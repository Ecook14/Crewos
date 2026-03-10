#!/usr/bin/env bash
# Repair script for broken Docker Snap in WSL
set -e

echo "==========================================="
echo "   CrewOS Docker Recovery Script"
echo "==========================================="

# 1. Remove the broken snap version
echo "[1/4] Removing corrupted Snap Docker..."
sudo snap remove docker --purge || true

# 2. Clean up any leftover snap artifacts
echo "[2/4] Cleaning snap mounts..."
sudo umount /snap/docker/* 2>/dev/null || true
sudo rm -rf /var/snap/docker /snap/docker

# 3. Install native Docker (Reliable for WSL)
echo "[3/4] Installing native Docker (apt) + Buildx..."
sudo apt update
sudo apt install -y docker.io docker-buildx

# 4. Start the Docker service
echo "[4/4] Restarting Docker service..."
sudo service docker start

# 5. Verify
echo "==========================================="
echo "Verification:"
docker --version
sudo docker run --rm hello-world
echo "==========================================="
echo "Success! You can now run 'make all' to finish OCI packaging."
