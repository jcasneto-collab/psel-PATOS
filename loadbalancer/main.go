package main

import (
	"bufio"
	"io"
	"load-balancer/internal/httpcore"
	"log"
	"net"
	"strconv"
)

func forwardbody (reader *bufio.Reader, connectionBackend net.Conn, contentLength string) error{
	numbytes, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Println("Error while converting")
		returncd ..
	}
	buffer := make([]byte,numbytes)
	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		log.Println("error while take bytes")
		err = httpcore.SendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
			if err != nil {
				log.Println("Error of Bad Requesting")
				return
			}
	}
	_, err = connectionBackend.Write(buffer)
	if err != nil {
		log.Println("error sending resposne response of body")
		err = httpcore.ParseRequest(connection, 400, "Bad Request", "text/html", "400 Bad Request")
		if err != nil {
			log.Println("error converting")
			return
		}
	}
	_, err = io.Copy(connection, connectionBackend)
	if err != nil{
		log.Println("error sending response")
		return
	}
}


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
		numbytes, err := strconv.Atoi(Request.Header["Content-Length"])
		if err != nil {
			log.Println("Error while converting")
			err = httpcore.SendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
			if err != nil {
				log.Println("error converting", err)
				return
		}
		buffer := make([]byte, numbytes)
		_, err = io.ReadFull(reader, buffer)
		if err != nil {

			log.Println("error while take bytes", err)
			err = httpcore.SendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
			if err != nil {
				log.Println("error of Bad Request")
				return
			}
			
		}

	}
		_, err = connectionBackend.Write(buffer)
		if err != nil {
			log.Println("error sending response of body", err)
			err = httpcore.SendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
			if err != nil {
				log.Println("error converting", err)
				return
			}

		}	
	

	_, err = io.Copy(connection, connectionBackend)
	if err != nil {
		log.Println("error sending response")
		return
	}
}
