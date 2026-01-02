package handler

import (
	"log"
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
	IP        string `json:"ip" binding:"required" example:"1.1.1.1"`
	Endpoint  string `json:"endpoint" binding:"required" example:"/orders"`
	Timestamp string `json:"timestamp" binding:"required" example:"2026-01-02T14:00:00+07:00"`
}

// CreateLog godoc
// @Summary Record API usage log
// @Description Ingest API usage logs (high throughput)
// @Tags Logs
// @Accept json
// @Produce json
// @Param X-API-Key header string true "API Key" example(ab90c204707f8251cf370a7a60e3b31f72dd58fe783c6d480c163835f8cefd77)
// @Param body body logReq true "Log payload"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/logs [post]
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

	if err := h.Record.Execute(c.Request.Context(), usecase.HitIn{
		ClientID:  cid,
		IP:        req.IP,
		Endpoint:  req.Endpoint,
		Timestamp: ts,
	}); err != nil {
		log.Printf("enqueue api hit failed: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
