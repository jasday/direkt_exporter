package models

type DecoderVideo struct {
	Format DecoderVideoFormat `json:"format"`
	Codec  DecoderVideoCodec  `json:"codec"`
}

type DecoderVideoFormat struct {
	Framerate         float64 `json:"framerate"`
	DisplayAspect     string  `json:"display_aspect"`
	TopFieldFirst     bool    `json:"top_field_first"`
	ChromaSubsampling string  `json:"chroma_subsampling"`
	Width             int     `json:"width"`
	ForcedAspect      bool    `json:"forced_aspect"`
	BitDepth          int     `json:"bit_depth"`
	Interlaced        bool    `json:"interlaced"`
	PixelAspect       string  `json:"pixel_aspect"`
	Height            int     `json:"height"`
}

type DecoderVideoCodec struct {
	Level           string `json:"level"`
	Profile         string `json:"profile"`
	Name            string `json:"name"`
	AdaptiveBitrate bool   `json:"adaptive_bitrate"`
	Bitrate         int64  `json:"bitrate"`
}

type DecoderAudioFormat struct {
	BitDepth   int `json:"bit_depth"`
	SampleRate int `json:"sample_rate"`
	Channels   int `json:"channels"`
}

type DecoderAudioCodec struct {
	Name            string `json:"name"`
	AdaptiveBitrate bool   `json:"adaptive_bitrate"`
}

type Audio struct {
	Codec  DecoderAudioCodec  `json:"codec"`
	Format DecoderAudioFormat `json:"format"`
}

type EndToEndDelay struct {
	Delay  float64 `json:"delay"`
	Target float64 `json:"target"`
}

type Buffers struct {
	Reception float64 `json:"reception"`
	Target    float64 `json:"target"`
	Decoder   float64 `json:"decoder"`
}

type Program struct {
	Video         DecoderVideo  `json:"video"`
	Audio         []Audio       `json:"audio"`
	Messages      []string      `json:"messages"`
	Thumbnail     string        `json:"thumbnail"`
	Number        int           `json:"number"`
	ID            string        `json:"id"`
	EndToEndDelay EndToEndDelay `json:"end_to_end_delay"`
	Buffers       Buffers       `json:"buffers"`
}

type DecoderBondingPath struct {
	Address  string   `json:"address"`
	Messages []string `json:"messages"`
}

type DecoderBonding struct {
	Buffer   float64              `json:"buffer"`
	Protocol string               `json:"protocol"`
	Paths    []DecoderBondingPath `json:"paths"`
}

type FEC struct {
	Buffer     float64 `json:"buffer"`
	PacketLoss float64 `json:"packet_loss"`
}

type Sender struct {
	Serial   string `json:"serial"`
	Verified bool   `json:"verified"`
}

type NetworkSource struct {
	Programs   []Program      `json:"programs"`
	Encrypted  bool           `json:"encrypted"`
	FEC        FEC            `json:"fec"`
	Bonding    DecoderBonding `json:"bonding"`
	Bitrate    int64          `json:"bitrate"`
	SourceType string         `json:"source_type"`
	Address    string         `json:"address"`
	Sender     Sender         `json:"sender"`
	PacketLoss float64        `json:"packet_loss"`
}

type DestinationsClients struct {
	Clients []string `json:"clients"`
}

type DecoderDestinations struct {
	Basic        []string            `json:"basic"`
	TCPOnRequest DestinationsClients `json:"tcp_on_request"`
	SRTOnRequest DestinationsClients `json:"srt_on_request"`
	RTMP         []string            `json:"rtmp"`
}

type NetworkInputStatus struct {
	Description   string              `json:"description"`
	NetworkSource NetworkSource       `json:"network_source"`
	Recording     interface{}         `json:"recording"`
	Messages      []string            `json:"messages"`
	Links         []Link              `json:"_links"`
	Active        bool                `json:"active"`
	Destinations  DecoderDestinations `json:"destinations"`
}
