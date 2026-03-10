#!/bin/sh
# CrewOS Runtime Control Plane (CrewAI-Go Executor)
# Executed by /init as PID 1 or secondary launcher

echo "[CrewOS-Runtime] Initializing generic execution mesh..."

# Ensure required run directories exist
mkdir -p /run/crewos
mkdir -p /var/log/crewos

# Run GPU PCI-E Detection if we are on the GPU variant
if [ "$OS_VARIANT" = "gpu" ] && [ -x /usr/bin/gpu_detect ]; then
    /usr/bin/gpu_detect
fi

# Determine startup mode
if [ -z "$CREW_STARTUP_MODE" ]; then
    CREW_STARTUP_MODE="sandbox"
fi

echo "[CrewOS-Runtime] Mode: $CREW_STARTUP_MODE"

# Launch Agent Process (Placeholder until Go binary is built)
if [ -x /usr/bin/crewai ]; then
    echo "[CrewOS-Runtime] Starting CrewAI..."
    exec /usr/bin/crewai --mode "$CREW_STARTUP_MODE"
else
    echo "[CrewOS-Runtime] Warning: /usr/bin/crewai not found."
    echo "[CrewOS-Runtime] Dropping to fallback shell..."
    exec /bin/sh
fi
