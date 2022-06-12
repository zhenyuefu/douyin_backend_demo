package constants

import "os"

const (
	MySQLDefaultDSN  = "root:password@tcp(host.containers.internal:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	UserTableName    = "user"
	VideoTableName   = "video"
	CommentTableName = "comment"
)

var JwtSignKey = []byte("douyin")
var VideoURLPrefix string

func Init() {
	VideoURLPrefix = os.Getenv("VIDEO_URL_PREFIX")
}
