package user

type ServiceInterface interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	DeleteUserByID(id int64) error
}

func NewService(repository RepositoryInterface) ServiceInterface {
	return &Service{repository}
}

type Service struct {
	repository RepositoryInterface
}

func (s *Service) CreateUser(user *User) error {
	return s.repository.Create(user)
}

func (s *Service) GetUserByUsername(username string) (*User, error) {
	return s.repository.GetByUsername(username)
}

func (s *Service) DeleteUserByID(id int64) error {
	return s.repository.DeleteByID(id)
}
