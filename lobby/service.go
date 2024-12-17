package lobby

type userRepo interface {
	GetUserByUUID(uuid string)
}

type Service struct {
	repo     *LobbyRepo
	userRepo userRepo
}

func NewService(repo *LobbyRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (ls *Service) CreateLobby(name string, userId string) (string, error) {
	return ls.repo.CreateLobby(name, userId)
}

func (ls *Service) DeleteLobby(lobbyId string, userId string) (bool, error) {
	rows, err := ls.repo.DeleteLobby(lobbyId, userId)
	return rows >= 1, err
}
