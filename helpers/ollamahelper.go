package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	types "github.com/asabla/gogent/types"

	"github.com/ollama/ollama/api"
)

type RequestTool struct {
	Name             string
	Description      string
	Parameters       types.ToolFunctionParameter
	FunctionCallBack func(api.ToolCallFunctionArguments) string
}

type RequestSettings struct {
	Model       string        // Name of the model you want to use
	Format      string        // Supports only "json" for now. If empty, defaults to "json"
	Stream      bool          // If the response should be streamed or not, tools require false. Will default to false
	KeepAlive   time.Duration // How long given model should be kept in memory
	Temperature float32       // Temperature for the model
}

type OllamaHandler struct {
	Client *api.Client
}

var Handler *OllamaHandler

func CreateClient() (*OllamaHandler, error) {
	if Handler == nil {
		newClient, err := api.ClientFromEnvironment()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		// TODO: This is quite ugly, make it better
		Handler = &OllamaHandler{
			Client: newClient,
		}

		return Handler, nil
	}

	return Handler, nil
}

func (c *OllamaHandler) SendWithMessage(messages []api.Message, reqSettings RequestSettings, tools *[]RequestTool) (string, error) {
	// TODO: is it necessary to attach this method in this way?
	var client *api.Client
	if Handler == nil {
		return "", fmt.Errorf("Client not initialized")
	} else {
		client = Handler.Client
	}

	if reqSettings.Model == "" {
		return "", fmt.Errorf("Model is required")
	}

	if tools == nil {
		tools = new([]RequestTool)
	}

	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    reqSettings.Model,
		Messages: messages,
		Stream:   &reqSettings.Stream,
		Format:   reqSettings.Format,
		KeepAlive: &api.Duration{
			Duration: reqSettings.KeepAlive,
		},
		Tools: getTools(tools),
		Options: map[string]interface{}{
			"temperature": reqSettings.Temperature,
		},
	}

	var respString string = ""

	respFunc := func(resp api.ChatResponse) error {
		// Content should be empty if there are tool calls and if streaming is false
		if resp.Message.Content == "" && len(resp.Message.ToolCalls) > 0 {
			for _, t := range *tools {

				for _, tool := range resp.Message.ToolCalls {
					if tool.Function.Name == t.Name {
						toolResponse := t.FunctionCallBack(tool.Function.Arguments)
						messages = append(messages, api.Message{
							Role:    "tool",
							Content: toolResponse,
						})
						break
					}
				}
			}

			// Send back the tool response, and let current model handle it
			c.SendWithMessage(messages, reqSettings, tools)

		} else {
			log.Println("Response: ", resp.Message.Content)
			respString += resp.Message.Content
		}

		return nil
	}

	err := client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return respString, nil
}

func (t *RequestTool) getToolDefinition() api.Tool {
	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
		},
	}
}

func getTools(t *[]RequestTool) []api.Tool {
	var tools []api.Tool

	for _, tool := range *t {
		tools = append(tools, tool.getToolDefinition())
	}

	return tools
}
