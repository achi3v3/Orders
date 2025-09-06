package subs

import (
	"context"
	"orders/internal/models"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	repo     Repository
	cache    map[string]cacheEntry
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
	logger   *logrus.Logger
}

type cacheEntry struct {
	order     models.OrderJson
	expiresAt time.Time
}

func NewService(repo *Repository, logger *logrus.Logger) *Service {
	service := &Service{
		repo:     *repo,
		cache:    make(map[string]cacheEntry),
		cacheTTL: 2 * time.Minute,
		logger:   logger,
	}

	go service.startCacheCleaner()

	return service
}
func (s *Service) Create(ctx context.Context, orderJson *models.OrderJson) error {
	err := s.repo.Create(ctx, orderJson)
	if err != nil {
		return err
	}

	s.invalidateCache(orderJson.OrderUID)
	return nil
}
func (s *Service) GetOrder(ctx context.Context, orderUID string) (*models.OrderJson, error) {
	if order, found := s.getFromCache(orderUID); found {
		s.logger.Info("Service.GetOrder: Get Order From CACHE")
		return order, nil
	}

	order, err := s.repo.GetOrder(ctx, orderUID)
	s.logger.Info("Service.GetOrder: Get Order From DB")
	if err != nil {
		s.logger.Errorf("Service.GetOrder: %v", err)
		return nil, err
	}
	if order == nil {
		return nil, nil
	}
	s.setToCache(orderUID, order)
	s.logger.Info("Service.GetOrder: Get Order From DB and cached")
	return order, nil
}

func (s *Service) WarmUpCache(ctx context.Context) error {
	orders, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Warnf("Service.WarmUpCache: %v", err)
		return err
	}
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	now := time.Now()
	for _, order := range orders {
		s.cache[order.OrderUID] = cacheEntry{
			order:     order,
			expiresAt: now.Add(s.cacheTTL),
		}
	}
	s.logger.Infof("Cache warmed up with %d orders", len(orders))
	return nil
}

func (s *Service) startCacheCleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanExpiredCache()
	}
}

func (s *Service) cleanExpiredCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	now := time.Now()
	for order_uid, entry := range s.cache {
		if now.After(entry.expiresAt) {
			delete(s.cache, order_uid)
			s.logger.Infof("Cache expired: %s", order_uid)
		}
	}
}

func (s *Service) invalidateCache(orderUID string) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	delete(s.cache, orderUID)
	s.logger.Infof("Cache invalidated for order: %s", orderUID)
}

func (s *Service) setToCache(orderUID string, order *models.OrderJson) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache[orderUID] = cacheEntry{
		order:     *order,
		expiresAt: time.Now().Add(s.cacheTTL),
	}
}
func (s *Service) getFromCache(orderUID string) (*models.OrderJson, bool) {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()

	entry, found := s.cache[orderUID]
	if !found || time.Now().After(entry.expiresAt) {
		return nil, false
	}
	return &entry.order, true
}
