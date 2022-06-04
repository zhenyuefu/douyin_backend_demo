package comment

import (
	"fmt"
	"utils/errmsg"

	"gorm.io/gorm"
)

type Comments struct {
	Id        uint64      `json:"id"`
	Userid    int32       `json:"userid"`
	Content   string      `json:"content"`
	ArticleId string      `json:"article_id"`
	ParentId  uint64      `json:"parent_id"`
	Child     []*Comments `json:"child" gorm:"-"`
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
func GetComments(data []*Comments) []*Comments {
	if data == nil || len(data) == 0 {
		return []*Comments{}
	}
	fmt.Println("这边")

	mapComment := make(map[uint64]*Comments, len(data))
	ret := []*Comments{}

	for _, c := range data {
		c.Child = []*Comments{}
		if c.ParentId == 0 {
			ret = append(ret, c)
		}
		mapComment[c.Id] = c
	}

	for _, c := range data {
		if c.ParentId != 0 {
			parent := mapComment[c.ParentId]
			c.Child = []*Comments{}
			parent.Child = append(parent.Child, c)

		}
	}
	return ret
}

// GetComment 查询单个评论
func GetComment(id int) (Comments, int) {
	var comment Comments
	err = db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return comment, errmsg.ERROR
	}
	return comment, errmsg.SUCCSE
}

// GetCommentList 后台所有获取评论列表
func GetCommentList(pageSize int, pageNum int) ([]Comments, int64, int) {

	var commentList []Comments
	var total int64
	db.Find(&commentList).Count(&total)
	err = db.Model(&commentList).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Select("comment.id, article.title,user_id,article_id, user.username, comment.content, comment.status,comment.created_at,comment.deleted_at").Joins("LEFT JOIN article ON comment.article_id = article.id").Joins("LEFT JOIN user ON comment.user_id = user.id").Scan(&commentList).Error
	if err != nil {
		return commentList, 0, errmsg.ERROR
	}
	return commentList, total, errmsg.SUCCSE
}

//// GetCommentCount 获取评论数量
func GetCommentCount(id int) int64 {
	var comment Comments
	var total int64
	db.Find(&comment).Where("article_id = ?", id).Where("status = ?", 1).Count(&total)
	return total
}

// GetCommentListFront 展示页面获取评论列表
func GetCommentListFront(id int, pageSize int, pageNum int) ([]Comments, int64, int) {
	var commentList []Comments
	var total int64
	db.Find(&Comments{}).Where("article_id = ?", id).Where("status = ?", 1).Count(&total)
	err = db.Model(&Comments{}).Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Select("comment.id, article.title, user_id, article_id, user.username, comment.content, comment.status,comment.created_at,comment.deleted_at").Joins("LEFT JOIN article ON comment.article_id = article.id").Joins("LEFT JOIN user ON comment.user_id = user.id").Where("article_id = ?",
		id).Where("status = ?", 1).Scan(&commentList).Error
	if err != nil {
		return commentList, 0, errmsg.ERROR
	}
	return commentList, total, errmsg.SUCCSE
}

// 编辑评论

// DeleteComment 删除评论
func DeleteComment(id uint) int {
	var comment Comments
	err = db.Where("id = ?", id).Delete(&comment).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
