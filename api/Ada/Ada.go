package Ada

import (
	"Ada/api/Actions"
	"Ada/api/Common"
	"Ada/api/Ollama"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var audioBuffer = new(bytes.Buffer)
var mu = new(sync.Mutex)

func RegisterAdaRoutes(router *gin.Engine) {
	fmt.Print("Registering Ada routes")
	router.Group("/ada")
	{
		router.POST("/ada", Post)
	}
}

func Post(c *gin.Context) {
	audioByteArray, err := c.GetRawData()
	if err != nil {
		log.Fatalf("Error getting raw data: %v", err)
	}
	transcription := transcribeAudio(audioByteArray)
	intent := determineIntent(transcription)

	adaContext := Common.AdaContext{
		Audio:          audioByteArray,
		Transcription:  transcription,
		Intent:         intent,
		OllamaResponse: Common.OllamaIntentResponse{},
	}

	respondToIntent(adaContext)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func respondToIntent(ctx Common.AdaContext) {
	handlers := map[Common.Intent]func(){
		// Light Control
		Common.Lightson:  TurnLightsOn,
		Common.Lightsoff: func() { fmt.Println("Turning off the lights") },

		// Alarm Control
		Common.SetAlarm: func() { fmt.Println("Setting the alarm") },

		// Spotify Control

		Common.CheckWeather: func() { fmt.Println("Checking the weather") },
		Common.Talk: func() {
			Actions.Talk(ctx)
		},

		// Fallback
		Common.NotForMe: func() { fmt.Println("not meant for Ada") },
		Common.Unknown:  func() { fmt.Println("Unknown intent") },
	}

	if handler, found := handlers[ctx.Intent]; found {
		handler()
	} else {
		fmt.Println("Unknown intent")
	}
}

func TurnLightsOn() {
	fmt.Println("Turning on the lights")
}

func determineIntent(transcription string) Common.Intent {
	//schema := map[string]interface{}{
	//	"type": "object",
	//	"properties": map[string]interface{}{
	//		"evaluation": map[string]interface{}{
	//			"type": "string",
	//			"enum": Common.Intents,
	//		},
	//	},
	//	"required": []string{"evaluation"},
	//}

	schema := Ollama.Schema{
		Type: Ollama.Object,
		Properties: map[string]interface{}{
			"evaluation": map[string]interface{}{
				"type": "string",
				"enum": Common.Intents,
			},
		},
		Required: []string{"evaluation"},
	}

	prompt := generateIntentPrompt()

	//requestPayload := map[string]interface{}{
	//	"model":    viper.GetString("intent_model"),
	//	"messages": []map[string]string{{"role": "system", "content": prompt}, {"role": "user", "content": transcription}},
	//	"stream":   false,
	//	"format":   schema,
	//}

	requestPayload := Ollama.RequestPayload{
		Model: viper.GetString("intent_model"),
		Messages: []map[string]Ollama.Message{
			{
				"system": Ollama.Message{
					Role:    "system",
					Content: prompt,
				},
				"user": Ollama.Message{
					Role:    "user",
					Content: transcription,
				},
			},
		},
		Stream: false,
		Format: schema,
	}

	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		log.Fatalf("Error marshalling request payload: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from Ollama API: %s", respBody)
	}

	var ollamaResponse Common.OllamaIntentResponse
	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling response to OllamaIntentResponse: %v", err)
	}
	fmt.Printf("Parsed Response: %+v\n", ollamaResponse)

	ollamaResponse.Message.Content.UnmarshalJSON([]byte(ollamaResponse.Message.Content.Evaluation))
	fmt.Printf("Parsed Response: %+v\n", ollamaResponse.Message.Content)
	fmt.Printf("Response: %s\n", ollamaResponse.Message.Content.Evaluation)

	return ollamaResponse.Message.Content.Evaluation
}

func generateIntentPrompt() string {
	return fmt.Sprintf(
		viper.GetString("intent_prompt"),
		Common.Intents,
	)
}

func transcribeAudio(audio []byte) string {

	dirPath := "audio_files"
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}

	filePath := fmt.Sprintf("%s/audio.wav", dirPath)
	err = os.WriteFile(filePath, audio, 0644)
	if err != nil {
		log.Fatalf("Error saving audio to file: %v", err)
	}

	cmd := exec.Command("whisper", "audio.wav", "--model", "tiny", "--language", "en", "--output_dir", viper.GetString("transcribe_audio_directory_path"))

	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}

	// Print the output
	log.Printf("Command output: %s", output)
	return string(output)
}

func processAudioChunk(audioData []byte) error {
	mu.Lock()                              // Lock to ensure thread-safe appending
	defer mu.Unlock()                      // Unlock after appending
	_, err := audioBuffer.Write(audioData) // Append the byte array
	if err != nil {
		return fmt.Errorf("error writing to buffer: %v", err)
	}
	return nil
}
