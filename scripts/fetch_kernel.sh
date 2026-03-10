#!/usr/bin/env bash
# Fetch Linux Kernel Source

set -e

KERNEL_VERSION="6.6.21"
BUILD_DIR="${PWD}/build"
KERNEL_DIR="${BUILD_DIR}/linux"
KERNEL_MAJOR=$(echo $KERNEL_VERSION | cut -d. -f1)

echo "[CrewOS] Fetching Linux Kernel ${KERNEL_VERSION}..."

mkdir -p "$BUILD_DIR"

if [ -d "$KERNEL_DIR" ]; then
    echo "[CrewOS] Linux Kernel already exists at $KERNEL_DIR. Skipping download."
    exit 0
fi

# Download Linux Kernel
curl -L "https://cdn.kernel.org/pub/linux/kernel/v${KERNEL_MAJOR}.x/linux-${KERNEL_VERSION}.tar.xz" -o "$BUILD_DIR/linux.tar.xz"

echo "[CrewOS] Extracting Linux Kernel..."
tar -xf "$BUILD_DIR/linux.tar.xz" -C "$BUILD_DIR"
mv "$BUILD_DIR/linux-${KERNEL_VERSION}" "$KERNEL_DIR"
rm "$BUILD_DIR/linux.tar.xz"

echo "[CrewOS] Linux Kernel successfully fetched to $KERNEL_DIR"
exit 0
