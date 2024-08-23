package helpers

import (
	types "github.com/asabla/gogent/types"
	"github.com/ollama/ollama/api"
)

type OllamaTool struct {
	Name         string
	Description  string
	Parameters   types.ToolFunctionParameter
	ToolFunction func(api.ToolCallFunctionArguments) string
}

func (t *OllamaTool) GetToolDefinition() api.Tool {
	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
		},
	}
}
