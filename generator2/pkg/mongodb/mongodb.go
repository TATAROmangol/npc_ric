package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBD struct{
	client *mongo.Client
	db *mongo.Collection
}

func NewMongoDB(ctx context.Context, cfg Config) (*MongoBD,error){
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Addr()))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.DBName).Collection(cfg.CollectionName)

	return &MongoBD{client:client, db: db},nil
}

func (db *MongoBD) Close() error{
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return db.client.Disconnect(ctx)
}