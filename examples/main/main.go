package main

import (
	"context"
	"fmt"
	"llm-gopher/enums/presetname"
	"llm-gopher/params"
	"llm-gopher/router"
	"log"
)

func main() {
	env := getAppConfig()

	clientMap := getClientMap(env)
	presetMap := getPresetSettingsMap()
	router := router.NewRouter(clientMap, presetMap)

	prompt := params.NewSimplePrompt("You are a helpful assistant.", "Hello, how are you?")
	response, err := router.SendPrompt(context.Background(), presetname.DeepseekV3, prompt)
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	fmt.Println(response)
}
