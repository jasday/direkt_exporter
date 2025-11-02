package models

type VideoFormat struct {
	Interlaced        bool    `json:"interlaced"`
	BitDepth          int     `json:"bit_depth"`
	ForcedAspect      bool    `json:"forced_aspect"`
	PixelAspect       string  `json:"pixel_aspect"`
	Height            int     `json:"height"`
	TopFieldFirst     bool    `json:"top_field_first"`
	DisplayAspect     string  `json:"display_aspect"`
	ChromaSubsampling string  `json:"chroma_subsampling"`
	Framerate         float64 `json:"framerate"`
	Width             int     `json:"width"`
}

type VideoCodec struct {
	Bitrate                   int     `json:"bitrate"`
	AdaptiveBitrate           bool    `json:"adaptive_bitrate"`
	ConfiguredPerformanceMode string  `json:"configured_performance_mode"`
	BitrateBuffer             float64 `json:"bitrate_buffer"`
	Name                      string  `json:"name"`
	DefaultPerformanceMode    string  `json:"default_performance_mode"`
	Profile                   string  `json:"profile"`
	PerformanceMode           string  `json:"performance_mode"`
	Level                     string  `json:"level"`
}

type VideoSource struct {
	Source              string        `json:"source"`
	Audio               []AudioStream `json:"audio"`
	Video               VideoStream   `json:"video"`
	ProgramID           int           `json:"program_id"`
	Thumbnail           string        `json:"thumbnail"`
	Available           bool          `json:"available"`
	FallbackType        string        `json:"fallback_type"`
	FallbackDescription string        `json:"fallback_description"`
}
