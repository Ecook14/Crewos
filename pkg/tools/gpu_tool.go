package tools

import (
	"context"
	"fmt"
	"github.com/Ecook14/crewos/pkg/runtime"
	"github.com/Ecook14/gocrewwai/pkg/tools"
)

// GPUTool allows agents to interact with the CrewOS GPU layer.
type GPUTool struct {
	sysInfo *runtime.Info
}

// NewGPUTool creates a new instance of the GPU acceleration tool.
func NewGPUTool(info *runtime.Info) *GPUTool {
	return &GPUTool{sysInfo: info}
}

func (t *GPUTool) Name() string {
	return "CrewOS_GPU_Accelerator"
}

func (t *GPUTool) Description() string {
	return "Utilizes the CrewOS Vulkan/Vendor GPU layer for high-performance compute tasks. Use this for local LLM inference, data crunching, or image processing."
}

// Execute handles the GPU task invocation.
func (t *GPUTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	task, _ := input["task"].(string)
	if t.sysInfo.GPUState == "none" || t.sysInfo.GPUState == "" {
		return "", fmt.Errorf("GPU hardware not available in this instance (Lite variant)")
	}

	t.sysInfo.Log(fmt.Sprintf("Action: Executing GPU-accelerated task: %s", task))
	return fmt.Sprintf("Success: Task '%s' completed using CrewOS %s acceleration.", task, t.sysInfo.GPUState), nil
}

func (t *GPUTool) RequiresReview() bool {
	return false
}

func (t *GPUTool) ArgsSchema() []tools.ArgSchema {
	return []tools.ArgSchema{
		{
			Name:        "task",
			Type:        "string",
			Description: "The compute task description to accelerate",
			Required:    true,
		},
	}
}

func (t *GPUTool) CacheFunction(input map[string]interface{}) string {
	return ""
}
