#!/usr/bin/env bash
# Termux + PRoot Execution Script for CrewOS

unset LD_PRELOAD
command="proot"
command+=" --link2symlink"
command+=" -0"
command+=" -r crewos-fs"
command+=" -b /dev"
command+=" -b /proc"
command+=" -b /sys"
# GPU Passthrough (Adreno/Mali)
if [ -e /dev/kgsl-3d0 ]; then command+=" -b /dev/kgsl-3d0"; fi
if [ -d /dev/dri ]; then command+=" -b /dev/dri"; fi
# Shared Memory for Go Runtime
command+=" -b /dev/shm"
command+=" -b crewos-fs/root:/dev/shm"

# Execute CrewOS init
command+=" -w /root /init"

echo "Booting CrewOS via PRoot..."
$command
