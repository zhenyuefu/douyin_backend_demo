package controller

import (
	"net/http"
	"strconv"

	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
)

type FavoriteActionResponse struct {
	structs.Response
}

type FavoriteListResponse struct {
	structs.Response
	VideoList []structs.Video `json:"video_list"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	//Check token
	token := c.Query("token")
	_, exist := usersLoginInfo[token]
	if !exist {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	//Get uid,vid,action_type
	action_type := c.Query("action_type")
	cancel := false
	if action_type == "2" {
		cancel = true
	}
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	uid := uint(user_id)
	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	vid := uint(video_id)

	like := db.LikeModel{
		UID:    uid,
		VID:    vid,
		Cancel: cancel,
	}

	//Create likeModel
	if err := db.CreateLike(&like); err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, structs.Response{StatusCode: 0})
	return
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {

	//Get user_id
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	uid := uint(user_id)

	//Get LikeList by user_id
	likeList, err := db.GetLikeList(uid)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//Get videoList from LikeList
	var videoList []structs.Video
	for _, like := range likeList {
		var v db.VideoModel
		db.DB.Where("uid = ?", like.VID).First(&v)
		likes := db.DB.Model(&v).Where("uid = ?", uid).Association("Likes").Count()
		videoList = append(videoList, structs.Video{
			Id:            v.ID,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: db.DB.Model(&v).Association("Likes").Count(),
			CommentCount:  db.DB.Model(&v).Association("Comments").Count(),
			IsFavorite:    likes > 0,
			Title:         v.Title,
		})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: structs.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
