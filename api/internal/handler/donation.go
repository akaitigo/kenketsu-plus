package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akaitigo/kenketsu-plus/api/internal/model"
	"github.com/akaitigo/kenketsu-plus/api/internal/repository"
	"github.com/akaitigo/kenketsu-plus/api/internal/service"
)

type DonationHandler struct {
	repo       repository.DonationRepo
	calculator *service.DonationCalculator
}

func NewDonationHandler(repo repository.DonationRepo, calculator *service.DonationCalculator) *DonationHandler {
	return &DonationHandler{repo: repo, calculator: calculator}
}

func (h *DonationHandler) List(w http.ResponseWriter, _ *http.Request) {
	donations := h.repo.List()
	writeJSON(w, http.StatusOK, donations)
}

func (h *DonationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var donation model.Donation
	if err := json.NewDecoder(r.Body).Decode(&donation); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.repo.Create(&donation)
	if err != nil {
		writeRepoError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *DonationHandler) NextAvailable(w http.ResponseWriter, r *http.Request) {
	genderStr := r.URL.Query().Get("gender")
	if genderStr == "" {
		writeError(w, http.StatusBadRequest, "gender query parameter is required")
		return
	}

	gender := model.Gender(genderStr)
	switch gender {
	case model.GenderMale, model.GenderFemale:
	default:
		writeError(w, http.StatusBadRequest, "gender must be male or female")
		return
	}

	donations := h.repo.List()
	result := h.calculator.NextAvailableDate(donations, gender)
	writeJSON(w, http.StatusOK, result)
}
