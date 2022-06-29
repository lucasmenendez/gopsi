package session

type Session struct {
}

type Manager struct {
	sessions []*Session
}

func NewManager() *Manager {
	return &Manager{
		sessions: []*Session{},
	}
}
