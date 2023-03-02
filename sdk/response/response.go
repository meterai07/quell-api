package response

import "github.com/gin-gonic/gin"

func Response(c *gin.Context, status int, message string, data interface{}) {
	switch status / 100 {
	case 2:
		c.JSON(status, gin.H{
			"message": message,
			"data":    data,
		})
	case 4:
		c.JSON(status, gin.H{
			"message": message,
		})
	case 5:
		c.JSON(status, gin.H{
			"message": message,
		})
	}
}
