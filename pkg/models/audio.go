package models

type AudioFormat struct {
	SampleRate int `json:"sample_rate"`
	Channels   int `json:"channels"`
	BitDepth   int `json:"bit_depth"`
}

type AudioCodec struct {
	Name            string `json:"name"`
	AdaptiveBitrate bool   `json:"adaptive_bitrate"`
	Bitrate         int    `json:"bitrate"`
}
