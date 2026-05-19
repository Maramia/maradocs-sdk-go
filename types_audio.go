package maradocs

// AudioValidateRequest validates uploaded audio.
type AudioValidateRequest struct {
	UnvalidatedFileHandle UnvalidatedFileHandle `json:"unvalidated_file_handle"`
}

type audioValidateInner struct {
	ClassName   string        `json:"class_name"`
	AudioHandle *AudioHandle  `json:"audio_handle,omitempty"`
	Error       string        `json:"error,omitempty"`
	Virus       string        `json:"virus,omitempty"`
}

// AudioValidateResponse is the polled response for /audio/validate.
type AudioValidateResponse struct {
	ClassName            string           `json:"class_name"`
	SourceAudioMetadata *AudioMetadata   `json:"source_audio_metadata"`
	Response             audioValidateInner `json:"response"`
}

// OkAudio extracts the audio handle from a validation response.
func OkAudio(v AudioValidateResponse) (AudioHandle, error) {
	switch v.Response.ClassName {
	case "AudioValidateResponseOk":
		if v.Response.AudioHandle == nil {
			return AudioHandle{}, &ValidationError{Message: "missing audio_handle"}
		}
		return *v.Response.AudioHandle, nil
	case "AudioValidateResponseError":
		return AudioHandle{}, &ValidationError{Message: v.Response.Error}
	case "AudioValidateResponseVirus":
		return AudioHandle{}, &ValidationVirus{Virus: v.Response.Virus}
	default:
		return AudioHandle{}, &ValidationError{Message: "unknown validation response type"}
	}
}
