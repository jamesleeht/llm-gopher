package main

import (
	"llm-gopher/examples/main/enums/modelname"
	"llm-gopher/examples/main/enums/presetname"
	"llm-gopher/params"
)

func getPresetSettingsMap() map[string]params.Settings {
	defaultTemperature := float64(0.5)
	presets := map[presetname.PresetName]params.Settings{
		presetname.DeepseekV3: {
			ModelName:       modelname.DeepseekDeepseekV3Turbo.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.NoThinkingBudget,
			IsSearchEnabled: false,
		},
		presetname.DeepseekV31: {
			ModelName:       modelname.DeepseekDeepseekV31.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.NoThinkingBudget,
			IsSearchEnabled: false,
		},
		presetname.Gemini25FlashNonThinking: {
			ModelName:       modelname.Gemini25Flash.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.NoThinkingBudget,
			IsSearchEnabled: false,
		},
		presetname.Gemini25FlashThinking: {
			ModelName:       modelname.Gemini25Flash.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.MediumThinkingBudget,
			IsSearchEnabled: false,
		},
		presetname.Gemini25FlashNonThinkingSearch: {
			ModelName:       modelname.Gemini25Flash.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.NoThinkingBudget,
			IsSearchEnabled: true,
		},
		presetname.Gemini25ProLow: {
			ModelName:       modelname.Gemini25Pro.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.SmallThinkingBudget,
			IsSearchEnabled: false,
		},
		presetname.Gemini25ProLowSearch: {
			ModelName:       modelname.Gemini25Pro.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.SmallThinkingBudget,
			IsSearchEnabled: true,
		},
		presetname.Gemini25ProHigh: {
			ModelName:       modelname.Gemini25Pro.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.LargeThinkingBudget,
			IsSearchEnabled: false,
		},
		presetname.Gemini25ProHighSearch: {
			ModelName:       modelname.Gemini25Pro.String(),
			Temperature:     &defaultTemperature,
			ThinkingBudget:  params.LargeThinkingBudget,
			IsSearchEnabled: true,
		},
		presetname.GPT4OMiniSearch: {
			ModelName:       modelname.Gpt4OMiniSearchPreview.String(),
			Temperature:     nil,
			ThinkingBudget:  params.NoThinkingBudget,
			IsSearchEnabled: true,
		},
		presetname.GPT4OSearch: {
			ModelName:       modelname.Gpt4OSearchPreview.String(),
			Temperature:     nil,
			ThinkingBudget:  params.NoThinkingBudget,
			IsSearchEnabled: true,
		},
	}

	result := make(map[string]params.Settings)
	for presetName, settings := range presets {
		result[presetName.String()] = settings
	}
	return result
}
