package params

type ThinkingBudget string

const (
	NoThinkingBudget      ThinkingBudget = "no"
	MinimalThinkingBudget ThinkingBudget = "minimal"
	SmallThinkingBudget   ThinkingBudget = "small"
	MediumThinkingBudget  ThinkingBudget = "medium"
	LargeThinkingBudget   ThinkingBudget = "large"
)

type Settings struct {
	ModelName       string
	Temperature     *float64
	ThinkingBudget  ThinkingBudget
	IsSearchEnabled bool
}
