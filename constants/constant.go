package constants

import (
	"os"
)

const (
	MySQLDefaultDSN  = "root:password@tcp(host.docker.internal:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	UserTableName    = "user"
	FollowTableName  = "follow"
	VideoTableName   = "video"
	CommentTableName = "comment"
	LikeTableName    = "like"
)

var JwtSignKey = []byte("douyin")
var VideoURLPrefix = os.Getenv("VIDEO_URL_PREFIX")
