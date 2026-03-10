# 📖 CrewOS Usage Guide

This guide covers the various ways to deploy and manage your CrewOS fleet.

## 🚀 Deployment Modes

### 1. Docker (Cloud & Local)
The easiest way to scale your agents.
```bash
# Build the images
make docker

# Run a lite node
docker run -d --name agent-node-01 -e SENTINEL_URL="http://your-ip:8080" crewos:lite
```

### 2. MicroVMs (High Isolation)
Perfect for multi-tenant or ultra-secure agent hosting.
```bash
# Generate the .ext4 images
make microvm-lite

# Boot with Firecracker (Example config in docker/firecracker-config.json)
firecracker --api-sock /tmp/firecracker.socket --config-file docker/firecracker-config.json
```

### 3. Android (Zero-Root Edge)
Run CrewOS nodes on mobile hardware.
1.  Transfer `build/crewos-lite-rootfs.tar.gz` to your device.
2.  Install [Termux](https://termux.dev/).
3.  Run the installer:
    ```bash
    pkg install wget -y
    chmod +x android/installer.sh
    ./android/installer.sh
    ```
4.  Launch: `./android/proot-launcher.sh`

---

## 🎮 Managing the Fleet (Sentinel)

Sentinel provides a centralized API and Dashboard.

- **Dashboard**: `http://localhost:8080`
- **Register Node**: `POST /api/register`
- **Telemetry**: Heartbeats include CPU, RAM, and GPU status.

---

## 🛠 Advanced Configuration

### OS Tunables
Edit `overlay/etc/crewos.conf` to adjust:
- `OS_VARIANT`: lite or gpu
- `MESH_PORT`: Port for internal agent communication
- `SENTINEL_URL`: The master control plane address

### Kernel Parameters
Modify `kernel/config-lite` or `kernel/config-gpu` then run `make prune && make all` to regenerate the OS with custom drivers or optimizations.
