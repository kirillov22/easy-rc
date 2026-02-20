package actions

type ActionResponse struct {
	Data []byte
	Err  error
}

func NewActionResponse(data []byte, err error) ActionResponse {
	return ActionResponse{Data: data, Err: err}
}

func NoopActionResponse() ActionResponse {
	return ActionResponse{nil, nil}
}
