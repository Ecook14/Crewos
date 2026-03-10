#!/usr/bin/env bash
# CrewOS - Termux Setup Script (Run inside Termux)
# Purpose: Prepare Environment for Zero-Root CrewOS Agent Deployment

set -e

echo "[CrewOS] Initializing Termux Environment..."

# 1. Update packages
pkg update -y
pkg upgrade -y

# 2. Install essential dependencies
pkg install -y proot proot-distro curl tar wget git

# 3. Setup Storage Access
termux-setup-storage

# 4. Create CrewOS workspace
mkdir -p ~/crewos
cd ~/crewos

echo "[CrewOS] Environment Ready."
echo "Next: Use android/proot-launcher.sh to start the agent."
