package direkt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/vividbroadcast/prometheus-direkt-exporter/pkg/metrics"
	"github.com/vividbroadcast/prometheus-direkt-exporter/pkg/models"
)

const (
	MetricEncoderVideoStatus                              = "encoder_video_config"
	MetricEncoderAudioStatus                              = "encoder_audio_config"
	MetricEncoderVideoInputStatus                         = "encoder_video_input_status"
	MetricEncoderAudioInputStatus                         = "encoder_audio_input_status"
	MetricEncoderTotalBitrate                             = "encoder_total_bitrate_bits"
	MetricEncoderDestinationBitrate                       = "encoder_basic_destination_bitrate_bits"
	MetricEncoderBasicDestinationPacketLoss               = "encoder_basic_destination_packet_loss"
	MetricEncoderBasicDestinationFECPacketLoss            = "encoder_basic_destination_fec_packet_loss"
	MetricEncoderBasicDestinationFECOverhead              = "encoder_basic_destination_fec_bitrate_overhead"
	MetricEncoderBasicDestinationUDPSmoothingBuffer       = "encoder_basic_destination_udp_smoothing_buffer_seconds"
	MetricEncoderBasicDestinationPathLatency              = "encoder_basic_destination_path_latency_seconds"
	MetricEncoderBasicDestinationPathLatencyHistorical    = "encoder_basic_destination_path_latency_historical_seconds"
	MetricEncoderBasicDestinationPathViable               = "encoder_basic_destination_path_viable"
	MetricEncoderBasicDestinationPathBitrate              = "encoder_basic_destination_path_bitrate_bits"
	MetricEncoderBasicDestinationPathPacketLoss           = "encoder_basic_destination_path_packet_loss"
	MetricEncoderBasicDestinationPathPacketLossHistorical = "encoder_basic_destination_path_packet_loss_historical"
	MetricEncoderBasicDestinationPathCapacity             = "encoder_basic_destination_path_estimated_capacity_bits"
	MetricEncoderBasicDestinationPathRedundancy           = "encoder_basic_destination_path_redundancy_bitrate_bits"
	MetricEncoderBasicDestinationFailoverActive           = "encoder_basic_destination_failover_active"
)

// Global label names
const (
	LabelEncoderIndex       = "encoder_index"
	LabelEncoderName        = "encoder_name"
	LabelDestination        = "destination"
	LabelDestinationIndex   = "destination_index"
	LabelBondingDestination = "bonding_destination"

	// Video format labels
	LabelInterlaced        = "interlaced"
	LabelChromaSubsampling = "chroma_subsampling"
	LabelFramerate         = "framerate"
	LabelBitDepth          = "bit_depth"
	LabelDisplayAspect     = "display_aspect"
	LabelWidth             = "width"
	LabelPixelAspect       = "pixel_aspect"
	LabelForcedAspect      = "forced_aspect"
	LabelHeight            = "height"
	LabelTopFieldFirst     = "top_field_first"
	LabelTargetBitrate     = "target_bitrate"

	// Video codec labels
	LabelCodecName                 = "codec_name"
	LabelAdaptiveBitrate           = "adaptive_bitrate"
	LabelConfiguredPerformanceMode = "configured_performance_mode"
	LabelBitrateBuffer             = "bitrate_buffer"
	LabelDefaultPerformanceMode    = "default_performance_mode"
	LabelProfile                   = "profile"
	LabelPerformanceMode           = "performance_mode"
	LabelLevel                     = "level"

	// Audio codec labels
	LabelSampleRate = "sample_rate"
	LabelChannels   = "channels"
)

var encoderMetrics = []metrics.Gauge{
	// --- VIDEO STATUS ---
	{
		Name: MetricEncoderVideoInputStatus,
		Desc: "Video input and encoding status. Value is 1 if source is available, 0 otherwise. Encoding and format info are exposed as labels.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelFramerate,
			LabelWidth,
			LabelHeight,
			LabelBitDepth,
			LabelInterlaced,
			LabelTopFieldFirst,
			LabelChromaSubsampling,
			LabelDisplayAspect,
			LabelPixelAspect,
			LabelForcedAspect,
		},
	},
	{
		Name: MetricEncoderAudioInputStatus,
		Desc: "Audio encoding status. Value is 1 if source is available, 0 otherwise. Audio properties are in labels.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelCodec,
			LabelAudioSampleRate,
			LabelAudioChannels,
		},
	},
	{
		Name: MetricEncoderVideoStatus,
		Desc: "Video encoder configuration. Exposes all static and configured parameters as labels. Value always 1 if encoder is active.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelCodec,
			LabelProfile,
			LabelCodecLevel,
			LabelTargetBitrate,
			LabelFramerate,
			LabelWidth,
			LabelHeight,
			LabelBitDepth,
			LabelInterlaced,
			LabelTopFieldFirst,
			LabelChromaSubsampling,
			LabelDisplayAspect,
			LabelPixelAspect,
			LabelForcedAspect,
		},
	},
	{
		Name: MetricEncoderAudioStatus,
		Desc: "Audio encoder configuration. Exposes all static and configured parameters as labels. Value always 1 if encoder is active.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelAudioIndex,
			LabelCodec,
			LabelTargetBitrate,
			LabelAudioSampleRate,
			LabelAudioChannels,
		},
	},
	{
		Name: MetricEncoderTotalBitrate,
		Desc: "Total encoder bitrate in bits per second (sum of video and audio streams).",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
		},
	},
	{
		Name: MetricEncoderDestinationBitrate,
		Desc: "Output bitrate to a destination in bits per second.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelDestination,
			LabelDestinationIndex,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPacketLoss,
		Desc: "Overall packet loss ratio for a destination (0-1).",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelDestination,
			LabelDestinationIndex,
		},
	},
	{
		Name: MetricEncoderBasicDestinationFECPacketLoss,
		Desc: "FEC packet loss for a destination (fractional).",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelDestination,
			LabelDestinationIndex,
		},
	},
	{
		Name: MetricEncoderBasicDestinationFECOverhead,
		Desc: "FEC bitrate overhead ratio for a destination (fractional).",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelDestination,
			LabelDestinationIndex,
		},
	},
	{
		Name: MetricEncoderBasicDestinationUDPSmoothingBuffer,
		Desc: "UDP smoothing buffer duration for a destination in seconds.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelDestination,
			LabelDestinationIndex,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathLatency,
		Desc: "Current latency in seconds for a destination path.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathLatencyHistorical,
		Desc: "Historical average latency in seconds for a destination path.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathViable,
		Desc: "1 if the destination path is viable, 0 otherwise.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathBitrate,
		Desc: "Current bitrate on a specific destination path (bits per second).",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathPacketLoss,
		Desc: "Current packet loss ratio (0-1) on a destination path.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathPacketLossHistorical,
		Desc: "Historical average packet loss ratio (0-1) on a destination path.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathCapacity,
		Desc: "Estimated capacity in bits per second for a destination path.",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationPathRedundancy,
		Desc: "Configured redundancy bitrate for a destination path (bits per second).",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelBondingDestination,
			LabelDestination,
			LabelDestinationIndex,
			LabelNetworkInterface,
		},
	},
	{
		Name: MetricEncoderBasicDestinationFailoverActive,
		Desc: "Configured redundancy bitrate for a destination path (bits per second).",
		Labels: []string{
			LabelEncoderIndex,
			LabelEncoderName,
			LabelDestination,
			LabelDestinationIndex,
		},
	},
}

func encoders(ctx context.Context, l zerolog.Logger, registry prometheus.Registerer, doReq func(l zerolog.Logger, request *http.Request) ([]byte, error), id string) error {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s/encoders", url, unitEndpoint, id), nil)
	if err != nil {
		return err
	}

	res, err := doReq(l, request)
	if err != nil {
		return err
	}
	var encoders models.EncodersResponse
	err = json.Unmarshal(res, &encoders)
	if err != nil {
		return err
	}

	mtrcs := metrics.NewGaugeMap(encoderMetrics)
	for _, metric := range mtrcs {
		registry.MustRegister(metric)
		metric.Reset()
	}

	for _, encoder := range encoders.Encoders {
		request, err = http.NewRequest("GET", fmt.Sprintf("%s%s%s/encoders/%d/status", url, unitEndpoint, id, encoder.Index), nil)
		if err != nil {
			l.Err(err).Int("encoder_index", encoder.Index).Msg("Error creating encoder request, skipping")
			continue
		}
		res, err := doReq(l, request)
		if err != nil {
			l.Err(err).Int("encoder_index", encoder.Index).Msg("Error getting encoder metrics, skipping")
			continue
		}
		var e models.EncoderStatus
		err = json.Unmarshal(res, &e)
		if err != nil {
			l.Err(err).Int("encoder_index", encoder.Index).Msg("Error decoding encoder status response, skipping")
			continue
		}

		encoderIdx := strconv.Itoa(encoder.Index)

		mtrcs[MetricEncoderVideoInputStatus].WithLabelValues(
			encoderIdx,
			e.Description,
			fmt.Sprintf("%.2f", e.VideoSource.Video.Format.Framerate),
			strconv.Itoa(e.VideoSource.Video.Format.Width),
			strconv.Itoa(e.VideoSource.Video.Format.Height),
			strconv.Itoa(e.VideoSource.Video.Format.BitDepth),
			metrics.BoolToString(e.VideoSource.Video.Format.Interlaced),
			metrics.BoolToString(e.VideoSource.Video.Format.TopFieldFirst),
			e.VideoSource.Video.Format.ChromaSubsampling,
			e.VideoSource.Video.Format.DisplayAspect,
			e.VideoSource.Video.Format.PixelAspect,
			metrics.BoolToString(e.VideoSource.Video.Format.ForcedAspect),
		).Set(metrics.BoolToFloat64(e.VideoSource.Available))

		// TODO
		// mtrcs[MetricEncoderAudioInputStatus].WithLabelValues(
		// 	encoderIdx,
		// 	e.Description,
		// 	fmt.Sprintf("%.2f", e.VideoSource.Video.Format.Framerate),
		// 	strconv.Itoa(e.VideoSource.Video.Format.Width),
		// 	strconv.Itoa(e.VideoSource.Video.Format.Height),
		// 	strconv.Itoa(e.VideoSource.Video.Format.BitDepth),
		// 	metrics.BoolToString(e.VideoSource.Video.Format.Interlaced),
		// 	metrics.BoolToString(e.VideoSource.Video.Format.TopFieldFirst),
		// 	e.VideoSource.Video.Format.ChromaSubsampling,
		// 	e.VideoSource.Video.Format.DisplayAspect,
		// 	e.VideoSource.Video.Format.PixelAspect,
		// 	metrics.BoolToString(e.VideoSource.Video.Format.ForcedAspect),
		// ).Set(metrics.BoolToFloat64(e.VideoSource.Available))

		mtrcs[MetricEncoderVideoStatus].WithLabelValues(
			encoderIdx,
			e.Description,
			e.Encoding.Video.Codec.Name,
			e.Encoding.Video.Codec.Profile,
			e.Encoding.Video.Codec.Level,
			strconv.Itoa(e.Encoding.Video.Codec.Bitrate),
			fmt.Sprintf("%.2f", e.Encoding.Video.Format.Framerate),
			strconv.Itoa(e.Encoding.Video.Format.Width),
			strconv.Itoa(e.Encoding.Video.Format.Height),
			strconv.Itoa(e.Encoding.Video.Format.BitDepth),
			metrics.BoolToString(e.Encoding.Video.Format.Interlaced),
			metrics.BoolToString(e.Encoding.Video.Format.TopFieldFirst),
			e.Encoding.Video.Format.ChromaSubsampling,
			e.Encoding.Video.Format.DisplayAspect,
			e.Encoding.Video.Format.PixelAspect,
			metrics.BoolToString(e.Encoding.Video.Format.ForcedAspect),
		).Set(metrics.BoolToFloat64(e.Active))

		for i, audio := range e.Encoding.Audio {
			mtrcs[MetricEncoderAudioStatus].WithLabelValues(
				encoderIdx,
				e.Description,
				strconv.Itoa(i),
				audio.Codec.Name,
				strconv.Itoa(audio.Codec.Bitrate),
				strconv.Itoa(audio.Format.SampleRate),
				strconv.Itoa(audio.Format.Channels),
			).Set(metrics.BoolToFloat64(e.Active))
		}

		mtrcs[MetricEncoderTotalBitrate].WithLabelValues(
			encoderIdx,
			e.Description,
		).Set(e.Encoding.TotalBitrate)

		for i, basic := range e.Destinations.Basic {
			destinationIdx := strconv.Itoa(i)
			mtrcs[MetricEncoderDestinationBitrate].WithLabelValues(
				encoderIdx,
				e.Description,
				basic.Bonding.Destination,
				destinationIdx,
			).Set(basic.Bitrate)

			mtrcs[MetricEncoderBasicDestinationPacketLoss].WithLabelValues(
				encoderIdx,
				e.Description,
				basic.Bonding.Destination,
				destinationIdx,
			).Set(basic.PacketLoss)

			mtrcs[MetricEncoderBasicDestinationFECPacketLoss].WithLabelValues(
				encoderIdx,
				e.Description,
				basic.Bonding.Destination,
				destinationIdx,
			).Set(basic.FEC.PacketLoss)

			mtrcs[MetricEncoderBasicDestinationFECOverhead].WithLabelValues(
				encoderIdx,
				e.Description,
				basic.Bonding.Destination,
				destinationIdx,
			).Set(basic.FEC.BitrateOverhead)

			mtrcs[MetricEncoderBasicDestinationUDPSmoothingBuffer].WithLabelValues(
				encoderIdx,
				e.Description,
				basic.Bonding.Destination,
				destinationIdx,
			).Set(basic.UDPSmoothingBuffer)

			mtrcs[MetricEncoderBasicDestinationFailoverActive].WithLabelValues(
				encoderIdx,
				e.Description,
				basic.Bonding.Destination,
				destinationIdx,
			).Set(metrics.BoolToFloat64(basic.Bonding.FailoverActive))

			for _, path := range basic.Bonding.Paths {
				mtrcs[MetricEncoderBasicDestinationPathLatency].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(path.Latency)

				mtrcs[MetricEncoderBasicDestinationPathLatencyHistorical].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(path.LatencyHistory)

				mtrcs[MetricEncoderBasicDestinationPathViable].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(metrics.BoolToFloat64(path.Viable))

				mtrcs[MetricEncoderBasicDestinationPathBitrate].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(path.Bitrate)

				mtrcs[MetricEncoderBasicDestinationPathPacketLoss].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(path.PacketLoss)

				mtrcs[MetricEncoderBasicDestinationPathPacketLossHistorical].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(path.PacketLossHistory)

				mtrcs[MetricEncoderBasicDestinationPathCapacity].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(path.EstimatedCapacity)

				mtrcs[MetricEncoderBasicDestinationPathRedundancy].WithLabelValues(
					encoderIdx,
					e.Description,
					basic.Bonding.Destination,
					path.Destination,
					destinationIdx,
					simplifyNetworkInterface(path.NetworkInterface),
				).Set(path.RedundancyBitrate)
			}
		}

	}
	return err
}
