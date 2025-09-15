package main

import (
	"llm-gopher/client"
	"llm-gopher/examples/basic/enums/modelname"
	"log"
)

func getClientMap(env appConfig) map[string][]*client.Client {
	chatGPTClient, err := client.NewClient(client.ClientConfig{
		BaseURL: "https://api.openai.com/v1",
		APIKey:  env.apiKey,
	}, client.ClientTypeOpenAI)
	if err != nil {
		log.Fatalf("failed to create chatgpt client: %v", err)
	}

	novitaClient, err := client.NewClient(client.ClientConfig{
		APIKey:  env.novitaKey,
		BaseURL: "https://api.novita.ai/v3/openai",
	}, client.ClientTypeOpenAI)
	if err != nil {
		log.Fatalf("failed to create novita client: %v", err)
	}

	vertexAIClient, err := client.NewClient(client.ClientConfig{
		ProjectID:       env.vertexAIProjectID,
		Location:        env.vertexAILocation,
		CredentialsPath: env.vertexAICredentialsPath,
	}, client.ClientTypeVertex)
	if err != nil {
		log.Fatalf("failed to create vertex ai client: %v", err)
	}

	clientMap := map[modelname.ModelName][]*client.Client{
		modelname.DeepseekDeepseekV3Turbo: []*client.Client{novitaClient},
		modelname.DeepseekDeepseekV31:     []*client.Client{novitaClient},
		modelname.Gemini20Flash:           []*client.Client{vertexAIClient},
		modelname.Gemini25Flash:           []*client.Client{vertexAIClient},
		modelname.Gemini25Pro:             []*client.Client{vertexAIClient},
		modelname.Gpt4OSearchPreview:      []*client.Client{chatGPTClient},
		modelname.Gpt4OMiniSearchPreview:  []*client.Client{chatGPTClient},
	}

	result := make(map[string][]*client.Client)
	for modelName, clients := range clientMap {
		result[modelName.String()] = clients
	}
	return result
}
