package hardware

import (
	"fmt"
	"os"
	"strings"
)

// GPUInfo represents the detected graphics hardware state.
type GPUInfo struct {
	Vendor      string
	Model       string
	VulkanReady bool
}

// DetectGPU scans the system PCI and DRM subsystems for graphics hardware.
func DetectGPU() (*GPUInfo, error) {
	info := &GPUInfo{
		Vendor:      "none",
		VulkanReady: false,
	}

	// 1. Check /sys/class/drm for active rendering nodes
	if _, err := os.Stat("/sys/class/drm/renderD128"); err == nil {
		info.VulkanReady = true
	}

	// 2. Mobile GPU Detection (Android/ARM64)
	if _, err := os.Stat("/dev/kgsl-3d0"); err == nil {
		info.Vendor = "adreno (qualcomm)"
		info.VulkanReady = true
	} else if _, err := os.Stat("/dev/mali0"); err == nil {
		info.Vendor = "mali (arm)"
		info.VulkanReady = true
	}

	// 3. Simple PCI scan via sysfs (for x86/WSL/Cloud)
	// For simplicity, we parse lspci output if available, or fallback to sysfs
	data, err := os.ReadFile("/sys/bus/pci/devices/0000:00:02.0/vendor") // Common Intel path in VM/PC
	if err == nil {
		vendorID := strings.TrimSpace(string(data))
		switch vendorID {
		case "0x8086":
			info.Vendor = "intel"
		case "0x10de":
			info.Vendor = "nvidia"
		case "0x1002":
			info.Vendor = "amd"
		}
	}

	// 3. Optional: Read vendor specific files if needed
	return info, nil
}

// String returns a human-readable summary of the GPU state.
func (g *GPUInfo) String() string {
	if g.Vendor == "none" {
		return "CPU (No Acceleration)"
	}
	readyStr := "Vulkan Disabled"
	if g.VulkanReady {
		readyStr = "Vulkan Ready"
	}
	return fmt.Sprintf("%s (%s)", g.Vendor, readyStr)
}
