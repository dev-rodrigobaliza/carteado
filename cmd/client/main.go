package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dev-rodrigobaliza/carteado/domain/request"
	"github.com/dev-rodrigobaliza/carteado/domain/response"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn    *websocket.Conn
	sendBuf chan []byte
	message chan response.WSResponse
}

func (c *Client) Stop() {
	if c.conn != nil {
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
	}
}

func (c *Client) Write(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	c.sendBuf <- data

	return nil
}

func (c *Client) listen() {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.Stop()
			fmt.Printf("cannot read websocket message: %s\n", err.Error())
			break
		}

		var message response.WSResponse
		err = message.FromBytes(data)
		if err != nil {
			c.Stop()
			fmt.Printf("cannot decode websocket message: %s\n", err.Error())
			break
		}

		c.message <- message
	}
}

func (c *Client) listenListen() {
	requestID := uint64(1)
	var tableID string

	for {
		message := <-c.message
		fmt.Printf("message received: %v\n", message)

		if message.Status != "info" && message.Status != "success" {
			fmt.Println("wrong expected status")
			break
		}

		requestID++
		switch message.Message {
		case "welcome stranger":
			data := make(map[string]interface{})
			data["token"] = "v4.local.cEPMuKYqQlP7YfAO6Gey8ilVVWYHDYyC3OTG3LFEbw8JcGKg-p28KcJ2UbRKZEljtArFNepFT_DLIf-c2BXSiziJbJBxnGhA5Mcq_bYEPBo3apVghqHw3Zfk0HIw-3kIcKLy4cUFPGrEP54DIIx9TAwIHQ0"

			req := &request.WSRequest{
				RequestID: requestID,
				Service:   "auth",
				Resource:  "login",
			}
			req.Data = data

			c.sendBuf <- req.ToBytes()

		case "hello user":
			data := make(map[string]interface{})
			data["game_mode"] = "blackjack"
			data["min_players"] = 1
			data["max_players"] = 5
			data["allow_bots"] = true

			req := &request.WSRequest{
				RequestID: requestID,
				Service:   "table",
				Resource:  "create",
			}
			req.Data = data

			c.sendBuf <- req.ToBytes()

		case "table create":
			var ok bool
			tableID, ok = message.Data["table"].(map[string]interface{})["id"].(string)
			if !ok {
				fmt.Println("table id not processed")
				break
			}

			data := make(map[string]interface{})
			data["table_id"] = tableID
			data["action"] = "bot"
			data["quantity"] = 10

			req := &request.WSRequest{
				RequestID: requestID,
				Service:   "table",
				Resource:  "group",
			}
			req.Data = data

			c.sendBuf <- req.ToBytes()

		case "bot table group":
			time.Sleep(time.Minute * 1)
			data := make(map[string]interface{})
			data["table_id"] = tableID
			data["action"] = "start"

			req := &request.WSRequest{
				RequestID: requestID,
				Service:   "table",
				Resource:  "game",
			}
			req.Data = data

			c.sendBuf <- req.ToBytes()

		case "start table game":
			fmt.Println("game started")

		default:
			fmt.Println("message not processed")
		}
	}
}

func (c *Client) listenWrite() {
	for data := range c.sendBuf {
		err := c.conn.WriteMessage(
			websocket.TextMessage,
			data,
		)
		if err != nil {
			fmt.Printf("cannot listenWrite websocket message: %s\n", err.Error())
			break
		}
		fmt.Printf("send: %s\n", data)
	}
}

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8081",
		Path:   "ws",
	}
	urlStr := u.String()

	client := &Client{
		sendBuf: make(chan []byte, 1),
		message: make(chan response.WSResponse, 10),
	}

	ws, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		fmt.Printf("cannot connect to websocket: %s - [%s]\n", urlStr, err.Error())
		return
	}
	client.conn = ws

	go client.listen()
	go client.listenListen()
	go client.listenWrite()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	client.Stop()
	close(client.sendBuf)
}
