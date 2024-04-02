package main

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

var initializers []initializer

type initializer func(t *Tools)

func register(i initializer) {
	initializers = append(initializers, i)
}

func NewTools(conf *Config) *Tools {
	tools := &Tools{
		conf: conf,
		log:  conf.Logger,
	}

	for _, fn := range initializers {
		fn(tools)
	}

	return tools
}

type Tools struct {
	conf      *Config
	log       *log.Logger
	functions []FunctionDefinition
}

func (t *Tools) String() string {
	bs, _ := json.MarshalIndent(t.functions, "", "  ")

	return string(bs)
}

type Tool interface {
	Execute(c *Call) (llms.MessageContent, bool)
}

type ToolFunc func(c *Call) (llms.MessageContent, bool)

func (fn ToolFunc) Execute(c *Call) (llms.MessageContent, bool) {
	return fn(c)
}

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
	Response    any    `json:"response"`
	Tool        Tool   `json:"-"`
}

func (t *Tools) push(fd FunctionDefinition) {
	t.functions = append(t.functions, fd)
}

func (t *Tools) execute(c *Call) (llms.MessageContent, bool) {
	var def *FunctionDefinition

	for _, v := range t.functions {
		if v.Name == c.Tool {
			def = &v
		}
	}

	if def == nil {
		t.conf.Logger.Errorf("invalid function call: %#v", c)

		return llms.TextParts(schema.ChatMessageTypeHuman,
			"Tool does not exist. Please try again."), true
	}

	return def.Tool.Execute(c)
}
