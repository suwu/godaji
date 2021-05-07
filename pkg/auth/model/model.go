package model

import (
	"gorm.io/gorm"
	dbcore "mitaitech.com/oa/pkg/common/db"
)

func init() {
	dbcore.RegisterInjector(func(db *gorm.DB) {
		dbcore.SetupTableModel(db, &Pet{})
	})
}

type Pet struct {
	dbcore.BaseModel
}
