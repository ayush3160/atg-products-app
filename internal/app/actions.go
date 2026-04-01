package app

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) AdjustProductStock(ctx context.Context, id primitive.ObjectID, delta int) (Product, error) {
	product, err := s.GetProduct(ctx, id)
	if err != nil {
		return Product{}, err
	}

	nextStock := product.Stock + delta
	if nextStock < 0 {
		return Product{}, fmt.Errorf("%w: stock cannot be negative", ErrValidation)
	}

	product.Stock = nextStock
	product.UpdatedAt = nowUTC()

	if _, err := s.products.ReplaceOne(ctx, bson.M{"_id": product.ID}, product); err != nil {
		return Product{}, mapMongoError(err)
	}

	return product, nil
}

func (s *Store) CancelOrder(ctx context.Context, id primitive.ObjectID) (OrderCancellationResult, error) {
	order, err := s.GetOrder(ctx, id)
	if err != nil {
		return OrderCancellationResult{}, err
	}

	now := nowUTC()
	order.Status = OrderStatusCancelled
	order.UpdatedAt = now

	if _, err := s.orders.ReplaceOne(ctx, bson.M{"_id": order.ID}, order); err != nil {
		return OrderCancellationResult{}, mapMongoError(err)
	}

	payments, err := s.ListPayments(ctx, &id, nil, "")
	if err != nil {
		return OrderCancellationResult{}, err
	}

	var refundedPayments int64
	for _, payment := range payments {
		switch payment.Status {
		case PaymentStatusRefunded, PaymentStatusFailed:
			continue
		}

		payment.Status = PaymentStatusRefunded
		payment.UpdatedAt = now
		if _, err := s.payments.ReplaceOne(ctx, bson.M{"_id": payment.ID}, payment); err != nil {
			return OrderCancellationResult{}, mapMongoError(err)
		}
		refundedPayments++
	}

	return OrderCancellationResult{
		Order:            order,
		RefundedPayments: refundedPayments,
	}, nil
}
