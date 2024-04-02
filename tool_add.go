package main

import (
	"encoding/json"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func init() {
	register(func(t *Tools) {
		t.push(FunctionDefinition{
			Name:        "add",
			Description: "Add two numbers",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"v1": {"type": "integer", "description": "The first number to add"},
					"v2": {"type": "integer", "description": "The second number to add"}
				},
				"required": ["v1", "v2"]
			}`),
			Response: json.RawMessage(`{
				"type": "object",
				"properties": {
					"result": {"type": "integer", "description": "The result of the addition"}
				}
			}`),
			Tool: ToolFunc(func(c *Call) (llms.MessageContent, bool) {
				v1, ok := c.Input["v1"].(int)
				if !ok {
					t.log.Errorf("function call: %#v", c)
					t.log.Fatal("invalid input")
				}

				v2, ok := c.Input["v2"].(int)
				if !ok {
					t.log.Errorf("function call: %#v", c)
					t.log.Fatal("invalid input")
				}

				t.log.Infof("Call tool: %s(%v, %v)", c.Tool, v1, v2)

				return llms.TextParts(schema.ChatMessageTypeHuman, fmt.Sprintf(`{"result": %d}`, v1+v2)), true
			}),
		})
	})
}
