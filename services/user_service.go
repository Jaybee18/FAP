package services

import (
	"fap-server/models"
	"fap-server/pkg"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	users    map[string]models.User    // In-memory store
	sessions map[string]models.Session // In-memory session store
	mu       sync.RWMutex              // For thread safety
}

func NewUserService() *UserService {
	return &UserService{
		users:    make(map[string]models.User),
		sessions: make(map[string]models.Session),
	}
}

func (u *UserService) UserExists(username string) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	_, ok := u.users[username]
	return ok
}

// Adds a user but only if a user with the same login name doesn't already
// exist. Returns true if the user was created and false otherwise
func (s *UserService) AddUser(user models.User) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.LoginName]; exists {
		return false
	}

	s.users[user.LoginName] = user

	return true
}

// Only returns a valid session id when user exists and credentials are correct.
// Otherwise an empty string is returned
func (s *UserService) Login(loginName, password string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[loginName]
	if !exists {
		return ""
	}

	if user.Password.Password != password {
		return ""
	}

	sessionID := uuid.New().String()
	s.sessions[sessionID] = models.Session{
		ID:        sessionID,
		UserID:    user.LoginName,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return sessionID
}

func (s *UserService) GetAllUsers() map[string]models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.users
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

func (s *UserService) ValidSession(username string, sessionID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	return exists && session.UserID == username && time.Now().Before(session.ExpiresAt)
}

func (u *UserService) GetStandortOfUser(username string) (*models.Location, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	var user models.User
	for name, u := range u.users {
		if name == username {
			user = u
			break
		}
	}

	if user.Location != nil {
		return user.Location, nil
	}

	// User currently has no entered location, so take his home adress instead
	if user.Country == nil || user.ZipCode == nil || user.City == nil || user.Street == nil {
		return nil, fmt.Errorf("user has no location set")
	}
	location, err := pkg.GetLocationByAdressGeoapify(*user.Country, *user.ZipCode, *user.City, *user.Street)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (u *UserService) SetStandortOfUser(username string, location models.Location) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	user, ok := u.users[username]
	if !ok {
		return fmt.Errorf("no user exists with name %s", username)
	}

	user.Location = &location
	u.users[username] = user
	return nil
}
