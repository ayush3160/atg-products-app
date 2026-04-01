package app

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	db        *mongo.Database
	users     *mongo.Collection
	addresses *mongo.Collection
	products  *mongo.Collection
	orders    *mongo.Collection
	payments  *mongo.Collection
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		db:        db,
		users:     db.Collection("users"),
		addresses: db.Collection("addresses"),
		products:  db.Collection("products"),
		orders:    db.Collection("orders"),
		payments:  db.Collection("payments"),
	}
}

func (s *Store) EnsureIndexes(ctx context.Context) error {
	indexes := []struct {
		coll *mongo.Collection
		mdl  []mongo.IndexModel
	}{
		{
			coll: s.users,
			mdl: []mongo.IndexModel{
				{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
			},
		},
		{
			coll: s.addresses,
			mdl: []mongo.IndexModel{
				{Keys: bson.D{{Key: "userId", Value: 1}}},
			},
		},
		{
			coll: s.products,
			mdl: []mongo.IndexModel{
				{Keys: bson.D{{Key: "sku", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "active", Value: 1}}},
			},
		},
		{
			coll: s.orders,
			mdl: []mongo.IndexModel{
				{Keys: bson.D{{Key: "orderNumber", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "userId", Value: 1}}},
				{Keys: bson.D{{Key: "status", Value: 1}}},
			},
		},
		{
			coll: s.payments,
			mdl: []mongo.IndexModel{
				{Keys: bson.D{{Key: "orderId", Value: 1}}},
				{Keys: bson.D{{Key: "userId", Value: 1}}},
				{Keys: bson.D{{Key: "status", Value: 1}}},
			},
		},
	}

	for _, entry := range indexes {
		if len(entry.mdl) == 0 {
			continue
		}
		if _, err := entry.coll.Indexes().CreateMany(ctx, entry.mdl); err != nil {
			return err
		}
	}

	return nil
}

func getByID[T any](ctx context.Context, coll *mongo.Collection, id primitive.ObjectID) (T, error) {
	return findOne[T](ctx, coll, bson.M{"_id": id})
}

func listAll[T any](ctx context.Context, coll *mongo.Collection, filter interface{}, sort bson.D) ([]T, error) {
	findOpts := options.Find()
	if len(sort) > 0 {
		findOpts.SetSort(sort)
	}

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items := make([]T, 0)
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func findOne[T any](ctx context.Context, coll *mongo.Collection, filter interface{}) (T, error) {
	var item T
	err := coll.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		return item, mapMongoError(err)
	}
	return item, nil
}

func mapMongoError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNotFound
	}
	if isDuplicateKeyError(err) {
		return fmt.Errorf("%w: duplicate key", ErrConflict)
	}
	return err
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "e11000") || strings.Contains(msg, "duplicate key")
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

func normalizeString(value string) string {
	return strings.TrimSpace(value)
}

func normalizeStatus(value, fallback string) string {
	status := strings.ToLower(strings.TrimSpace(value))
	if status == "" {
		return fallback
	}
	return status
}

func normalizeCurrency(value string) string {
	currency := strings.ToUpper(strings.TrimSpace(value))
	if currency == "" {
		return "USD"
	}
	return currency
}

func roundMoney(value float64) float64 {
	return math.Round(value*100) / 100
}

func snapshotUser(user User) UserSnapshot {
	return UserSnapshot{
		ID:    user.ID,
		Name:  strings.TrimSpace(user.FirstName + " " + user.LastName),
		Email: user.Email,
		Phone: user.Phone,
	}
}

func snapshotAddress(address Address) AddressSnapshot {
	return AddressSnapshot{
		ID:        address.ID,
		Label:     address.Label,
		Line1:     address.Line1,
		Line2:     address.Line2,
		City:      address.City,
		State:     address.State,
		PostalCode: address.PostalCode,
		Country:   address.Country,
	}
}

func shortObjectID(id primitive.ObjectID) string {
	hex := id.Hex()
	if len(hex) <= 8 {
		return strings.ToUpper(hex)
	}
	return strings.ToUpper(hex[len(hex)-8:])
}

func generateOrderNumber(id primitive.ObjectID) string {
	return fmt.Sprintf("ORD-%s-%s", time.Now().UTC().Format("20060102"), shortObjectID(id))
}

func generatePaymentRef(orderNumber string, id primitive.ObjectID) string {
	return fmt.Sprintf("PAY-%s-%s", orderNumber, shortObjectID(id))
}

func validOrderStatus(status string) bool {
	switch status {
	case OrderStatusDraft, OrderStatusPending, OrderStatusPaid, OrderStatusCancelled, OrderStatusFulfilled, OrderStatusRefunded:
		return true
	default:
		return false
	}
}

func validPaymentStatus(status string) bool {
	switch status {
	case PaymentStatusPending, PaymentStatusSucceeded, PaymentStatusFailed, PaymentStatusRefunded:
		return true
	default:
		return false
	}
}

func userName(user User) string {
	return strings.TrimSpace(user.FirstName + " " + user.LastName)
}
