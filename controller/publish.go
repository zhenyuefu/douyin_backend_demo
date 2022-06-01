package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/constants"
	"github.com/RaymondCode/simple-demo/db"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/structs"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"net/http"
	"path/filepath"
)

type VideoListResponse struct {
	structs.Response
	VideoList []structs.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	var user db.UserModel
	var uid uint
	if claim, exist := c.Get("user"); exist {
		uid = claim.(*middleware.Claims).ID
	} else {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}
	err := db.GetUser(&user, uid)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	imageName := fmt.Sprintf("%s.jpg", finalName)
	imageUrl := fmt.Sprintf("%s.jpg", saveFile)
	if err := ffmpeg.Input(saveFile).
		Output(imageUrl, ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		Run(); err != nil {
		log.Printf("ffmpeg error: %v", err.Error())
		imageUrl = ""
	}
	video := db.VideoModel{
		Author:   user,
		PlayUrl:  constants.VideoURLPrefix + finalName,
		CoverUrl: constants.VideoURLPrefix + imageName,
		Title:    c.PostForm("title"),
	}
	if err := db.CreateVideo(&video); err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, structs.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var uid uint
	if claim, exist := c.Get("user"); exist {
		uid = claim.(*middleware.Claims).ID
	} else {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}

	videos, err := db.GetVideoList(uid)
	if err != nil {
		c.JSON(http.StatusOK, structs.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	var videoList []structs.Video
	for _, v := range videos {
		like := db.DB.Model(&v).Where("uid = ?", uid).Association("Likes").Count()
		videoList = append(videoList, structs.Video{
			Id:            v.ID,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: db.DB.Model(&v).Association("Likes").Count(),
			CommentCount:  db.DB.Model(&v).Association("Comments").Count(),
			IsFavorite:    like > 0,
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
