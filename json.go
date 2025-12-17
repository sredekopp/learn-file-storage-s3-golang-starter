package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func getVideoAspectRatio(filePath string) (string, error) {

	type VideoProbe struct {
		Streams []struct {
			Index              int    `json:"index"`
			CodecName          string `json:"codec_name,omitempty"`
			CodecLongName      string `json:"codec_long_name,omitempty"`
			Profile            string `json:"profile,omitempty"`
			CodecType          string `json:"codec_type"`
			CodecTagString     string `json:"codec_tag_string"`
			CodecTag           string `json:"codec_tag"`
			Width              int    `json:"width,omitempty"`
			Height             int    `json:"height,omitempty"`
			CodedWidth         int    `json:"coded_width,omitempty"`
			CodedHeight        int    `json:"coded_height,omitempty"`
			ClosedCaptions     int    `json:"closed_captions,omitempty"`
			FilmGrain          int    `json:"film_grain,omitempty"`
			HasBFrames         int    `json:"has_b_frames,omitempty"`
			SampleAspectRatio  string `json:"sample_aspect_ratio,omitempty"`
			DisplayAspectRatio string `json:"display_aspect_ratio,omitempty"`
			PixFmt             string `json:"pix_fmt,omitempty"`
			Level              int    `json:"level,omitempty"`
			ColorRange         string `json:"color_range,omitempty"`
			ColorSpace         string `json:"color_space,omitempty"`
			ColorTransfer      string `json:"color_transfer,omitempty"`
			ColorPrimaries     string `json:"color_primaries,omitempty"`
			ChromaLocation     string `json:"chroma_location,omitempty"`
			FieldOrder         string `json:"field_order,omitempty"`
			Refs               int    `json:"refs,omitempty"`
			IsAvc              string `json:"is_avc,omitempty"`
			NalLengthSize      string `json:"nal_length_size,omitempty"`
			ID                 string `json:"id"`
			RFrameRate         string `json:"r_frame_rate"`
			AvgFrameRate       string `json:"avg_frame_rate"`
			TimeBase           string `json:"time_base"`
			StartPts           int    `json:"start_pts"`
			StartTime          string `json:"start_time"`
			DurationTs         int    `json:"duration_ts"`
			Duration           string `json:"duration"`
			BitRate            string `json:"bit_rate,omitempty"`
			BitsPerRawSample   string `json:"bits_per_raw_sample,omitempty"`
			NbFrames           string `json:"nb_frames"`
			ExtradataSize      int    `json:"extradata_size"`
			SampleFmt          string `json:"sample_fmt,omitempty"`
			SampleRate         string `json:"sample_rate,omitempty"`
			Channels           int    `json:"channels,omitempty"`
			ChannelLayout      string `json:"channel_layout,omitempty"`
			BitsPerSample      int    `json:"bits_per_sample,omitempty"`
			InitialPadding     int    `json:"initial_padding,omitempty"`
			Disposition        struct {
				Default         int `json:"default"`
				Dub             int `json:"dub"`
				Original        int `json:"original"`
				Comment         int `json:"comment"`
				Lyrics          int `json:"lyrics"`
				Karaoke         int `json:"karaoke"`
				Forced          int `json:"forced"`
				HearingImpaired int `json:"hearing_impaired"`
				VisualImpaired  int `json:"visual_impaired"`
				CleanEffects    int `json:"clean_effects"`
				AttachedPic     int `json:"attached_pic"`
				TimedThumbnails int `json:"timed_thumbnails"`
				NonDiegetic     int `json:"non_diegetic"`
				Captions        int `json:"captions"`
				Descriptions    int `json:"descriptions"`
				Metadata        int `json:"metadata"`
				Dependent       int `json:"dependent"`
				StillImage      int `json:"still_image"`
			} `json:"disposition"`
			Tags struct {
				Language    string `json:"language"`
				HandlerName string `json:"handler_name"`
				VendorID    string `json:"vendor_id"`
				Encoder     string `json:"encoder"`
				Timecode    string `json:"timecode"`
			} `json:"tags,omitempty"`
		} `json:"streams"`
	}
	probe := VideoProbe{}

	buff := bytes.Buffer{}
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	cmd.Stdout = &buff
	err := cmd.Run()
	if err != nil {
		return "", errors.New("error probing video attributes")
	}

	err = json.Unmarshal(buff.Bytes(), &probe)
	if err != nil {
		return "", errors.New("error decoding video attributes")
	}

	width := 0
	height := 0
	for _, stream := range probe.Streams {
		if stream.Width > 0 && stream.Height > 0 {
			width = stream.Width
			height = stream.Height
			break
		}
	}

	aspect := ""
	if width == 0 || height == 0 {
		aspect = "other"
	} else if width/height > 0 {
		aspect = "16:9"
	} else {
		aspect = "9:16"
	}

	return aspect, nil
}

func processVideoForFastStart(filePath string) (string, error) {
	processingPath := filePath + ".processing"

	buff := bytes.Buffer{}
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", processingPath)
	cmd.Stdout = &buff
	err := cmd.Run()
	if err != nil {
		return "", errors.New("error probing video attributes")
	}

	return processingPath, nil
}
