package auth

import (
	"context"

	"mitaitech.com/oa/pkg/common/dbcore"
)

type authDomain struct{}

func NewAuthDomain() *authDomain {
	return &authDomain{}
}

func (a *authDomain) UserDao(ctx context.Context) IUserDao {
	return &userDao{dbcore.GetDB(ctx)}
}
