package clienttype

//go:generate go-enum --noprefix --marshal --sqlint

/*
ENUM(
openai
vertex
)
*/
type ClientType string
