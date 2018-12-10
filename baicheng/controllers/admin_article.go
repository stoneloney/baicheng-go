package controllers

import (
	"strconv"
	"net/http"
	"html/template"

	"common/method"
	"common/models"
	"common/controllers"
	"common/system"

	"github.com/gin-gonic/gin"
)

func AdminArticle(c *gin.Context) {
	H := commonControllers.DefaultH(c)
	H["Title"] = "Admin article"
	H["AdminPath"] = system.GetAdminPath()

	page, _ := strconv.Atoi(c.Query("page"))
	pageNum := 20
	count := 0

	var articles []models.Article
	db := models.GetDB()

	// 总数
	db.Find(&articles).Count(&count)
	H["count"] = count
	H["pagenum"] = pageNum

	// 获取数据列表
	db.Limit(pageNum).Offset(page*pageNum).Order("createtime desc").Find(&articles)
	H["Articles"] = articles

	c.HTML(http.StatusOK, "admin/article", H)
}

// 资讯页面
func AdminArticlePage(c *gin.Context) {
	H := commonControllers.DefaultH(c)
	H["Title"] = "Admin article"
	H["AdminPath"] = system.GetAdminPath()

	// 获取频道
	var channels []models.Channel
	db := models.GetDB()
	db.Order("id desc").Find(&channels)
	H["Channels"] = channels

	c.HTML(http.StatusOK, "admin/article_content", H)
}

// 添加资讯
func AdminArticleAdd(c *gin.Context) {
	article := &models.Article{}
	// 标题
	title := c.PostForm("title")
	if len(title) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"标题不能为空"})
		return
	}
	title = template.HTMLEscapeString(title)
	article.Title = title

	// 简介
	desc := c.PostForm("desc")
	if len(desc) > 0 {
		desc = template.HTMLEscapeString(desc)
		article.Desc = desc
	}

	// 内容
	content := c.PostForm("editorValue")
	if len(content) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"内容不能为空"})
		return
	}
	// content = template.HTMLEscapeString(content)
	article.Content = content

	// 所属频道
	channel := c.PostForm("channel")
	if len(channel) == 0 || !method.IsNumeric(channel) {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"选择资讯所属频道"})
		return
	}
	article.Channel = channel

	// 作者
	author := c.PostForm("author")
	if len(author) > 0 {
		author = template.HTMLEscapeString(author)
		article.Author = author
	}

	// 缩略图
	thumburl := c.PostForm("thumburl")
	if len(thumburl) > 0 {
		thumburl = template.HTMLEscapeString(thumburl)
		article.Thumburl = thumburl
	}

	// 资讯状态
	status := c.PostForm("status")
	if len(status) == 0 || !method.IsNumeric(status) {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":"状态参数错误"})
		return
	}
	article.Status = status

	now := method.Now()
	article.Modifytime = now

	db := models.GetDB()

	aid := c.PostForm("aid")
	if len(aid) > 0 {
		aid, _ := strconv.ParseInt(aid, 10, 64) 
		article.ID = aid
		if err := db.Model(&models.Article{}).Updates(article).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code":5, "msg":err.Error()})
			return
		}
	} else {
		article.Createtime = now
		if err := db.Create(&article).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code":4, "msg":err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}

// 获取资讯
func AdminArticleEdit(c *gin.Context) {
	id := c.Query("id");
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	db := models.GetDB()
	article := models.Article{}
	db.First(&article, id)
	if article.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"资讯不存在"})
		return
	}

	H := commonControllers.DefaultH(c)
	H["Title"] = "Admin article"
	H["AdminPath"] = system.GetAdminPath()
	H["Article"] = article

	// 获取频道
	var channels []models.Channel
	db.Order("id desc").Find(&channels)
	H["Channels"] = channels

	c.HTML(http.StatusOK, "admin/article_content", H)
}

// 删除资讯
func AdminArticleDelete(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{"code":1, "msg":"参数错误"})
		return
	}
	db := models.GetDB()
	article := models.Article{}
	db.First(&article, id)
	if article.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"code":2, "msg":"资讯不存在"})
		return
	}
	if err := db.Delete(&article).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code":3, "msg":err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code":0, "msg":""})
}



