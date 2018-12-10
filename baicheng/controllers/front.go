package controllers

import (
	//"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"` 
	}
	resp.Code = 0
	resp.Msg = "success"
	c.JSON(http.StatusOK, resp)
}