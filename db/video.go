package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/gorm"
)

type VideoModel struct {
	gorm.Model
	AuthorID uint `gorm:"index"`
	PlayUrl  string
	CoverUrl string
	Likes    []LikeModel    `gorm:"foreignKey:VID"`
	Comments []CommentModel `gorm:"foreignKey:VID"`
}

func (v *VideoModel) TableName() string {
	return constants.VideoTableName
}
