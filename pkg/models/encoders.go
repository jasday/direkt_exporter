package models

type EncodersResponse struct {
	Encoders []Encoder `json:"encoders"`
	Links    []Link    `json:"_links"`
}

type Encoder struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	Href        string `json:"href"`
	Index       int    `json:"index"`
	Type        string `json:"type"`
	Links       []Link `json:"_links"`
}
