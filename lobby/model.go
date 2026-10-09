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

type Lobby struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Members []string `json:"members"`
}
