package storage

import (
	"context"
	"generator/pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct{
	db *mongo.Collection
}

func NewStorage(db *mongo.Collection) *Storage{
	return &Storage{db: db}
}

func (s *Storage) DeleteTemplate(ctx context.Context, id int) error{
	filter := bson.M{"id": id}
	_, err := s.db.DeleteOne(ctx, filter)
	if err != nil{
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to delete template", err)
		return err
	}
	return nil
}
func (s *Storage) UploadTemplate(ctx context.Context, id int, data []byte) error{
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"data": data}}
	opts := options.Update().SetUpsert(true)

	_, err := s.db.UpdateOne(ctx, filter, update, opts)
	if err != nil{
		logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to update template", err)
		return err
	}
  
	return nil
}
func (s *Storage) GetTemplate(ctx context.Context, id int) ([]byte ,error){
	var result struct {
        Data []byte `bson:"data"`
    }

    err := s.db.FindOne(ctx, bson.M{"id": id}).Decode(&result)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            logger.GetFromCtx(ctx).ErrorContext(ctx, "template not found", err)
            return nil, err
        }
        logger.GetFromCtx(ctx).ErrorContext(ctx, "failed to find template", err)
        return nil, err
    }

    return result.Data, nil
}