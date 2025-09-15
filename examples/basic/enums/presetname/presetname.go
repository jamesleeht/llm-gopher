package presetname

import "slices"

//go:generate go-enum --noprefix --marshal --sqlint

/*
ENUM(
Gemini 2.5 Flash Non-Thinking
Gemini 2.5 Flash Thinking
Gemini 2.5 Flash Non-Thinking Search
Gemini 2.5 Pro Low
Gemini 2.5 Pro Low Search
Gemini 2.5 Pro High
Gemini 2.5 Pro High Search
GPT 4o Search
GPT 4o Mini Search
Deepseek V3
Deepseek V3.1
)
*/
type PresetName string

func GetAllOptions() []PresetName {
	var presetOptions []PresetName
	for presetName := range _PresetNameValue {
		presetOptions = append(presetOptions, PresetName(presetName))
	}
	slices.Sort(presetOptions)
	return presetOptions
}
