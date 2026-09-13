package main

import (
	"bufio"
	"io"
	"load-balancer/internal/httpcore"
	"log"
	"net"
	"os"
	"strconv"
)

func forwardbody(reader *bufio.Reader, connectionBackend net.Conn, contentLength string) error {
	numbytes, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Println("Error while converting")
		return err
	}
	buffer := make([]byte, numbytes)
	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		log.Println("error while take bytes")
		return err
	}
	_, err = connectionBackend.Write(buffer)
	if err != nil {
		log.Println("error sending resposne response of body")
		return err
	}
	return nil
}

func forwardToBackend(Request httpcore.Request, connection net.Conn, reader *bufio.Reader) {
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

	defer connectionBackend.Close()

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

	if Request.Method == "POST" {
		err = forwardbody(reader, connectionBackend, Request.Header["Content-Length"])
		if err != nil {
			log.Println("Error sending response")
			err = httpcore.SendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
			if err != nil {
				log.Println("Error sending response")
				return
			}
		}

	}

	_, err = io.Copy(connection, connectionBackend)
	if err != nil {
		log.Println("error sending response")
		return
	}
}

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Println("Error", err)
		os.Exit(1)
	}
	defer listen.Close()

	for {
		connection, err := listen.Accept()
		if err != nil {
			log.Println("Error stablishing connection", err)
			continue
		}
		go httpcore.Accepting_con(connection, forwardToBackend)
	}

}
