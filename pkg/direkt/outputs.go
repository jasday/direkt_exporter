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
	OutputVideoActive          = "output_video_active"
	OutputAudioActive          = "output_audio_active"
	OutputVideoSourceAvailable = "output_video_source_available"
	OutputAudioSourceAvailable = "output_audio_source_available"
)

const (
	LabelSourceAvailable = "source_available"
	LabelOutputIndex     = "output_index"
	LabelOutputName      = "output_name"
	LabelSourceIndex     = "source_index"
	LabelSourceName      = "source_name"
	LabelAudioCodecName  = "audio_codec_name"
	LabelAudioChannels   = "audio_channels"
	LabelAudioSampleRate = "audio_sample_rate"
	LabelAudioBitDepth   = "audio_bit_depth"
	LabelAudioIndex      = "audio_index"
)

var videoMetrics = []metrics.Gauge{
	{
		Name: OutputVideoActive,
		Desc: "Indicates if the video output is active (1=active, 0=inactive) with video format properties as labels",
		Labels: []string{
			LabelOutputIndex,
			LabelOutputName,
			LabelWidth,
			LabelHeight,
			LabelFramerate,
			LabelBitDepth,
			LabelInterlaced,
			LabelChromaSubsampling,
			LabelPixelAspect,
			LabelDisplayAspect,
			LabelTopFieldFirst,
		},
	},
	{
		Name: OutputAudioActive,
		Desc: "Indicates if the audio output is active (1=active, 0=inactive) with audio format properties as labels",
		Labels: []string{
			LabelOutputIndex,
			LabelOutputName,
			LabelAudioIndex,
			LabelAudioChannels,
			LabelAudioSampleRate,
			LabelAudioBitDepth,
		},
	},
	{
		Name: OutputVideoSourceAvailable,
		Desc: "Indicates if the video source is available (1=active, 0=inactive) with video format properties as labels",
		Labels: []string{
			LabelSourceIndex,
			LabelSourceName,
			LabelCodecName,
			LabelCodecBitrate,
			LabelProfile,
			LabelLevel,
			LabelWidth,
			LabelHeight,
			LabelFramerate,
			LabelBitDepth,
			LabelInterlaced,
			LabelChromaSubsampling,
			LabelPixelAspect,
			LabelDisplayAspect,
			LabelTopFieldFirst,
		},
	},
	{
		Name: OutputAudioSourceAvailable,
		Desc: "Indicates if the audio source is available (1=active, 0=inactive) with audio format properties as labels",
		Labels: []string{
			LabelSourceIndex,
			LabelSourceName,
			LabelAudioIndex,
			LabelAudioCodecName,
			LabelAudioChannels,
			LabelAudioSampleRate,
			LabelAudioBitDepth,
		},
	},
}

func outputs(ctx context.Context, l zerolog.Logger, registry prometheus.Registerer, doReq func(l zerolog.Logger, request *http.Request) ([]byte, error), id string) error {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s/video_outputs", url, unitEndpoint, id), nil)
	if err != nil {
		return err
	}

	res, err := doReq(l, request)
	if err != nil {
		return err
	}

	var outputs models.VideoOutputsResponse
	err = json.Unmarshal(res, &outputs)
	if err != nil {
		return err
	}

	mtrcs := metrics.NewGaugeMap(videoMetrics)
	for _, metric := range mtrcs {
		registry.MustRegister(metric)
		metric.Reset()
	}

	for _, output := range outputs.VideoOutputs {
		request, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s/video_outputs/%d/status", url, unitEndpoint, id, output.Index), nil)
		if err != nil {
			l.Err(err).Int("output_index", output.Index).Msg("Error creating output request, skipping")
			continue
		}
		res, err := doReq(l, request)
		if err != nil {
			l.Err(err).Int("output_index", output.Index).Msg("Error getting output metrics, skipping")
			continue
		}

		var e models.VideoOutputStatus
		err = json.Unmarshal(res, &e)
		if err != nil {
			l.Err(err).Int("output_index", output.Index).Msg("Error decoding output status response, skipping")
			continue
		}

		mtrcs[OutputVideoActive].WithLabelValues(
			strconv.Itoa(output.Index),
			output.Description,
			strconv.Itoa(e.VideoOut.Video.Format.Width),
			strconv.Itoa(e.VideoOut.Video.Format.Height),
			fmt.Sprintf("%.2f", e.VideoOut.Video.Format.Framerate),
			strconv.Itoa(e.VideoOut.Video.Format.BitDepth),
			metrics.BoolToString(e.VideoOut.Video.Format.Interlaced),
			e.VideoOut.Video.Format.ChromaSubsampling,
			e.VideoOut.Video.Format.PixelAspect,
			e.VideoOut.Video.Format.DisplayAspect,
			metrics.BoolToString(e.VideoOut.Video.Format.TopFieldFirst),
		).Set(metrics.BoolToFloat64(e.Active))

		for i, audio := range e.VideoOut.Audio {
			mtrcs[OutputAudioActive].WithLabelValues(
				strconv.Itoa(output.Index),
				output.Description,
				strconv.Itoa(i),
				strconv.Itoa(audio.Format.Channels),
				strconv.Itoa(audio.Format.SampleRate),
				strconv.Itoa(audio.Format.BitDepth),
			).Set(metrics.BoolToFloat64(e.Active))
		}

		mtrcs[OutputVideoSourceAvailable].WithLabelValues(
			strconv.Itoa(output.Index),
			output.Description,
			e.VideoSource.Video.Codec.Name,
			strconv.Itoa(e.VideoSource.Video.Codec.Bitrate),
			e.VideoSource.Video.Codec.Profile,
			e.VideoSource.Video.Codec.Level,
			strconv.Itoa(e.VideoSource.Video.Format.Width),
			strconv.Itoa(e.VideoSource.Video.Format.Height),
			fmt.Sprintf("%.2f", e.VideoSource.Video.Format.Framerate),
			strconv.Itoa(e.VideoSource.Video.Format.BitDepth),
			metrics.BoolToString(e.VideoSource.Video.Format.Interlaced),
			e.VideoSource.Video.Format.ChromaSubsampling,
			e.VideoSource.Video.Format.PixelAspect,
			e.VideoSource.Video.Format.DisplayAspect,
			metrics.BoolToString(e.VideoSource.Video.Format.TopFieldFirst),
		).Set(metrics.BoolToFloat64(e.VideoSource.Available))

		for i, audio := range e.VideoOut.Audio {
			mtrcs[OutputAudioSourceAvailable].WithLabelValues(
				strconv.Itoa(output.Index),
				output.Description,
				strconv.Itoa(i),
				audio.Codec.Name,
				strconv.Itoa(audio.Format.Channels),
				strconv.Itoa(audio.Format.SampleRate),
				strconv.Itoa(audio.Format.BitDepth),
			).Set(metrics.BoolToFloat64(e.VideoSource.Available))
		}

	}

	return err
}
