package requests

type Response struct {
	Code    bool        `json:"code"`
	Message string      `json:"message"`
	Errors  error       `json:"errors"`
	Data    interface{} `json:"data"`
}
