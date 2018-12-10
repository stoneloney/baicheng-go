package models

type Article struct {
	Model
	Title       string  `form:"title"`
	Desc        string  `form:"desc"`
	Content     string  `form:"content"`
	Channel     string  `form:"channel"`
	Author      string  `form:"author"`
	Thumburl    string  `form:"thumburl"`
	Status      string  `form:"status"`
	Createtime  string  `form:"createtime"`
	Modifytime  string  `form:"modifytime"`
}

// 获取列表
func ArticleList() []Article {
	var articles []Article
	db := GetDB()
	db.Order("id desc").Find(&articles)
	return articles
}