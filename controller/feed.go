package controller

import (
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
	var videos []structs.Video
	var err error
	log.Printf("uid: %d, nextTimeBefore: %d", uid, nextTime)
	err = db.GetVideosBefore(&videos, time.Unix(nextTime, 0), uid)
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].UpdatedAt.Unix()
	}
	//log.Printf("%+v", videos)
	log.Printf("uid: %d, nextTimeAfter: %d", uid, nextTime)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  structs.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  nextTime,
	})
}
