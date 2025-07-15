package ai

type AIAdaptor struct {
	api_key string
}

func NewAIAdaptor(key string) *AIAdaptor {
	return &AIAdaptor{
		api_key: key,
	}
}
