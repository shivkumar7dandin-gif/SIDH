package repository

import (
	"context"

	"github.com/shivkumar7dandin-gif/students-api/internal/user/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user model.User,
) error {

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*model.User, error) {

	var user model.User

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"username": username,
		},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UsernameExists(
	ctx context.Context,
	username string,
) (bool, error) {

	count, err := r.collection.CountDocuments(
		ctx,
		bson.M{
			"username": username,
		},
	)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
