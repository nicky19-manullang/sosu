package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sosu-backend/internal/model"
	"sosu-backend/internal/repository"
	"sosu-backend/internal/service"
)

type AuthHandler struct {
	userRepo    *repository.UserRepository
	authService *service.AuthService
}

func NewAuthHandler(userRepo *repository.UserRepository, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, authService: authService}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email valid dan password minimal 6 karakter wajib diisi"})
		return
	}

	existing, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email sudah terdaftar"})
		return
	}

	hash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memproses password"})
		return
	}

	user := &model.User{Email: req.Email, PasswordHash: hash}
	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat sesi login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "email": user.Email})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email dan password wajib diisi"})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil || !h.authService.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password salah"})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat sesi login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "email": user.Email})
}