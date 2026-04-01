package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) CreateAddress(ctx context.Context, req CreateAddressRequest) (Address, error) {
	userID, err := primitive.ObjectIDFromHex(normalizeString(req.UserID))
	if err != nil {
		return Address{}, fmt.Errorf("%w: invalid userId", ErrValidation)
	}

	if _, err := s.GetUser(ctx, userID); err != nil {
		return Address{}, err
	}

	label := normalizeString(req.Label)
	line1 := normalizeString(req.Line1)
	city := normalizeString(req.City)
	state := normalizeString(req.State)
	postalCode := normalizeString(req.PostalCode)
	country := normalizeString(req.Country)

	if label == "" || line1 == "" || city == "" || state == "" || postalCode == "" || country == "" {
		return Address{}, fmt.Errorf("%w: address fields are required", ErrValidation)
	}

	now := nowUTC()
	address := Address{
		ID:                primitive.NewObjectID(),
		UserID:            userID,
		Label:             label,
		Line1:             line1,
		Line2:             normalizeString(req.Line2),
		City:              city,
		State:             state,
		PostalCode:        postalCode,
		Country:           country,
		IsDefaultShipping: req.IsDefaultShipping,
		IsDefaultBilling:  req.IsDefaultBilling,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if _, err := s.addresses.InsertOne(ctx, address); err != nil {
		return Address{}, mapMongoError(err)
	}

	return address, nil
}

func (s *Store) ListAddresses(ctx context.Context, userID *primitive.ObjectID) ([]Address, error) {
	filter := bson.M{}
	if userID != nil {
		filter["userId"] = *userID
	}
	return listAll[Address](ctx, s.addresses, filter, bson.D{{Key: "createdAt", Value: -1}})
}

func (s *Store) GetAddress(ctx context.Context, id primitive.ObjectID) (Address, error) {
	return getByID[Address](ctx, s.addresses, id)
}

func (s *Store) UpdateAddress(ctx context.Context, id primitive.ObjectID, req UpdateAddressRequest) (Address, error) {
	existing, err := s.GetAddress(ctx, id)
	if err != nil {
		return Address{}, err
	}

	userID := existing.UserID
	if normalizeString(req.UserID) != "" {
		userID, err = primitive.ObjectIDFromHex(normalizeString(req.UserID))
		if err != nil {
			return Address{}, fmt.Errorf("%w: invalid userId", ErrValidation)
		}
		if _, err := s.GetUser(ctx, userID); err != nil {
			return Address{}, err
		}
	}

	label := normalizeString(req.Label)
	line1 := normalizeString(req.Line1)
	city := normalizeString(req.City)
	state := normalizeString(req.State)
	postalCode := normalizeString(req.PostalCode)
	country := normalizeString(req.Country)

	if label == "" || line1 == "" || city == "" || state == "" || postalCode == "" || country == "" {
		return Address{}, fmt.Errorf("%w: address fields are required", ErrValidation)
	}

	updated := Address{
		ID:                existing.ID,
		UserID:            userID,
		Label:             label,
		Line1:             line1,
		Line2:             normalizeString(req.Line2),
		City:              city,
		State:             state,
		PostalCode:        postalCode,
		Country:           country,
		IsDefaultShipping: req.IsDefaultShipping,
		IsDefaultBilling:  req.IsDefaultBilling,
		CreatedAt:         existing.CreatedAt,
		UpdatedAt:         nowUTC(),
	}

	if _, err := s.addresses.ReplaceOne(ctx, bson.M{"_id": existing.ID}, updated); err != nil {
		return Address{}, mapMongoError(err)
	}

	return updated, nil
}

func (s *Store) DeleteAddress(ctx context.Context, id primitive.ObjectID) error {
	if _, err := s.GetAddress(ctx, id); err != nil {
		return err
	}

	_, err := s.addresses.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
