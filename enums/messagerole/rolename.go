package messagerole

//go:generate go-enum --noprefix --marshal --sqlint

/*
ENUM(
assistant
user
)
*/
type MessageRole string
