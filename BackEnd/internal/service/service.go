package service

import "forum/internal/storage"

type Service struct {
	userService *UserService
	postService *PostService
	storage     *storage.Storage
	session     *Session
}

func (s *Service) User() *UserService {
	if s.userService != nil {
		return s.userService
	}

	s.userService = &UserService{
		storage: s.storage,
	}

	return s.userService
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
	if s.postService != nil {
		return s.postService
	}
	s.postService = &PostService{
		storage: s.storage,
	}
	return s.postService
}

func New(s *storage.Storage) *Service {
	return &Service{storage: s}
}
