package controller

import (
	"utils/errmsg"
	"gorm.io/driver/sqlite"
	//"github.com/RaymondCode/simple-demo/structs"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	//"net/http"
)

type Userinfo struct {
	gorm.Model
	Username string `json:"username" gorm:"column:'username'"`
	Password string `json:"password"`
}


type Comments struct {
	Id uint64 `json:"id"`
	Userid int32 `json:"userid"`
	Content string `json:"content"`
	ArticleId string `json:"article_id"`
	ParentId uint64 `json:"parent_id"`
	Child []*Comments `json:"child" gorm:"-"`
}


var db *gorm.DB
var err error
// AddComment 新增评论
func AddComment(data *Comments) int {
	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//序列化评论
func GetComments(data []Comments) []*Comments {
	fmt.Println(data)
	if data == nil || len(data) == 0 {
		return []*Comments{}
	}
	fmt.Println("这边")

	mapComment := make(map[uint64]*Comments,len(data))
	ret := []*Comments{}

	for _, cm := range data{
		fmt.Println("过来了")
		cm.Child = []*Comments{}
		fmt.Println(cm.Content)
		//if cm.ParentId == 0 {
		//	ret = append(ret, cm)
		//}
		//mapComment[cm.Id] =
	}

	for _, c := range data {
		if c.ParentId != 0 {
			parent := mapComment[c.ParentId]
			c.Child = []*Comments{}
			parent.Child = append(parent.Child)

		}
	}
	return ret
}




//type CommentListResponse struct {
//	structs.Response
//	CommentList []structs.Comment `json:"comment_list,omitempty"`
//}
//
//type CommentActionResponse struct {
//	structs.Response
//	Comment structs.Comment `json:"comment,omitempty"`
//}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	token := c.Query("token")
	actionType := c.Query("action_type")
	fmt.Println(token)
	fmt.Println(actionType)

	//if user, exist := usersLoginInfo[token]; exist {
	//	if actionType == "1" {
	//		text := c.Query("comment_text")
	//		c.JSON(http.StatusOK, CommentActionResponse{Response: structs.Response{StatusCode: 0},
	//			Comment: structs.Comment{
	//				Id:         1,
	//				User:       user,
	//				Content:    text,
	//				CreateDate: "05-01",
	//			}})
	//		return
	//	}
	//	c.JSON(http.StatusOK, structs.Response{StatusCode: 0})
	//} else {
	//	c.JSON(http.StatusOK, structs.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	//}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	var data []Comments
	db.Model(&Comments{}).Find(&data)
	if data == nil || len(data) == 0 {
		c.JSON(200,gin.H{
			"data": "没有数据",
		})
		return
	}
	//s := GetComments(data)
	c.JSON(200,data)
}

func CreateComment(ctx *gin.Context)  {
	fmt.Println("请求过来了")
	var comment Comments
	err = ctx.ShouldBindJSON(&comment)
	if err != nil{
		ctx.JSON(200,gin.H{
			"data": "错误",
			"error": err.Error(),
		})
		return
	}
	fmt.Println(comment.Id)
	db.Create(&comment)

	ctx.JSON(200,"成功")
}

func main()  {
	db, err = gorm.Open(sqlite.Open("test1.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Comments{},&Userinfo{})
	r := gin.Default()
	r.POST("createcomment",CreateComment)
	r.GET("listcomment",CommentList)
	r.Run()
}