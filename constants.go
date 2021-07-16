package goreplica

// SpeechExtension : This is as type used to represent the allowed extension type of the Replica API.
type SpeechExtension string

const (
	WAV          SpeechExtension = "wav"
	MP3          SpeechExtension = "mp3"
	OGG          SpeechExtension = "ogg"
	FLAC         SpeechExtension = "flac"
	CLIENTID     string          = "client_id"
	CLIENTSECRET string          = "secret"
	UUID         string          = "uuid"
	NAME         string          = "name"
	TEXT         string          = "txt"
	SPEAKERID    string          = "speaker_id"
	EXTENSION    string          = "extension"
	BITRATE      string          = "bit_rate"
	SAMPLERATE   string          = "sample_rate"
)
