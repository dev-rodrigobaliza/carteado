package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"

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
			fmt.Printf("Cannot read websocket message: %s\n", err.Error())
			break
		}

		var message response.WSResponse
		err = message.FromBytes(data)
		if err != nil {
			c.Stop()
			fmt.Printf("Cannot decode websocket message: %s\n", err.Error())
			break
		}

		c.message <- message
	}
}

func (c *Client) listenListen() {
	requestID := uint64(1)

	for {
		message := <-c.message
		fmt.Printf("message received: %v\n", message)
		requestID++
		switch message.Message {
		case "welcome player":
			data := make(map[string]interface{})
			data["token"] = "v4.local.TqajZDVRzA-wGa2yhzzfnNErivkuofatcxSGr1RyniA1P06m0NhbZUC9n3BfIu4eSgxupyDDgZ2a447sbM8D3dtvGdpiyioSmd2VjYveBbAZu8SXo9ODI1YiYA4jgBVxZvebKNl9DOptvZ_lKwi19xvrtp_UdW0"

			req := &request.WSRequest{
				RequestID: requestID,
				Service:   "auth",
				Resource:  "login",
			}
			req.Data = data

			c.sendBuf <- req.ToBytes()

		case "authenticated":
			data := make(map[string]interface{})
			data["game_type"] = "blackjack"
			data["min_players"] = 2
			data["max_players"] = 5
			data["allow_bots"] = true
			data["secret"] = "secret"

			req := &request.WSRequest{
				RequestID: requestID,
				Service:   "game",
				Resource:  "create",
			}
			req.Data = data

			c.sendBuf <- req.ToBytes()

		case "game created":
			gameID, ok := message.Data["game_id"].(string)
			if !ok {
				fmt.Printf("message received and not processed: %v\n", message)
				continue
			}

			data := make(map[string]interface{})
			data["game_id"] = gameID

			req := &request.WSRequest{
				RequestID: requestID,
				Service:   "game",
				Resource:  "status",
			}
			req.Data = data

			c.sendBuf <- req.ToBytes()

		// case "game status", "game removed":
		// 	gameID, ok := message.Data["game_id"].(string)
		// 	if !ok {
		// 		fmt.Printf("message received and not processed: %v\n", message)
		// 		continue
		// 	}

		// 	data := make(map[string]interface{})
		// 	data["game_id"] = gameID

		// 	req := &request.WSRequest{
		// 		RequestID: requestID,
		// 		Service:   "game",
		// 		Resource:  "remove",
		// 	}
		// 	req.Data = data

		// 	c.sendBuf <- req.ToBytes()
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
		fmt.Printf("Cannot connect to websocket: %s - [%s]\n", urlStr, err.Error())
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
