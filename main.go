package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		s.Join("chat")
		server.BroadcastToRoom("", "chat", "reply", "A user has joined the room!")
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Printf("Server side log ---> (%s) %s\n", s.ID(), msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		server.BroadcastToRoom("", "chat", "reply", fmt.Sprintf("(%s) %s", s.ID(), msg))
		return "recv " + msg
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		server.BroadcastToRoom("", "chat", "reply", fmt.Sprintf("(%s) has left the room :(", s.ID()))
		log.Println("closed", reason)
	})

	server.OnError("/", func(s socketio.Conn, err error) {
		log.Println(err)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
