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
	Seed        int32         // Seed for the model
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

	// TODO: Need to check if messages receieved have a single or multiple tasks in them
	// as well as the need of running tools.
	// 1. Have a user/system message which describes how to extract tasks from a single message
	// 		- maybe a tool which can call the model?
	// 2. For each task, run a chatrequest and then validate the output if it needs to be re-run
	//		and if so, then reformat the instruction of the task and run the chatrequest again.
	//		repeat N amount of times until the output is correct.
	// 3. Include information about how each task was solved (specifically if a tool was used or not)
	// 		- might need to be able to attach which request has run a tool or not, and which ones
	// 4. Merge the response to each task into a single response and return it
	//		- alternatively, summarize the response to each task and return it

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
			"seed":        reqSettings.Seed,
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

			// TODO: add a flag if this should run or not.
			// default should be true

			// Send back the tool response, and let current model handle it
			c.SendWithMessage(messages, reqSettings, tools)

		} else {
			log.Println("Response:\n", resp.Message.Content)
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
