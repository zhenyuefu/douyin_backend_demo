package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	structs.Response
	VideoList []structs.Video `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	token := c.Query("token")
	var uid uint
	if token != "" {
		claims, err := middleware.ParseToken(token)
		if err != nil {
			log.Printf(err.Error())
		}
		uid = claims.ID
	}
	timeStamp := c.Query("latest_time")
	var nextTime int64
	if timeStamp == "" {
		nextTime = time.Now().Unix()
	} else {
		if len(timeStamp) > 10 {
			timeStamp = timeStamp[:10]
		}
		nextTime, _ = strconv.ParseInt(timeStamp, 10, 64)
	}
	videos, err := db.GetVideosBefore(time.Unix(nextTime, 0))
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	var videoList []structs.Video
	for _, v := range videos {
		var like int64
		if token != "" {
			like = db.DB.Model(&v).Where("uid = ?", uid).Association("Likes").Count()
		}
		videoList = append(videoList, structs.Video{
			Id:       v.ID,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Author: structs.User{
				Id:            v.Author.ID,
				Name:          v.Author.Identifier,
				FollowCount:   db.DB.Model(&v.Author).Association("Followers").Count(),
				FollowerCount: db.DB.Model(&v.Author).Association("Followings").Count(),
			},
			FavoriteCount: db.DB.Model(&v).Association("Likes").Count(),
			CommentCount:  db.DB.Model(&v).Association("Comments").Count(),
			IsFavorite:    like > 0,
			Title:         v.Title,
		})
	}
	log.Println(videoList)
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].CreatedAt.Unix()
	}
	fmt.Print(videoList)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  structs.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
