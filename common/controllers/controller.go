package commonControllers

import (
	//"fmt"
	//"net/http"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func DefaultH(c *gin.Context) gin.H {
	return gin.H{
		"Title":  "",
		"Context":  c,
		"Csrf": csrf.GetToken(c),
	}
}



