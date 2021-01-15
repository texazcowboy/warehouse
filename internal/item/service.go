package item

type ServiceInterface interface {
	CreateItem(item *Item) error
	GetItem(id int64) (*Item, error)
	GetAllItems() ([]*Item, error)
	UpdateItem(item *Item) (int64, error)
	DeleteItem(id int64) (int64, error)
}

type Service struct {
	repository RepositoryInterface
}

func (s *Service) CreateItem(item *Item) error {
	return s.repository.Create(item)
}

func (s *Service) GetItem(id int64) (*Item, error) {
	return s.repository.GetByID(id)
}

func (s *Service) GetAllItems() ([]*Item, error) {
	return s.repository.GetAll()
}

func (s *Service) UpdateItem(item *Item) (int64, error) {
	return s.repository.Update(item)
}

func (s *Service) DeleteItem(id int64) (int64, error) {
	return s.repository.DeleteByID(id)
}

func NewService(repository RepositoryInterface) ServiceInterface {
	return &Service{repository}
}
