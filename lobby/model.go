package lobby

type CreateLobbyRequest struct {
	Name string `binding:"required"`
}

type CreateLobbyResponse struct {
	Id string
}

type GetLobbyResponse struct {
	Lobby Lobby
}

type UpdateLobbyRequest struct {
	Name  string `binding:"required"`
	Owner string `binding:"required"`
}

type AddUserToLobbyRequest struct {
	UserId string `binding:"required"`
}

type Lobby struct {
	Id      string
	Name    string
	Owner   string
	Members []string
}
