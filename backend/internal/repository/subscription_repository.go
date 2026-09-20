package repository

import (
	"time"

	"gorm.io/gorm"

	"sosu-backend/internal/model"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(s *model.Subscription) error {
	return r.db.Create(s).Error
}

func (r *SubscriptionRepository) FindAll() ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.Find(&subs).Error
	return subs, err
}

func (r *SubscriptionRepository) FindByUserID(userID uint) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.Where("user_id = ?", userID).Find(&subs).Error
	return subs, err
}

func (r *SubscriptionRepository) UpdateLastNotified(id uint, status string) error {
	now := time.Now()
	return r.db.Model(&model.Subscription{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_notified_status": status,
		"last_notified_at":     now,
	}).Error
}

func (r *SubscriptionRepository) DeleteByUserAndLocation(userID uint, lat, lng float64) error {
	return r.db.Where("user_id = ? AND latitude = ? AND longitude = ?", userID, lat, lng).Delete(&model.Subscription{}).Error
}
func (r *SubscriptionRepository) GroupedByUser() (map[uint][]model.Subscription, error) {
	subs, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	grouped := make(map[uint][]model.Subscription)
	for _, sub := range subs {
		grouped[sub.UserID] = append(grouped[sub.UserID], sub)
	}
	return grouped, nil
}