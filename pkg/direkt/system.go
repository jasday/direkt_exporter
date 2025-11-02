package direkt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/vividbroadcast/prometheus-direkt-exporter/pkg/metrics"
	"github.com/vividbroadcast/prometheus-direkt-exporter/pkg/models"
)

// Metric names
const (
	SystemInfo            = "system_info"
	CPUUtilisationPercent = "cpu_utilisation_percent"
	MemoryTotalBytes      = "memory_total_bytes"
	MemoryAvailableBytes  = "memory_available_bytes"
	BondingPathRTTSeconds = "bonding_path_rtt_seconds"
	BondingPathRxBitrate  = "bonding_path_rx_bitrate_bytes_per_second"
	BondingPathTxBitrate  = "bonding_path_tx_bitrate_bytes_per_second"
	BondingPathHealth     = "bonding_path_health"
)

// Label names
const (
	LabelActiveFirmwareVersion  = "active_firmware_version"
	LabelBackupFirmwareVersion  = "backup_firmware_verison"
	LabelDefaultFirmwareVersion = "default_firmware_version"
	LabelNetworkInterface       = "network_interface"
)

var sysMetrics = []metrics.Gauge{
	{
		Name:   SystemInfo,
		Desc:   "Provides informtion on system uptime and statistics",
		Labels: []string{LabelActiveFirmwareVersion, LabelBackupFirmwareVersion, LabelDefaultFirmwareVersion},
	},
	{
		Name:   CPUUtilisationPercent,
		Desc:   "Percentage of CPU utilisation",
		Labels: []string{},
	},
	{
		Name:   MemoryTotalBytes,
		Desc:   "Total amount of memory in bytes",
		Labels: []string{},
	},
	{
		Name:   MemoryAvailableBytes,
		Desc:   "Amount of memory available in bytes",
		Labels: []string{},
	},
	{
		Name:   BondingPathRTTSeconds,
		Desc:   "Round trip time in seconds for bonding path",
		Labels: []string{LabelNetworkInterface},
	},
	{
		Name:   BondingPathRxBitrate,
		Desc:   "Receive bitrate in bytes per second for bonding path",
		Labels: []string{LabelNetworkInterface},
	},
	{
		Name:   BondingPathTxBitrate,
		Desc:   "Transmit bitrate in bytes per second for bonding path",
		Labels: []string{LabelNetworkInterface},
	},
	{
		Name:   BondingPathHealth,
		Desc:   "Network management bonding path health",
		Labels: []string{LabelNetworkInterface},
	},
}

func system(ctx context.Context, l zerolog.Logger, registry prometheus.Registerer, doReq func(l zerolog.Logger, request *http.Request) ([]byte, error), id string) error {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s/system/status", url, unitEndpoint, id), nil)
	if err != nil {
		return err
	}

	mtrcs := metrics.NewGaugeMap(sysMetrics)
	for _, metric := range mtrcs {
		registry.MustRegister(metric)
		metric.Reset()
	}

	var info models.SystemResponse
	var success float64 = 0

	res, err := doReq(l, request)
	if err == nil {
		err = json.Unmarshal(res, &info)
		if err == nil {
			l.Trace().Msg("Successfully retrieved metrics for system status")
			success = 1
			mtrcs[CPUUtilisationPercent].WithLabelValues().Set(info.CPU.Usage)
			mtrcs[MemoryAvailableBytes].WithLabelValues().Set(float64(info.Memory.Available))
			mtrcs[MemoryTotalBytes].WithLabelValues().Set(float64(info.Memory.Total))
			for _, path := range info.RemoteManagement.Bonding.Paths {
				ni := simplifyNetworkInterface(path.NetworkInterface)
				mtrcs[BondingPathRTTSeconds].WithLabelValues(ni).Set(path.RTT)
				mtrcs[BondingPathRxBitrate].WithLabelValues(ni).Set(float64(path.RxBitrate))
				mtrcs[BondingPathTxBitrate].WithLabelValues(ni).Set(float64(path.TxBitrate))
				mtrcs[BondingPathHealth].WithLabelValues(ni).Set(float64(metrics.StringBoolToInt(path.Health)))
			}
		}
	}

	// Set the gauge with labels from the parsed info
	mtrcs[SystemInfo].WithLabelValues(
		info.Firmware.Running.Version,
		info.Firmware.Recovery.Version,
		info.Firmware.Default.Version,
	).Set(float64(success))

	return err
}

// simplifyNetworkInterface extracts 0 out of "/api/v1/units/D02018/network_interfaces/0"
func simplifyNetworkInterface(fullPath string) string {
	parts := strings.Split(fullPath, "/")
	if len(parts) < 7 {
		return fullPath
	}

	return parts[6]
}

// var prettyJSON bytes.Buffer
// json.Indent(&prettyJSON, res, "", "  ")
// fmt.Println(prettyJSON.String())
