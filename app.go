package main

import (
	"encoding/json"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type App struct {
	conf      *Config
	tools     *Tools
	userQuery string
	history   []llms.MessageContent
	errors    int
}

func NewApp() *App {
	conf := LoadConfig()

	userQuery := "What is the weather like in Beijing? Prefer celsius."
	if len(os.Args) > 1 {
		userQuery = os.Args[1]
	}

	app := &App{
		conf:      conf,
		tools:     NewTools(conf),
		userQuery: userQuery,
	}

	app.push(app.systemMessage(userQuery))

	return app
}

func (a *App) banner() {
	a.conf.Logger.Infof("Model: %v", a.conf.Model)
	a.conf.Logger.Infof("Query: %v", a.userQuery)
	a.conf.Logger.Debugf("Available tools: %v", a.tools)
}

func (a *App) push(msg llms.MessageContent) {
	a.conf.Logger.Debug("push history", "pos", len(a.history), "role", msg.Role, "contents", msg.Parts)

	a.history = append(a.history, msg)
}

func (a *App) systemMessage(question string) llms.MessageContent {
	bs, err := json.Marshal(a.tools.functions)
	if err != nil {
		a.conf.Logger.Fatal(err)
	}

	return makeSystemMessage(string(bs), question)
}

func (a *App) handle(input string) bool {
	a.push(llms.TextParts(schema.ChatMessageTypeAI, input))

	if calls := a.unmarshalCall(input); calls != nil {
		var text string

		done := false

		for _, call := range calls {
			msg, ok := a.tools.execute(&call)
			if !ok {
				done = true
			}

			text += partsToString(msg.Parts)
		}

		a.push(makeResponse(text))

		return !done
	} else {
		a.conf.Logger.Infof("Not a call: %v", input)

		a.errors++

		return a.errors < 10
	}
}

func (a *App) unmarshalCall(input string) []Call {
	cs, err := ParseCalls(input)
	if err != nil {
		a.conf.Logger.Errorf("parse calls failed: %v", err)

		return nil
	}

	for i, c := range cs {
		a.conf.Logger.Debug("call", "index", i, "function", c)
	}

	return cs
}
