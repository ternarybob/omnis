package omnis

type ApiResponse struct {
	Version       string            `json:"version"`
	Name          string            `json:"name"`
	Support       string            `json:"support"`
	Status        int               `json:"status"`
	Scope         string            `json:"scope"`
	CorrelationId string            `json:"correlationid"`
	Err           string            `json:"error,omitempty"`
	Stack         []string          `json:"stack,omitempty"`
	Request       map[string]string `json:"request,omitempty"`
	Log           map[string]string `json:"log,omitempty"`
	Result        interface{}       `json:"result"`
}

type ApiTypedResponse[T any] struct {
	ApiResponse
	Result T `json:"result"`
}

type ApiPagedResult[T any] struct {
	Total    int   `json:"total"`
	Page     int   `json:"page"`
	Size     int   `json:"size"`
	Duration int64 `json:"duration"`
	Data     []*T  `json:"data"`
}

type ApiPagedResponse[T any] struct {
	ApiResponse
	Result ApiPagedResult[T] `json:"result"`
}
