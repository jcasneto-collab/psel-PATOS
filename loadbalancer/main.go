package main

import (
	"bufio"
	"load-balancer/internal/httpcore"
	"log"
	"net"
)

func forwardToBackend(connection net.Conn, Request httpcore.Request, reader *bufio.Reader) {
	connectionBackend, err := net.Dial("tcp", "localhost:8001")
	if err != nil {
		log.Println("error sending connection to backend")
		err = httpcore.SendResponse(connection, 500, "Internal Server Error", "text/html", "500 Internal Server Error")
		if err != nil {
			log.Println("Error of server")
			return
		}
		return
	}

	message := Request.Method + " " + Request.Uri + " " + Request.Version + "\r\n"

	for chave, valor := range Request.Header {
		message += chave + ": " + valor + "\r\n"

	}

	message += "\r\n"

	messageb := []byte(message)

	_, err = connectionBackend.Write(messageb)
	if err != nil {
		log.Println("Error sending answer")
		return
	}

}
