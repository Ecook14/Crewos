#!/usr/bin/env bash
# Package CrewOS for Firecracker MicroVM Execution
set -e

echo "[CrewOS] Packaging MicroVM Image..."

if [ -z "$1" ]; then
    echo "Usage: ./package_microvm.sh <lite|gpu>"
    exit 1
fi

VARIANT=$1
ROOTFS_TARBALL="build/crewos-${VARIANT}-rootfs.tar.gz"
EXT4_IMAGE="build/crewos-${VARIANT}.ext4"

if [ ! -f "$ROOTFS_TARBALL" ]; then
    echo "Error: RootFS tarball $ROOTFS_TARBALL not found."
    exit 1
fi

# Cleanup any stale mounts or loop devices from previous failed runs
echo "Cleaning up stale mounts and loop devices..."
if mountpoint -q /tmp/crewos-mount; then
    sudo fuser -km /tmp/crewos-mount || true
    sudo umount -l /tmp/crewos-mount || true
fi

# Ensure the image file itself isn't locked by any process
if [ -f "$EXT4_IMAGE" ]; then
    sudo fuser -k "$EXT4_IMAGE" || true
fi

# Detach ALL loop devices associated with this file
for dev in $(losetup -j "$EXT4_IMAGE" | cut -d: -f1); do
    sudo losetup -d "$dev" || true
done

# Final settle to let the kernel catch up
sudo udevadm settle || true

echo "Creating 512MB ext4 filesystem image..."
dd if=/dev/zero of="$EXT4_IMAGE" bs=1M count=512
mkfs.ext4 -F "$EXT4_IMAGE"

mkdir -p /tmp/crewos-mount
# Require sudo for loopback mount
sudo mount -o loop "$EXT4_IMAGE" /tmp/crewos-mount

echo "Extracting RootFS onto ext4 image..."
sudo tar -xf "$ROOTFS_TARBALL" -C /tmp/crewos-mount

sudo umount /tmp/crewos-mount
rm -rf /tmp/crewos-mount

echo "[CrewOS] MicroVM ext4 Image created: $EXT4_IMAGE"
echo "(Can be booted via Firecracker or QEMU)"
exit 0
