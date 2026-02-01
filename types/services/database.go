package services

import (
	"fmt"
	"sort"
	"time"

	"github.com/statping-next/statping-next/database"
	"github.com/statping-next/statping-next/types/errors"
	"github.com/statping-next/statping-next/types/metrics"
	"github.com/statping-next/statping-next/utils"
)

var (
	db          database.Database
	log         = utils.Log.WithField("type", "service")
	allServices map[int64]*Service
)

func (s *Service) Validate() error {
	if s.Name == "" {
		return errors.New("missing service name")
	} else if s.Domain == "" && s.Type != "static" {
		return errors.New("missing domain name")
	} else if s.Type == "" {
		return errors.New("missing service type")
	} else if s.Interval == 0 && s.Type != "static" {
		return errors.New("missing check interval")
	}
	return nil
}

func (s *Service) BeforeCreate() error {
	return s.Validate()
}

func (s *Service) BeforeUpdate() error {
	return s.Validate()
}

func (s *Service) AfterFind() {
	// Removed automatic loading of related data to improve performance
	// Related data (Incidents, Messages, Checkins) should be loaded explicitly when needed
	// This prevents N+1 query problems and significantly improves main page load times
	metrics.Query("service", "find")
}

func (s *Service) AfterCreate() error {
	s.prevOnline = true
	allServices[s.Id] = s
	metrics.Query("service", "create")
	return nil
}

func (s *Service) AfterUpdate() {
	metrics.Query("service", "update")
}

func (s *Service) AfterDelete() {
	metrics.Query("service", "delete")
}

func init() {
	allServices = make(map[int64]*Service)
}

func Services() map[int64]*Service {
	return allServices
}

func SetDB(database database.Database) {
	db = database.Model(&Service{})
}

func Find(id int64) (*Service, error) {
	srv := allServices[id]
	if srv == nil {
		return nil, errors.Missing(&Service{}, id)
	}
	db.First(&srv, id)
	return srv, nil
}

func all() []*Service {
	var services []*Service
	db.Find(&services)
	return services
}

func All() map[int64]*Service {
	return allServices
}

func AllInOrder() []Service {
	var services []Service
	for _, service := range allServices {
		service.UpdateStats()
		services = append(services, *service)
	}
	sort.Sort(ServiceOrder(services))
	return services
}

// AllInOrderLite returns services in order without calculating expensive stats
// This is optimized for the main status page which only needs up/down status
// For offline services, it loads the downtime start time (first failure after last success)
func AllInOrderLite() []Service {
	var services []Service
	for _, service := range allServices {
		// For offline services, we need to know when downtime started
		// If DowntimeStarted is not set (e.g., service was offline on startup),
		// find the first failure after the last success
		if !service.Online && service.DowntimeStarted.IsZero() {
			// Get the last success time
			lastHit := service.LastHit()
			lastCheckinHit := service.AllCheckinHits().First()
			var lastSuccessTime time.Time
			if lastHit != nil {
				lastSuccessTime = lastHit.CreatedAt
			}
			if lastCheckinHit != nil && lastCheckinHit.CreatedAt.After(lastSuccessTime) {
				lastSuccessTime = lastCheckinHit.CreatedAt
			}
			
			// Find failures after the last success and get the earliest one
			if !lastSuccessTime.IsZero() {
				failuresSince := service.FailuresSince(lastSuccessTime).List()
				if len(failuresSince) > 0 {
					// Find the earliest failure by time (not ID)
					earliest := failuresSince[0]
					for _, f := range failuresSince {
						if f.CreatedAt.Before(earliest.CreatedAt) {
							earliest = f
						}
					}
					service.DowntimeStarted = earliest.CreatedAt
				} else {
					// If no failure found after last success, use the most recent failure
					lastFailure := service.AllFailures().Last()
					if lastFailure != nil {
						service.DowntimeStarted = lastFailure.CreatedAt
					}
				}
			} else {
				// No success found, use the most recent failure as fallback
				lastFailure := service.AllFailures().Last()
				if lastFailure != nil {
					service.DowntimeStarted = lastFailure.CreatedAt
				}
			}
		}
		// Skip UpdateStats() - only load basic service data
		// This avoids expensive queries for hits/failures counts, uptime percentages, etc.
		services = append(services, *service)
	}
	sort.Sort(ServiceOrder(services))
	return services
}

func (s *Service) Create() error {
	err := db.Create(s)
	if err.Error() != nil {
		log.Errorln(fmt.Sprintf("Failed to create service %v #%v: %v", s.Name, s.Id, err))
		return err.Error()
	}
	return nil
}

func (s *Service) Update() error {
	q := db.Update(s)
	s.Close()
	allServices[s.Id] = s
	s.SleepDuration = s.Duration()
	go ServiceCheckQueue(allServices[s.Id], true)
	return q.Error()
}

func (s *Service) Delete() error {
	s.Close()
	if err := s.AllFailures().DeleteAll(); err != nil {
		return err
	}
	if err := s.AllHits().DeleteAll(); err != nil {
		return err
	}
	if err := s.DeleteCheckins(); err != nil {
		return err
	}
	db.Model(s).Association("Checkins").Clear()
	if err := s.DeleteIncidents(); err != nil {
		return err
	}
	db.Model(s).Association("Incidents").Clear()
	if err := s.DeleteMessages(); err != nil {
		return err
	}
	db.Model(s).Association("Messages").Clear()

	delete(allServices, s.Id)
	q := db.Model(&Service{}).Delete(s)
	return q.Error()
}

func (s *Service) DeleteMessages() error {
	for _, m := range s.Messages {
		if err := m.Delete(); err != nil {
			return err
		}
	}
	db.Model(s).Association("messages").Clear()
	return nil
}

func (s *Service) DeleteCheckins() error {
	for _, c := range s.Checkins {
		if err := c.Delete(); err != nil {
			return err
		}
	}
	db.Model(s).Association("checkins").Clear()
	return nil
}
