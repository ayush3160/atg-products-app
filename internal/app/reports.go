package app

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Store) GetSummary(ctx context.Context) (SummaryReport, error) {
	users, err := s.users.CountDocuments(ctx, bson.M{})
	if err != nil {
		return SummaryReport{}, err
	}

	activeUsers, err := s.users.CountDocuments(ctx, bson.M{"status": UserStatusActive})
	if err != nil {
		return SummaryReport{}, err
	}

	addresses, err := s.addresses.CountDocuments(ctx, bson.M{})
	if err != nil {
		return SummaryReport{}, err
	}

	products, err := s.products.CountDocuments(ctx, bson.M{})
	if err != nil {
		return SummaryReport{}, err
	}

	activeProducts, err := s.products.CountDocuments(ctx, bson.M{"active": true})
	if err != nil {
		return SummaryReport{}, err
	}

	orders, err := s.orders.CountDocuments(ctx, bson.M{})
	if err != nil {
		return SummaryReport{}, err
	}

	ordersByStatus, err := countStatuses(ctx, s.orders, "status", []string{
		OrderStatusDraft,
		OrderStatusPending,
		OrderStatusPaid,
		OrderStatusCancelled,
		OrderStatusFulfilled,
		OrderStatusRefunded,
	})
	if err != nil {
		return SummaryReport{}, err
	}

	payments, err := s.payments.CountDocuments(ctx, bson.M{})
	if err != nil {
		return SummaryReport{}, err
	}

	paymentsByStatus, err := countStatuses(ctx, s.payments, "status", []string{
		PaymentStatusPending,
		PaymentStatusSucceeded,
		PaymentStatusFailed,
		PaymentStatusRefunded,
	})
	if err != nil {
		return SummaryReport{}, err
	}

	grossRevenue, err := sumField(ctx, s.payments, bson.M{"status": PaymentStatusSucceeded}, "amount")
	if err != nil {
		return SummaryReport{}, err
	}

	pendingOrderValue, err := sumField(ctx, s.orders, bson.M{"status": OrderStatusPending}, "total")
	if err != nil {
		return SummaryReport{}, err
	}

	return SummaryReport{
		Users:             users,
		ActiveUsers:       activeUsers,
		Addresses:         addresses,
		Products:          products,
		ActiveProducts:    activeProducts,
		Orders:            orders,
		OrdersByStatus:    ordersByStatus,
		Payments:          payments,
		PaymentsByStatus:  paymentsByStatus,
		GrossRevenue:      roundMoney(grossRevenue),
		PendingOrderValue: roundMoney(pendingOrderValue),
		GeneratedAt:       nowUTC(),
	}, nil
}

func countStatuses(ctx context.Context, coll *mongo.Collection, field string, statuses []string) (map[string]int64, error) {
	counts := make(map[string]int64, len(statuses))
	for _, status := range statuses {
		count, err := coll.CountDocuments(ctx, bson.M{field: status})
		if err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, nil
}

func sumField(ctx context.Context, coll *mongo.Collection, match bson.M, field string) (float64, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$" + field},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var rows []struct {
		Total float64 `bson:"total"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return 0, err
	}
	if len(rows) == 0 {
		return 0, nil
	}
	return rows[0].Total, nil
}
