package models

import (
	"cloud.google.com/go/vertexai/genai"
)

func BuildVertexModel(client *genai.Client) genai.GenerativeModel {
	generativeModel := client.GenerativeModel("gemini-1.5-pro-001")
	generativeModel.SetMaxOutputTokens(8192)
	generativeModel.SetTemperature(1.0)
	generativeModel.SetTopP(0.95)

	generativeModel.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}
	generativeModel.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(`
		You work as an extraction service mainly focusing on event detail extraction from different type of media or text.
		Input are text, image or video.
		text input has format as following struct:
		type Info struct {
			SuppliedBy  string 
			SuppliedFor string 
			Time        time.Time 
			Text        string 
		}
		Along with this text input there could be a file of any type.
		You extract following information from the input.
		Generate the output in JSON Array of objects to fit following struct:
		
		type Event struct {
			SuppliedBy       string    
			SuppliedFor      string    
			Time             time.Time 
			EventName        string    
			PeopleInvolved   []string 
			EventType        string    
			Activities       []string  
			Vibe             string    
			ThingsToRemember []string  
		}

		`)},
	}
	return *generativeModel
}
