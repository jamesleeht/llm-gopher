package oai

import (
	"reflect"

	"github.com/jamesleeht/llm-gopher/params"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

func mapPromptToMessages(prompt params.Prompt) []openai.ChatCompletionMessageParamUnion {
	messages := []openai.ChatCompletionMessageParamUnion{}
	if prompt.SystemMessage != "" {
		messages = append(messages, openai.SystemMessage(prompt.SystemMessage))
	}

	if prompt.Messages != nil {
		for _, message := range prompt.Messages {
			messages = append(messages, openai.UserMessage(message.Content))
		}
	}
	return messages
}

func mapPromptToResponseFormat(prompt params.Prompt) *openai.ChatCompletionNewParamsResponseFormatUnion {
	if prompt.ResponseFormatName == "" || prompt.ResponseFormat == nil {
		return nil
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   prompt.ResponseFormatName,
		Schema: generateSchemaFromType(reflect.TypeOf(prompt.ResponseFormat)),
		Strict: openai.Bool(true),
	}

	return &openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
			JSONSchema: schemaParam,
		},
	}
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

func mapSettingsToParams(settings params.Settings) openai.ChatCompletionNewParams {
	reasoningEffort := getReasoningEffortFromThinkingBudget(settings.ThinkingBudget)
	params := openai.ChatCompletionNewParams{
		Model:               shared.ChatModel(string(settings.ModelName)),
		ReasoningEffort:     reasoningEffort,
		MaxCompletionTokens: openai.Int(16000),
	}

	if settings.Temperature != nil {
		params.Temperature = openai.Float(*settings.Temperature)
	}

	return params
}

func getReasoningEffortFromThinkingBudget(thinkingBudget params.ThinkingBudget) shared.ReasoningEffort {
	switch thinkingBudget {
	case params.MinimalThinkingBudget:
		return shared.ReasoningEffortMinimal
	case params.SmallThinkingBudget:
		return shared.ReasoningEffortLow
	case params.MediumThinkingBudget:
		return shared.ReasoningEffortMedium
	case params.LargeThinkingBudget:
		return shared.ReasoningEffortHigh
	default:
		return ""
	}
}
