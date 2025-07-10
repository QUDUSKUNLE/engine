package ports

type AIService interface {
	InterpretLabResult(prompt string) (interface{}, error)
}
