package database

import (
	"campaign/internal/models"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	Database() *mongo.Database
}

type service struct {
	db     *mongo.Database
	client *mongo.Client
}

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	password = os.Getenv("DB_ROOT_PASSWORD")
	user     = os.Getenv("DB_USERNAME")
)

func New() Service {
	slog.Info("Connecting to database")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)))

	if err != nil {
		slog.Error("Error connecting to database", "Error", err)
		log.Fatal(err)

	}

	slog.Info("Connected to database successfully", "Host", host, "Port", port)

	db := client.Database("campaign")

	setupDBUniqueIndex(db)

	return &service{
		db:     db,
		client: client,
	}
}

func setupDBUniqueIndex(db *mongo.Database) {
	indexes := []mongo.IndexModel{{
		Keys: bson.M{
			"email": 1,
		},
		Options: options.Index().SetUnique(true).SetName("email"),
	},
	}

	_, err := db.Collection(string(models.UsersCollection)).Indexes().CreateMany(context.Background(), indexes)

	if err != nil {
		slog.Error("Error creating index: ", "error", err)
	}

}

func (s *service) Health() map[string]string {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) Database() *mongo.Database {
	return s.db
}

type Database interface {
	SetCollection(collection models.Collections)
	InsertOne(document bson.M) error
	InsertMany(documents []interface{}) error
	FindOne(filter bson.M, result interface{}) error
	FindMany(filter bson.M, result interface{}) error
	AggregateMany(pipeline []bson.M, result interface{}) error
	UpdateOne(filter bson.M, update bson.M) error
	DeleteOne(filter bson.M) error
}

type databaseService struct {
	ctx        context.Context
	db         *mongo.Database
	collection models.Collections
}

func NewDatabaseService(ctx context.Context, client *mongo.Database, collection models.Collections) Database {
	return &databaseService{
		ctx:        ctx,
		db:         client,
		collection: collection,
	}
}

func (s *databaseService) SetCollection(collection models.Collections) {
	s.collection = collection
}

func (s *databaseService) InsertOne(document bson.M) error {
	c := s.db.Collection(string(s.collection))
	_, err := c.InsertOne(s.ctx, document)

	return err
}

func (s *databaseService) InsertMany(documents []interface{}) error {
	c := s.db.Collection(string(s.collection))
	_, err := c.InsertMany(s.ctx, documents)
	return err
}

func (s *databaseService) FindOne(filter bson.M, result interface{}) error {
	c := s.db.Collection(string(s.collection))
	err := c.FindOne(s.ctx, filter).Decode(result)

	return err
}

func (s *databaseService) FindMany(filter bson.M, result interface{}) error {
	c := s.db.Collection(string(s.collection))
	w, err := c.Find(s.ctx, filter)

	if err != nil {
		return err
	}

	err = w.All(s.ctx, result)

	return err
}

func (s *databaseService) UpdateOne(filter bson.M, update bson.M) error {
	c := s.db.Collection(string(s.collection))
	_, err := c.UpdateOne(s.ctx, filter, update)

	return err
}

func (s *databaseService) DeleteOne(filter bson.M) error {
	c := s.db.Collection(string(s.collection))
	_, err := c.DeleteOne(s.ctx, filter)

	return err
}

func (s *databaseService) AggregateMany(pipeline []bson.M, result interface{}) error {
	c := s.db.Collection(string(s.collection))
	cursor, err := c.Aggregate(s.ctx, pipeline)

	if err != nil {
		return err
	}

	err = cursor.All(s.ctx, result)

	return err
}
