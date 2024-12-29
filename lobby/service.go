package lobby

type userRepo interface {
	GetUserByUUID(uuid string)
}

type Service struct {
	repo     *LobbyRepo
	userRepo userRepo
}

func (ls *Service) AddUserToLobby(lobbyId string, userId string) error {
	return ls.repo.AddUserToLobby(lobbyId, userId)
}

func NewService(repo *LobbyRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (ls *Service) CreateLobby(name string, userId string) (string, error) {
	lobbyId, err := ls.repo.CreateLobby(name, userId)
	if err != nil {
		return "", err
	}
	err = ls.repo.AddMemberToLobby(lobbyId, userId)
	if err != nil {
		return "", err
	}
	return lobbyId, err
}

func (ls *Service) DeleteLobby(lobbyId string, userId string) (bool, error) {
	rows, err := ls.repo.DeleteLobby(lobbyId, userId)
	return rows >= 1, err
}

func (ls *Service) GetLobby(id string) (*Lobby, error) {
	lobby, err := ls.repo.GetLobby(id)
	if err != nil {
		return nil, err
	} else if lobby == nil {
		return nil, nil
	}

	members, err := ls.repo.GetLobbyMembers(id)
	lobby.Members = members
	return lobby, err
}
