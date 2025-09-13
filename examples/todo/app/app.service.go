package app

type Service struct{}

func NewService() *Service  { return &Service{} }

func (s *Service) Greet() string { return  "hello from Grove!" }
