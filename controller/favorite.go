package controller

import (
	"github.com/RaymondCode/simple-demo/middleware"
	"log"
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
	//获取uid
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	uid := user.(*middleware.Claims).ID
	//获取vid
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "Video id is not valid"})
		return
	}

	actionType := c.Query("action_type")
	switch actionType {
	case "1":
		{

			like := db.LikeModel{
				UID: uid,
				VID: uint(vid),
			}

			//Create likeModel
			if err := db.Like(&like); err != nil {
				c.JSON(http.StatusOK, structs.Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, structs.Response{StatusCode: 0})
			return
		}
	case "2":
		err := db.UnLike(uid, uint(vid))
		if err != nil {
			c.JSON(http.StatusOK, structs.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return
		}
	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {

	//Get user_id
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	uid := uint(userId)

	var videos []structs.Video
	err := db.GetLikeVideos(&videos, uid)
	log.Printf("%+v", videos)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: structs.Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
