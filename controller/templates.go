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
		ID       int    `json:"id"`
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
		ID       int    `json:"id"`
		Username string `json:"username"`
	} `json:"result"`
}

type templateUserData struct {
	Status bool `json:"status"`
	Result struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
		Surname  string `json:"surname"`
	} `json:"result"`
}

var templateGameResults = map[string]any{
	"userid1": 0.0,
	"userid2": 0.0,
	"score1":  0.0,
	"score2":  0.0,
}

type templateScoreUpdateSuccess struct {
	Status bool `json:"status"`
	User1  struct {
		ID       int64   `json:"id"`
		Username string  `json:"username"`
		Score    float64 `json:"score"`
	} `json:"user1"`
	User2 struct {
		ID       int64   `json:"id"`
		Username string  `json:"username"`
		Score    float64 `json:"score"`
	} `json:"user2"`
}

var templateLeaderboardRequest = map[string]any{
	"page":  0.0,
	"count": 0.0,
}

type leaderboardPlacement struct {
	ID    int     `json:"id"`
	Score float64 `json:"score"`
	Rank  int     `json:"rank"`
}

type templateLeaderboardPage struct {
	Status bool                   `json:"status"`
	Result []leaderboardPlacement `json:"result"`
}

var templateSimulationParams = map[string]any{
	"usercount": 0.0,
}

type userData struct {
	ID       int     `json:"id"`
	Score    float64 `json:"score"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Surname  string  `json:"surname"`
}

type templateSimulationResult struct {
	Status bool       `json:"status"`
	Result []userData `json:"result"`
}
