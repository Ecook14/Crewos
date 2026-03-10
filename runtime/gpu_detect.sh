#!/bin/sh
# CrewOS GPU Autodetection Script

echo "[GPU-Detect] Scanning PCI bus for graphics devices..."

GPU_MODE="none"

if lspci | grep -i 'VGA\|3D\|Display' | grep -qi 'NVIDIA'; then
    echo "[GPU-Detect] Detected NVIDIA Device."
    GPU_MODE="vendor_nvidia"
elif lspci | grep -i 'VGA\|3D\|Display' | grep -qi 'AMD\|Radeon'; then
    echo "[GPU-Detect] Detected AMD Device."
    GPU_MODE="vendor_amd"
elif lspci | grep -i 'VGA\|3D\|Display' | grep -qi 'Intel'; then
    echo "[GPU-Detect] Detected Intel Device."
    GPU_MODE="vendor_intel"
else
    echo "[GPU-Detect] No specific vendor graphics found. Defauting to generic Vulkan."
    GPU_MODE="vulkan"
fi

# Export mode to config or runtime state
echo "export GPU_MODE=$GPU_MODE" > /run/crewos/gpu_state.env
echo "[GPU-Detect] Selected Mode: $GPU_MODE"
