package lobby

type CreateLobbyRequest struct {
	Name string
}

type CreateLobbyResponse struct {
	Id string
}

type Lobby struct {
	Id    string
	Name  string
	Owner string
}
