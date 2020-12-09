# EventTube

EventTube is a simple real-time and bidirectional messaging server between FrontEnd and BackEnd.

## Running from source

This way requires a working Go development environment.
The [GettingStarted](http://golang.org/doc/install) page describes how to install the
development environment.

Once you have Go up and running, you can download, build and run the example
using the following commands.

    $ git clone git@github.com:colber/eventtube-server.git your_dir
    $ cd your_dir
    $ go get github.com/gorilla/websocket
    $ go get github.com/prometheus/client_golang/prometheus
    $ go get github.com/prometheus/client_golang/prometheus/promauto
    $ go get github.com/prometheus/client_golang/prometheus/promhttp
    $ go run main.go

## Running in the Docker
We recommend using Docker to run EventTube:

    $ docker pull ptimofeev/eventtube:latest
    $ docker run --name eventtube --rm -p 9000:9000 ptimofeev/eventtube


## Server
By default the server start on the `localhost:9000`
The server has got slim SDK (6Kb) for work with him. 
Also the server has got thin JS SDK (6 Kb) for work with him from the FrontEnd. 
You can find it on the path `http://localhost:9000/sdk/js`


## On the FrontEnd

The first include SDK:
In the header of the yor web page add following commands

    $ <script src="http://localhost:9000/sdk/js" async onload="initEventTube()"></script>

The seccond connect to the server:

    $ function initEventTube(){
    $ var options={
    $     connection:{
    $     host:'localhost',
    $     port:'9000'
    $     }
    $ }
    $ var eventTube=new EventTube(options);
    $ window.EventTube=eventTube;
    $ window.EventTube.connect();
}
