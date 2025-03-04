package Common

import (
	"encoding/json"
	"fmt"
)

type Intent string

const (
	Lightson     Intent = "turn_on_light"
	Lightsoff    Intent = "turn_off_light"
	SetAlarm     Intent = "set_alarm"
	CheckWeather Intent = "check_weather"
	Talk         Intent = "talk"
	Unknown      Intent = "unknown"
	NotForMe     Intent = "notforme"
)

var Intents = []Intent{Lightson, Lightsoff, SetAlarm, CheckWeather, Talk, Unknown, NotForMe}

type OllamaIntentResponse struct {
	CreatedAt          string        `json:"created_at"`
	Done               bool          `json:"done"`
	DoneReason         string        `json:"done_reason"`
	EvalCount          int           `json:"eval_count"`
	EvalDuration       float64       `json:"eval_duration"`
	LoadDuration       float64       `json:"load_duration"`
	Message            NestedMessage `json:"message"`
	Model              string        `json:"model"`
	PromptEvalCount    int           `json:"prompt_eval_count"`
	PromptEvalDuration float64       `json:"prompt_eval_duration"`
	TotalDuration      float64       `json:"total_duration"`
}

type NestedMessage struct {
	Content ContentField `json:"content"`
	Role    string       `json:"role"`
}

type ContentField struct {
	Evaluation Intent
	Raw        string
}

func (c *ContentField) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal as a struct { evaluation: string }
	var obj struct {
		Evaluation string `json:"evaluation"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		c.Evaluation = Intent(obj.Evaluation)
		return nil
	}

	// Fallback: Assume it's a plain string
	var raw string
	if err := json.Unmarshal(data, &raw); err == nil {
		c.Raw = raw
		return nil
	}

	return fmt.Errorf("invalid content format")
}
