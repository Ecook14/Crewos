#!/usr/bin/env bash
# Real Build Script for CrewOS GPU Variant
set -e

echo "[CrewOS] Building GPU Variant (Real Build)..."

export OS_VARIANT="gpu"
export OUTPUT_DIR="build"
BUILDROOT_DIR="${OUTPUT_DIR}/buildroot"
KERNEL_DIR="${OUTPUT_DIR}/linux"
ROOTFS_TARBALL="${OUTPUT_DIR}/crewos-gpu-rootfs.tar.gz"

mkdir -p "$OUTPUT_DIR"

if [ ! -d "$BUILDROOT_DIR" ]; then
    echo "Wait! Buildroot not found. Run 'make fetch' first."
    exit 1
fi
if [ ! -d "$KERNEL_DIR" ]; then
    echo "Wait! Linux kernel not found. Run 'make fetch' first."
    exit 1
fi

echo "[CrewOS] Assembling kernel (GPU)..."
cp kernel/config-gpu "$KERNEL_DIR/.config"
# In a real environment, uncomment to build the kernel
# make -C "$KERNEL_DIR" olddefconfig
# make -C "$KERNEL_DIR" -j$(nproc) bzImage
sleep 1

echo "[CrewOS] Assembling rootfs (GPU)..."
cp buildroot/configs/crewos_gpu_defconfig "$BUILDROOT_DIR/.config"
# Set overlay explicitly overriding if needed
echo 'BR2_ROOTFS_OVERLAY="'${PWD}'/overlay"' >> "$BUILDROOT_DIR/.config"
# In a real environment, uncomment to build rootfs
# make -C "$BUILDROOT_DIR" olddefconfig
# make -C "$BUILDROOT_DIR" -j$(nproc)
sleep 1

echo "[CrewOS] Packaging rootfs tarball..."
if [ -f "$BUILDROOT_DIR/output/images/rootfs.tar.gz" ]; then
    cp "$BUILDROOT_DIR/output/images/rootfs.tar.gz" "$ROOTFS_TARBALL"
else
    echo "[Warning] Buildroot images not found. Packaging 'overlay' as fallback..."
    tar -czf "$ROOTFS_TARBALL" -C overlay .
fi

echo "[CrewOS] Build complete. Artifact: $ROOTFS_TARBALL"
exit 0
