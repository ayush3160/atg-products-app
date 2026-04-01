package app

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) CreateProduct(ctx context.Context, req CreateProductRequest) (Product, error) {
	sku := strings.ToUpper(normalizeString(req.SKU))
	name := normalizeString(req.Name)
	description := normalizeString(req.Description)
	currency := normalizeCurrency(req.Currency)

	if sku == "" || name == "" {
		return Product{}, fmt.Errorf("%w: sku and name are required", ErrValidation)
	}
	if req.Price <= 0 {
		return Product{}, fmt.Errorf("%w: price must be greater than zero", ErrValidation)
	}
	if req.Stock < 0 {
		return Product{}, fmt.Errorf("%w: stock cannot be negative", ErrValidation)
	}

	now := nowUTC()
	product := Product{
		ID:          primitive.NewObjectID(),
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       roundMoney(req.Price),
		Currency:    currency,
		Stock:       req.Stock,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if req.Active != nil {
		product.Active = *req.Active
	}

	if _, err := s.products.InsertOne(ctx, product); err != nil {
		return Product{}, mapMongoError(err)
	}

	return product, nil
}

func (s *Store) ListProducts(ctx context.Context, active *bool) ([]Product, error) {
	filter := bson.M{}
	if active != nil {
		filter["active"] = *active
	}
	return listAll[Product](ctx, s.products, filter, bson.D{{Key: "createdAt", Value: -1}})
}

func (s *Store) GetProduct(ctx context.Context, id primitive.ObjectID) (Product, error) {
	return getByID[Product](ctx, s.products, id)
}

func (s *Store) UpdateProduct(ctx context.Context, id primitive.ObjectID, req UpdateProductRequest) (Product, error) {
	existing, err := s.GetProduct(ctx, id)
	if err != nil {
		return Product{}, err
	}

	sku := strings.ToUpper(normalizeString(req.SKU))
	name := normalizeString(req.Name)
	description := normalizeString(req.Description)
	currency := normalizeCurrency(req.Currency)

	if sku == "" || name == "" {
		return Product{}, fmt.Errorf("%w: sku and name are required", ErrValidation)
	}
	if req.Price <= 0 {
		return Product{}, fmt.Errorf("%w: price must be greater than zero", ErrValidation)
	}
	if req.Stock < 0 {
		return Product{}, fmt.Errorf("%w: stock cannot be negative", ErrValidation)
	}

	updated := Product{
		ID:          existing.ID,
		SKU:         sku,
		Name:        name,
		Description: description,
		Price:       roundMoney(req.Price),
		Currency:    currency,
		Stock:       req.Stock,
		Active:      existing.Active,
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   nowUTC(),
	}
	if req.Active != nil {
		updated.Active = *req.Active
	}

	if _, err := s.products.ReplaceOne(ctx, bson.M{"_id": existing.ID}, updated); err != nil {
		return Product{}, mapMongoError(err)
	}

	return updated, nil
}

func (s *Store) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	if _, err := s.GetProduct(ctx, id); err != nil {
		return err
	}

	_, err := s.products.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
