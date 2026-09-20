package handler

import (
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"sosu-backend/internal/model"
	"sosu-backend/internal/repository"
	"sosu-backend/internal/service"
)

type SubscriptionHandler struct {
	repo          *repository.SubscriptionRepository
	userRepo      *repository.UserRepository
	emailService  *service.EmailService
	statusService *service.StatusService
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository, userRepo *repository.UserRepository, emailService *service.EmailService, statusService *service.StatusService) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo, userRepo: userRepo, emailService: emailService, statusService: statusService}
}

type subscribeRequest struct {
	Latitude      *float64 `json:"latitude" binding:"required"`
	Longitude     *float64 `json:"longitude" binding:"required"`
	LocationLabel string   `json:"location_label"`
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req subscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude dan longitude wajib diisi"})
		return
	}

	if !validCoordinates(*req.Latitude, *req.Longitude) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude atau longitude tidak valid"})
		return
	}

	if req.LocationLabel == "" {
		req.LocationLabel = "lokasi pilihan"
	}

	sub := &model.Subscription{
		UserID:        userID,
		Latitude:      *req.Latitude,
		Longitude:     *req.Longitude,
		LocationLabel: req.LocationLabel,
	}

	if err := h.repo.Create(sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err == nil && user != nil {
		go func() {
			if err := h.emailService.SendConfirmation(user.Email, req.LocationLabel); err != nil {
				log.Println("[subscribe] failed to send confirmation email:", err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"status": "subscribed"})
}

func validCoordinates(latitude, longitude float64) bool {
	return !math.IsNaN(latitude) && !math.IsInf(latitude, 0) && latitude >= -90 && latitude <= 90 &&
		!math.IsNaN(longitude) && !math.IsInf(longitude, 0) && longitude >= -180 && longitude <= 180
}

type subscriptionWithStatus struct {
	ID            uint    `json:"id"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	LocationLabel string  `json:"location_label"`
	FinalStatus   string  `json:"final_status"`
	AQICategory   string  `json:"aqi_category"`
	AQIAvailable  bool    `json:"aqi_available"`
	HotspotCount  int64   `json:"hotspot_count"`
}

func (h *SubscriptionHandler) MySubscriptions(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	subs, err := h.repo.FindByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]subscriptionWithStatus, 0, len(subs))
	for _, sub := range subs {
		statusResult, err := h.statusService.GetStatus(sub.Latitude, sub.Longitude)
		item := subscriptionWithStatus{
			ID:            sub.ID,
			Latitude:      sub.Latitude,
			Longitude:     sub.Longitude,
			LocationLabel: sub.LocationLabel,
		}
		if err == nil {
			item.FinalStatus = string(statusResult.FinalStatus)
			item.AQICategory = statusResult.AQICategory
			item.AQIAvailable = statusResult.AQIAvailable
			item.HotspotCount = statusResult.HotspotCount
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"subscriptions": result})
}

type unsubscribeRequest struct {
	Latitude  *float64 `json:"latitude" binding:"required"`
	Longitude *float64 `json:"longitude" binding:"required"`
}

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req unsubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude dan longitude wajib diisi"})
		return
	}

	if !validCoordinates(*req.Latitude, *req.Longitude) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude atau longitude tidak valid"})
		return
	}

	if err := h.repo.DeleteByUserAndLocation(userID, *req.Latitude, *req.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "unsubscribed"})
}
