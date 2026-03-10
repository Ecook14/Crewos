package workload

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Manager handles the execution of custom AI workloads.
type Manager struct{}

// NewManager creates a new workload manager.
func NewManager() *Manager {
	return &Manager{}
}

// DeployImage pulls and runs an OCI image using nerdctl.
func (m *Manager) DeployImage(ctx context.Context, image string) error {
	log.Printf("[Workload] Pulling image: %s", image)
	pullCmd := exec.CommandContext(ctx, "nerdctl", "pull", image)
	if err := pullCmd.Run(); err != nil {
		return fmt.Errorf("failed to pull image: %v", err)
	}

	log.Printf("[Workload] Running container: %s", image)
	runCmd := exec.CommandContext(ctx, "nerdctl", "run", "-d", "--name", "ai-workload", image)
	return runCmd.Run()
}

// DeployGit clones or updates a repository and attempts to run it.
func (m *Manager) DeployGit(ctx context.Context, repoURL string) error {
	cacheDir := "/var/cache/crewos/workloads"
	os.MkdirAll(cacheDir, 0755)

	log.Printf("[Workload] Processing Git repo: %s", repoURL)
	
	// Check if already cloned
	repoPath := fmt.Sprintf("%s/current", cacheDir)
	if _, err := os.Stat(repoPath); err == nil {
		log.Printf("[Workload] Updating existing cache...")
		updateCmd := exec.CommandContext(ctx, "git", "-C", repoPath, "pull")
		updateCmd.Run()
	} else {
		log.Printf("[Workload] Cloning fresh copy...")
		cloneCmd := exec.CommandContext(ctx, "git", "clone", repoURL, repoPath)
		if err := cloneCmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repo: %v", err)
		}
	}

	// Simple heuristic: if there's a main.go, try to run it
	log.Printf("[Workload] Attempting to run Git workload...")
	runCmd := exec.CommandContext(ctx, "go", "run", repoPath+"/main.go")
	return runCmd.Start() // Run in background
}
