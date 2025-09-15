package router

import (
	"context"
	"fmt"
	"llm-gopher/client"
	"llm-gopher/params"
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
	for modelName := range presetMap {
		modelsFromPresetMap[modelName] = true
	}

	for modelName := range modelsFromPresetMap {
		if !modelsFromClientMap[modelName] {
			return fmt.Errorf("model %s defined in preset map but not in client map", modelName)
		}
	}

	for modelName := range modelsFromClientMap {
		if !modelsFromPresetMap[modelName] {
			return fmt.Errorf("model %s defined in client map but not in preset map", modelName)
		}
	}
	return nil
}

func (r *Router) SendPrompt(ctx context.Context,
	presetName string,
	prompt params.Prompt) (string, error) {
	preset, exists := r.presetMap[presetName]
	if !exists {
		return "", fmt.Errorf("preset %s not found", presetName)
	}

	client, err := r.GetClientForModelName(preset.ModelName)
	if err != nil {
		return "", fmt.Errorf("failed to get client for model %s: %w", preset.ModelName, err)
	}

	response, err := client.SendMessage(ctx, prompt, preset)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
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
