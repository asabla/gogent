package main

import (
	"fmt"
	"log"
	"time"

	helper "github.com/asabla/gogent/helpers"
	types "github.com/asabla/gogent/types"

	"github.com/ollama/ollama/api"
)

var (
	Model                     = "llama3.1:8b"
	Prompt                    = ""
	Format                    = ""
	Stream      bool          = false
	Temperature float32       = 0.0
	KeepAlive   time.Duration = 30 * time.Minute
)

func main() {
	client, err := helper.CreateClient()
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.SendWithMessage(
		[]api.Message{
			{
				// TODO: make these into arguments and/or read from a config file
				Role:    "system",
				Content: "Hello, how are you?",
			}, {
				Role:    "user",
				Content: "Fetch the contents of this page https://raw.githubusercontent.com/ollama/ollama/main/docs/faq.md and summarize it",
			}, {
				Role:    "user",
				Content: "Do you know what time it is?",
			},
		},
		helper.RequestSettings{
			Model:       Model,
			Format:      Format,
			Stream:      Stream,
			Temperature: Temperature,
			KeepAlive:   KeepAlive,
		},
		&[]helper.RequestTool{
			// TODO: after each tool call is done, a system message needs to be generated
			// before continueing with the next tool call. Otherwise lesser models (smaller)
			// will easily forget about previous ones and only answer the last question
			getCurrentTimeTool(),
			getContentsOfUrlTool(),
		})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

// TODO: refactor this into an automated process while running
// embedded lua scripts. Make sure to expose a register tool
// function to the lua environment
func getCurrentTimeTool() helper.RequestTool {
	return helper.RequestTool{
		Name:        "get_current_time",
		Description: "Gets the current time for given city",
		Parameters: types.ToolFunctionParameter{
			Type: "object",
			Properties: types.ParameterProperties{
				"city": {
					Type:        "string",
					Description: "Which city to fetch time for",
				},
			},
		},
		FunctionCallBack: func(args api.ToolCallFunctionArguments) string {
			return time.Now().String()
		},
	}
}

func getContentsOfUrlTool() helper.RequestTool {
	return helper.RequestTool{
		Name:        "get_contents_of_url",
		Description: "Fetches the content of a given URL",
		Parameters: types.ToolFunctionParameter{
			Type: "object",
			Properties: types.ParameterProperties{
				"url": {
					Type:        "string",
					Description: "URL to fetch content from",
				},
			},
		},
		FunctionCallBack: func(args api.ToolCallFunctionArguments) string {
			log.Println("Fetching content from: ", args["url"].(string))
			responseString, err := helper.FetchContent(args["url"].(string))
			if err != nil {
				return fmt.Sprintf("Error fetching content: %v", err)
			}

			return responseString
		},
	}
}
