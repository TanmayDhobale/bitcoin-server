package server

import (
	"log"
	"sync"
	"github.com/gorilla/websocket"
	"encoding/json"
	"net/http"
)

type CentralServer struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

func NewCentralServer() *CentralServer {
	return &CentralServer{
			clients:    make(map[*websocket.Conn]bool),
			broadcast:  make(chan []byte),
			register:   make(chan *websocket.Conn),
			unregister: make(chan *websocket.Conn),
	}
}

func (s *CentralServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client] = true
			s.mutex.Unlock()
			log.Println("New client connected")
		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				client.Close()
			}
			s.mutex.Unlock()
			log.Println("Client disconnected")
		case message := <-s.broadcast:
			s.mutex.Lock()
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error broadcasting message: %v", err)
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mutex.Unlock()
		}
	}
}

func (s *CentralServer) HandleConnection(conn *websocket.Conn) {
	s.register <- conn
	defer func() {
		s.unregister <- conn
		conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		s.broadcast <- message
	}
}

func (s *CentralServer) BroadcastUpdate(updateType string, data interface{}) {
	update := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: updateType,
		Data: data,
	}

	jsonData, err := json.Marshal(update)
	if err != nil {
		log.Printf("Error marshaling update: %v", err)
		return
	}

	s.broadcast <- jsonData
}


func (s *CentralServer) HandleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := generateNewWalletAddress()

	response := struct {
		Address string `json:"address"`
	}{
		Address: address,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}


func generateNewWalletAddress() string {
			return "dummy_wallet_address"
}

func (s *CentralServer) HandleGetWalletBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")

	balance := getWalletBalance(address)

	response := struct {
		Balance int `json:"balance"`
	}{
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func getWalletBalance(address string) int {
	return 0
}

func (s *CentralServer) HandleGetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")

	history := getTransactionHistory(address)

	response := struct {
		History []string `json:"history"`
	}{
		History: history,
	}
}

func getTransactionHistory(address string) []string {
	return []string{}
}

func (s *CentralServer) HandleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	blockchain := getBlockchain()


	response := struct {
		Blockchain []string `json:"blockchain"`
	}{
		Blockchain: blockchain,
	}
}

func getBlockchain() []string {
	// TODO: Implement this function to get the blockchain
	return []string{}	
}				

func (s *CentralServer) HandleGetPendingTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	transactions := getPendingTransactions()

	response := struct {
		Transactions []string `json:"transactions"`
	}{
		Transactions: transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func getPendingTransactions() []string {
	// TODO: Implement this function to get the pending transactions
	return []string{}
}

func (s *CentralServer) HandleGetMinerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := getMinerStatus()

	response := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
	

func getMinerStatus() string {



	return "dummy_miner_status"
}

func (s *CentralServer) HandleGetMinerStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := getMinerStats()

	response := struct {
		Stats []string `json:"stats"`
	}{
		Stats: stats,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func getMinerStats() []string {
      
	return []string{}
}

func (s *CentralServer) HandleGetMinerPerformance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	performance := getMinerPerformance()

	response := struct {
		Performance []string `json:"performance"`
	}{
		Performance: performance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}


