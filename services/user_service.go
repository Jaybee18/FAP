package services

import (
	"errors"
	"fap-server/models"
	"sync"
)

type UserService struct {
	users map[string]models.User // In-memory store
	mu    sync.RWMutex           // For thread safety
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]models.User),
	}
}

func (s *UserService) AddUser(user models.User) (models.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	if _, exists := s.users[user.LoginName]; exists {
		return models.Response{
			Result:  false,
			Message: "User already exists",
		}, errors.New("user already exists")
	}

	// Store the user
	s.users[user.LoginName] = user

	return models.Response{
		Result:  true,
		Message: "User created successfully",
	}, nil
}

func (s *UserService) GetUsers() map[string]models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return a copy to prevent external modifications
	usersCopy := make(map[string]models.User)
	for k, v := range s.users {
		usersCopy[k] = v
	}
	return usersCopy
}