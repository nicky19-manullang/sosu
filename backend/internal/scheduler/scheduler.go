package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"

	"sosu-backend/internal/fetcher"
	"sosu-backend/internal/repository"
	"sosu-backend/internal/service"
	"time"
)

type Scheduler struct {
	cron             *cron.Cron
	firmsFetcher     *fetcher.FIRMSFetcher
	openaqFetcher    *fetcher.OpenAQFetcher
	statusService    *service.StatusService
	subscriptionRepo *repository.SubscriptionRepository
	userRepo         *repository.UserRepository
	emailService     *service.EmailService
	hotspotRepo      *repository.HotspotRepository
}

func New(
	firmsFetcher *fetcher.FIRMSFetcher,
	openaqFetcher *fetcher.OpenAQFetcher,
	statusService *service.StatusService,
	subscriptionRepo *repository.SubscriptionRepository,
	userRepo *repository.UserRepository,
	emailService *service.EmailService,
	hotspotRepo *repository.HotspotRepository,
) *Scheduler {
	return &Scheduler{
		cron:             cron.New(),
		firmsFetcher:     firmsFetcher,
		openaqFetcher:    openaqFetcher,
		statusService:    statusService,
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
		emailService:     emailService,
		hotspotRepo:      hotspotRepo,
	}
}

func (s *Scheduler) checkSubscriptions() {
	subs, err := s.subscriptionRepo.FindAll()
	if err != nil {
		log.Println("[scheduler] failed to load subscriptions:", err)
		return
	}

	for _, sub := range subs {
		result, err := s.statusService.GetStatus(sub.Latitude, sub.Longitude)
		if err != nil {
			log.Println("[scheduler] failed to get status for subscription:", sub.ID, err)
			continue
		}

		currentStatus := string(result.FinalStatus)
		isAlertWorthy := currentStatus == "Waspada" || currentStatus == "Tidak Aman"
		alreadyNotifiedThisStatus := sub.LastNotifiedStatus == currentStatus

		if isAlertWorthy && !alreadyNotifiedThisStatus {
			user, err := s.userRepo.FindByID(sub.UserID)
			if err != nil || user == nil {
				continue
			}

			err = s.emailService.SendAlert(user.Email, sub.LocationLabel, currentStatus, result.AQICategory, result.HotspotCount)
			if err != nil {
				log.Println("[scheduler] failed to send alert email:", err)
				continue
			}
			s.subscriptionRepo.UpdateLastNotified(sub.ID, currentStatus)
			log.Printf("[scheduler] alert sent to %s for status %s\n", user.Email, currentStatus)
		} else if !isAlertWorthy && sub.LastNotifiedStatus != "" {
			s.subscriptionRepo.UpdateLastNotified(sub.ID, "")
		}
	}
}

func (s *Scheduler) Start() {
	_, err := s.cron.AddFunc("*/30 * * * *", func() {
		log.Println("[scheduler] running FIRMS fetch...")
		if err := s.firmsFetcher.Run(); err != nil {
			log.Println("[scheduler] FIRMS fetch failed:", err)
		}
	})
	if err != nil {
		log.Fatal("failed to schedule FIRMS fetch: ", err)
	}

	_, err = s.cron.AddFunc("0 * * * *", func() {
		log.Println("[scheduler] running OpenAQ fetch...")
		if err := s.openaqFetcher.Run(); err != nil {
			log.Println("[scheduler] OpenAQ fetch failed:", err)
		}
	})
	if err != nil {
		log.Fatal("failed to schedule OpenAQ fetch: ", err)
	}

	_, err = s.cron.AddFunc("0 */3 * * *", func() {
		log.Println("[scheduler] running FIRMS global fetch...")
		if err := s.firmsFetcher.RunGlobal(); err != nil {
			log.Println("[scheduler] FIRMS global fetch failed:", err)
		}
	})
	if err != nil {
		log.Fatal("failed to schedule FIRMS global fetch: ", err)
	}

	_, err = s.cron.AddFunc("15,45 * * * *", func() {
		log.Println("[scheduler] checking subscriptions...")
		s.checkSubscriptions()
	})
	if err != nil {
		log.Fatal("failed to schedule subscription check: ", err)
	}
	_, err = s.cron.AddFunc("0 3 * * *", func() {
	log.Println("[scheduler] running daily cleanup...")
	s.cleanupOldData()
	})
	if err != nil {
		log.Fatal("failed to schedule cleanup: ", err)
	}
	_, err = s.cron.AddFunc("0 8 * * 1", func() {
	log.Println("[scheduler] sending weekly digest...")
	s.sendWeeklyDigest()
	})
	if err != nil {
		log.Fatal("failed to schedule weekly digest: ", err)
	}
	s.cron.Start()
	log.Println("[scheduler] started: FIRMS every 30m, OpenAQ hourly, FIRMS global every 3h, subscription check every 30m, daily cleanup at 3am, weekly digest Monday 8am")
}

func (s *Scheduler) RunInitialFetch() {
	log.Println("[scheduler] running initial fetch on startup...")
	if err := s.firmsFetcher.Run(); err != nil {
		log.Println("[scheduler] initial FIRMS fetch failed:", err)
	}
	if err := s.openaqFetcher.Run(); err != nil {
		log.Println("[scheduler] initial OpenAQ fetch failed:", err)
	}
}

func (s *Scheduler) cleanupOldData() {
	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	if err := s.hotspotRepo.DeleteOlderThan(cutoff); err != nil {
		log.Println("[scheduler] cleanup failed:", err)
		return
	}
	log.Println("[scheduler] cleaned up hotspots older than 7 days")
}
func (s *Scheduler) sendWeeklyDigest() {
	grouped, err := s.subscriptionRepo.GroupedByUser()
	if err != nil {
		log.Println("[scheduler] failed to load subscriptions for digest:", err)
		return
	}

	for userID, subs := range grouped {
		user, err := s.userRepo.FindByID(userID)
		if err != nil || user == nil {
			continue
		}

		var summaries []service.DigestLocationSummary
		for _, sub := range subs {
			result, err := s.statusService.GetStatus(sub.Latitude, sub.Longitude)
			if err != nil {
				continue
			}
			summaries = append(summaries, service.DigestLocationSummary{
				LocationLabel: sub.LocationLabel,
				Status:        string(result.FinalStatus),
				AQICategory:   result.AQICategory,
				AQIValue:      result.AQIValue,
				HotspotCount:  result.HotspotCount,
			})
		}

		if len(summaries) == 0 {
			continue
		}

		if err := s.emailService.SendWeeklyDigest(user.Email, summaries); err != nil {
			log.Println("[scheduler] failed to send digest to", user.Email, err)
			continue
		}
		log.Println("[scheduler] weekly digest sent to", user.Email)
	}
}
func (s *Scheduler) TriggerDigestNow() {
	go s.sendWeeklyDigest()
}
func (s *Scheduler) Stop() {
	s.cron.Stop()
}