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

type GetLobbyResponse struct {
	Lobby Lobby
}

type UpdateLobbyRequest struct {
	Name  string
	Owner string
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
