package models

type EncoderStatus struct {
	VideoSource  VideoSource         `json:"video_source"`
	Destinations EncoderDestinations `json:"destinations"`
	Links        []Link              `json:"_links"`
	Description  string              `json:"description"`
	Active       bool                `json:"active"`
	Recording    map[string]any      `json:"recording"`
	Messages     []any               `json:"messages"`
	Encoding     EncodingStatus      `json:"encoding"`
}

type EncoderDestinations struct {
	TCPOnRequest TCPOnRequest  `json:"tcp_on_request"`
	SRTOnRequest SRTOnRequest  `json:"srt_on_request"`
	RTMP         []any         `json:"rtmp"`
	Basic        []BasicOutput `json:"basic"`
}

type TCPOnRequest struct {
	Clients []any `json:"clients"`
}

type SRTOnRequest struct {
	Clients []any `json:"clients"`
}

type BasicOutput struct {
	Messages           []any       `json:"messages"`
	UDPSmoothingBuffer float64     `json:"udp_smoothing_buffer"`
	FEC                FECStatus   `json:"fec"`
	Bonding            BondingInfo `json:"bonding"`
	PacketLoss         float64     `json:"packet_loss"`
	ID                 string      `json:"id"`
	Bitrate            float64     `json:"bitrate"`
}

type FECStatus struct {
	PacketLoss      float64 `json:"packet_loss"`
	BitrateOverhead float64 `json:"bitrate_overhead"`
}

type BondingInfo struct {
	Paths             []EncoderBondingPath `json:"paths"`
	Destinations      []string             `json:"destinations"`
	Bitrate           float64              `json:"bitrate"`
	EstimatedCapacity float64              `json:"estimated_capacity"`
	EstimateIsMax     bool                 `json:"estimate_is_max"`
	Destination       string               `json:"destination"`
	FailoverActive    bool                 `json:"failover_active"`
}

type EncoderBondingPath struct {
	LatencyHistory    float64 `json:"latency_history"`
	EstimateIsMax     bool    `json:"estimate_is_max"`
	PacketLateHistory float64 `json:"packet_late_history"`
	Messages          []any   `json:"messages"`
	Destination       string  `json:"destination"`
	RedundancyBitrate float64 `json:"redundancy_bitrate"`
	PacketLossHistory float64 `json:"packet_loss_history"`
	Bitrate           float64 `json:"bitrate"`
	Viable            bool    `json:"viable"`
	EstimatedCapacity float64 `json:"estimated_capacity"`
	PacketLate        int     `json:"packet_late"`
	NetworkInterface  string  `json:"network_interface"`
	Latency           float64 `json:"latency"`
	PacketLoss        float64 `json:"packet_loss"`
}

type EncodingStatus struct {
	Audio        []AudioStream `json:"audio"`
	TotalBitrate float64       `json:"total_bitrate"`
	Video        VideoStream   `json:"video"`
}
