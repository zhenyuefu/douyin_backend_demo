package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/gorm"
)

type LikeModel struct {
	gorm.Model
	UID    uint `gorm:"index"`
	VID    uint `gorm:"index"`
	Cancel bool `gorm:"default:false"`
}

func (l *LikeModel) TableName() string {
	return constants.LikeTableName
}
