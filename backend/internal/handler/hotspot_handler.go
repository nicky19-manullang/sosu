package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sosu-backend/internal/repository"
)

type HotspotHandler struct {
	repo *repository.HotspotRepository
}

func NewHotspotHandler(repo *repository.HotspotRepository) *HotspotHandler {
	return &HotspotHandler{repo: repo}
}
func (h *HotspotHandler) GetHotspotsInBBox(c *gin.Context) {
	minLat, err1 := strconv.ParseFloat(c.Query("min_lat"), 64)
	minLng, err2 := strconv.ParseFloat(c.Query("min_lng"), 64)
	maxLat, err3 := strconv.ParseFloat(c.Query("max_lat"), 64)
	maxLng, err4 := strconv.ParseFloat(c.Query("max_lng"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter min_lat, min_lng, max_lat, max_lng wajib diisi"})
		return
	}

	since := time.Now().Add(-72 * time.Hour)

	hotspots, err := h.repo.FindWithinBBox(minLat, minLng, maxLat, maxLng, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": len(hotspots), "hotspots": hotspots})
}

func (h *HotspotHandler) GetHotspots(c *gin.Context) {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter lat wajib diisi"})
		return
	}
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter lng wajib diisi"})
		return
	}

	radiusKm := 50.0
	if r := c.Query("radius_km"); r != "" {
		if parsed, err := strconv.ParseFloat(r, 64); err == nil {
			radiusKm = parsed
		}
	}
	if radiusKm < 0 {
		radiusKm = 0
	}
	if radiusKm > 50 {
		radiusKm = 50
	}

	since := time.Now().Add(-72 * time.Hour)
	if days := c.Query("days"); days != "" {
		if parsed, err := strconv.Atoi(strings.TrimSpace(days)); err == nil {
			since = time.Now().Add(-time.Duration(parsed) * 24 * time.Hour)
		}
	}

	hotspots, err := h.repo.FindWithinRadius(lat, lng, radiusKm*1000, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": len(hotspots), "hotspots": hotspots})
}
