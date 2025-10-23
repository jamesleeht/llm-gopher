package oai

import (
	"reflect"

	"github.com/jamesleeht/llm-gopher/params"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

func mapPromptToMessages[T any](prompt params.Prompt[T]) []openai.ChatCompletionMessageParamUnion {
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

func mapPromptToResponseFormat[T any]() *openai.ChatCompletionNewParamsResponseFormatUnion {
	// Get the type of T
	var zero T
	t := reflect.TypeOf(zero)

	// If T is interface{} (any) or nil, don't use structured output
	if t == nil || t.Kind() == reflect.Interface {
		return nil
	}

	// Get the type name from the struct
	schemaName := t.Name()

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   schemaName,
		Schema: generateSchemaFromType(t),
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

	// If it's a pointer type, get the element type
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
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
