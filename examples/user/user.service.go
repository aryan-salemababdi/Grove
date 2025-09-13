package user

type UserService struct{}

func NewService() *UserService {
	return &UserService{}
}

func (s *UserService) FindAll() []string {
	return []string{"Aryan", "Jasmin"}
}
