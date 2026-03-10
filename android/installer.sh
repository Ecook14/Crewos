#!/usr/bin/env bash
# Termux + PRoot Installer Script for CrewOS Android Support

echo "Downloading CrewOS RootFS..."
# curl -LO https://github.com/Ecook14/crewos/releases/latest/download/crewos-lite-rootfs.tar.gz

echo "Extracting RootFS for PRoot environment..."
mkdir -p crewos-fs
tar -xf crewos-lite-rootfs.tar.gz -C crewos-fs

echo "Installation complete. Run ./proot-launcher.sh to start CrewOS."
