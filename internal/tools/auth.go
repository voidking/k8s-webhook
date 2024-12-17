package tools

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Authorize(allowedUsers gin.Accounts) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, _ := c.Request.BasicAuth()
		logrus.Debugf("user: %s", user)
		if _, ok := allowedUsers[user]; ok {
			if allowedUsers[user] == pass {
				c.Next() // 调用后续的处理函数
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed!"})
		c.Abort() // 终止后续的处理
	}
}
