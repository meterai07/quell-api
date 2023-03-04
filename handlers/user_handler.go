package handlers

import (
	"net/http"
	"quell-api/entity"
	"quell-api/sdk/response"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	user := c.MustGet("user").(entity.User)
	response.Response(c, http.StatusOK, "success", user)
}
