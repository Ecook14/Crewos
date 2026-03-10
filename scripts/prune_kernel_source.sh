#!/usr/bin/env bash
# Aggressive Kernel Source Pruning for CrewOS
# Physically removes code to ensure it's not present during build

set -e

KERNEL_DIR="build/linux"

if [ ! -d "$KERNEL_DIR" ]; then
    echo "Error: KERNEL_DIR ($KERNEL_DIR) not found. Run 'make fetch' first."
    exit 1
fi

echo "[CrewOS] Starting physical kernel source pruning..."

# 1. Prune unused architectures
echo "Pruning architectures..."
find "$KERNEL_DIR/arch" -maxdepth 1 -mindepth 1 -not -name "x86" -not -name "Kconfig" -not -name "um" -exec rm -rf {} +

# 2. Prune Documentation and Samples
echo "Pruning Documentation and samples..."
rm -rf "$KERNEL_DIR/Documentation"
rm -rf "$KERNEL_DIR/samples"

# 3. Prune unused drivers (Broad sweep)
echo "Pruning drivers..."
# We keep: base, block (virtio), char (virtio), virtio, net (virtio), gpu (virtio for GPU variant)
UNUSED_DRIVERS="accessibility atm auxdisplay bluetooth board_arm counter crypto dax dca dma-buf edac firewire firmware fmc fpga gnss gpu/drm/amd gpu/drm/nouveau gpu/drm/radeon hid iio infiniband isdn leds macintosh media memstick message mfd misc mmc mtd net/arcnet net/can net/ethernet net/hamradio net/irda net/usb nfc ntb nvme parport pcmcia platform/x86 pnp power pps pwm rapidio regmap remotedev s390 scsi slimbus sound staging target tty/serial/8250 uio usb vhost video w1 watchdog"

for driver in $UNUSED_DRIVERS; do
    if [ -d "$KERNEL_DIR/drivers/$driver" ]; then
        rm -rf "$KERNEL_DIR/drivers/$driver"
    fi
done

# 4. Prune unused filesystems
echo "Pruning filesystems..."
# Keep ONLY the bare essentials for CrewOS
KEEP_FS="overlayfs squashfs proc sysfs tmpfs devpts"
find "$KERNEL_DIR/fs" -maxdepth 1 -mindepth 1 -type d | while read -r fs_dir; do
    fs_name=$(basename "$fs_dir")
    keep=0
    for k in $KEEP_FS; do
        if [ "$fs_name" = "$k" ]; then keep=1; break; fi
    done
    if [ $keep -eq 0 ]; then
        rm -rf "$fs_dir"
    fi
done

# 5. Prune unused networking
echo "Pruning networking protocols..."
# Keep: core, ipv4, unix, packet, netlink
KEEP_NET="core ipv4 unix packet netlink sched ethernet"
find "$KERNEL_DIR/net" -maxdepth 1 -mindepth 1 -type d | while read -r net_dir; do
    net_name=$(basename "$net_dir")
    keep=0
    for k in $KEEP_NET; do
        if [ "$net_name" = "$k" ]; then keep=1; break; fi
    done
    if [ $keep -eq 0 ]; then
        rm -rf "$net_dir"
    fi
done

# 6. Prune entire subsystems (Sound, ISDN, etc.)
echo "Pruning top-level subsystems..."
rm -rf "$KERNEL_DIR/sound"
rm -rf "$KERNEL_DIR/drivers/isdn"

echo "[CrewOS] Extreme Kernel source pruning complete."
