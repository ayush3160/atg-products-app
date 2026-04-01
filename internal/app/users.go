package app

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) CreateUser(ctx context.Context, req CreateUserRequest) (User, error) {
	firstName := normalizeString(req.FirstName)
	lastName := normalizeString(req.LastName)
	email := strings.ToLower(normalizeString(req.Email))
	phone := normalizeString(req.Phone)

	if firstName == "" || lastName == "" || email == "" {
		return User{}, fmt.Errorf("%w: first name, last name, and email are required", ErrValidation)
	}

	now := nowUTC()
	user := User{
		ID:        primitive.NewObjectID(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := s.users.InsertOne(ctx, user); err != nil {
		return User{}, mapMongoError(err)
	}

	return user, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]User, error) {
	return listAll[User](ctx, s.users, bson.M{}, bson.D{{Key: "createdAt", Value: -1}})
}

func (s *Store) GetUser(ctx context.Context, id primitive.ObjectID) (User, error) {
	return getByID[User](ctx, s.users, id)
}

func (s *Store) UpdateUser(ctx context.Context, id primitive.ObjectID, req UpdateUserRequest) (User, error) {
	existing, err := s.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}

	firstName := normalizeString(req.FirstName)
	lastName := normalizeString(req.LastName)
	email := strings.ToLower(normalizeString(req.Email))
	phone := normalizeString(req.Phone)
	status := normalizeStatus(req.Status, existing.Status)

	if firstName == "" || lastName == "" || email == "" {
		return User{}, fmt.Errorf("%w: first name, last name, and email are required", ErrValidation)
	}
	if status != UserStatusActive && status != UserStatusDisabled {
		return User{}, fmt.Errorf("%w: invalid user status", ErrValidation)
	}

	updated := User{
		ID:        existing.ID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Status:    status,
		CreatedAt: existing.CreatedAt,
		UpdatedAt: nowUTC(),
	}

	if _, err := s.users.ReplaceOne(ctx, bson.M{"_id": existing.ID}, updated); err != nil {
		return User{}, mapMongoError(err)
	}

	return updated, nil
}

func (s *Store) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	if _, err := s.GetUser(ctx, id); err != nil {
		return err
	}

	if _, err := s.payments.DeleteMany(ctx, bson.M{"userId": id}); err != nil {
		return err
	}

	if _, err := s.orders.DeleteMany(ctx, bson.M{"userId": id}); err != nil {
		return err
	}

	if _, err := s.addresses.DeleteMany(ctx, bson.M{"userId": id}); err != nil {
		return err
	}

	_, err := s.users.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
