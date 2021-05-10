package auth

import (
	"context"

	"gorm.io/gorm"
	"mitaitech.com/oa/pkg/common/dbcore"
)

func init() {
	dbcore.RegisterInjector(func(db *gorm.DB) {
		dbcore.SetupTableModel(db, &User{})
	})
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(ctx context.Context) *userDao {
	return &userDao{dbcore.GetDB(ctx)}
}

func (dao *userDao) List(query *User, offset, limit int) ([]*User, error) {
	var r []*User

	db := dbcore.WithOffsetLimit(dao.db, offset, limit)

	err := db.Where(query).Find(&r).Error
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (dao *userDao) Get(id string) (*User, error) {
	var r User
	err := dao.db.Where("id = ?", id).First(&r).Error
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (dao *userDao) Create(in *User) (*User, error) {
	err := dao.db.Create(in).Error
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (dao *userDao) Update(in *User) (*User, error) {
	err := dao.db.Updates(in).Error
	if err != nil {
		return nil, err
	}

	return in, nil
}

func (dao *userDao) Delete(in *User) error {
	err := dao.db.Where(in).Delete(&User{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (dao *userDao) GetByName(name string) (*User, error) {
	var r User
	err := dao.db.Where("name = ?", name).First(&r).Error
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (dao *userDao) GetByTel(tel string) (*User, error) {
	var r User
	err := dao.db.Where("tel = ?", tel).First(&r).Error
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (dao *userDao) GetByEmail(email string) (*User, error) {
	var r User
	err := dao.db.Where("email = ?", email).First(&r).Error
	if err != nil {
		return nil, err
	}

	return &r, nil
}
