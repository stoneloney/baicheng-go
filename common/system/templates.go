package system

import (
	"fmt"
	"strings"
	"os"
	"time"
	"path/filepath"
	"html/template"

	"common/models"

	"github.com/gin-gonic/gin"
)

var tmpl *template.Template

func LoadTemplates() {
	tmpl = template.New("").Funcs(template.FuncMap{
		"now":   			now,
		"activeUsername":   activeUsername,
		"formatDateTime":   formatDateTime,
		"isActiveLink":	    isActiveLink,
		//"channels":         channels,
		//"articles":         articles,
	})
	fn := func (path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			var err error
			tmpl, err = tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if err := filepath.Walk("views", fn); err != nil {
		panic(err)
	}
}

func GetTemplates() *template.Template {
	return tmpl
}

func now() time.Time {
	return time.Now()
}

// 时间格式化
func formatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// 当前用户名
func activeUsername(c *gin.Context) string {
	if c != nil {
		u, _ := c.Get("User")
		if user, ok := u.(*models.User); ok {
			return user.Name
		}
	}
	return ""
}

// 是否当前链接
func isActiveLink(c *gin.Context, uri string) string {
	if c != nil {
		var uris = strings.Split(c.Request.RequestURI, "/")
		var action string
		if len(uris[2]) == 0 {
			action = "/"
		} else {
			action = uris[2]
		}
		if action == uri {
			return "active"
		}
	}
	return ""
}

// 频道列表
/*
func channels(c *gin.Context, desc string) []models.Channel {
	return models.ChannelList()
}

// 资讯列表
func articles(c *gin.Context, id int, number int) []models.Article {
	return models.ArticleList()
}
*/

