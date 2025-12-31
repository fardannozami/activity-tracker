package handler

import (
	"net/http"

	"github.com/fardannozami/activity-tracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	Register *usecase.RegisterClientUC
}

func NewClientHandler(r *usecase.RegisterClientUC) *ClientHandler {
	return &ClientHandler{Register: r}
}

type registerReq struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func (client *ClientHandler) RegisterClient(ctx *gin.Context) {
	var req registerReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := client.Register.Execute(ctx, req.Name, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, out)
}
