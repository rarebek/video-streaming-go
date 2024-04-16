package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/", nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	outputFile, err := os.Create("output.mp4")
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outputFile.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		log.Println("Received termination signal. Closing connection and exiting...")
		conn.Close()
		os.Exit(0)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from WebSocket:", err)
			break
		}
		_, err = outputFile.Write(message)
		if err != nil {
			log.Println("Error writing to output file:", err)
			break
		}
	}

	fmt.Println("Video streaming complete. Output saved to output.mp4")
}
