package auth

import "context"

type IGuard interface {
	Check() bool
	CurrentUser() (IAuthenticatable, error)
	Guest() bool
	Id() string
	Validate(username, password string) bool
	SetUser(user IAuthenticatable)
}

type sessionGuard struct {
	provider IUserProvider
}

func NewSessionGuard() *sessionGuard {
	return &sessionGuard{NewDatabaseUserProvider(context.Background())}
}

func (g *sessionGuard) Check() bool {
	user, err := g.CurrentUser()
	if err != nil {
		return false
	}
}
