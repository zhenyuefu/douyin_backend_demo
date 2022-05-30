package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"gorm.io/gorm"
	"time"
)

type VideoModel struct {
	gorm.Model
	AuthorID uint `gorm:"index"`
	PlayUrl  string
	CoverUrl string
	Likes    []LikeModel    `gorm:"foreignKey:VID"`
	Comments []CommentModel `gorm:"foreignKey:VID"`
	Title    string
}

func (v *VideoModel) TableName() string {
	return constants.VideoTableName
}

// CreateVideo 创建视频
func CreateVideo(video *VideoModel) error {
	return DB.Create(video).Error
}

// GetVideoList 获取视频
func GetVideoList(uid uint) ([]VideoModel, error) {
	var videos []VideoModel
	result := DB.Where("author_id = ?", uid).Find(&videos)
	return videos, result.Error
}

// GetVideosBefore 获取视频
func GetVideosBefore(time time.Time) ([]VideoModel, error) {
	var videos []VideoModel
	result := DB.Where("updated_at < ?", time.Format("2006-01-02 15:04:05")).Limit(30).Order("updated_at desc").Find(&videos)
	return videos, result.Error
}
