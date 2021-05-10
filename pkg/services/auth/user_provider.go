package auth

import (
	"context"
	"errors"

	"mitaitech.com/oa/pkg/common/hasher"
)

type IAuthenticatable interface {
	GetId() string
	GetPassword() string
	GetRememberToken() string
	SetRememberToken(token string) error
}

type IUserProvider interface {
	RetrieveById(id string) (IAuthenticatable, error)
	RetrieveByToken(id, token string) (IAuthenticatable, error)
	RetrieveByCredentials(username, password string) (IAuthenticatable, error)
}

type databaseUserProvider struct {
	dao *userDao
}

func NewDatabaseUserProvider(ctx context.Context) *databaseUserProvider {
	return &databaseUserProvider{NewUserDao(ctx)}
}

func (p *databaseUserProvider) RetrieveById(id string) (IAuthenticatable, error) {
	user, err := p.dao.Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *databaseUserProvider) RetrieveByToken(id, token string) (IAuthenticatable, error) {
	user, err := p.dao.Get(id)
	if err != nil {
		return nil, err
	}
	if user.RememberToken != token {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (p *databaseUserProvider) RetrieveByCredentials(username, password string) (IAuthenticatable, error) {
	user, err := p.dao.GetByName(username)
	if err != nil {
		return nil, err
	}

	valid, err := hasher.Check(password, user.password)
	if err != nil {
		return nil, errors.New("credentials wrong.")
	}

	if !valid {
		return nil, errors.New("password wrong.")
	}

	return user, nil
}
