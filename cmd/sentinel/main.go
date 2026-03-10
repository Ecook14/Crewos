package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// AgentStatus represents the real-time state of a CrewOS agent.
type AgentStatus struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	OSVariant string    `json:"os_variant"`
	GPUState  string    `json:"gpu_state"`
	CPUUsage  float64   `json:"cpu_usage"`
	MemUsage  float64   `json:"mem_usage"`
	NetIn     float64   `json:"net_in"`
	NetOut    float64   `json:"net_out"`
	Heartbeat time.Time `json:"last_heartbeat"`
	Status    string    `json:"status"` // online, idle, busy
}

var (
	agentsMap = make(map[string]*AgentStatus)
	agentsMu  sync.RWMutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients = make(map[*websocket.Conn]bool)
	broadcast = make(chan []byte)
)

func main() {
	fmt.Println("===============================================")
	fmt.Println("      CrewOS Sentinel Fleet Manager")
	fmt.Println("      Enterprise Control Plane v1.0")
	fmt.Println("===============================================")

	http.HandleFunc("/api/register", handleRegister)
	http.HandleFunc("/api/agents", handleListAgents)
	http.HandleFunc("/api/deploy/image", handleDeployImage)
	http.HandleFunc("/api/deploy/git", handleDeployGit)
	http.HandleFunc("/api/updates/check", handleCheckUpdates)
	http.HandleFunc("/ws", handleWebSocket)

	go handleMessages()
	go checkDeadAgents()

	addr := ":8080"
	log.Printf("[Sentinel] Starting API server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type DeployRequest struct {
	AgentID string `json:"agent_id"`
	Image   string `json:"image,omitempty"`
	Source  string `json:"source,omitempty"` // Git URL
}

func handleDeployImage(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Sentinel] Deploying Image %s to %s", req.Image, req.AgentID)
	// Broadast command to the agent via WebSocket/Mesh
	msg, _ := json.Marshal(map[string]interface{}{
		"type":    "deploy_image",
		"target":  req.AgentID,
		"image":   req.Image,
	})
	broadcast <- msg
	w.WriteHeader(http.StatusAccepted)
}

func handleDeployGit(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[Sentinel] Deploying Git Source %s to %s", req.Source, req.AgentID)
	msg, _ := json.Marshal(map[string]interface{}{
		"type":    "deploy_git",
		"target":  req.AgentID,
		"source":  req.Source,
	})
	broadcast <- msg
	w.WriteHeader(http.StatusAccepted)
}

func handleCheckUpdates(w http.ResponseWriter, r *http.Request) {
	// Mock implementation of update server
	update := map[string]string{
		"version":  "v1.1.0",
		"checksum": "a1b2c3d4e5f6...", // SHA256 sum
		"url":      "http://localhost:8080/static/updates/crewos-v1.1.0.tar.gz",
	}
	json.NewEncoder(w).Encode(update)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var status AgentStatus
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status.Heartbeat = time.Now()
	status.Status = "online"

	agentsMu.Lock()
	agentsMap[status.ID] = &status
	agentsMu.Unlock()

	log.Printf("[Sentinel] Agent Registered: %s (%s)", status.Role, status.ID)
	msg, _ := json.Marshal(map[string]interface{}{"type": "agent_update", "agent": status})
	broadcast <- msg

	w.WriteHeader(http.StatusOK)
}

func handleListAgents(w http.ResponseWriter, r *http.Request) {
	agentsMu.RLock()
	defer agentsMu.RUnlock()

	var list []*AgentStatus
	for _, a := range agentsMap {
		list = append(list, a)
	}

	json.NewEncoder(w).Encode(list)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Printf("[Sentinel] New Dashboard Connection established.")

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("[Error] WS error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func checkDeadAgents() {
	for {
		time.Sleep(10 * time.Second)
		agentsMu.Lock()
		for id, agent := range agentsMap {
			// 1. Heartbeat Check
			if time.Since(agent.Heartbeat) > 30*time.Second {
				log.Printf("[Sentinel] Agent Timed Out: %s", agent.Role)
				agent.Status = "offline"
				msg, _ := json.Marshal(map[string]interface{}{"type": "agent_offline", "id": id})
				broadcast <- msg
				delete(agentsMap, id)
				continue
			}

			// 2. Resource/OOM Guard (Self-Healing)
			if agent.MemUsage > 90.0 {
				log.Printf("[Sentinel] WARNING: Agent %s (%s) High Memory Usage: %.1f%%", agent.Role, id, agent.MemUsage)
				msg, _ := json.Marshal(map[string]interface{}{"type": "resource_alert", "id": id, "reason": "OOM Prevention"})
				broadcast <- msg
				
				// Optional: Trigger remote restart if this were a container
				// go triggerRestart(id)
			}
		}
		agentsMu.Unlock()
	}
}
