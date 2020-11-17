package gorilla

import (
	"fmt"
	"../../models"
	"flag"
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Runner ...
type Runner struct {
	Config		*models.Config
	Listeners 	map[string][]int
	Upgrader 	websocket.Upgrader
	Hub 		*Hub
}

// NewRunner ...
func NewRunner(config *models.Config) (*Runner, error) {	
	var upgrader = websocket.Upgrader{} // use default options
	upgrader.CheckOrigin = func(r *http.Request) bool { 
		return true 
	}

	hub:= newHub()

	runner := &Runner{
		Config:   	config,
		Upgrader:	upgrader,
		Hub: hub,
	}
	return runner, nil
}

// Start ...
func (srv *Runner) Start() {
	
	var addr = flag.String("addr", fmt.Sprintf("%s:%s",srv.Config.Service.Host,srv.Config.Service.Port), "http service address")
	
	flag.Parse()

	go srv.Hub.run()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/sdk/js", srv.serveSdk)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(srv.Hub, w, r)
	})

	if srv.Config.Service.Mode=="prod"{
		err:= http.ListenAndServeTLS(*addr, "/app/certs/fullchain.pem", "/app/certs/privkey.pem", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}else{
		err:= http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
	
	
}
// func (srv *Runner) serveHome(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
	
// 	http.ServeFile(w, r, "./ui/dist/index.html")
// }

func (srv *Runner) serveSdk(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "sdk.js")
}
