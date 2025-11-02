package models

type StatusResponse struct {
	Status []InterfaceStatus `json:"status"`
	Links  []Link            `json:"_links"`
}

type InterfaceStatus struct {
	TestingInternetAccess bool      `json:"testing_internet_access"`
	RxBitrate             float64   `json:"rx_bitrate"`
	Ethernet              Ethernet  `json:"ethernet"`
	InternetAccess        bool      `json:"internet_access"`
	PrimaryInterface      bool      `json:"primary_interface"`
	TxBitrate             float64   `json:"tx_bitrate"`
	IP                    IPAddress `json:"ip"`
}

type Ethernet struct {
	Link    float64 `json:"link"`    // link speed in bps
	Duplex  string  `json:"duplex"`  // e.g., "full"
	Address string  `json:"address"` // MAC address
}

type IPAddress struct {
	Address string `json:"address"`
	Netmask string `json:"netmask"`
}
