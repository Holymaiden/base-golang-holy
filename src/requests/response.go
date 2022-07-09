package requests

type Response struct {
	Code    bool        `json:"code"`
	Message string      `json:"message"`
	Errors  error       `json:"errors"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}
