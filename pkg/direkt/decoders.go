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
	NetworkInputVideoStatus       = "network_input_video_status"
	NetworkInputAudioStatus       = "network_input_audio_status"
	NetworkInputVideoBitrate      = "network_input_video_bitrate"
	NetworkInputAudioBitrateBytes = "network_input_audio_bitrate"
	NetworkInputBitrate           = "network_input_bitrate"
	NetworkInputPacketLoss        = "network_input_packet_loss"
	NetworkInputEndToEndDelay     = "network_input_end_to_end_delay_seconds"
	NetworkInputBuffersReception  = "network_input_buffers_reception_seconds"
	NetworkInputBuffersDecoder    = "network_input_buffers_decoder_seconds"
	NetworkInputBuffersTarget     = "network_input_buffers_target_seconds"
	NetworkInputFecBuffer         = "network_input_fec_buffer_seconds"
	NetworkInputFecPacketLoss     = "network_input_fec_packet_loss"
	NetworkInputBondingBuffer     = "network_input_bonding_buffer_seconds"
	NetworkInputBondingPaths      = "network_input_bonding_paths"
	NetworkInputActive            = "network_input_active"
)

const (
	LabelInputIndex     = "input_index"
	LabelInputName      = "input_name"
	LabelActive         = "active"
	LabelCodec          = "codec"
	LabelCodecBitrate   = "codec_bitrate"
	LabelCodecLevel     = "codec_level"
	LabelSourceType     = "source_type"
	LabelTarget         = "target"
	LabelProtocol       = "protocol"
	LabelProgramIndex   = "program_number"
	LabelAddress        = "address"
	LabelSenderSerial   = "sender_serial"
	LabelSenderVerified = "sender_verified"
	LevelAudioIndex     = "audio_index"
)

var networkInputMetrics = []metrics.Gauge{
	{
		Name:   NetworkInputVideoStatus,
		Desc:   "Video input status (1=active, 0=inactive)",
		Labels: []string{LabelInputIndex, LabelInputName, LabelCodec, LabelProfile, LabelCodecLevel, LabelChromaSubsampling, LabelFramerate, LabelWidth, LabelHeight, LabelBitDepth, LabelInterlaced, LabelTopFieldFirst, LabelDisplayAspect, LabelPixelAspect, LabelForcedAspect, LabelProgramIndex},
	},
	{
		Name:   NetworkInputAudioStatus,
		Desc:   "Audio input status (1=active, 0=inactive)",
		Labels: []string{LabelInputIndex, LabelInputName, LabelCodec, LabelChannels, LabelSampleRate, LabelBitDepth, LabelProgramIndex, LevelAudioIndex},
	},
	{
		Name:   NetworkInputVideoBitrate,
		Desc:   "Video codec bitrate in bits per second",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProgramIndex},
	},
	// {
	// 	Name:   NetworkInputAudioBitrateBytes,
	// 	Desc:   "Audio codec bitrate in bytes",
	// 	Labels: []string{LabelInputIndex, LabelInputName, LabelCodec, LabelProgramNumber},
	// },
	{
		Name:   NetworkInputBitrate,
		Desc:   "Total network input bitrate in bits per second",
		Labels: []string{LabelInputIndex, LabelInputName, LabelSourceType, LabelSenderSerial},
	},
	{
		Name:   NetworkInputPacketLoss,
		Desc:   "Network input packet loss",
		Labels: []string{LabelInputIndex, LabelInputName},
	},
	{
		Name:   NetworkInputEndToEndDelay,
		Desc:   "End-to-end delay for the input in seconds",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProgramIndex, LabelTarget},
	},
	{
		Name:   NetworkInputBuffersReception,
		Desc:   "Reception buffer duration in seconds",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProgramIndex},
	},
	{
		Name:   NetworkInputBuffersDecoder,
		Desc:   "Decoder buffer duration in seconds",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProgramIndex},
	},
	{
		Name:   NetworkInputBuffersTarget,
		Desc:   "Target buffer duration in seconds",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProgramIndex},
	},
	{
		Name:   NetworkInputFecBuffer,
		Desc:   "FEC buffer duration in seconds",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProgramIndex},
	},
	{
		Name:   NetworkInputFecPacketLoss,
		Desc:   "FEC packet loss",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProgramIndex},
	},
	{
		Name:   NetworkInputBondingBuffer,
		Desc:   "Bonding buffer duration in seconds",
		Labels: []string{LabelInputIndex, LabelInputName, LabelProtocol},
	},
	{
		Name:   NetworkInputActive,
		Desc:   "Indicates if the network input is active (1=active, 0=inactive)",
		Labels: []string{LabelInputIndex, LabelInputName, LabelSourceType, LabelAddress, LabelSenderSerial, LabelSenderVerified},
	},
}

func decoders(ctx context.Context, l zerolog.Logger, registry prometheus.Registerer, doReq func(l zerolog.Logger, request *http.Request) ([]byte, error), id string) error {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s/network_inputs", url, unitEndpoint, id), nil)
	if err != nil {
		return err
	}

	res, err := doReq(l, request)
	if err != nil {
		return err
	}

	var decoders models.NetworkInputsResponse
	err = json.Unmarshal(res, &decoders)
	if err != nil {
		return err
	}

	mtrcs := metrics.NewGaugeMap(networkInputMetrics)
	for _, metric := range mtrcs {
		registry.MustRegister(metric)
		metric.Reset()
	}

	for _, decoder := range decoders.NetworkInputs {
		request, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s/network_inputs/%d/status", url, unitEndpoint, id, decoder.Index), nil)
		if err != nil {
			l.Err(err).Int("encoder_index", decoder.Index).Msg("Error creating encoder request, skipping")
			continue
		}
		res, err := doReq(l, request)
		if err != nil {
			l.Err(err).Int("encoder_index", decoder.Index).Msg("Error getting encoder metrics, skipping")
			continue
		}

		var e models.NetworkInputStatus
		err = json.Unmarshal(res, &e)
		if err != nil {
			l.Err(err).Int("encoder_index", decoder.Index).Msg("Error decoding encoder status response, skipping")
			continue
		}

		decoderIdx := strconv.Itoa(decoder.Index)

		mtrcs[NetworkInputActive].WithLabelValues(
			decoderIdx,
			e.Description,
			e.NetworkSource.SourceType,
			e.NetworkSource.Address,
			e.NetworkSource.Sender.Serial,
			metrics.BoolToString(e.NetworkSource.Sender.Verified),
		).Set(metrics.BoolToFloat64(e.Active))

		mtrcs[NetworkInputBitrate].WithLabelValues(
			decoderIdx,
			e.Description,
			e.NetworkSource.SourceType,
			e.NetworkSource.Sender.Serial,
		).Set(float64(e.NetworkSource.Bitrate))

		mtrcs[NetworkInputPacketLoss].WithLabelValues(
			decoderIdx,
			e.Description,
		).Set(e.NetworkSource.PacketLoss)

		mtrcs[NetworkInputFecBuffer].WithLabelValues(decoderIdx, e.Description, "0").Set(e.NetworkSource.FEC.Buffer)
		mtrcs[NetworkInputFecPacketLoss].WithLabelValues(decoderIdx, e.Description, "0").Set(e.NetworkSource.FEC.PacketLoss)

		if e.NetworkSource.Bonding.Protocol != "" {
			mtrcs[NetworkInputBondingBuffer].WithLabelValues(
				decoderIdx,
				e.Description,
				e.NetworkSource.Bonding.Protocol,
			).Set(e.NetworkSource.Bonding.Buffer)
		}

		for progIndex, prog := range e.NetworkSource.Programs {
			progIdxStr := fmt.Sprintf("%d", progIndex)

			// Video status
			video := prog.Video
			mtrcs[NetworkInputVideoStatus].WithLabelValues(
				decoderIdx,
				e.Description,
				video.Codec.Name,
				video.Codec.Profile,
				video.Codec.Level,
				video.Format.ChromaSubsampling,
				fmt.Sprintf("%.2f", video.Format.Framerate),
				strconv.Itoa(video.Format.Width),
				strconv.Itoa(video.Format.Height),
				strconv.Itoa(video.Format.BitDepth),
				metrics.BoolToString(video.Format.Interlaced),
				metrics.BoolToString(video.Format.TopFieldFirst),
				video.Format.DisplayAspect,
				video.Format.PixelAspect,
				metrics.BoolToString(video.Format.ForcedAspect),
				progIdxStr,
			).Set(metrics.BoolToFloat64(e.Active))

			// Input bitrate
			mtrcs[NetworkInputVideoBitrate].WithLabelValues(
				decoderIdx,
				e.Description,
				progIdxStr,
			).Set(float64(video.Codec.Bitrate))

			// Audio status
			for audioIdx, audio := range prog.Audio {
				audioIdxStr := fmt.Sprintf("%d", audioIdx)
				mtrcs[NetworkInputAudioStatus].WithLabelValues(
					decoderIdx,
					e.Description,
					audio.Codec.Name,
					strconv.Itoa(audio.Format.Channels),
					strconv.Itoa(audio.Format.SampleRate),
					strconv.Itoa(audio.Format.BitDepth),
					progIdxStr,
					audioIdxStr,
				).Set(metrics.BoolToFloat64(e.Active))
			}

			// Buffers
			mtrcs[NetworkInputBuffersReception].WithLabelValues(decoderIdx, e.Description, progIdxStr).Set(prog.Buffers.Reception)
			mtrcs[NetworkInputBuffersDecoder].WithLabelValues(decoderIdx, e.Description, progIdxStr).Set(prog.Buffers.Decoder)
			mtrcs[NetworkInputBuffersTarget].WithLabelValues(decoderIdx, e.Description, progIdxStr).Set(prog.Buffers.Target)

			// End-to-end delay
			mtrcs[NetworkInputEndToEndDelay].WithLabelValues(decoderIdx, e.Description, progIdxStr, fmt.Sprintf("%.4f", prog.EndToEndDelay.Target)).Set(prog.EndToEndDelay.Delay)
		}
	}

	return err
}
