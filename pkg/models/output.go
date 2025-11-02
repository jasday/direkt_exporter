package models

type VideoOutputStatus struct {
	Active      bool        `json:"active"`
	VideoSource VideoSource `json:"video_source"`
	Links       []Link      `json:"_links"`
	Messages    []string    `json:"messages"`
	VideoOut    VideoOut    `json:"video_out"`
	Description string      `json:"description"`
}

type AudioStream struct {
	Format AudioFormat `json:"format"`
	Codec  AudioCodec  `json:"codec,omitempty"` // Some entries may not have codec
}

type VideoStream struct {
	Format VideoFormat `json:"format"`
	Codec  VideoCodec  `json:"codec,omitempty"` // Some entries may not have codec
}

type VideoOut struct {
	ConnectorName string        `json:"connector_name"`
	Audio         []AudioStream `json:"audio"`
	Video         VideoStream   `json:"video"`
}
