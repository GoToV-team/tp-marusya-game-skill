package game

import (
	"encoding/json"
	"github.com/evrone/go-clean-template/pkg/marusia"
)

type ClosedClient interface {
	Close()
}

type HubResult struct {
	Result
	WorkedClient ClosedClient
}

type Client struct {
	hub      *ScriptHub
	op       Director
	clientId string
}

func NewClient(hub *ScriptHub, sessionId string, op Director) *Client {
	return &Client{
		hub:      hub,
		clientId: sessionId,
		op:       op,
	}
}

func (c *Client) Close() {
	c.hub.UnregisterClient(c)
}

func (c *Client) PlayScene(msg sceneMessage) Result {
	res := c.op.PlayScene(SceneRequest{
		Command:      msg.rq.command,
		FullUserText: msg.rq.fullUserText,
		WasButton:    msg.rq.wasButton,
		Payload:      msg.rq.payload,
		Info: UserInfo{
			UserId:    msg.userId,
			SessionId: msg.sessionId,
		},
	})

	return res
}

type ScriptHub struct {
	Clients    map[string]*Client
	broadcast  chan *sceneMessage
	register   chan *Client
	unregister chan *Client
	stopHub    chan bool
}

type request struct {
	command      string
	fullUserText string
	payload      json.RawMessage
	wasButton    bool
}

func fromMarusiaRequest(rqm marusia.RequestIn) *request {
	return &request{
		command:      rqm.Command,
		fullUserText: rqm.OriginalUtterance,
		payload:      rqm.Payload,
		wasButton:    rqm.Type == marusia.ButtonPressed,
	}
}

type sceneMessage struct {
	userId    string
	sessionId string
	answer    chan HubResult
	rq        request
}

func NewHub() *ScriptHub {
	return &ScriptHub{
		broadcast:  make(chan *sceneMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
		stopHub:    make(chan bool),
	}
}

func (h *ScriptHub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *ScriptHub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *ScriptHub) RunScene(rq marusia.Request) chan HubResult {
	answer := make(chan HubResult)
	h.broadcast <- &(sceneMessage{
		userId:    rq.Session.UserID,
		sessionId: rq.Session.SessionID,
		rq:        *fromMarusiaRequest(rq.Request),
		answer:    answer,
	})
	return answer
}

func (h *ScriptHub) StopHub() {
	h.stopHub <- true
}

func (h *ScriptHub) unregisterAll() {
	for key, _ := range h.Clients {
		delete(h.Clients, key)
	}
}

func (h *ScriptHub) runScene(msg *sceneMessage) {
	if client, ok := h.Clients[msg.sessionId]; ok {
		go func(ans chan HubResult, client *Client) {
			ans <- HubResult{
				Result:       client.PlayScene(*msg),
				WorkedClient: client,
			}
		}(msg.answer, client)
	}
}

func (h *ScriptHub) unregisterClient(client *Client) {
	if _, ok := h.Clients[client.clientId]; ok {
		delete(h.Clients, client.clientId)
	}
}

func (h *ScriptHub) Run() {
	for {
		select {
		case client, ok := <-h.register:
			if ok {
				h.Clients[client.clientId] = client
			}
			break
		case client, ok := <-h.unregister:
			if ok {
				h.unregisterClient(client)
			}
			break
		case msg, ok := <-h.broadcast:
			if ok {
				h.runScene(msg)
			}
			break
		case <-h.stopHub:
			h.unregisterAll()
			return
		default:
			break
		}
	}
}
