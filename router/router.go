package router

import (
	"context"
	"fmt"

	"github.com/jamesleeht/llm-gopher/client"
	"github.com/jamesleeht/llm-gopher/params"
)

type Router struct {
	clientMap ClientMap
	presetMap PresetMap
}

type ClientMap map[string][]*client.Client
type PresetMap map[string]params.Settings

func NewRouter(clients ClientMap, presetMap PresetMap) (*Router, error) {
	err := validateAllModelsDefined(clients, presetMap)
	if err != nil {
		return nil, err
	}

	return &Router{
		clientMap: clients,
		presetMap: presetMap,
	}, nil
}

func validateAllModelsDefined(clientMap ClientMap, presetMap PresetMap) error {
	modelsFromClientMap := make(map[string]bool)
	for modelName := range clientMap {
		modelsFromClientMap[modelName] = true
	}

	modelsFromPresetMap := make(map[string]bool)
	for _, preset := range presetMap {
		modelsFromPresetMap[preset.ModelName] = true
	}

	for modelName := range modelsFromPresetMap {
		if !modelsFromClientMap[modelName] {
			return fmt.Errorf("model %s defined in preset map but not in client map", modelName)
		}
	}
	return nil
}

// SendPrompt sends a prompt using the specified preset and returns a response.
// For unstructured text responses, use SendPrompt with Prompt[any].
// For structured JSON responses, use SendPromptTyped[T] instead.
func (r *Router) SendPrompt(ctx context.Context,
	presetName string,
	prompt params.Prompt[any]) (*client.Response[any], error) {
	preset, exists := r.presetMap[presetName]
	if !exists {
		return nil, fmt.Errorf("preset %s not found", presetName)
	}

	c, err := r.GetClientForModelName(preset.ModelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get client for model %s: %w", preset.ModelName, err)
	}

	response, err := client.SendMessage[any](ctx, c, prompt, preset)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return response, nil
}

// SendPromptTyped sends a prompt using the specified preset and returns a typed response.
// Use this when you want structured JSON output.
func SendPromptTyped[T any](ctx context.Context,
	r *Router,
	presetName string,
	prompt params.Prompt[T]) (*client.Response[T], error) {
	preset, exists := r.presetMap[presetName]
	if !exists {
		return nil, fmt.Errorf("preset %s not found", presetName)
	}

	c, err := r.GetClientForModelName(preset.ModelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get client for model %s: %w", preset.ModelName, err)
	}

	response, err := client.SendMessage[T](ctx, c, prompt, preset)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return response, nil
}

func (r *Router) GetClientForModelName(modelName string) (*client.Client, error) {
	clients, exists := r.clientMap[modelName]
	if !exists {
		return nil, fmt.Errorf("model name: %s not defined in client map", modelName)
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("no clients found for model name: %s", modelName)
	}

	// TODO: load balance between different clients
	return clients[0], nil
}
