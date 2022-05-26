package constants

const (
	MySQLDefaultDSN  = "root:password@tcp(localhost:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	UserTableName    = "user"
	FollowTableName  = "follow"
	VideoTableName   = "video"
	CommentTableName = "comment"
	LikeTableName    = "like"
)

var JwtSignKey = []byte("douyin")
