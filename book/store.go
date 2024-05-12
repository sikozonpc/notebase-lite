package book

import (
	"context"

	t "github.com/sikozonpc/notebase/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "notebase"
	CollName = "books"
)

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {
	return &Store{db: db}
}

func (s *Store) GetByISBN(ctx context.Context, isbn string) (*t.Book, error) {
	col := s.db.Database(DbName).Collection(CollName)

	oID, _ := primitive.ObjectIDFromHex(isbn)

	var b t.Book
	err := col.FindOne(ctx, bson.M{
		"isbn":        oID,
	}).Decode(&b)

	return &b, err
}

func (s *Store) Create(ctx context.Context, b *t.CreateBookRequest) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)

	newBook, err := col.InsertOne(ctx, b)

	id := newBook.InsertedID.(primitive.ObjectID)
	return id, err
}
