package highlight

import (
	"context"

	t "github.com/sikozonpc/notebase/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "notebase"
	CollName = "highlights"
)

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserHighlights(ctx context.Context, userID primitive.ObjectID) ([]*t.Highlight, error) {
	col := s.db.Database(DbName).Collection(CollName)

	cursor, err := col.Find(ctx, bson.M{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}

	var highlights []*t.Highlight
	if err = cursor.All(ctx, &highlights); err != nil {
		return nil, err
	}

	return highlights, nil
}

func (s *Store) CreateHighlight(ctx context.Context, h *t.CreateHighlightRequest) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)

	newHighlight, err := col.InsertOne(ctx, h)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id := newHighlight.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (s *Store) GetHighlightByID(ctx context.Context, id primitive.ObjectID, userID primitive.ObjectID) (*t.Highlight, error) {
	col := s.db.Database(DbName).Collection(CollName)

	var h t.Highlight
	err := col.FindOne(ctx, bson.M{
		"_id":    id,
		"userId": userID,
	}).Decode(&h)

	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (s *Store) DeleteHighlight(ctx context.Context, id primitive.ObjectID) error {
	col := s.db.Database(DbName).Collection(CollName)

	_, err := col.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetRandomHighlights(ctx context.Context, userID primitive.ObjectID, limit int) ([]*t.Highlight, error) {
	col := s.db.Database(DbName).Collection(CollName)

	cursor, err := col.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{Key: `$match`, Value: bson.M{
				"userId": userID,
			}},
		},
		bson.D{
			{Key: "$sample", Value: bson.M{
				"size": limit,
			}},
		},
	})
	if err != nil {
		return nil, err
	}

	var highlights []*t.Highlight
	if err = cursor.All(ctx, &highlights); err != nil {
		return nil, err
	}

	return highlights, nil
}
