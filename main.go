package main

import (
  "fmt"
  "github.com/go-redis/redis"
  "github.com/gorilla/websocket"
  "net/http"
  "strconv"
  "flag"
  "os"
)

var upgrader = websocket.Upgrader{}
var updateChannel = make(chan *redis.Message)
var client *redis.Client

var flags = CaptureFlags()

type Flags struct {
  RedisHost string
  RedisPort string
  SocketHost string
  SocketPort string
  SocketPath string
  Channels string
  LogLevel string
}

func Getenv(variable string, defaultValue string) string {
  val, ok := os.LookupEnv(variable)
  if ok {
    return val
  } else {
    return defaultValue
  }
}

func CaptureFlags() *Flags {
  var flags Flags

  redisHost := flag.String("redisHost", Getenv("REDIS_HOST", "localhost"), "Redis host name")
  redisPort := flag.String("redisPort", Getenv("REDIS_PORT", "6379"), "Redis port")
  socketHost := flag.String("socketHost", Getenv("SOCKET_HOST", "localhost"), "Websocket host name")
  socketPort := flag.String("socketPort", Getenv("SOCKET_PORT", "8081"), "Websocket port")
  socketPath := flag.String("socketPath", Getenv("SOCKET_PATH", "/commands"), "Websocket http path")
  channels := flag.String("channels", Getenv("REDIS_CHANNELS", "*"), "Channels to subscribe to")
  logLevel := flag.String("logLevel", "error", "Log level: debug or error")
  flag.Parse()

  flags.RedisHost = *redisHost
  flags.RedisPort = *redisPort
  flags.SocketHost = *socketHost
  flags.SocketPort = *socketPort
  flags.SocketPath = *socketPath
  flags.Channels = *channels
  flags.LogLevel = *logLevel
  return &flags
}

func receiveEvents(c *websocket.Conn) {
	defer c.Close()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err)
			return
		}
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

  http.HandleFunc(flags.SocketPath, commands)
	fmt.Println("starting socket listener...")
  go http.ListenAndServe(flags.SocketHost + ":" + flags.SocketPort, nil)

  client = redis.NewClient(&redis.Options{
  		Addr:     flags.RedisHost + ":" + flags.RedisPort,
  		DB:       0,
  	})

  pong, err := client.Ping().Result()
  fmt.Println(pong, err)

  pubsub := client.PSubscribe(flags.Channels)
  defer pubsub.Close()

  for {
    msg, err := pubsub.ReceiveMessage()
    if err != nil {
      panic(err)
    }
    if flags.LogLevel == "debug" {
      fmt.Println("Message: ", msg)
    }
    select {
    case updateChannel <- msg:
    default:
    }
  }
}
