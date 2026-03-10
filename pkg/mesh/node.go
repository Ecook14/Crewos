package mesh

import (
	"fmt"
	"sync"
)

// Node represents a single agent node in the CrewOS mesh.
type Node struct {
	mu       sync.Mutex
	running  bool
	peers    []string
}

// NewNode creates a new mesh node instance.
func NewNode() *Node {
	return &Node{
		peers: make([]string, 0),
	}
}

// Start begins the discovery process and joins the mesh.
func (n *Node) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.running {
		return fmt.Errorf("node already running")
	}

	fmt.Println("[Mesh] Initializing discovery protocol (UDP/mDNS)...")
	// Placeholder for actual discovery logic
	n.running = true

	// Initialize A2A Listener (Placeholder for cross-platform delegation)
	go n.listenForA2ATasks()
	
	return nil
}

// listenForA2ATasks listens for task delegations from other Gocrewwai nodes.
func (n *Node) listenForA2ATasks() {
	fmt.Println("[Mesh] A2A Bridge active: Ready to receive delegated tasks from Gocrew fleets.")
	// In a real implementation, this would open a gRPC or WebSocket listener
}

// Stop gracefully leaves the mesh.
func (n *Node) Stop() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !n.running {
		return
	}

	fmt.Println("[Mesh] Notifying peers of departure...")
	n.running = false
}
