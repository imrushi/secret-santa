package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a single user connection
type Client struct {
	Conn   *websocket.Conn
	Name   string
	Avatar string
	Room   *Room
}

// Room holds the state of a Secret Santa room
type Room struct {
	ID         string
	Clients    map[*Client]bool
	Host       *Client
	IsStarted  bool
	Mutex      sync.Mutex
	Broadcast  chan []byte
	Unregister chan *Client
}

// Message is the standard JSON structure for communication
type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// Participant is used for sending the user list to frontend
type Participant struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	IsHost bool   `json:"isHost"`
}

// A map to hold all active rioms: RoomID -> Room
var rooms = make(map[string]*Room)
var roomsMutex sync.Mutex

func main() {
	mux := http.NewServeMux()

	// The single WebSocket endpoint handling Create, Join and Game Loop
	mux.HandleFunc("/ws", handleWebSocket)
	// Start the server
	fmt.Println("🎅 Santa Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func newRoom(id string) *Room {
	return &Room{
		ID:         id,
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Unregister: make(chan *Client),
	}
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	roomID := query.Get("room")
	name := query.Get("name")
	avatar := query.Get("avatar")
	action := query.Get("action") // "create" or "join"

	if name == "" || roomID == "" {
		http.Error(w, "Missing name or room ID", http.StatusBadRequest)
		return
	}

	// upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}

	client := &Client{Conn: conn, Name: name, Avatar: avatar}

	// Handle Room logic
	roomsMutex.Lock()
	room, exists := rooms[roomID]
	if action == "create" {
		if exists {
			// Room already exists, cannot overwrite
			roomsMutex.Unlock()
			sendError(client, "Room already exists")
			conn.Close()
			return
		}
		// Create new room
		room = newRoom(roomID)
		rooms[roomID] = room
		room.Host = client // First User is the host
		go room.run()      // Start rooms event loop
	} else { // join
		if !exists {
			roomsMutex.Unlock()
			sendError(client, "Room does not exist")
			conn.Close()
			return
		}
		if room.IsStarted {
			roomsMutex.Unlock()
			sendError(client, "Game already started")
			conn.Close()
			return
		}
	}

	// Add client to room
	room.Clients[client] = true
	client.Room = room
	roomsMutex.Unlock()

	// Listen for message (Blocking loop)
	// Notify everyone about new participant
	broadcastUserList(room)

	defer func() {
		room.Unregister <- client
		conn.Close()
	}()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			break // Disconnected
		}

		if msg.Type == "START_GAME" {
			// Only host can start the game
			if client == room.Host {
				performDraw(room)
			}
		}
	}
}

// run handles the concurrent broadcasting for the room
func (r *Room) run() {
	for {
		select {
		case client := <-r.Unregister:
			// User disconnected
			r.Mutex.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				// If host leaves, room logic handles it assign new host or close room
				if client == r.Host {
					// Assign new host if possible
					for c := range r.Clients {
						r.Host = c
						break
					}
				} else if len(r.Clients) == 0 {
					// Clean up empty room
					roomsMutex.Lock()
					delete(rooms, r.ID)
					roomsMutex.Unlock()
					r.Mutex.Unlock()
					return
				}
			}
			r.Mutex.Unlock()
			broadcastUserList(r)
		case message := <-r.Broadcast:
			// Send message to all clients
			r.Mutex.Lock()
			for client := range r.Clients {
				client.Conn.WriteMessage(websocket.TextMessage, message)
			}
			r.Mutex.Unlock()
		}
	}
}

func sendError(c *Client, errorMsg string) {
	c.Conn.WriteJSON(Message{Type: "ERROR", Payload: errorMsg})
}

func broadcastUserList(r *Room) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	participants := []Participant{}
	for client := range r.Clients {
		participants = append(participants, Participant{
			Name:   client.Name,
			IsHost: (client == r.Host),
			Avatar: client.Avatar,
		})
	}

	msg := Message{Type: "UPDATE_PARTICIPANTS", Payload: participants}
	data, _ := json.Marshal(msg)

	go func() {
		r.Broadcast <- data
	}()
}

func performDraw(r *Room) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	if len(r.Clients) < 2 {
		// Not enough participants
		return
	}
	r.IsStarted = true

	// Extract Clients into a slice
	clients := make([]*Client, 0, len(r.Clients))
	for c := range r.Clients {
		clients = append(clients, c)
	}

	// Shuffle (Fisher-Yates)
	rs := rand.New(rand.NewSource(time.Now().UnixNano()))
	rs.Shuffle(len(clients), func(i, j int) {
		clients[i], clients[j] = clients[j], clients[i]
	})

	// Assign Matches (Circular Linked List approach)
	// A-->B, B->C, ..., N->A
	for i := 0; i < len(clients); i++ {
		giver := clients[i]
		receiver := clients[(i+1)%len(clients)]

		// Send PRIVATE message to giver
		// We do NOT broadcast the whole list
		resultMsg := Message{
			Type:    "MATCH_RESULT",
			Payload: receiver.Name, // The name they have to gift
		}
		giver.Conn.WriteJSON(resultMsg)
	}
}
