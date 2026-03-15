package domain

type PromptElement struct {
	Type           string             `json:"type"`                      // fixed | embedding | category | event | route
	RawPrompt      string             `json:"raw_prompt,omitempty"`
	Coords         []float64          `json:"coords,omitempty"`
	Name           string             `json:"name,omitempty"`
	Categories     map[string]float64 `json:"categories,omitempty"`
	IsPivotPoint   *bool              `json:"isPivotPoint,omitempty"`
	GeneratedPrompt string            `json:"generated_prompt,omitempty"`
	ParsedElements []PromptElement    `json:"parsed_elements,omitempty"`
	Time           float64            `json:"time,omitempty"`
}
