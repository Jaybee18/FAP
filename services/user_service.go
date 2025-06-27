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

// UserExists tries accessing the entry in the user list
// with the given username as the key. When this operation
// succeeds, the user exists. Otherwise he doesn't exist.
func (s *UserService) UserExists(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.users[username]
	return ok
}

// Adds a user but only if a user with the same login name doesn't already
// exist. Returns true if the user was created and false otherwise
func (s *UserService) AddUser(user models.User) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.UserExists(user.LoginName) {
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

// Logout checks if the given username / session combination is valid
// and deletes the session from the session store if that is the case.
// It returns true if the deletion was successfull and false if the
// deletion fails for any reason
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

// CleanupSessions captures the current time and iterates through the
// stored list of sessions to delete (aka. invalidate) every session
// that has an expiry timestamp that lies in the past
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

// ValidSession verifies if a given username / session combination is valid.
// A session is valid for a given username if it exists and the username
// matches the username stored in the session and the session is not expired
func (s *UserService) ValidSession(username string, sessionID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sessionID]
	return exists && session.UserID == username && time.Now().Before(session.ExpiresAt)
}

// GetStandortOfUser determines the current location of the given user. It does
// so by first looking up if the user has a current stored location. If he doesn't
// have a stored location, his location is determined by requesting the coordinates
// of his stored home address via api.
func (s *UserService) GetStandortOfUser(username string) (*models.Location, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var user models.User
	for name, u := range s.users {
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

// SetStandortOfUser assigns the given location to the given user
func (s *UserService) SetStandortOfUser(username string, location models.Location) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[username]
	if !ok {
		return fmt.Errorf("no user exists with name %s", username)
	}

	user.Location = &location
	s.users[username] = user
	return nil
}
