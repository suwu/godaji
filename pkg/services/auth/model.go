package auth

import (
	"context"
	"time"

	"mitaitech.com/oa/pkg/common/dbcore"
)

type User struct {
	dbcore.BaseModel
	Name            string
	Tel             string
	Email           string
	EmailVerifiedAt time.Time
	Password        string
	RememberToken   string
}

type IUserDao interface {
	List(query *User, offset, limit int) ([]*User, error)
	Get(id string) (*User, error)
	Create(in *User) (*User, error)
	Update(in *User) (*User, error)
	Delete(in *User) error
	GetByName(name string) (*User, error)
	GetByTel(tel string) (*User, error)
	GetByEmail(email string) (*User, error)
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) GetRememberToken() string {
	return u.RememberToken
}

func (u *User) SetRememberToken(token string) error {
	u.RememberToken = token
	_, err := NewUserDao(context.Background()).Update(u)
	if err != nil {
		return err
	}
	return nil
}

// $table->id();
// $table->string('name');
// $table->string('email')->unique();
// $table->timestamp('email_verified_at')->nullable();
// $table->string('password');
// $table->rememberToken();
// $table->foreignId('current_team_id')->nullable();
// $table->text('profile_photo_path')->nullable();
// $table->timestamps();
