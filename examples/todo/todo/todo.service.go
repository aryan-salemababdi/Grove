package todo

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) FindAll() []string {
	return []string{"buy milk", "write Grove example", "read Go docs"}
}
