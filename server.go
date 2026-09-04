package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type request struct {
	method  string
	uri     string
	header  map[string]string
	version string
	body    string
}

// declarar variáveis csonstantes que serão o host, a porta e o tipo de conexao (TCP)
const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

//criar um listener, via função net.Listen
//for para escutar as conexoes

func sendResponse(connection net.Conn, statuscode int, statusmessage string, contentType string, body string) error {

	bodyb := len(body)
	bodys := strconv.Itoa(bodyb)
	statuscodeString := strconv.Itoa(statuscode)
	message := "HTTP/1.1" + " " + statuscodeString + " " + statusmessage + "\r\n" + "Content-type: " + contentType + "\r\n" + "Content-Length: " + bodys + "\r\n" + "\r\n" + body
	bodybyte := []byte(message)
	_, err := connection.Write(bodybyte)
	if err != nil {
		log.Println("error sending answer")
		return err
	}

	return nil

}

func HandleConnection(connection net.Conn) {
	var Request request
	defer connection.Close()

	reader := bufio.NewReader(connection) //mesma coisa que o buffer fazia

	requestline, err := reader.ReadString('\n')
	if err != nil {
		log.Println("error while reading request line")
		return
	}

	requestline = strings.Trim(requestline, "\r\n")
	var fields []string
	fields = strings.Split(requestline, " ")

	if len(fields) != 3 {
		log.Println("error spliting string")
		err = sendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
		if err != nil {
			log.Println("Malformated request", err)
		}

		return
	}

	Request.method = fields[0]
	Request.uri = fields[1]
	Request.version = fields[2]

	var line string
	Request.header = make(map[string]string)

	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			log.Println("error while read headers", err)
			return
		}

		line = strings.Trim(line, "\r\n")

		if line == "" {
			break
		}

		var partes []string
		partes = strings.SplitN(line, ":", 2)

		if len(partes) != 2 {
			log.Println("Erros spliting parts")
			return
		}

		partes[1] = strings.Trim(partes[1], " ")

		Request.header[partes[0] /*chave*/] = partes[1] /*valor*/

	}

	if Request.method == "GET" || Request.method == "HEAD" {
		log.Println("Method:", Request.method)
		log.Println("Uri:", Request.uri)
		log.Println("Version: ", Request.version)
		log.Println("Headers:", Request.header)

		caminho := "files" + Request.uri
		if strings.Contains(Request.uri, "..") {
			err = sendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
			if err != nil {
				log.Println("path transversal")
				return
			}
			return
		} else {
			caminhob, err := os.ReadFile(caminho)
			caminhoext := filepath.Ext(caminho)
			content, existe := mimemaps[caminhoext]
			if existe == false {
				content = "application/octet-stream"
			}

			if err != nil {
				log.Println("not found archive", err)

				err = sendResponse(connection, 404, "not found", "text/html", "404 not found")
				if err != nil {
					log.Println("error requesting")
					return
				}
				return
			} else {
				caminhoString := string(caminhob)
				log.Println(caminhoString)
				err = sendResponse(connection, 200, "OK", content, caminhoString)
				if err != nil {
					log.Println("error sending archive")
				}
			}
		}

	} else if Request.method == "POST" {
		numbytes, err := strconv.Atoi(Request.header["Content-Length"])
		if err != nil {
			log.Println("Error while converting")
			//adicionar erro 400
			err = sendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
			if err != nil {
				log.Println("Error in POST request")
			}
			return
		} else {
			buffer := make([]byte, numbytes)
			_, err := io.ReadFull(reader, buffer)
			if err != nil {

				log.Println("error while take bytes", err)
				err = sendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
				if err != nil {
					log.Println("Bad Request")
				}
				return
			} else if numbytes != 0 {
				//Request.body = string(buffer)
				caminhoPost := "files" + Request.uri
				if Request.uri == "" || strings.Contains(Request.uri, "..") == true {
					err = sendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
					if err != nil {
						log.Println("error of Bad Request")
						return
					}
					return
				}
				err = os.WriteFile(caminhoPost, buffer, 0644)
				if err != nil {
					log.Println("error sending archive", err)
					err = sendResponse(connection, 500, "Internal Server Error", "text/html", "500 Internal Server Error")
					if err != nil {
						log.Println("error of Service Unavaliable")
						return
					}
					return
				} else {
					err = sendResponse(connection, 201, "Resource created", "text/html", "201 Resource Created")
					if err != nil {
						log.Println("error creating resource")
						return
					}
					return
				}
			} else if numbytes == 0 {
				err = sendResponse(connection, 400, "Bad Request", "text/html", "400 Bad Request")
				if err != nil {
					log.Println("Content-Length = 0")
					return
				}
				return
			}

		}

		log.Println("Method:", Request.method)
		log.Println("Uri:", Request.uri)
		log.Println("Version: ", Request.version)
		log.Println("Headers:", Request.header)
		log.Println("Body: ", Request.body)
	} else {
		err = sendResponse(connection, 501, "Not Implemented", "text/html", "501 Not implemented")
		if err != nil {
			log.Println("Not implemented request")
		}
		return

	}

}

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Print("Error", err)
		os.Exit(1)
	}
	defer listen.Close()

	for {
		connection, err := listen.Accept()
		if err != nil {
			log.Print("Error stablishing connection", err)
			continue
		}
		go HandleConnection(connection)

	}

}

//criar um handle connection
