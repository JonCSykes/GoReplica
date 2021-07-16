package goreplica

// IReplicaClient : This is the interface for the GoReplica client code.
type IReplicaClient interface {
	Auth(clientID, secret string)
	GetVoices()
	GetSpeech(text, speakerID string, bitRate, sampleRate int, ext SpeechExtension)
}
