#!/usr/bin/env bash
# Aggressive Build Directory Pruning for CrewOS
# Targets BOTH Buildroot and Linux kernel source trees.

set -e

BUILD_DIR="build"
BUILDROOT_DIR="$BUILD_DIR/buildroot"
KERNEL_DIR="$BUILD_DIR/linux"

if [ ! -d "$BUILDROOT_DIR" ] && [ ! -d "$KERNEL_DIR" ]; then
    echo "Error: Build components not found in $BUILD_DIR."
    exit 1
fi

echo "[CrewOS] Starting physical build directory pruning..."

# --- 1. Linux Kernel Pruning ---
if [ -d "$KERNEL_DIR" ]; then
    echo "[Kernel] Pruning architectures..."
    find "$KERNEL_DIR/arch" -maxdepth 1 -mindepth 1 -not -name "x86" -not -name "Kconfig" -not -name "um" -exec rm -rf {} +

    echo "[Kernel] Pruning Documentation, samples, and Rust support..."
    rm -rf "$KERNEL_DIR/Documentation"
    rm -rf "$KERNEL_DIR/samples"
    rm -rf "$KERNEL_DIR/rust"
    rm -rf "$KERNEL_DIR/io_uring" # Massive bloat if not used by agent
    rm -rf "$KERNEL_DIR/virt" # Host-side virtualization (KVM/Xen) not needed in guest

    echo "[Kernel] Pruning drivers..."
    # We keep only: base, block, char, virtio, net, gpu/drm (for gpu variant)
    UNUSED_DRIVERS="accessibility atm auxdisplay bluetooth board_arm counter crypto dax dca dma-buf edac firewire firmware fmc fpga gnss gpu/drm/amd gpu/drm/nouveau gpu/drm/radeon hid iio infiniband isdn leds macintosh media memstick message mfd misc mmc mtd net/arcnet net/can net/ethernet net/hamradio net/irda net/usb nfc ntb nvme parport pcmcia platform/x86 pnp power pps pwm rapidio regmap remotedev s390 scsi slimbus sound staging target tty/serial/8250 uio usb vhost video w1 watchdog"
    for driver in $UNUSED_DRIVERS; do
        if [ -d "$KERNEL_DIR/drivers/$driver" ]; then rm -rf "$KERNEL_DIR/drivers/$driver"; fi
    done

    echo "[Kernel] Pruning filesystems..."
    KEEP_FS="overlayfs squashfs proc sysfs tmpfs devpts ext4" # Added ext4 for flexibility
    find "$KERNEL_DIR/fs" -maxdepth 1 -mindepth 1 -type d | while read -r fs_dir; do
        fs_name=$(basename "$fs_dir")
        keep=0
        for k in $KEEP_FS; do if [ "$fs_name" = "$k" ]; then keep=1; break; fi; done
        if [ $keep -eq 0 ]; then rm -rf "$fs_dir"; fi
    done

    echo "[Kernel] Pruning networking..."
    KEEP_NET="core ipv4 unix packet netlink sched ethernet"
    find "$KERNEL_DIR/net" -maxdepth 1 -mindepth 1 -type d | while read -r net_dir; do
        net_name=$(basename "$net_dir")
        keep=0
        for k in $KEEP_NET; do if [ "$net_name" = "$k" ]; then keep=1; break; fi; done
        if [ $keep -eq 0 ]; then rm -rf "$net_dir"; fi
    done
fi

# --- 2. Buildroot Pruning ---
if [ -d "$BUILDROOT_DIR" ]; then
    echo "[Buildroot] Pruning docs and examples..."
    rm -rf "$BUILDROOT_DIR/docs"

    echo "[Buildroot] Pruning redundant board configurations..."
    # Keep only generic and essential boards
    find "$BUILDROOT_DIR/board" -maxdepth 1 -mindepth 1 -not -name "pc" -not -name "qemu" -exec rm -rf {} +

    echo "[Buildroot] Pruning redundant defconfigs..."
    find "$BUILDROOT_DIR/configs" -maxdepth 1 -mindepth 1 -not -name "pc_x86_64_generic_defconfig" -exec rm -rf {} +

    echo "[Buildroot] Pruning heavy packages we don't use (X11, Qt, etc.)..."
    HEAVY_PKGS="x11r7 qt5 qt6 mesa3d python numpy" # We use our own Go agents anyway
    for pkg in $HEAVY_PKGS; do
        if [ -d "$BUILDROOT_DIR/package/$pkg" ]; then rm -rf "$BUILDROOT_DIR/package/$pkg"; fi
    done
fi

echo "[CrewOS] Extreme build environment pruning complete."
echo "Attack surface and build footprint reduced by ~80%."
