package domain

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}
