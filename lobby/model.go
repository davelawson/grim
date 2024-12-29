package lobby

type CreateLobbyRequest struct {
	Name string
}

type CreateLobbyResponse struct {
	Id string
}

type DeleteLobbyRequest struct {
	Id string
}

type GetLobbyRequest struct {
	Id string
}

type GetLobbyResponse struct {
	Lobby Lobby
}

type UpdateLobbyRequest struct {
	Lobby Lobby
}

type AddUserToLobbyRequest struct {
	UserId string
}

type Lobby struct {
	Id      string
	Name    string
	Owner   string
	Members []string
}
