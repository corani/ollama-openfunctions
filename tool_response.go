package main

import (
	"encoding/json"

	"github.com/tmc/langchaingo/llms"
)

func init() {
	register(func(t *Tools) {
		t.push(FunctionDefinition{
			Name:        "finalResponse",
			Description: "Provide the final response to the user query",
			Parameters: json.RawMessage(`{
				"type": "object", 
				"properties": {
					"response": {"type": "string", "description": "The final response to the user query"}
				}, 
				"required": ["response"]
			}`),
			Tool: ToolFunc(func(c *Call) (llms.MessageContent, bool) {
				text, ok := c.Input["response"].(string)
				if !ok {
					t.log.Errorf("function call: %#v", c)
					t.log.Fatal("invalid input")
				}

				t.log.Debugf("Call tool: %s(%q)", c.Tool, text)
				t.log.Infof("Final response: %v", text)

				return llms.MessageContent{}, false
			}),
		})
	})
}
