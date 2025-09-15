package modelname

//go:generate go-enum --noprefix --marshal --sqlint

/*
ENUM(
deepseek/deepseek-v3-turbo
deepseek/deepseek-v3.1
gemini-2.0-flash
gemini-2.5-flash
gemini-2.5-pro
gpt-4o-search-preview
gpt-4o-mini-search-preview
)
*/
type ModelName string

func (modelName ModelName) IsGemini() bool {
	return modelName == Gemini20Flash || modelName == Gemini25Flash || modelName == Gemini25Pro
}
