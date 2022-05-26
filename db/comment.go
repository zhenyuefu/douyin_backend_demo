package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/gorm"
)

type CommentModel struct {
	gorm.Model
	UID     uint
	VID     uint `gorm:"index"`
	Content string
}

func (c *CommentModel) TableName() string {
	return constants.CommentTableName
}
