package campaignservice

import (
	"campaign/internal/database"
	"campaign/internal/models"
	"campaign/internal/utils/jwt"
	"context"
	"errors"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type CampaignService interface {
	CreateCampaign(c models.Campaign) error
	GetCampaigns() ([]models.Campaign, error)
	GetCampaignByID(id string) (models.Campaign, error)
	UpdateCampaign(id string, c models.Campaign) error
	DeleteCampaign(id string) error
}

type service struct {
	ctx context.Context
	db  database.Database
}

func NewService(ctx context.Context, db database.Database) CampaignService {
	return &service{ctx: ctx, db: db}
}

func (s *service) CreateCampaign(c models.Campaign) error {

	userID, err := jwt.GetAuthContext(s.ctx)

	if err != nil {
		slog.Error("Error getting auth context", "error", err)

		return errors.New("error creating campaign")
	}

	s.db.SetCollection(models.CampaignsCollection)

	c.CreatedBy = userID.Sub
	err = s.db.InsertOne(bson.M{
		"name":        c.Name,
		"description": c.Description,
		"start_date":  c.StartDate,
		"end_date":    c.EndDate,
		"banner_url":  c.BannerURL,
		"created_by":  c.CreatedBy,
		"created_at":  c.CreatedAt,
		"updated_at":  c.UpdatedAt,
	})

	if err != nil {
		slog.Error("Error creating campaign", "error", err)

		return errors.New("error creating campaign")

	}

	return nil
}

func (s *service) GetCampaigns() ([]models.Campaign, error) {
	var campaigns []models.Campaign

	user, err := jwt.GetAuthContext(s.ctx)

	if err != nil {
		slog.Error("Error getting auth context", "error", err)

		return campaigns, errors.New("error getting campaigns")
	}

	s.db.SetCollection(models.CampaignsCollection)

	err = s.db.FindMany(bson.M{"created_by": user.Sub}, &campaigns)

	if err != nil {
		slog.Error("Error getting campaigns", "error", err)

		return campaigns, errors.New("error getting campaigns")
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(id string) (models.Campaign, error) {
	campaign := models.Campaign{}

	user, err := jwt.GetAuthContext(s.ctx)

	if err != nil {
		slog.Error("Error getting auth context", "error", err)

		return campaign, errors.New("error getting campaign")
	}

	s.db.SetCollection(models.CampaignsCollection)

	err = s.db.FindOne(bson.M{"_id": id, "created_by": user.Sub}, &campaign)

	if err != nil {
		slog.Error("Error getting campaign", "error", err)

		return campaign, errors.New("error getting campaign")
	}

	return campaign, nil
}

func (s *service) UpdateCampaign(id string, c models.Campaign) error {
	user, err := jwt.GetAuthContext(s.ctx)

	if err != nil {
		slog.Error("Error getting auth context", "error", err)

		return errors.New("error updating campaign")
	}

	s.db.SetCollection(models.CampaignsCollection)

	err = s.db.UpdateOne(bson.M{"_id": id, "created_by": user.Sub}, bson.M{
		"name":        c.Name,
		"description": c.Description,
		"start_date":  c.StartDate,
		"end_date":    c.EndDate,
		"banner_url":  c.BannerURL,
		"updated_at":  time.Now(),
	})

	if err != nil {
		slog.Error("Error updating campaign", "error", err)

		return errors.New("error updating campaign")
	}

	return nil
}

func (s *service) DeleteCampaign(id string) error {
	user, err := jwt.GetAuthContext(s.ctx)

	if err != nil {
		slog.Error("Error getting auth context", "error", err)

		return errors.New("error deleting campaign")
	}

	s.db.SetCollection(models.CampaignsCollection)

	err = s.db.DeleteOne(bson.M{"_id": id, "created_by": user.Sub})

	if err != nil {
		slog.Error("Error deleting campaign", "error", err)

		return errors.New("error deleting campaign")
	}

	return nil
}
