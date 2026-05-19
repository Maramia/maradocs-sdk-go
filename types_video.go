package maradocs

// VideoValidateRequest validates uploaded video.
type VideoValidateRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

type videoValidateInner struct {
	ClassName   string        `json:"class_name"`
	VideoHandle *VideoHandle  `json:"video_handle,omitempty"`
	Error       string        `json:"error,omitempty"`
	Virus       string        `json:"virus,omitempty"`
}

// VideoValidateResponse is the polled response for /video/validate.
type VideoValidateResponse struct {
	ClassName            string           `json:"class_name"`
	SourceVideoMetadata *VideoMetadata   `json:"source_video_metadata"`
	SourceAudioMetadata *AudioMetadata   `json:"source_audio_metadata"`
	Response             videoValidateInner `json:"response"`
}

// OkVideo extracts the video handle from a validation response.
func OkVideo(v VideoValidateResponse) (VideoHandle, error) {
	switch v.Response.ClassName {
	case "VideoValidateResponseOk":
		if v.Response.VideoHandle == nil {
			return VideoHandle{}, &ValidationError{Message: "missing video_handle"}
		}
		return *v.Response.VideoHandle, nil
	case "VideoValidateResponseError":
		return VideoHandle{}, &ValidationError{Message: v.Response.Error}
	case "VideoValidateResponseVirus":
		return VideoHandle{}, &ValidationVirus{Virus: v.Response.Virus}
	default:
		return VideoHandle{}, &ValidationError{Message: "unknown validation response type"}
	}
}
