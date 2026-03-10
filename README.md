# 🚀 CrewOS: The Universal Agentic Operating System

**Small. Light. Powerful. Hardened.**

CrewOS is a minimalist, enterprise-grade Linux distribution designed specifically to host and manage **Gocrewwai** agent meshes. It transforms edge devices, cloud VMs, and even Android phones into high-performance, SDK-agnostic agent hosts.

---

## 💎 Elite Features

- **Extreme Minimalism**: ~16MB RootFS (Lite) / ~32MB Docker image. Reduced attack surface by ~80%.
- **SDK-Agnostic Runtime**: Decoupled core platform (`crew-node`) handles telemetry and orchestration, while agents run as separate payloads.
- **GPU-Ready**: Automatic discovery and passthrough for NVIDIA, Intel, AMD, and Mobile (Adreno/Mali) hardware.
- **Universal Deployment**: Runs on Bare Metal (x86/ARM), MicroVMs (Firecracker/QEMU), OCI Containers, and Zero-Root Android (Termux/PRoot).
- **Glassmorphic Fleet Dashboard**: Real-time observability for CPU, RAM, GPU, and Mesh Network health.
- **Immutable Boot**: dm-verity hashes and hardened kernel configurations for production security.

---

## 🏗 Architecture

CrewOS follows a clean separation of concerns:
1. **The Engine**: A pruned, hardened Linux kernel and minimalist RootFS.
2. **The Node Manager (`crew-node`)**: A Go-based platform manager that handles mesh heartbeats and workload hosting.
3. **The Control Plane (Sentinel)**: A centralized manager for fleet orchestration and real-time visualization.

---

## 🛠 Quick Start

### 1. Build the OS (WSL2/Linux)
```bash
# Clone and build entire stack
make all
```

### 2. Launch the Fleet Manager
```bash
./build/sentinel
```

### 3. Deploy a Node (Docker)
```bash
docker run -d --name crew-node-01 crewos:lite
```

---

## 📖 Documentation

- **[Installation & Usage](USAGE.md)**: Detailed guides for MicroVMs, Android, and Cloud.
- **[Architecture Deep-Dive](docs/ARCHITECTURE.md)**: Understanding the "Extreme Pruning" and Mesh logic.
- **[Contributing](CONTRIBUTING.md)**: Help us push the boundaries of agentic infrastructure.

---

## 🌍 Open Source Strategy
CrewOS is designed for the community. The "Extreme Pruning" scripts and minimalist configurations are open for audit and improvement. We believe the future of AI belongs at the edge—small, fast, and secure.

---

Developed with ❤️.