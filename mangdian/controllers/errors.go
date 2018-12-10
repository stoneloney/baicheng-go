package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	ShowErrorPage(c, http.StatusNotFound, nil)
}

func MethodNotAllowed(c *gin.Context) {
	ShowErrorPage(c, http.StatusMethodNotAllowed, nil)
}

func ShowErrorPage(c *gin.Context, code int, err error) {
	H := DefaultH(c)
	H["Error"] = err
	c.HTML(code, fmt.Sprintf("err_%d.html", code), H)
}