package lobby

type CreateLobbyRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateLobbyResponse struct {
	Id string `json:"id"`
}

type GetLobbyResponse struct {
	Lobby Lobby `json:"lobby"`
}

type UpdateLobbyRequest struct {
	Name  string `json:"name" binding:"required"`
	Owner string `json:"owner" binding:"required"`
}

type AddUserToLobbyRequest struct {
	UserId string `json:"userid" binding:"required"`
}

type ReadyRequest struct {
	Ready *bool `json:"ready" binding:"required"`
}

type LaunchRequest struct {
	RequestID string `json:"requestid" binding:"required,uuid"`
}

type Lobby struct {
	Id        string          `json:"id"`
	Name      string          `json:"name"`
	Owner     string          `json:"owner"`
	Members   []string        `json:"members"`
	Status    string          `json:"status"`
	MatchID   *string         `json:"matchid"`
	Readiness map[string]bool `json:"readiness"`
}
