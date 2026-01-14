package auth

import (
	"net/http"

	auth_dto "github.com/MAJIAXIT/projname/api/internal/dto/auth"
	"github.com/MAJIAXIT/projname/api/internal/service/auth"
	"github.com/MAJIAXIT/projname/api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (h *handler) Login(c *gin.Context) {
	var req auth_dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.transactionsMiddleware.GetTx(c)

	user, err := h.authService.Authenticate(tx, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	accessToken, refreshToken, err := h.authService.GenerateTokens(
		tx, user.ID, user.Login, userAgent, ip)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, auth_dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 minutes
	})
}

func (h *handler) Signup(c *gin.Context) {
	var req auth_dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.transactionsMiddleware.GetTx(c)

	// Check if user already exists
	_, err := h.usersService.UserBeforeCreateExistsCheck(
		tx, req.Login)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Create user
	usr, err := h.usersService.CreateUser(tx, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	// Remove password from response
	usr.EncrPassword = ""
	c.JSON(http.StatusOK, usr)
}

func (h *handler) Refresh(c *gin.Context) {
	var req auth_dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.transactionsMiddleware.GetTx(c)

	claims, err := h.authService.ValidateToken(req.RefreshToken, auth.RefreshToken)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	_, err = h.authService.ValidateSessionByToken(tx, req.RefreshToken)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	if err := h.authService.DeleteSessionByToken(tx, req.RefreshToken); err != nil {
		utils.HandleError(c, err)
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ip := c.ClientIP()

	accessToken, refreshToken, err := h.authService.GenerateTokens(
		tx, claims.UserID, claims.Login, userAgent, ip)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, auth_dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    15 * 60, // 15 minutes
	})
}

func (h *handler) Logout(c *gin.Context) {
	var req auth_dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.transactionsMiddleware.GetTx(c)

	if err := h.authService.DeleteSessionByToken(tx, req.RefreshToken); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *handler) LogoutAll(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not present in claims"})
		return
	}

	tx := h.transactionsMiddleware.GetTx(c)

	if err := h.authService.DeleteSessionsByUserID(tx, userID.(uint)); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out from all devices successfully"})
}

func (h *handler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not present in claims"})
		return
	}

	tx := h.transactionsMiddleware.GetTx(c)

	usr, err := h.usersService.GetUserByID(tx, userID.(uint))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	usr.EncrPassword = ""
	c.JSON(http.StatusOK, usr)
}
