// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gorilla

import (
	"../../models"
	"log"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client
	lastClientId int32

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Subscribe client to the event.
	subscribe chan *Subscription

	// Unsubscribe client from the event.
	unsubscribe chan *Subscription

	// publish event for subscribes. 
	publish chan *models.Message

	// Registered event listeners.
	listeners map[string]map[string]*Client
}

func newHub() *Hub {
	return &Hub{
		clients:   		make(map[string]*Client),
		register:   	make(chan *Client),
		unregister: 	make(chan *Client),
		subscribe:   	make(chan *Subscription),
		unsubscribe:	make(chan *Subscription),
		publish: 		make(chan *models.Message),
		broadcast:  	make(chan []byte),
		listeners: 		make(map[string]map[string]*Client),
	}
}

func (h *Hub) nextClientId() string {
	clientId:=fmt.Sprintf("%d",h.lastClientId)
	h.lastClientId++
	return clientId
}


func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			
			h.clients[client.id] = client
			message:=&models.Message{
				Event:"CLIENT_CONNECTED",
				Data:client.id,
			}
			client.send <- message.Encode()
			log.Println("Client Id:",client.id,"registred")

		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				
				for event,_ := range h.listeners{
					h.UnSub(event,client.id)
				}

				// удаляю клиента из общего списка
				delete(h.clients, client.id)
				close(client.send)

				log.Println("Client Id:",client.id,"unregistred")
				//h.DebugState()
			}else{
				log.Println("Client Id:",client.id,"d'nt registred")
			}
		case  subscribe := <-h.subscribe:
			h.Sub(subscribe)
		case  subscribe := <-h.unsubscribe:
		 	h.UnSub(subscribe.Event,subscribe.Client.id)

		case  message := <-h.publish:
			fmt.Println("event:",message.Data)
			for _,client := range h.listeners[message.Event] {
				select {
					case client.send <- message.Encode():
					default:
						close(client.send)
						delete(h.clients, client.id)
				}
			
			}
			
		case message := <-h.broadcast:
			// log.Println("broadcast")
			for clientId := range h.clients {
				client:=h.clients[clientId]
				select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, clientId)
				}
			}
		}
	}
}
func (h *Hub) Sub(s *Subscription){

	log.Println("Sub:",s.Event,s.Client.id)
	if client, ok := h.clients[s.Client.id]; ok {
		if _, ok := h.listeners[s.Event]; ok {
			if _, ok := h.listeners[s.Event][s.Client.id]; ok {
				log.Println("Client ID",s.Client.id,"was sbscribed to ",s.Event," early")
			}else{
				h.listeners[s.Event][s.Client.id]=client
				log.Println("Client ID",s.Client.id,"sbscribed to ",s.Event)
			}

		}else{
			listeners:=make(map[string]*Client)
			listeners[client.id]=client

			h.listeners[s.Event]=listeners

			//это первая подписка на событие,
			//подписываюсь на одноименное событие в NATS 
			// sub,err:=h.mq.Subscribe(r.Event,func(event,data string) {
			// 	//log.Println(event,data)
			// 	m:=&models.Message{
			// 		Event:event,
			// 		Data:data,
			// 	}
			// 	h.publish <- m
			// 	return
			// })
			// if err!=nil{
			// 	log.Println("Subscribe to NATS event ",r.Event," has err",err)
			// }else{
			// 	h.mqSubscribes[r.Event]=sub
			// }
		}
	}else{
		log.Println("Client ID",s.Client.id,"d'nt registred")
	}
}

func (h *Hub) UnSub(event,clientId string){

	// удаляю клиента из списка подписчиков
	delete(h.listeners[event], clientId);

	log.Println("Client Id:",clientId,"unsubscribed of event ",event)

	//если у события больше нет подписчиков
	//отписываюсь от NATS
	if len(h.listeners[event])==0{
		delete(h.listeners, event);

		// sub:=h.mqSubscribes[eventId]
		// err:=sub.Unsubscribe()
		// if err!=nil{
		// 	log.Println("Unsubscribe of NATS event",eventId,"has err:",err)
		// }
		// delete(h.mqSubscribes, eventId);
	}
}

