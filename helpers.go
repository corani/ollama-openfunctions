package main

import (
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func makeResponse(text string) llms.MessageContent {
	response := fmt.Sprintf(`
### Instruction: 
%s 
### Response:`,
		text)

	return llms.TextParts(schema.ChatMessageTypeHuman, response)
}

func makeSystemMessage(tools, question string) llms.MessageContent {
	text := fmt.Sprintf(`
You are an AI programming assistant, utilizing the Gorilla LLM model, developed by Gorilla LLM, and
you only answer questions related to computer science. For politically sensitive questions, security
and privacy issues, and other non-computer science questions, you will refuse to answer.

### Instruction: 
<<function>> %s
<<question>> %s 
### Response:
`, tools, question)

	return llms.TextParts(schema.ChatMessageTypeHuman, text)
}

func partsToString(parts []llms.ContentPart) string {
	var text string

	for _, part := range parts {
		if v, ok := part.(llms.TextContent); ok {
			text += v.Text + "\n"
		}
	}

	return text
}
