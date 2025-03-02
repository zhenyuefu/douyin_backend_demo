package structs

import "time"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	ID            uint      `json:"id,omitempty"`
	UpdatedAt     time.Time `json:"-"`
	Author        User      `json:"author" gorm:"foreignKey:ID"`
	PlayUrl       string    `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount  int64     `json:"comment_count,omitempty"`
	IsFavorite    bool      `json:"is_favorite,omitempty"`
	Title         string    `json:"title,omitempty"`
}

type Comment struct {
	ID         uint   `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	ID            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	TotalFavorite int64  `json:"total_favorited,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	WorkCount     int64  `json:"work_count,omitempty"`
}
