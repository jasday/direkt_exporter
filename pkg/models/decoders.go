package models

type NetworkInputsResponse struct {
	NetworkInputs []NetworkInput `json:"network_inputs"`
	Links         []Link         `json:"_links"`
}

type NetworkInput struct {
	Links       []Link `json:"_links"`
	Active      bool   `json:"active"`
	Name        string `json:"name"`
	Href        string `json:"href"`
	Index       int    `json:"index"`
	Type        string `json:"type"`
	Description string `json:"description"`
}
