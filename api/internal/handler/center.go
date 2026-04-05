package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
)

type CenterHandler struct {
	repo repository.CenterRepo
}

func NewCenterHandler(repo repository.CenterRepo) *CenterHandler {
	return &CenterHandler{repo: repo}
}

func (h *CenterHandler) List(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		centers := h.repo.List()
		writeJSON(w, http.StatusOK, centers)
		return
	}

	h.listByDistance(w, r, latStr, lngStr)
}

func (h *CenterHandler) listByDistance(w http.ResponseWriter, r *http.Request, latStr, lngStr string) {
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid lat parameter")
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid lng parameter")
		return
	}

	radius := 10.0
	if radiusStr := r.URL.Query().Get("radius"); radiusStr != "" {
		radius, err = strconv.ParseFloat(radiusStr, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid radius parameter")
			return
		}
	}

	if radius <= 0 || radius > 500 {
		writeError(w, http.StatusBadRequest, "radius must be greater than 0 and at most 500 km")
		return
	}

	centers := h.repo.ListByDistance(lat, lng, radius)
	writeJSON(w, http.StatusOK, centers)
}

func (h *CenterHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	center, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, center)
}

func (h *CenterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var center model.DonationCenter
	if err := json.NewDecoder(r.Body).Decode(&center); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.repo.Create(&center)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
