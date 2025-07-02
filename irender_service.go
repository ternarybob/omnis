package omnis

type IRenderService interface {
	AsModel(code int, output interface{})
	AsResult(code int, payload interface{})
	AsResultWithError(code int, payload interface{}, err error)
	AsError(code int, err interface{})
}
