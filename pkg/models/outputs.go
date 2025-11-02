package models

type VideoOutputsResponse struct {
	VideoOutputs []VideoOutput `json:"video_outputs"`
	Links        []Link        `json:"_links"`
}

type VideoOutput struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	VideoPort   VideoPort `json:"video_port"`
	Name        string    `json:"name"`
	Links       []Link    `json:"_links"`
	Active      bool      `json:"active"`
	Href        string    `json:"href"`
	Index       int       `json:"index"`
}

type VideoPort struct {
	VideoCard string `json:"video_card"`
	PortIndex int    `json:"port_index"`
	Usable    bool   `json:"usable"`
}
