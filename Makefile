# Top-level Makefile for CrewOS

BUILD_DIR := $(PWD)/build
BUILDROOT_DIR := $(BUILD_DIR)/buildroot
KERNEL_DIR := $(BUILD_DIR)/linux

.PHONY: all default fetch fetch-buildroot fetch-kernel prune node lite gpu sentinel clean docker microvm-lite microvm-gpu

default: fetch prune node lite sentinel

all: fetch prune node lite gpu sentinel microvm-lite microvm-gpu docker

fetch: fetch-buildroot fetch-kernel

fetch-buildroot:
	@chmod +x scripts/fetch_buildroot.sh
	@./scripts/fetch_buildroot.sh

fetch-kernel:
	@chmod +x scripts/fetch_kernel.sh
	@./scripts/fetch_kernel.sh

prune:
	@chmod +x scripts/prune_build_env.sh
	@./scripts/prune_build_env.sh

node:
	@echo "Building CrewOS Node Manager..."
	@go build -o overlay/usr/bin/crew-node ./cmd/node

sentinel:
	@echo "Building Sentinel Fleet Manager..."
	@GOEXPERIMENT=aliastypeparams go build -o build/sentinel ./cmd/sentinel

lite: node
	@chmod +x scripts/build_lite.sh
	@./scripts/build_lite.sh

gpu:
	@chmod +x scripts/build_gpu.sh
	@./scripts/build_gpu.sh

docker:
	@chmod +x scripts/package_docker.sh
	@./scripts/package_docker.sh

microvm-lite: lite
	@chmod +x scripts/package_microvm.sh
	@./scripts/package_microvm.sh lite

microvm-gpu: gpu
	@chmod +x scripts/package_microvm.sh
	@./scripts/package_microvm.sh gpu

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)/*
	@echo "Clean complete."
