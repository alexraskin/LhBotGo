package database

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Guess struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	LhGuess   string        `bson:"lhguess"`
	GuessedBy string        `bson:"guessedBy"`
	GuessedAt time.Time     `bson:"guessedAt"`
}

type Guesses struct {
	Guesses []Guess `bson:"guesses"`
}

type MongoClient interface {
	GetGuesses(ctx context.Context) ([]Guess, error)
	AddGuess(ctx context.Context, guess Guess) error
	CountGuesses(ctx context.Context) (int64, error)
	GetGuess(ctx context.Context, guess string) (Guess, error)
	GetLatestGuesses(ctx context.Context, limit int) ([]Guess, error)
	Disconnect(ctx context.Context) error
}

type client struct {
	mongo *mongo.Client
}

func New(ctx context.Context, uri string) (MongoClient, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		slog.Error("failed to connect to mongo", "error", err)
		return nil, err
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		slog.Error("failed to ping mongo", "error", err)
		_ = mongoClient.Disconnect(ctx)
		return nil, err
	}

	return &client{mongo: mongoClient}, nil
}

func (c *client) GetGuesses(ctx context.Context) ([]Guess, error) {
	collection := c.mongo.Database("lhbot").Collection("lhbot_collection")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var guesses []Guess
	if err := cursor.All(ctx, &guesses); err != nil {
		return nil, err
	}
	return guesses, nil
}

func (c *client) AddGuess(ctx context.Context, guess Guess) error {
	collection := c.mongo.Database("lhbot").Collection("lhbot_collection")
	_, err := collection.InsertOne(ctx, guess)
	return err
}

func (c *client) CountGuesses(ctx context.Context) (int64, error) {
	collection := c.mongo.Database("lhbot").Collection("lhbot_collection")
	return collection.CountDocuments(ctx, bson.M{})
}

func (c *client) GetGuess(ctx context.Context, guess string) (Guess, error) {
	collection := c.mongo.Database("lhbot").Collection("lhbot_collection")
	var g Guess
	err := collection.FindOne(ctx, bson.M{"guess": guess}).Decode(&g)
	return g, err
}

func (c *client) GetLatestGuesses(ctx context.Context, limit int) ([]Guess, error) {
	collection := c.mongo.Database("lhbot").Collection("lhbot_collection")
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"_id": -1}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}

	var guesses []Guess
	if err := cursor.All(ctx, &guesses); err != nil {
		return nil, err
	}
	return guesses, nil
}

func (c *client) Disconnect(ctx context.Context) error {
	if c.mongo == nil {
		return nil
	}
	slog.Info("Disconnecting from MongoDB")
	return c.mongo.Disconnect(ctx)
}
