package main

import (
	"context"

	"github.com/charmbracelet/log"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	app := NewApp()

	// TODO(daniel): check if JSON format can be made to work, to avoid the custom
	// function call parser.
	llm, err := ollama.New(
		ollama.WithModel(app.conf.Model),
		// ollama.WithFormat("json"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	app.banner()

	for {
		resp, err := llm.GenerateContent(ctx, app.history)
		if err != nil {
			app.conf.Logger.Fatal(err)
		}

		if len(resp.Choices) == 0 {
			app.conf.Logger.Fatalf("Got no choices: %v", resp)
		}

		if !app.handle(resp.Choices[0].Content) {
			break
		}
	}
}
