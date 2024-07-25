package campaignservice

import (
	"campaign/internal/database"
	"campaign/internal/models"
	"campaign/internal/utils/jwt"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	err = s.db.InsertOne(bson.M{
		"name":        c.Name,
		"description": c.Description,
		"start_date":  c.StartDate.Local(),
		"end_date":    c.EndDate.Local(),
		"banner_url":  c.BannerURL,
		"created_by":  userID.Sub,
		"created_at":  time.Now().Local(),
		"updated_at":  time.Now().Local(),
	})

	if err != nil {
		slog.Error("Error creating campaign", "error", err)

		return errors.New("error creating campaign")

	}

	return nil
}

func (s *service) GetCampaigns() ([]models.Campaign, error) {
	campaigns := []models.Campaign{}

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

	objid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		slog.Error("Error converting id to object id", "error", err)

		return campaign, fmt.Errorf("invalid campaign id: %s", id)

	}

	err = s.db.FindOne(bson.M{
		"_id":        objid,
		"created_by": user.Sub}, &campaign)

	if err != nil {
		slog.Error("Error getting campaign", "error", err)

		return campaign, fmt.Errorf("no campaigns with id: %s found", id)
	}

	return campaign, nil
}

func (s *service) UpdateCampaign(id string, c models.Campaign) error {
	user, err := jwt.GetAuthContext(s.ctx)

	if err != nil {
		slog.Error("Error getting auth context", "error", err)

		return errors.New("error updating campaign")
	}

	objid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		slog.Error("Error converting id to object id", "error", err)

		return fmt.Errorf("invalid campaign id: %s", id)

	}

	s.db.SetCollection(models.CampaignsCollection)

	err = s.db.UpdateOne(bson.M{"_id": objid, "created_by": user.Sub}, bson.M{
		"name":        c.Name,
		"description": c.Description,
		"start_date":  c.StartDate.Unix(),
		"end_date":    c.EndDate.Unix(),
		"banner_url":  c.BannerURL,
		"updated_at":  time.Now(),
	})

	if err != nil {
		slog.Error("Error updating campaign", "error", err)

		return fmt.Errorf("could not update campaign with id: %s", id)
	}

	return nil
}

func (s *service) DeleteCampaign(id string) error {
	user, err := jwt.GetAuthContext(s.ctx)

	if err != nil {
		slog.Error("Error getting auth context", "error", err)

		return errors.New("error deleting campaign")
	}

	objid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		slog.Error("Error converting id to object id", "error", err)

		return fmt.Errorf("invalid campaign id: %s", id)

	}

	s.db.SetCollection(models.CampaignsCollection)

	err = s.db.DeleteOne(bson.M{"_id": objid, "created_by": user.Sub})

	if err != nil {
		slog.Error("Error deleting campaign", "error", err)

		return fmt.Errorf("could not delete campaign with id: %s", id)
	}

	return nil
}
