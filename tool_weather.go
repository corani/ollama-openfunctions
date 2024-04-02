package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

// see https://www.weatherapi.com/api-explorer.aspx
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func NewWeatherTool(t *Tools) *WeatherTool {
	return &WeatherTool{
		conf: t.conf,
		log:  t.conf.Logger,
	}
}

type WeatherTool struct {
	conf *Config
	log  *log.Logger
}

func (t *WeatherTool) Execute(c *Call) (llms.MessageContent, bool) {
	// TODO(daniel): would be nice to dynamically parse out the arguments based on the schema.
	loc, ok := c.Input["location"].(string)
	if !ok {
		t.log.Errorf("function call: %#v", c)
		t.log.Fatal("invalid input")
	}
	unit, ok := c.Input["unit"].(string)
	if !ok {
		t.log.Errorf("function call: %#v", c)
		t.log.Fatal("invalid input")
	}

	t.log.Infof("Call tool: %s(%q, %q)", c.Tool, loc, unit)

	weather, err := t.getCurrentWeather(loc, unit)
	if err != nil {
		t.log.Errorf("function call: %#v", c)
		t.log.Fatal(err)
	}

	return llms.TextParts(schema.ChatMessageTypeHuman, weather), true
}

func (t *WeatherTool) getCurrentWeather(location string, unit string) (string, error) {
	if t.conf.WeatherAPI == "" {
		return "", fmt.Errorf("please set WEATHERAPI_KEY before use")
	}

	endpoint := fmt.Sprintf(
		"http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		t.conf.WeatherAPI, url.QueryEscape(location),
	)
	res, err := http.Get(endpoint)
	if err != nil {
		t.log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		bs, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("Weather API not available: code=%v, body=%v", res.StatusCode, string(bs))
	}

	var weather Weather
	if err := json.NewDecoder(res.Body).Decode(&weather); err != nil {
		t.log.Fatal(err)
	}

	weatherInfo := map[string]any{
		"location":  location,
		"condition": weather.Current.Condition.Text,
	}

	switch unit {
	case "fahrenheit":
		weatherInfo["temperature"] = weather.Current.TempF
		weatherInfo["unit"] = "fahrenheit"
	default:
		weatherInfo["temperature"] = weather.Current.TempC
		weatherInfo["unit"] = "celcius"
	}

	b, err := json.Marshal(weatherInfo)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func init() {
	register(func(t *Tools) {
		t.push(FunctionDefinition{
			Name:        "getCurrentWeather",
			Description: "Get the current weather in a given location",
			Parameters: json.RawMessage(`{
				"type": "object", 
				"properties": {
					"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, 
					"unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
				}, 
				"required": ["location", "unit"]
			}`),
			Response: json.RawMessage(`{
				"type": "object",
				"properties": {
					"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, 
					"unit": {"type": "string", "enum": ["celsius", "fahrenheit"]},
					"temperature": {"type": "integer", "description": "The current temperature"},
					"condition": {"type": "string", "description": "The current weather conditions"}
				}
			}`),
			Tool: NewWeatherTool(t),
		})
	})
}
