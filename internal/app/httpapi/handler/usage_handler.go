package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fardannozami/activity-tracker/internal/repo/cache"
	"github.com/fardannozami/activity-tracker/internal/repo/postgres"
	"github.com/gin-gonic/gin"
)

type UsageHandler struct {
	Cache cache.Cache
	Usage *postgres.UsageRepo
}

func NewUsageHandler(ca cache.Cache, ur *postgres.UsageRepo) *UsageHandler {
	return &UsageHandler{
		Cache: ca,
		Usage: ur,
	}
}

// GetDailyUsage godoc
// @Summary Get daily usage
// @Tags Usage
// @Security BearerAuth
// @Produce json
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} any
// @Router /api/usage/daily [get]
func (h *UsageHandler) Daily(c *gin.Context) {
	clientId, _ := c.Get("client_id")
	cid, _ := clientId.(string)
	log.Println(cid)
	ver, _ := h.getVer(c, cid)
	key := fmt.Sprintf("usage:daily:%s:v%d", cid, ver)
	if s, ok, _ := h.Cache.Get(c.Request.Context(), key); ok {
		c.Data(http.StatusOK, "application/json", []byte(s))
		return
	}

	rows, err := h.Usage.GetDailyLast7(c.Request.Context(), cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	out := make([]gin.H, 0, len(rows))
	for _, r := range rows {
		out = append(out, gin.H{
			"day":   r.Day.Format("2006-01-02"),
			"total": r.Total,
		})
	}

	b, _ := json.Marshal(gin.H{"client_id": cid, "data": out})
	_ = h.Cache.Set(c.Request.Context(), key, string(b), 1*time.Hour)
	c.Data(http.StatusOK, "application/json", b)
}

// GetTopUsage godoc
// @Summary Get Top usage
// @Tags Usage
// @Security BearerAuth
// @Produce json
// @Success 200 {object} any
// @Router /api/usage/top [get]
func (h *UsageHandler) Top(c *gin.Context) {
	// Global ranking, cache key based on time bucket (or global version)
	key := "usage:top:last24h"

	if s, ok, _ := h.Cache.Get(c.Request.Context(), key); ok {
		c.Data(http.StatusOK, "application/json", []byte(s))
		return
	}

	rows, err := h.Usage.GetTopLast24Hours(c.Request.Context(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	b, _ := json.Marshal(gin.H{"data": rows})
	_ = h.Cache.Set(c.Request.Context(), key, string(b), 1*time.Hour)
	c.Data(http.StatusOK, "application/json", b)
}

func (h *UsageHandler) getVer(c *gin.Context, cid string) (int64, error) {
	verKey := fmt.Sprintf("usage:ver:%s", cid)
	v, err := h.Cache.Incr(c.Request.Context(), verKey, 0, 7*24*time.Hour)
	if err != nil {
		return 0, err
	}

	return v, nil
}
