package vertex

import (
	"fmt"
	"llm-gopher/params"
	"reflect"

	"github.com/invopop/jsonschema"
	"google.golang.org/genai"
)

func mapPromptToMessages(prompt params.Prompt) []*genai.Content {
	messages := []*genai.Content{}
	for _, message := range prompt.Messages {
		var role genai.Role
		switch message.Role {
		case params.MessageRoleAssistant:
			role = genai.RoleModel
		case params.MessageRoleUser:
			role = genai.RoleUser
		}
		messages = append(messages, genai.NewContentFromText(message.Content, role))
	}
	return messages
}

func mapSettingsToVertexSettings(prompt params.Prompt, settings params.Settings) (*genai.GenerateContentConfig, error) {
	if settings.IsSearchEnabled && prompt.ResponseFormatName != "" {
		return nil, fmt.Errorf("gemini - response format is not supported when search is enabled")
	}

	safetySettings := []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryUnspecified,
			Threshold: genai.HarmBlockThresholdOff,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockThresholdOff,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockThresholdOff,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockThresholdOff,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockThresholdOff,
		},
	}

	var tools []*genai.Tool
	if settings.IsSearchEnabled {
		tools = []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
			{URLContext: &genai.URLContext{}},
		}
	}

	respFormat := mapResponseFormatToVertexResponseFormat(prompt)
	respMimeType := ""
	if respFormat != nil {
		respMimeType = "application/json"
	}

	var genaiThinkingConfig *genai.ThinkingConfig
	thinkingBudget := getThinkingBudget(settings.ThinkingBudget)
	if thinkingBudget > 0 {
		genaiThinkingConfig = &genai.ThinkingConfig{
			IncludeThoughts: false,
			ThinkingBudget:  &thinkingBudget,
		}
	}

	var temperature *float32
	if settings.Temperature != nil {
		v := float32(*settings.Temperature)
		temperature = &v
	}

	var systemInstruction *genai.Content
	if prompt.SystemMessage != "" {
		systemInstruction = &genai.Content{Parts: []*genai.Part{{Text: prompt.SystemMessage}}}
	}

	return &genai.GenerateContentConfig{
		SystemInstruction:  systemInstruction,
		Temperature:        temperature,
		SafetySettings:     safetySettings,
		ThinkingConfig:     genaiThinkingConfig,
		Tools:              tools,
		ResponseJsonSchema: respFormat,
		ResponseMIMEType:   respMimeType,
	}, nil
}

func mapResponseFormatToVertexResponseFormat(prompt params.Prompt) any {
	// Return nil if no response format is specified
	if prompt.ResponseFormatName == "" || prompt.ResponseFormat == nil {
		return nil
	}
	// Create a zero value of the type and generate schema
	return generateSchemaFromType(reflect.TypeOf(prompt.ResponseFormat))
}

func generateSchemaFromType(t reflect.Type) interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}

	// Create a zero value of the type
	v := reflect.New(t).Elem().Interface()
	schema := reflector.Reflect(v)
	return schema
}

func getThinkingBudget(thinkingBudget params.ThinkingBudget) int32 {
	switch thinkingBudget {
	case params.NoThinkingBudget:
		return 0
	case params.MinimalThinkingBudget:
		return 512
	case params.SmallThinkingBudget:
		return 1024
	case params.MediumThinkingBudget:
		return 2048
	case params.LargeThinkingBudget:
		return 4096
	default:
		return 0
	}
}
