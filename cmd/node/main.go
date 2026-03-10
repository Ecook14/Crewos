package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ecook14/crewos/pkg/mesh"
	"github.com/Ecook14/crewos/pkg/ota"
	"github.com/Ecook14/crewos/pkg/runtime"
)

func main() {
	fmt.Println("===============================================")
	fmt.Println("      CrewOS Universal Node Manager")
	fmt.Println("      Status: SDK-Agnostic Platform")
	fmt.Println("===============================================")

	// 1. Load CrewOS Runtime Info
	sysInfo, err := runtime.GetSystemInfo()
	if err != nil {
		log.Printf("[Warning] Failed to read system info: %v", err)
	} else {
		sysInfo.Log("Runtime initialized.")
	}

	// 2. Initialize Mesh Node
	node := mesh.NewNode()
	if err := node.Start(); err != nil {
		log.Fatalf("[Error] Failed to join mesh: %v", err)
	}

	fmt.Printf("[Node] Online on %s (GPU: %s)\n", sysInfo.OSVariant, sysInfo.GPUState)

	// 3. Register with Sentinel Fleet Manager & OTA Client
	go registerWithSentinel(sysInfo)
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

	// 4. Workload Manager (Handles dynamic deployments of Gocrewwai agents)
	// mgr := workload.NewManager() - Not used in this basic loop yet
	
	// 5. Keep alive & Signal Handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("[Node] Ready to host workloads. Listening for commands via Mesh/Sentinel...")

	<-sigs
	fmt.Println("\n[Node] Shutting down...")
	node.Stop()
}

func registerWithSentinel(info *runtime.Info) {
	sentinelURL := "http://localhost:8080/api/register"
	ticker := time.NewTicker(10 * time.Second)
	
	id := fmt.Sprintf("node-%d", time.Now().UnixNano())

	for {
		// Refresh system info for real-time telemetry
		sInfo, _ := runtime.GetSystemInfo()
		
		payload := map[string]interface{}{
			"id":         id,
			"role":       "CrewOS Node Host",
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
