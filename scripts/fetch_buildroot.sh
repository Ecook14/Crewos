#!/usr/bin/env bash
# Fetch Buildroot Source

set -e

BUILDROOT_VERSION="2023.11.1"
BUILD_DIR="${PWD}/build"
BUILDROOT_DIR="${BUILD_DIR}/buildroot"

echo "[CrewOS] Fetching Buildroot ${BUILDROOT_VERSION}..."

mkdir -p "$BUILD_DIR"

if [ -d "$BUILDROOT_DIR" ]; then
    echo "[CrewOS] Buildroot already exists at $BUILDROOT_DIR. Skipping download."
    exit 0
fi

# Download Buildroot
curl -L "https://buildroot.org/downloads/buildroot-${BUILDROOT_VERSION}.tar.gz" -o "$BUILD_DIR/buildroot.tar.gz"

echo "[CrewOS] Extracting Buildroot..."
tar -xf "$BUILD_DIR/buildroot.tar.gz" -C "$BUILD_DIR"
mv "$BUILD_DIR/buildroot-${BUILDROOT_VERSION}" "$BUILDROOT_DIR"
rm "$BUILD_DIR/buildroot.tar.gz"

echo "[CrewOS] Buildroot successfully fetched to $BUILDROOT_DIR"
exit 0
