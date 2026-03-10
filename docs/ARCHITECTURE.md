# 🏗 CrewOS Architecture Deep-Dive

CrewOS is built on the principle of **Infrastructure as a Scalpel**. Every component is precision-engineered for the specific task of hosting AI agents.

## 1. The "Extreme Pruning" Strategy
Most "minimal" OSs still include thousands of lines of code for hardware support that will never be used in a cloud or microVM environment.
- **Kernel Level**: We physically remove drivers for floppy disks, ancient network cards, and non-essential filesystems from the source.
- **Buildroot Level**: We strip out documentation, example binaries, and heavy support libraries (X11, Qt, Python) to keep the RootFS under 20MB.
- **Security Impact**: By removing ~80% of the build environment, we eliminate 80% of the potential vulnerabilities associated with common Linux distributions.

## 2. Universal Node Host (`crew-node`)
Instead of building the agent logic directly into the OS, CrewOS uses a **Node Host** architecture.
- **Role**: `crew-node` acts as the "Manager" for the physical or virtual hardware.
- **Mesh Connection**: It establishes an encrypted websocket tunnel to **Sentinel**.
- **Telemetry**: It streams live hardware metrics (CPU speed, Memory pressure, GPU utilization) without needing a heavy agent SDK.

## 3. Workload Isolation
CrewOS supports high-density agent deployments through three layers:
- **OCI (Docker)**: Standard isolation for cloud and local dev.
- **MicroVMs (Firecracker/QEMU)**: Bare-metal performance with hardware-level isolation.
- **PRoot (Android)**: User-space isolation for running on mobile devices without root access.

## 4. Hardware Acceleration (GPU)
CrewOS treats the GPU as a first-class citizen.
- **Universal Vulkan**: Native support for Vulkan-based compute out of the box.
- **Vendor Passthrough**: Pre-configured hooks for NVIDIA CUDA and Intel Level Zero, ensuring agents can tap into hardware power as soon as they are "Pushed" to the node.

---

## 🔒 Security Posture
- **Immutable Boot**: The RootFS can be mounted as read-only with dm-verity verification.
- **Hardened Kernel**: SECCOMP, AppArmor, and Kernel Lockdown Mode are enabled by default.
- **Minimal Surface Area**: No SSH, no bash (at runtime), no unused ports.
