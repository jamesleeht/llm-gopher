package params

import "llm-gopher/enums/modelname"

type ThinkingBudget string

const (
	NoThinkingBudget      ThinkingBudget = "no"
	MinimalThinkingBudget ThinkingBudget = "minimal"
	SmallThinkingBudget   ThinkingBudget = "small"
	MediumThinkingBudget  ThinkingBudget = "medium"
	LargeThinkingBudget   ThinkingBudget = "large"
)

type Settings struct {
	ModelName       modelname.ModelName
	Temperature     *float64
	ThinkingBudget  ThinkingBudget
	IsSearchEnabled bool
}
