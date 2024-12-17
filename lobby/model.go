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

type Lobby struct {
	Id    string
	Name  string
	Owner string
}
