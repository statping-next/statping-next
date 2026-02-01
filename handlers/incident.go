package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/statping-next/statping-next/types/errors"
	"github.com/statping-next/statping-next/types/incidents"
	"github.com/statping-next/statping-next/utils"
)

func findIncident(r *http.Request) (*incidents.Incident, int64, error) {
	vars := mux.Vars(r)
	if utils.NotNumber(vars["id"]) {
		return nil, 0, errors.NotNumber
	}
	id := utils.ToInt(vars["id"])
	if id == 0 {
		return nil, id, errors.IDMissing
	}
	checkin, err := incidents.Find(id)
	if err != nil {
		return nil, id, errors.Missing(&incidents.Incident{}, id)
	}
	return checkin, id, nil
}

func apiAllIncidentsHandler(w http.ResponseWriter, r *http.Request) {
	all := incidents.AllWithUpdates()
	returnJson(all, w, r)
}

// apiAllIncidentsHandlerScoped returns all incidents (including archived) for admin; non-archived for public.
// ProcessAutoArchive runs for every request so delay=0 archives run whether the user is public or authenticated.
func apiAllIncidentsHandlerScoped(r *http.Request) interface{} {
	incidents.ProcessAutoArchive()
	if IsFullAuthenticated(r) {
		return incidents.AllWithUpdates()
	}
	return incidents.AllPublic()
}

// apiArchivedIncidentsHandler returns archived incidents with updates (public incident history page).
func apiArchivedIncidentsHandler(r *http.Request) interface{} {
	return incidents.AllArchivedWithUpdates()
}

func apiServiceIncidentsHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	serviceIncidents := incidents.FindByServicePublic(service.Id)
	returnJson(serviceIncidents, w, r)
}

func apiIncidentUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	incid, _, err := findIncident(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	returnJson(incid.Updates, w, r)
}

func apiCreateIncidentUpdateHandler(w http.ResponseWriter, r *http.Request) {
	incid, _, err := findIncident(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}

	var update *incidents.IncidentUpdate
	if err := DecodeJSON(r, &update); err != nil {
		sendErrorJson(err, w, r)
		return
	}

	update.IncidentId = incid.Id

	if err := update.Create(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	// Run auto-archive so delay=0 (immediate) archives right after a Resolved update
	incidents.ProcessAutoArchive()
	sendJsonAction(update, "create", w, r)
}

func apiCreateIncidentHandler(w http.ResponseWriter, r *http.Request) {
	service, err := findService(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	var incident *incidents.Incident
	if err := DecodeJSON(r, &incident); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	incident.ServiceId = service.Id
	if err := incident.Create(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(incident, "create", w, r)
}

// apiCreateGlobalIncidentHandler creates an incident with service=0 (global).
func apiCreateGlobalIncidentHandler(w http.ResponseWriter, r *http.Request) {
	var incident *incidents.Incident
	if err := DecodeJSON(r, &incident); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	incident.ServiceId = 0
	if err := incident.Create(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(incident, "create", w, r)
}

func apiIncidentUpdateHandler(w http.ResponseWriter, r *http.Request) {
	incident, _, err := findIncident(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	var body struct {
		Title                   *string `json:"title"`
		Description             *string `json:"description"`
		AutoArchiveEnabled      *bool   `json:"auto_archive_enabled"`
		AutoArchiveDelayMinutes *int   `json:"auto_archive_delay_minutes"`
	}
	if err := DecodeJSON(r, &body); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if body.Title != nil {
		incident.Title = *body.Title
	}
	if body.Description != nil {
		incident.Description = *body.Description
	}
	if body.AutoArchiveEnabled != nil {
		incident.AutoArchiveEnabled = *body.AutoArchiveEnabled
	}
	if body.AutoArchiveDelayMinutes != nil {
		incident.AutoArchiveDelayMinutes = *body.AutoArchiveDelayMinutes
	}
	if err := incident.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(incident, "update", w, r)
}

func apiArchiveIncidentHandler(w http.ResponseWriter, r *http.Request) {
	incident, _, err := findIncident(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	var body struct {
		Archived *bool `json:"archived"`
	}
	if err := DecodeJSON(r, &body); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if body.Archived != nil {
		incident.Archived = *body.Archived
		// When user explicitly unarchives, disable auto-archive so ProcessAutoArchive
		// does not immediately re-archive it. User can re-enable in the edit modal.
		if !*body.Archived {
			incident.AutoArchiveEnabled = false
		}
	}
	if err := incident.Update(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(incident, "update", w, r)
}

func apiDeleteIncidentHandler(w http.ResponseWriter, r *http.Request) {
	incident, _, err := findIncident(r)
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := incident.Delete(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(incident, "delete", w, r)
}

func apiDeleteIncidentUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	update, err := incidents.FindUpdate(utils.ToInt(vars["uid"]))
	if err != nil {
		sendErrorJson(err, w, r)
		return
	}
	if err := update.Delete(); err != nil {
		sendErrorJson(err, w, r)
		return
	}
	sendJsonAction(update, "delete", w, r)
}
