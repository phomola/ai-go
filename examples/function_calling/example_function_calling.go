package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/phomola/ai-go/gemini/ai"
)

type (
	period struct {
		Number          int       `json:"number"`
		Name            string    `json:"name"`
		Start           time.Time `json:"startTime"`
		End             time.Time `json:"endTime"`
		Temperature     float32   `json:"temperature"`
		TemperatureUnit string    `json:"temperatureUnit"`
		Detailed        string    `json:"detailedForecast"`
	}

	weatherRequest struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	weatherResponse struct {
		City    string   `json:"city"`
		State   string   `json:"state"`
		Periods []period `json:"periods"`
	}

	toolOutput struct {
		Periods []modelPeriod `json:"periods"`
		Error   string        `json:"error" jsonschema:"An error occurred while getting the forecast."`
	}

	modelPeriod struct {
		Name     string `json:"name" jsonschema:"The name of the period."`
		Forecast string `json:"forecast" jsonschema:"The forecast for the period."`
	}
)

func main() {
	ctx := context.Background()

	cl, err := ai.NewClient(ctx, ai.Gemini3FlashPreview)
	if err != nil {
		log.Fatal(err)
	}

	var weatherTool ai.Tool
	if err := ai.AddFunction(&weatherTool, "weatherTool", "Provides weather forecasts for the US.", func(in *weatherRequest) (*toolOutput, error) {
		b, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://arax.ee/weather/forecast", bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		var cl http.Client
		resp, err := cl.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			return &toolOutput{Error: "The location couldn't be found. The tool only provides data for the US."}, nil
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
		}
		var out weatherResponse
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return nil, err
		}
		fmt.Printf("got weather forecast for %s/%s\n", out.City, out.State)
		periods := make([]modelPeriod, 0, 5)
		for i, p := range out.Periods {
			if i == 5 {
				break
			}
			periods = append(periods, modelPeriod{p.Name, p.Detailed})
		}
		return &toolOutput{Periods: periods}, nil
	}); err != nil {
		log.Fatal(err)
	}

	resp, err := cl.GenerateText(ctx, ai.NewText("What's the weather forecast for Seattle?"), []*ai.Tool{&weatherTool})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
