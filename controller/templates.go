package controller

type templateErrorResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

var templateRegisterRequest = map[string]any{
	"username": "",
	"password": "",
	"name":     "",
	"surname":  "",
}

type templateRegisterSuccess struct {
	Status bool `json:"status"`
	Result struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	} `json:"result"`
}

var templateLogInRequest = map[string]any{
	"username": "",
	"password": "",
}

type templateLoginSuccess struct {
	Status bool `json:"status"`
	Result struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
	} `json:"result"`
}

type templateUserData struct {
	Status bool `json:"status"`
	Result struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Surname  string `json:"surname"`
	} `json:"result"`
}
