package models

import "time"

type SystemResponse struct {
	Memory              Memory           `json:"memory"`
	Firmware            Firmware         `json:"firmware"`
	UpgradeMediaPresent bool             `json:"upgrade_media_present"`
	CPU                 CPU              `json:"cpu"`
	UpgradeServer       UpgradeServer    `json:"upgrade_server"`
	Datetime            time.Time        `json:"datetime"`
	RemoteManagement    RemoteManagement `json:"remote_management"`
}

// Memory holds memory info
type Memory struct {
	Available uint64 `json:"available"`
	Total     uint64 `json:"total"`
}

// Firmware contains firmware details
type Firmware struct {
	Running  FirmwareVersion `json:"running"`
	Recovery FirmwareVersion `json:"recovery"`
	Default  FirmwareVersion `json:"default"`
}

// FirmwareVersion holds firmware version info
type FirmwareVersion struct {
	Version  string    `json:"version"`
	Datetime time.Time `json:"datetime"`
}

// CPU holds CPU usage info
type CPU struct {
	Usage float64 `json:"usage"`
}

// UpgradeServer holds upgrade server info
type UpgradeServer struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// RemoteManagement holds remote management details
type RemoteManagement struct {
	Bonding           Bonding `json:"bonding"`
	StatusDescription string  `json:"status_description"`
	NetworkInterface  string  `json:"network_interface"`
	Connected         bool    `json:"connected"`
	ViaHTTP           bool    `json:"via_http"`
	Address           string  `json:"address"`
}

// Bonding holds bonding info for remote management
type Bonding struct {
	Paths    []BondingPath `json:"paths"`
}

// BondingPath holds data about each bonding path
type BondingPath struct {
	HttpsConnectivityStatus string  `json:"https_connectivity_status"`
	NetworkInterface        string  `json:"network_interface"`
	SilenceTime             float64 `json:"silence_time"`
	RTT                     float64 `json:"rtt"`
	RxBitrate               int     `json:"rx_bitrate"`
	ViaHTTPS                bool    `json:"via_https"`
	Health                  string  `json:"health"`
	TxBitrate               int     `json:"tx_bitrate"`
}
