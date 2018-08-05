package main

import (
  "fmt"
  "github.com/go-redis/redis"
  "github.com/gorilla/websocket"
  "net/http"
  "strconv"
)

// TODO: Pull the port from environment or something similar
var addr = "localhost:8081"
var upgrader = websocket.Upgrader{}
var updateChannel = make(chan *redis.Message)
var client *redis.Client

func receiveEvents(c *websocket.Conn) {
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err)
			return
		}

    fmt.Println("mt: ", mt)
		fmt.Println("recv: ", message)

    // TODO: Implement the redis command set here. Basically just pass through the
    // command and return the result to the client
  }
}

func writeEvents(c *websocket.Conn) {
	defer c.Close()
	for {
		message := <-updateChannel
    output := "{" + strconv.Quote("channel") + ":" + strconv.Quote(message.Channel) + ", " + strconv.Quote("payload") + ": " + message.Payload + "}"
		err := c.WriteMessage(websocket.TextMessage, []byte(output))
		if err != nil {
			fmt.Println("error: ", err)
			return
		}
	}
}

func commands(w http.ResponseWriter, r *http.Request) {

	// connect the websocket
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade: ", err)
		return
	}

	go receiveEvents(c)
	writeEvents(c)
}

func main() {
  fmt.Println("Starting socket-redis...")

  // TODO: make the url here configurable
  http.HandleFunc("/commands", commands)
	fmt.Println("starting socket listener...")
  go http.ListenAndServe(addr, nil)

  // TODO: also configure what port to use and IP to bind to
  client = redis.NewClient(&redis.Options{
  		Addr:     "localhost:6379",
  		DB:       0,
  	})

  pong, err := client.Ping().Result()
  fmt.Println(pong, err)

  // TODO: Make the client actualy subscribe instead of doing this automatically
  pubsub := client.PSubscribe("*")
  defer pubsub.Close()

  for {
    msg, err := pubsub.ReceiveMessage()
    if err != nil {
      panic(err)
    }
    fmt.Println("Message: ", msg)
    select {
    case updateChannel <- msg:
    default:
    }
  }
}
