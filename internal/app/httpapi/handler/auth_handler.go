package handler

import (
	"net/http"

	"github.com/fardannozami/activity-tracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct{ Issue *usecase.IssueTokenUC }

func NewAuthHandler(uc *usecase.IssueTokenUC) *AuthHandler { return &AuthHandler{Issue: uc} }

type tokenReq struct {
	Email  string `json:"email" binding:"required,email" example:"halo@halo.com"`
	APIKey string `json:"api_key" binding:"required" example:"ab90c204707f8251cf370a7a60e3b31f72dd58fe783c6d480c163835f8cefd77"`
}

// IssueToken godoc
// @Summary Issue JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body tokenReq true "Credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/auth/token [post]
func (h *AuthHandler) Token(c *gin.Context) {
	var req tokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwtStr, clientID, err := h.Issue.Execute(c.Request.Context(), req.Email, req.APIKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":       jwtStr,
		"client_id":          clientID,
		"token_type":         "Bearer",
		"expires_in_seconds": 7200,
	})
}
