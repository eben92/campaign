package campaign

import (
	"campaign/internal/database"
	"campaign/internal/models"
	campaignservice "campaign/internal/services/campaign"
	"campaign/internal/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type CampaignHandler interface {
	CreateCampaignHandler(w http.ResponseWriter, r *http.Request)
	GetCampaignsHandler(w http.ResponseWriter, r *http.Request)
	GetCampaignByIDHandler(w http.ResponseWriter, r *http.Request)
	UpdateCampaignHandler(w http.ResponseWriter, r *http.Request)
	DeleteCampaignHandler(w http.ResponseWriter, r *http.Request)
}

type campaignHandler struct {
	db *mongo.Database
}

func NewCampaignHandler(db *mongo.Database) CampaignHandler {
	return &campaignHandler{db: db}
}

func validateCampaign(c models.Campaign) error {
	if c.Name == "" {
		return errors.New("name is required")
	}

	if c.Description == "" {
		return errors.New("description is required")
	}

	if c.StartDate.IsZero() {
		return errors.New("start date is required")
	}

	if c.EndDate.IsZero() {
		return errors.New("end date is required")
	}

	if c.EndDate.Before(c.StartDate) {

		return errors.New("end date must be after start date")
	}

	if c.BannerURL == "" {
		return errors.New("banner url is required")
	}

	return nil
}

func (c *campaignHandler) CreateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := models.Campaign{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()

	if err != nil {
		res := utils.WrapInResponse("error decoding request body", nil)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)

		return
	}

	if err := validateCampaign(reqBody); err != nil {
		res := utils.WrapInResponse(err.Error(), nil)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)

		return
	}

	dbM := database.NewDatabaseService(r.Context(), c.db, models.CampaignsCollection)

	campaignService := campaignservice.NewService(r.Context(), dbM)

	err = campaignService.CreateCampaign(reqBody)

	if err != nil {
		res := utils.WrapInResponse(err.Error(), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(res)

		return
	}

	res := utils.WrapInResponse("campaign created successfully", nil)
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(res)

}

func (c *campaignHandler) GetCampaignsHandler(w http.ResponseWriter, r *http.Request) {
	dbM := database.NewDatabaseService(r.Context(), c.db, models.CampaignsCollection)

	campaignService := campaignservice.NewService(r.Context(), dbM)

	campaigns, err := campaignService.GetCampaigns()

	if err != nil {
		res := utils.WrapInResponse(err.Error(), nil)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(res)
		return
	}

	res := utils.WrapInResponse("campaigns retrieved successfully", campaigns)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)

}

func (c *campaignHandler) GetCampaignByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	dbM := database.NewDatabaseService(r.Context(), c.db, models.CampaignsCollection)
	campaignService := campaignservice.NewService(r.Context(), dbM)

	campaign, err := campaignService.GetCampaignByID(id)

	if err != nil {

		res := utils.WrapInResponse(err.Error(), nil)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write(res)
		return
	}

	res := utils.WrapInResponse("campaign retrieved successfully", campaign)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)

}

func (c *campaignHandler) UpdateCampaignHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	reqBody := models.Campaign{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	defer r.Body.Close()

	if err != nil {
		res := utils.WrapInResponse("error decoding request body", nil)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)

		return
	}

	if err := validateCampaign(reqBody); err != nil {
		res := utils.WrapInResponse(err.Error(), nil)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(res)

		return
	}

	dbM := database.NewDatabaseService(r.Context(), c.db, models.CampaignsCollection)

	campaignService := campaignservice.NewService(r.Context(), dbM)

	err = campaignService.UpdateCampaign(id, reqBody)

	if err != nil {
		res := utils.WrapInResponse(err.Error(), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(res)
		return
	}

	res := utils.WrapInResponse("campaign updated successfully", nil)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)

}

func (c *campaignHandler) DeleteCampaignHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	dbM := database.NewDatabaseService(r.Context(), c.db, models.CampaignsCollection)

	campaignService := campaignservice.NewService(r.Context(), dbM)

	err := campaignService.DeleteCampaign(id)

	if err != nil {
		res := utils.WrapInResponse(err.Error(), nil)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(res)
		return
	}

	res := utils.WrapInResponse("campaign deleted successfully", nil)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)

}
