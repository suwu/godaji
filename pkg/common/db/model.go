package db

import (
	"context"
	"time"
)

// tag按首字母排序
type BaseModel struct {
	Id        string    `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type Object interface {
	GetId() string
}

func (s *BaseModel) GetId() string {
	return s.Id
}

type ITransaction interface {
	Transaction(ctx context.Context, fc func(txctx context.Context) error) error
}

type NoopTransaction struct{}

func (*NoopTransaction) Transaction(ctx context.Context, fc func(txctx context.Context) error) error {
	return fc(ctx)
}
