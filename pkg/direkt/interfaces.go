package direkt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/vividbroadcast/prometheus-direkt-exporter/pkg/metrics"
	"github.com/vividbroadcast/prometheus-direkt-exporter/pkg/models"
)

const (
	InterfaceRxBitrate             = "interface_rx_bitrate_bytes_per_second"
	InterfaceTxBitrate             = "interface_tx_bitrate_bytes_per_second"
	InterfaceLinkSpeed             = "interface_link_speed_bits_per_second"
	InterfaceInternetAccess        = "interface_internet_access"
	InterfaceTestingInternetAccess = "interface_testing_internet_access"
)

const (
	LabelInterfaceMAC     = "interface_mac"
	LabelIPAddress        = "ip_address"
	LabelPrimaryInterface = "primary_interface"
)

var interfaceMetrics = []metrics.Gauge{
	{
		Name:   InterfaceRxBitrate,
		Desc:   "Receive bitrate in bits per second for the interface",
		Labels: []string{LabelInterfaceMAC, LabelIPAddress, LabelPrimaryInterface},
	},
	{
		Name:   InterfaceTxBitrate,
		Desc:   "Transmit bitrate in bits per second for the interface",
		Labels: []string{LabelInterfaceMAC, LabelIPAddress, LabelPrimaryInterface},
	},
	{
		Name:   InterfaceLinkSpeed,
		Desc:   "Ethernet link speed in bits per second. -1 if unknown",
		Labels: []string{LabelInterfaceMAC, LabelIPAddress, LabelPrimaryInterface},
	},
	{
		Name:   InterfaceInternetAccess,
		Desc:   "Boolean indicating if the interface has internet access (1 = yes, 0 = no)",
		Labels: []string{LabelInterfaceMAC, LabelIPAddress, LabelPrimaryInterface},
	},
	{
		Name:   InterfaceTestingInternetAccess,
		Desc:   "Boolean indicating if the interface is testing internet access (1 = yes, 0 = no)",
		Labels: []string{LabelInterfaceMAC, LabelIPAddress, LabelPrimaryInterface},
	},
}

func interfaces(ctx context.Context, l zerolog.Logger, registry prometheus.Registerer, doReq func(l zerolog.Logger, request *http.Request) ([]byte, error), id string) error {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s/network_interfaces/status", url, unitEndpoint, id), nil)
	if err != nil {
		return err
	}

	mtrcs := metrics.NewGaugeMap(interfaceMetrics)
	for _, metric := range mtrcs {
		registry.MustRegister(metric)
		metric.Reset()
	}

	res, err := doReq(l, request)
	if err != nil {
		return err
	}

	var interfaces models.StatusResponse
	err = json.Unmarshal(res, &interfaces)
	if err != nil {
		return err
	}

	l.Trace().Msg("Successfully retrieved metrics for network interfaces status")

	for _, nwint := range interfaces.Status {
		mac := nwint.Ethernet.Address
		ip := nwint.IP.Address
		isPri := nwint.PrimaryInterface
		mtrcs[InterfaceRxBitrate].WithLabelValues(mac, ip, metrics.BoolToString(isPri)).Set(nwint.RxBitrate)
		mtrcs[InterfaceTxBitrate].WithLabelValues(mac, ip, metrics.BoolToString(isPri)).Set(nwint.TxBitrate)
		mtrcs[InterfaceLinkSpeed].WithLabelValues(mac, ip, metrics.BoolToString(isPri)).Set(nwint.Ethernet.Link)
		mtrcs[InterfaceInternetAccess].WithLabelValues(mac, ip, metrics.BoolToString(isPri)).Set(metrics.BoolToFloat64(nwint.InternetAccess))
		mtrcs[InterfaceTestingInternetAccess].WithLabelValues(mac, ip, metrics.BoolToString(isPri)).Set(metrics.BoolToFloat64(nwint.TestingInternetAccess))
	}

	return err
}
