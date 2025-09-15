package main

import (
	"context"
	"fmt"
	"llm-gopher/examples/main/enums/presetname"
	"llm-gopher/params"
	"llm-gopher/router"
	"log"
)

func main() {
	env := getAppConfig()

	clientMap := getClientMap(env)
	presetMap := getPresetSettingsMap()
	router, err := router.NewRouter(clientMap, presetMap)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	prompt := params.NewSimplePrompt("You are a helpful assistant.", "Hello, how are you?")
	response, err := router.SendPrompt(context.Background(), presetname.DeepseekV3.String(), prompt)
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	fmt.Println(response)
}
