package runtime

import (
	"fmt"
	"os"
	"github.com/Ecook14/crewos/pkg/hardware"
)

// Info holds metadata about the CrewOS environment.
type Info struct {
	OSVariant string  `json:"os_variant"`
	GPUState  string  `json:"gpu_state"`
	CPUUsage  float64 `json:"cpu_usage"`
	MemUsage  float64 `json:"mem_usage"`
	NetIn     float64 `json:"net_in"`  // KB/s
	NetOut    float64 `json:"net_out"` // KB/s
}

// GetSystemInfo reads configuration and state from /etc and /run inside CrewOS.
func GetSystemInfo() (*Info, error) {
	variant := os.Getenv("OS_VARIANT")
	if variant == "" {
		variant = "unknown"
	}

	// Use the new hardware detection package
	gpuInfo, err := hardware.DetectGPU()
	gpuState := "none"
	if err == nil {
		gpuState = gpuInfo.String()
	}

	return &Info{
		OSVariant: variant,
		GPUState:  gpuState,
		CPUUsage:  12.5, // Mock initial
		MemUsage:  45.0, // Mock initial
		NetIn:     256.0, // Mock initial KB/s
		NetOut:    64.0,  // Mock initial KB/s
	}, nil
}

// Log logs a message with runtime context.
func (i *Info) Log(msg string) {
	fmt.Printf("[%s][%s] %s\n", i.OSVariant, i.GPUState, msg)
}
