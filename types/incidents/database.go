package incidents

import (
	"strings"
	"time"

	"github.com/statping-next/statping-next/database"
	"github.com/statping-next/statping-next/types/errors"
	"github.com/statping-next/statping-next/types/metrics"
	"github.com/statping-next/statping-next/utils"
)

var (
	db       database.Database
	dbUpdate database.Database
	log      = utils.Log.WithField("type", "service")
)

func SetDB(database database.Database) {
	db = database.Model(&Incident{})
	dbUpdate = database.Model(&IncidentUpdate{})
}

func (i *Incident) Validate() error {
	if i.Title == "" {
		return errors.New("missing title")
	}
	return nil
}

func (i *Incident) BeforeUpdate() error {
	return i.Validate()
}

func (i *Incident) BeforeCreate() error {
	return i.Validate()
}

func (i *Incident) AfterFind() {
	loadUpdatesAsc(i)
	metrics.Query("incident", "find")
}

// loadUpdatesAsc loads i.Updates ordered by id ASC (oldest-first) for display in public and dashboard.
func loadUpdatesAsc(i *Incident) {
	dbUpdate.Where("incident = ?", i.Id).Order("id ASC").Find(&i.Updates)
}

func (i *Incident) AfterCreate() {
	metrics.Query("incident", "create")
}

func (i *Incident) AfterUpdate() {
	metrics.Query("incident", "update")
}

func (i *Incident) AfterDelete() {
	metrics.Query("incident", "delete")
}

func (i *IncidentUpdate) Validate() error {
	if i.Message == "" {
		return errors.New("missing incident update title")
	}
	return nil
}

func (i *IncidentUpdate) BeforeUpdate() error {
	return i.Validate()
}

func (i *IncidentUpdate) BeforeCreate() error {
	return i.Validate()
}

func (i *IncidentUpdate) AfterFind() {
	metrics.Query("incident_update", "find")
}

func (i *IncidentUpdate) AfterCreate() {
	metrics.Query("incident_update", "create")
}

func (i *IncidentUpdate) AfterUpdate() {
	metrics.Query("incident_update", "update")
}

func (i *IncidentUpdate) AfterDelete() {
	metrics.Query("incident_update", "delete")
}

func FindUpdate(uid int64) (*IncidentUpdate, error) {
	var update IncidentUpdate
	q := dbUpdate.Where("id = ?", uid).Find(&update)
	return &update, q.Error()
}

func Find(id int64) (*Incident, error) {
	var incident Incident
	q := db.Where("id = ?", id).Find(&incident)
	return &incident, q.Error()
}

func FindByService(id int64) []*Incident {
	var incidents []*Incident
	db.Where("service = ?", id).Find(&incidents)
	return incidents
}

func All() []*Incident {
	var incidents []*Incident
	db.Order("created_at DESC").Find(&incidents)
	return incidents
}

// AllWithUpdates returns all incidents (including archived) with their updates loaded (for admin).
func AllWithUpdates() []*Incident {
	incidents := All()
	for _, i := range incidents {
		loadUpdatesAsc(i)
	}
	return incidents
}

// ProcessAutoArchive marks incidents as archived when auto_archive_enabled, last update is "resolved", and the delay has passed.
// Delay is from the last update's CreatedAt: archive when (lastUpdate.CreatedAt + delay) <= now.
// A delay of 0 means immediate: archive as soon as the last update is "resolved".
func ProcessAutoArchive() {
	var list []*Incident
	db.Where("auto_archive_enabled = ? AND archived = ?", true, false).Find(&list)
	now := time.Now()
	for _, i := range list {
		loadUpdatesAsc(i)
		if len(i.Updates) == 0 {
			continue
		}
		last := i.Updates[len(i.Updates)-1] // latest update (slice is id ASC, last = newest)
		if !strings.EqualFold(last.Type, "resolved") {
			continue
		}
		// delay 0 = immediate; otherwise wait N minutes after the resolved update
		delay := time.Duration(i.AutoArchiveDelayMinutes) * time.Minute
		if i.AutoArchiveDelayMinutes <= 0 {
			delay = 0
		}
		archiveAt := last.CreatedAt.Add(delay)
		if archiveAt.After(now) {
			continue // not yet time to archive
		}
		i.Archived = true
		_ = i.Update()
	}
}

// AllPublic returns non-archived incidents with updates (for public). Runs ProcessAutoArchive first.
func AllPublic() []*Incident {
	ProcessAutoArchive()
	var incidents []*Incident
	db.Where("archived = ?", false).Find(&incidents)
	for _, i := range incidents {
		loadUpdatesAsc(i)
	}
	return incidents
}

// AllArchivedWithUpdates returns archived incidents with updates (for public incident history). Newest first by creation date.
func AllArchivedWithUpdates() []*Incident {
	ProcessAutoArchive()
	var list []*Incident
	db.Where("archived = ?", true).Order("created_at DESC").Find(&list)
	for _, i := range list {
		loadUpdatesAsc(i)
	}
	return list
}

// FindByServicePublic returns non-archived incidents for a service (for public). Runs ProcessAutoArchive first.
func FindByServicePublic(serviceId int64) []*Incident {
	ProcessAutoArchive()
	var incidents []*Incident
	db.Where("service = ? AND archived = ?", serviceId, false).Find(&incidents)
	for _, i := range incidents {
		loadUpdatesAsc(i)
	}
	return incidents
}

func (i *Incident) Create() error {
	return db.Create(i).Error()
}

// Update persists the incident. Uses a map so zero values (0, false) are written to the DB.
func (i *Incident) Update() error {
	now := time.Now()
	i.UpdatedAt = now
	return db.Model(i).Updates(map[string]interface{}{
		"title":                      i.Title,
		"description":                i.Description,
		"archived":                   i.Archived,
		"auto_archive_enabled":       i.AutoArchiveEnabled,
		"auto_archive_delay_minutes": i.AutoArchiveDelayMinutes,
		"updated_at":                 now,
	}).Error()
}

func (i *Incident) Delete() error {
	for _, u := range i.Updates {
		if err := u.Delete(); err != nil {
			return err
		}
	}
	return db.Delete(i).Error()
}
