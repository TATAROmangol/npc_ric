package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBD struct{
	Client *mongo.Client
	DB *mongo.Collection
}

func NewMongoDB(ctx context.Context, cfg Config) (*MongoBD,error){
	ctx, cancel := context.WithTimeout(ctx, 2 * time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Addr()))
	if err != nil {
		return nil, err
	}

	db := client.Database(cfg.DBName).Collection(cfg.CollectionName)

	return &MongoBD{Client:client, DB: db},nil
}

func (db *MongoBD) Close() error{
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return db.Client.Disconnect(ctx)
}