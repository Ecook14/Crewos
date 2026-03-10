package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ecook14/gocrewwai/gocrew"
	"github.com/Ecook14/gocrewwai/pkg/llm"
	"github.com/Ecook14/crewos/pkg/mesh"
	"github.com/Ecook14/crewos/pkg/ota"
	"github.com/Ecook14/crewos/pkg/runtime"
	"github.com/Ecook14/crewos/pkg/tools"
	"github.com/Ecook14/crewos/pkg/workload"
)

func main() {
	fmt.Println("===============================================")
	fmt.Println("      CrewOS Universal Agent Runtime")
	fmt.Println("      Engine: Gocrewwai v0.9.0")
	fmt.Println("===============================================")

	// 1. Load CrewOS Runtime Info
	sysInfo, err := runtime.GetSystemInfo()
	if err != nil {
		log.Printf("[Warning] Failed to read system info: %v", err)
	} else {
		sysInfo.Log("Runtime initialized.")
	}

	// 2. Initialize Hardware-Optimized Tools
	gpuTool := tools.NewGPUTool(sysInfo)

	// 3. Initialize Mesh Node
	node := mesh.NewNode()
	if err := node.Start(); err != nil {
		log.Fatalf("[Error] Failed to join mesh: %v", err)
	}

	// 4. Build the Resident Agent using Gocrewwai SDK
	// This agent can handle tasks assigned via the mesh or local config
	model := llm.NewOpenAIClient(os.Getenv("OPENAI_API_KEY"))
	
	agent := gocrew.NewAgent(gocrew.AgentConfig{
		Role:      "CrewOS Resident Agent",
		Goal:      "Coordinate local hardware resources and execute mesh tasks",
		Backstory: fmt.Sprintf("You are a specialized agent running on %s. Your GPU state is %s.", sysInfo.OSVariant, sysInfo.GPUState),
		LLM:       model,
		Tools:     []gocrew.Tool{gpuTool}, // Inject our custom GPU tool
	})

	fmt.Printf("[Agent] %s is online and listening for tasks...\n", agent.Role)

	// 5. Register with Sentinel Fleet Manager & OTA Client
	go registerWithSentinel(agent, sysInfo)
	otaClient := ota.NewClient("v1.0.0", "http://localhost:8080")
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			if update, err := otaClient.CheckForUpdate(); err == nil && update != nil {
				log.Printf("[OTA] New version found: %s", update.Version)
				// otaClient.ApplyUpdate(update)
			}
		}
	}()

	// 6. Real-Time Observability (Telemetry) & Workload Manager
	mgr := workload.NewManager()
	events := gocrew.GlobalBus.Subscribe()
	go func() {
		for e := range events {
			log.Printf("[Event] %s: %v", e.Type, e.Payload)
			
			// Handle deployment commands from Sentinel
			if e.Type == "deploy_image" {
				image, _ := e.Payload["image"].(string)
				go mgr.DeployImage(context.Background(), image)
			} else if e.Type == "deploy_git" {
				source, _ := e.Payload["source"].(string)
				go mgr.DeployGit(context.Background(), source)
			}
		}
	}()

	// 7. Keep alive & Signal Handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	fmt.Println("\n[Agent] Shutting down...")
	node.Stop()
}

func registerWithSentinel(agent *gocrew.Agent, info *runtime.Info) {
	sentinelURL := "http://localhost:8080/api/register"
	ticker := time.NewTicker(10 * time.Second)
	
	id := fmt.Sprintf("agent-%d", time.Now().UnixNano())

	for {
		// Refresh system info for real-time telemetry
		sInfo, _ := runtime.GetSystemInfo()
		
		payload := map[string]interface{}{
			"id":         id,
			"role":       agent.Role,
			"os_variant": info.OSVariant,
			"gpu_state":  info.GPUState,
			"cpu_usage":  sInfo.CPUUsage,
			"mem_usage":  sInfo.MemUsage,
			"net_in":     sInfo.NetIn,
			"net_out":    sInfo.NetOut,
		}
		
		data, _ := json.Marshal(payload)
		resp, err := http.Post(sentinelURL, "application/json", bytes.NewBuffer(data))
		if err == nil {
			resp.Body.Close()
		} else {
			log.Printf("[Warning] Failed to heartbeat to Sentinel: %v", err)
		}
		
		<-ticker.C
	}
}
