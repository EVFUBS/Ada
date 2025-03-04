package Common

type AdaContext struct {
	Audio          []byte
	Transcription  string
	Intent         Intent
	OllamaResponse OllamaIntentResponse
}
