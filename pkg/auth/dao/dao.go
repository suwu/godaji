package dao

import (
	"github.com/win5do/go-lib/errx"
	"gorm.io/gorm"
)

type petDb struct {
	db *gorm.DB
}

func (s *petDb) List(query *petmodel.Pet, offset, limit int) ([]*petmodel.Pet, error) {
	var r []*petmodel.Pet

	db := dbcore.WithOffsetLimit(s.db, offset, limit)

	err := db.Where(query).Find(&r).Error
	if err != nil {
		return nil, errx.WithStackOnce(err)
	}

	return r, nil
}

func (s *petDb) Get(id string) (*petmodel.Pet, error) {
	var r petmodel.Pet
	err := s.db.Where("id = ?", id).First(&r).Error
	if err != nil {
		return nil, errx.WithStackOnce(err)
	}

	return &r, nil
}

func (s *petDb) Create(in *petmodel.Pet) (*petmodel.Pet, error) {
	err := s.db.Create(in).Error
	if err != nil {
		return nil, errx.WithStackOnce(err)
	}

	return in, nil
}

func (s *petDb) Update(in *petmodel.Pet) (*petmodel.Pet, error) {
	err := s.db.Updates(in).Error
	if err != nil {
		return nil, errx.WithStackOnce(err)
	}

	return in, nil
}

func (s *petDb) Delete(in *petmodel.Pet) error {
	err := s.db.Where(in).Delete(&petmodel.Pet{}).Error
	if err != nil {
		return errx.WithStackOnce(err)
	}

	return nil
}
