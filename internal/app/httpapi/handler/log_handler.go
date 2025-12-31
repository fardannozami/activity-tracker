package handler

import (
	"net/http"
	"time"

	"github.com/fardannozami/activity-tracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	Record *usecase.RecordHitUC
}

func NewLogHandler(uc *usecase.RecordHitUC) *LogHandler {
	return &LogHandler{
		Record: uc,
	}
}

type logReq struct {
	IP        string `json:"ip" binding:"required"`
	Endpoint  string `json:"endpoint" binding:"required"`
	Timestamp string `json:"timestamp" binding:"required"`
}

func (h *LogHandler) Create(c *gin.Context) {
	var req logReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "timestamp must be RFC3339"})
		return
	}

	ts, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "timestamp must be RFC3339"})
		return
	}

	clientID, _ := c.Get("client_id")
	cid, _ := clientID.(string)

	_ = h.Record.Execute(c.Request.Context(), usecase.HitIn{
		ClientID:  cid,
		IP:        req.IP,
		Endpoint:  req.Endpoint,
		Timestamp: ts,
	})

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
