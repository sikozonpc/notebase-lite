package user

import (
	"context"

	t "github.com/sikozonpc/notebase/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "notebase"
	CollName = "users"
)

type Store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, b t.RegisterRequest) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)

	newUser, err := col.InsertOne(ctx, b)

	id := newUser.InsertedID.(primitive.ObjectID)
	return id, err
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*t.User, error) {
	col := s.db.Database(DbName).Collection(CollName)

	var u t.User
	err := col.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&u)

	return &u, err
}

func (s *Store) GetUserByID(ctx context.Context, id string) (*t.User, error) {
	col := s.db.Database(DbName).Collection(CollName)

	oID, _ := primitive.ObjectIDFromHex(id)

	var u t.User
	err := col.FindOne(ctx, bson.M{
		"_id": oID,
	}).Decode(&u)

	return &u, err
}

func (s *Store) GetUsers(ctx context.Context) ([]*t.User, error) {
	col := s.db.Database(DbName).Collection(CollName)

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	users := make([]*t.User, 0)
	for cursor.Next(ctx) {
		var u t.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	return users, nil
}

func (s *Store) UpdateUser(ctx context.Context, u t.User) error {
	col := s.db.Database(DbName).Collection(CollName)

	_, err := col.UpdateOne(ctx, bson.M{
		"_id": u.ID,
	}, bson.M{
		"$set": bson.M{
			"firstName": u.FirstName,
			"lastName":  u.LastName,
			"email":     u.Email,
			"password":  u.Password,
			"isActive":  u.IsActive,
		},
	})

	return err
}
