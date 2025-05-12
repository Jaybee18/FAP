package services

import (
	"errors"
	"fap-server/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	users map[string]models.User // In-memory store
	sessions map[string]models.Session // In-memory session store
	mu    sync.RWMutex           // For thread safety
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]models.User),
		sessions: make(map[string]models.Session),
	}
}

func (s *UserService) AddUser(user models.User) (models.AddUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.LoginName]; exists {
		return models.AddUserResponse{
			Result:  false,
			Message: "User already exists",
		}, errors.New("user already exists")
	}

	s.users[user.LoginName] = user

	return models.AddUserResponse{
		Result:  true,
		Message: "User created successfully",
	}, nil
}

func (s *UserService) Login(loginName, password string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[loginName]
	if !exists {
		return "", errors.New("invalid credentials")
	}

	if user.Password.Password != password {
		return "", errors.New("invalid credentials")
	}

	sessionID := uuid.New().String()
	s.sessions[sessionID] = models.Session{
		ID: sessionID,
		UserID: user.LoginName,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return sessionID, nil
}

func (s *UserService) GetUser(loginName, sessionID string) (models.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    userID, err := s.ValidateSession(sessionID)
    if err != nil {
        return models.User{}, err
    }

    if loginName != "" && loginName != userID {
        return models.User{}, errors.New("unauthorized access")
    }

    user, exists := s.users[userID]
    if !exists {
        return models.User{}, errors.New("user not found")
    }

    return user, nil
}

func (s *UserService) Logout(sessionID, loginName string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Validate the session belongs to the user
    if session, exists := s.sessions[sessionID]; exists {
        if session.UserID == loginName {
            delete(s.sessions, sessionID)
            return true
        }
    }
    return false
}

func (s *UserService) CleanupSessions() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    now := time.Now()
    for id, session := range s.sessions {
        if now.After(session.ExpiresAt) {
            delete(s.sessions, id)
        }
    }
}

func (s *UserService) ValidateSession(sessionID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

    session, exists := s.sessions[sessionID]
    if !exists {
        return "", errors.New("invalid session")
    }

    if time.Now().After(session.ExpiresAt) {
        return "", errors.New("session expired")
    }

    return session.UserID, nil
}

func (s *UserService) NameTaken(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.users[name]; !ok {
		return false
	}

	return true
}