package repository

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
)

type CenterRepository struct {
	centers map[string]*model.DonationCenter
	mu      sync.RWMutex
	nextID  int
}

func NewCenterRepository() *CenterRepository {
	return &CenterRepository{
		centers: make(map[string]*model.DonationCenter),
	}
}

func (r *CenterRepository) List() []*model.DonationCenter {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.DonationCenter, 0, len(r.centers))
	for _, c := range r.centers {
		result = append(result, c)
	}
	return result
}

func (r *CenterRepository) ListByDistance(lat, lng, radiusKm float64) []*model.DonationCenter {
	all := r.List()
	var filtered []*model.DonationCenter
	for _, c := range all {
		if haversineKm(lat, lng, c.Lat, c.Lng) <= radiusKm {
			filtered = append(filtered, c)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		di := haversineKm(lat, lng, filtered[i].Lat, filtered[i].Lng)
		dj := haversineKm(lat, lng, filtered[j].Lat, filtered[j].Lng)
		return di < dj
	})
	return filtered
}

func (r *CenterRepository) GetByID(id string) (*model.DonationCenter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.centers[id]
	if !ok {
		return nil, fmt.Errorf("center not found: %s", id)
	}
	return c, nil
}

func (r *CenterRepository) Create(c *model.DonationCenter) (*model.DonationCenter, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	now := time.Now()
	c.ID = fmt.Sprintf("center-%d", r.nextID)
	c.CreatedAt = now
	c.UpdatedAt = now

	if c.Status == "" {
		c.Status = model.CenterStatusOpen
	}

	r.centers[c.ID] = c
	return c, nil
}

func haversineKm(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := degreesToRadians(lat2 - lat1)
	dLng := degreesToRadians(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}
