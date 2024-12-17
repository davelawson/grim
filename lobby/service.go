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
