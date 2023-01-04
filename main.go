package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(currentConn *websocket.Conn, conns *[]*websocket.Conn) {
	for {
		messageType, message, err := currentConn.ReadMessage()
		if err != nil {
			log.Println("ERROR: ", err)

			if strings.Contains(err.Error(), "websocket: close") {
				fmt.Println("closing connection...")
				var connsWithoutSelf []*websocket.Conn
				for _, conn := range *conns {
					if conn == currentConn {
						continue
					}
					connsWithoutSelf = append(connsWithoutSelf, conn)
				}
				*conns = connsWithoutSelf
				return
			}

		}

		for _, conn := range *conns {

			if err := conn.WriteMessage(messageType, message); err != nil {
				continue
			}

		}
	}
}

func setupRoutes() {
	var connections []*websocket.Conn

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}

		connections = append(connections, ws)
		fmt.Printf("ini adalah connections %+v \n", connections)

		go reader(ws, &connections)
	})
}

func main() {
	fmt.Println("hello world!")

	setupRoutes()
	// change this
	log.Fatal(http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), nil))
}
