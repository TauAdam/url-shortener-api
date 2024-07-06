package response

type Response struct {
	Alias  string `json:"alias,omitempty"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func Success() Response {
	return Response{Status: StatusSuccess}
}
func Error(message string) Response {
	return Response{
		Status: StatusError,
	}
}
