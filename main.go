package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		file, err := os.Open("video.mp4")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		bufferSize := 1024
		buffer := make([]byte, bufferSize)
		for {
			bytesRead, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Println(err)
				}
				break
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buffer[:bytesRead])
			if err != nil {
				log.Println(err)
				break
			}
		}

		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println(err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
