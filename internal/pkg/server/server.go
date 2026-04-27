package server

import (
	"encoding/csv"
	"encoding/json"
	"faro/internal/pkg/storage"
	"faro/internal/pkg/types"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					h.unregister <- client
				}
			}
			h.mu.Unlock()
		}
	}
}

// Server handles HTTP and WebSocket traffic.
type Server struct {
	Store storage.Store
	Hub   *Hub
}

func NewServer(store storage.Store) *Server {
	return &Server{
		Store: store,
		Hub:   NewHub(),
	}
}

func (s *Server) Start(addr string) error {
	go s.Hub.Run()

	http.HandleFunc("/ws", s.serveWs)
	http.HandleFunc("/api/stats", s.handleStats)
	http.HandleFunc("/api/duplicates", s.handleDuplicates)
	http.HandleFunc("/api/export", s.handleExport)
	http.HandleFunc("/api/resolve", s.handleResolve)

	// Serve static files from the assets directory
	http.Handle("/", http.FileServer(http.Dir("internal/pkg/server/assets")))

	fmt.Printf("Faro Dashboard starting on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	s.Hub.register <- conn
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	records, _ := s.Store.ListRecords()
	dups, _ := s.Store.GetDuplicates()

	var potential, resolved int
	for _, d := range dups {
		if d.Resolved {
			resolved++
		} else {
			potential++
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_records": len(records),
		"duplicates":    potential,
		"resolved":      resolved,
	})
}

func (s *Server) handleDuplicates(w http.ResponseWriter, r *http.Request) {
	dups, _ := s.Store.GetDuplicates()
	json.NewEncoder(w).Encode(dups)
}

func (s *Server) handleResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RecordA string `json:"record_a"`
		RecordB string `json:"record_b"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Find and update the duplicate
	dups, _ := s.Store.GetDuplicates()
	var found *types.SimilarityResult
	for _, d := range dups {
		if (d.RecordA == req.RecordA && d.RecordB == req.RecordB) ||
			(d.RecordA == req.RecordB && d.RecordB == req.RecordA) {
			found = &d
			break
		}
	}

	if found != nil {
		found.Resolved = true
		s.Store.PutDuplicate(*found)
	} else {
		// If not found, create a new resolved entry (unlikely but safe)
		s.Store.PutDuplicate(types.SimilarityResult{
			RecordA:     req.RecordA,
			RecordB:     req.RecordB,
			Resolved:    true,
			IsDuplicate: true,
			Algorithm:   "Manual",
		})
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=faro_report.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"Record A", "Record B", "Similarity Score", "Algorithm", "Status", "Raw Metadata"})
	
	// Sample data for export demo (updated with metadata)
	writer.Write([]string{"REC001", "REC002", "100.00%", "Levenshtein", "Potential", "{}"})
	writer.Write([]string{"PAT001", "PAT001-CON", "50.00%", "HierarchicalMetadata", "Conflict", `{"study_instance_uid":"1.2.3.4.5","series":[{"series_instance_uid":"1.2.3.4.5.1","modality":"CT"}]}`})
}

// BroadcastDiscovery sends a new discovery event to all connected dashboards.
func (s *Server) BroadcastDiscovery(res types.SimilarityResult) {
	data, _ := json.Marshal(res)
	s.Hub.broadcast <- data
}
