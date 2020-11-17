class EventTube {
    constructor(options){
        this.options = options;
        
        // Ограничение по числу подписок
        if(options.subLimit){
            this.subLimit = options.subLimit;
        }else{
            this.subLimit = 0;
        }
        
        this.listeners = [];        // зарегистрирванные подписчики в формате [event]:[handler1,handler2,..]}
        this.size = 0;              // суммарное число подписчиков (handle-ов)
        this.ws = null;             // экземпляр WS
        this.retryInterval = 0;     // стартовый интервал перед повторным подключением
        this.retryInc =  1000;      // с каждой неудачной попыткой переподключения +1000

        // Подключение к WS сервису
        this.connect = function () {
            var self=this;
            if (this.ws) {
                return false
            }
            
            var cnn=this.options.connection
            var uri='ws://'+cnn.host+':'+cnn.port+'/ws';
            
            this.ws = new WebSocket(uri);
            this.ws.onopen = function(e) {
                self.retryInterval=0; // сбрасываю интервал задержки перед повторым подключением
            }

            this.ws.onclose = function(e) {
                self.ws = null;
                // Автоматическое переподключение
                setTimeout(function(){
                    self.retryInterval=self.retryInterval+self.retryInc;
                    self.connect();
                },self.retryInterval);
            }

            this.ws.onmessage = function(evt) {
                
                var events=evt.data.split('\n');
                
                events.forEach(function(data){
                    var resp=JSON.parse(data);
                    if(resp.event=='CLIENT_CONNECTED'){
                        self.ws.clientId=resp.data;
                        
                        // После успешного подключения/переподключения
                        // если есть подписки, подписываюсь на них заново с новым clientId
                        Object.keys(self.listeners).forEach(function(event){
                            // Отправляю сообщение на сервер
                            var message={
                                event:event,
                                type:'subscribe',
                            }
                            // Отправляю сообщение на сервер
                            self.ws.send(JSON.stringify(message));
                        });
                        
                        
                    }else{
                        // Если пришло серверное событие
                        var listeners=self.listeners[resp.event]
                        if(listeners && listeners.handlers){
                            Object.keys(listeners.handlers).forEach(function(handler){
                                listeners.handlers[handler].call(this,JSON.parse(data));
                            });
                        }
                    }
                })
            }
            this.ws.onerror = function(e) {
                self.ws = null;
                // Автоматическое переподключение
                setTimeout(function(){
                    self.retryInterval=self.retryInterval+self.retryInc;
                    self.connect();
                },self.retryInterval);
            }
            return false;
        }

        // Подписка на событие
        this.sub = function(event,cb){
            if (!this.ws) {
                this.connect();
            }
            
            var self=this;
            return new Promise(function (resolve, reject) {
                if(self.subLimit==0 || self.size<self.subLimit){
                    if(!self.listeners[event]){
                        self.listeners[event]={
                            handlers:[],
                        }
                    }

                    // Формирую уникальный ID обработчика
                    var subId=Date.now();
                    
    
                    // Регистрирую подписку на клиенте
                    self.listeners[event].handlers[subId]=cb;
    
                    // Увеличиваю общее кол-во подписчиков на все события
                    self.size++;
    
                    // Считаю кол-во подписчиков на событие
                    // Если это первый подписчик на событие, регистрирую подписку на сервере
                    var eventListenerCount=Object.keys(self.listeners[event].handlers).length
                    if(eventListenerCount==1 && self.ws.clientId){
                        
                        // Отправляю сообщение на сервер
                        var message={
                            event:event,
                            type:'subscribe',
                        }
                        if (!self.ws) {
                            var err=new Error('ws is closed')
                            reject(err);
                        }
                        // Отправляю сообщение на сервер
                        self.ws.send(JSON.stringify(message));
                    }
    
                    // Возвращаю ID подписчика
                    resolve(subId);
                }else{
                    const err = new Error('превышен лимит подписки');
                    reject(err);
                }
            });
        }

        // Отписка от события
        this.unSub = function(eventName,subId){
            if (!this.ws) {
                this.connect();
            }
            var self=this;

            return new Promise(function (resolve, reject) {
                if(subId){
                    if(self.listeners[eventName] && self.listeners[eventName].handlers[subId]){
    
                        // Отписываюсь локально
                        var removed=self.listeners[eventName].handlers[subId];
                        delete self.listeners[eventName].handlers[subId]
                        
                        // Уменьшаю общее число подписчиков
                        self.size--;
        
                        // Считаю текущее кол-во подписчиков на событие
                        // Если локальных подписчиков больше нет, удаляю подписку на сервере
                        var eventListenerCount=Object.keys(self.listeners[eventName].handlers).length
                        if(eventListenerCount==0){
                            var message={
                                event:eventName,
                                type:'unsubscribe',
                            }
                            // Отправляю сообщение на сервер
                            self.ws.send(JSON.stringify(message));
                        }
    
                        // Возвращаю удаленный объект
                        resolve(removed);
                    }else{
                        const err = new Error('Incorrect subId');
                        reject(err);
                    }
                }else if(self.listeners[eventName]){
                    var message={
                        event:eventName,
                        type:'unsubscribe',
                    }
                    // Отправляю сообщение на сервер
                    self.ws.send(JSON.stringify(message));

                    
                    var removed=self.listeners[eventName].handlers
                    var rmSize=removed.length;
                    delete self.listeners[eventName];
                    self.size=self.size-rmSize;
                    
                    // Возвращаю удаленный объект
                    resolve(removed);
                }else{
                    const err = new Error('Incorrect eventName');
                    reject(err);
                }
            });
            
            
        }

        // Публикация события
        this.pub = function(event,data){
            if (!this.ws) {
                this.connect();
            }
            
            // Формирую сообщение
            var message={
                event:event,
                data:data,
                type:'publish',
            }
            if (!this.ws) {
                console.log('disconnected')
                return false;
            }
            // Отправляю сообщение на сервер
            this.ws.send(JSON.stringify(message));
            
        }
        
    }
}