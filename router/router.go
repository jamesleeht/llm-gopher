package router

import (
	"context"
	"fmt"
	"llm-gopher/client"
	"llm-gopher/enums/modelname"
	"llm-gopher/enums/presetname"
	"llm-gopher/params"
)

type Router struct {
	clientMap ClientMap
	presetMap PresetMap
}

type ClientMap map[modelname.ModelName][]*client.Client
type PresetMap map[presetname.PresetName]params.Settings

func NewRouter(clients ClientMap, presetMap PresetMap) *Router {
	return &Router{
		clientMap: clients,
		presetMap: presetMap,
	}
}

func (r *Router) SendMessage(ctx context.Context,
	presetName presetname.PresetName,
	prompt params.Prompt) (string, error) {
	presetSettings, exists := r.presetMap[presetName]
	if !exists {
		return "", fmt.Errorf("preset %s not found", presetName)
	}

	client, err := r.GetClientForModelName(presetSettings.ModelName)
	if err != nil {
		return "", fmt.Errorf("failed to get client for model %s: %w", presetSettings.ModelName, err)
	}

	response, err := client.SendMessage(ctx, prompt, presetSettings)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	return response, nil
}

func (r *Router) GetClientForModelName(modelName modelname.ModelName) (*client.Client, error) {
	clients := r.clientMap[modelName]
	if len(clients) == 0 {
		return nil, fmt.Errorf("no clients found for model %s", modelName)
	}

	// TODO: load balance between different clients
	return clients[0], nil
}
