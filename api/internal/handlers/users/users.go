package users

import (
	"net/http"

	users_dto "github.com/MAJIAXIT/api_base/api/internal/dto/users"
	"github.com/MAJIAXIT/api_base/api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *handler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not present in claims"})
		return
	}
	tx := h.transactionsMiddleware.GetTx(c)

	currentUser, err := h.usersService.GetUserByID(tx, userID.(uint))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	currentUser.EncrPassword = ""

	c.JSON(http.StatusOK, currentUser)
}

func (h *handler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not present in claims"})
		return
	}

	var req users_dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.transactionsMiddleware.GetTx(c)

	user, err := h.usersService.UpdateUser(tx, userID.(uint), &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
