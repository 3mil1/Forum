package service

import "forum/internal/storage"

type Service struct {
	UserService *UserService
	PostService *PostService
	storage     *storage.Storage
	session     *Session
}

func (s *Service) User() *UserService {
	if s.UserService != nil {
		return s.UserService
	}

	s.UserService = &UserService{
		storage: s.storage,
	}

	return s.UserService
}

func (s *Service) Session() *Session {
	if s.session != nil {
		return s.session
	}
	s.session = &Session{
		service: s,
		storage: s.storage,
	}
	return s.session
}

func (s *Service) Post() *PostService {
	if s.PostService != nil {
		return s.PostService
	}
	s.PostService = &PostService{
		storage: s.storage,
	}
	return s.PostService
}

func New(s *storage.Storage) *Service {
	return &Service{storage: s}
}
