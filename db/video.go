package db

import (
	"github.com/RaymondCode/simple-demo/constants"
	"github.com/RaymondCode/simple-demo/structs"
	"gorm.io/gorm"
	"log"
	"time"
)

type VideoModel struct {
	gorm.Model
	AuthorID uint      `gorm:"index"`
	Author   UserModel `gorm:"foreignKey:AuthorID"`
	PlayUrl  string
	CoverUrl string
	Likes    []UserModel    `gorm:"many2many:like;joinForeignKey:vid;joinReferences:uid"`
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
func GetVideoList(uid uint) ([]structs.Video, error) {
	var videos []structs.Video
	var VideoModels []VideoModel
	result := VideoSelect(&VideoModels, uid).
		Where(&VideoModel{AuthorID: uid}).Scan(&videos)
	log.Printf("%+v", videos)
	return videos, result.Error
}

// GetVideosBefore 获取视频
func GetVideosBefore(videos *[]structs.Video, time time.Time, uid uint) (error, int64) {
	var videoModels []VideoModel
	var nextTime int64
	result := VideoSelect(&videoModels, uid).
		Where("video.updated_at < ?", time.Format("2006-01-02 15:04:05")).
		Limit(30).
		Order("video.updated_at desc").
		Scan(&videos)
	log.Printf("%+v", videos)
	if len(videoModels) > 0 {
		nextTime = videoModels[len(videoModels)-1].CreatedAt.Unix()
	}
	return result.Error, nextTime
}

// GetLikeVideos 获取点赞的视频
func GetLikeVideos(videos *[]structs.Video, uid uint) error {
	var videoModels []VideoModel
	result := VideoSelect(&videoModels, uid).
		Joins("left join `like` on `like`.`vid` = `video`.`id`").
		Where("`like`.`uid` = ?", uid).
		Scan(&videos)
	return result.Error
}

func VideoSelect(videoModels *[]VideoModel, uid uint) *gorm.DB {
	var likeCountQuery = DB.Table("like").Select("vid,count(*) as favorite_count").Group("vid")
	var followCountQuery = DB.Table("follows").Select("uid,count(uid) as follow_count").Group("uid")
	var followerCountQuery = DB.Table("follows").Select("fid,count(fid) as follower_count").Group("fid")
	var totalFavoriteQuery = DB.Table("like").Select("author_id,count(*) as count").Joins("LEFT JOIN video v on v.id = `like`.vid").Group("author_id")
	var userLikeCountQuery = DB.Table("like").Select("uid,count(*) as favorite_count").Group("uid")
	var commentCountQuery = DB.Table("comment").Select("v_id,count(*) as comment_count").Where("deleted_at IS NULL").Group("v_id")
	var workCountQuery = DB.Table("video").Select("author_id,count(*) as work_count").Group("author_id")
	return DB.Model(&videoModels).
		Select("likeCount.favorite_count,"+
			"userLikeCount.favorite_count as Author__favorite_count,"+
			"followerCount.follower_count as Author__follower_count,"+
			"followCount.follow_count as Author__follow_count,"+
			"totalFavorite.count as Author__total_favorite,"+
			"workCount.work_count as Author__work_count,"+
			"commentCount.comment_count as comment_count,"+
			"EXISTS(SELECT 1 from `follows` WHERE fid=video.author_id and uid=? limit 1) as Author__is_follow,"+
			"EXISTS(SELECT 1 from `like` WHERE vid = video.id and uid = ? limit 1) as is_favorite,"+
			"`video`.`id`,"+
			"`video`.`deleted_at`,"+
			"`video`.`author_id`,"+
			"`video`.`play_url`,"+
			"`video`.`cover_url`,"+
			"`video`.`title`,"+
			"`Author`.`id` AS `Author__id`,"+
			"`Author`.`name` AS `Author__name`,"+
			"`Author`.`avatar` AS `Author__avatar`", uid, uid).
		Joins("Author").
		Joins("left join (?) as likeCount on video.id = likeCount.vid ", likeCountQuery).
		Joins("left join (?) as followCount on video.author_id = followCount.uid ", followCountQuery).
		Joins("left join (?) as followerCount on video.author_id = followerCount.fid ", followerCountQuery).
		Joins("left join (?) as userLikeCount on video.author_id = userLikeCount.uid ", userLikeCountQuery).
		Joins("left join (?) as totalFavorite on video.author_id = totalFavorite.author_id ", totalFavoriteQuery).
		Joins("left join (?) as workCount on video.author_id = workCount.author_id ", workCountQuery).
		Joins("left join (?) as commentCount on video.id = commentCount.v_id ", commentCountQuery)
}
