package app

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) CreatePayment(ctx context.Context, req CreatePaymentRequest) (Payment, error) {
	payment, err := s.assemblePayment(ctx, req, nil)
	if err != nil {
		return Payment{}, err
	}

	if _, err := s.payments.InsertOne(ctx, payment); err != nil {
		return Payment{}, mapMongoError(err)
	}

	if payment.Status == PaymentStatusSucceeded {
		order, err := s.GetOrder(ctx, payment.OrderID)
		if err != nil {
			return Payment{}, err
		}
		order.Status = OrderStatusPaid
		order.UpdatedAt = nowUTC()
		if _, err := s.orders.ReplaceOne(ctx, bson.M{"_id": order.ID}, order); err != nil {
			return Payment{}, mapMongoError(err)
		}
	}

	return payment, nil
}

func (s *Store) ListPayments(ctx context.Context, orderID, userID *primitive.ObjectID, status string) ([]Payment, error) {
	filter := bson.M{}
	if orderID != nil {
		filter["orderId"] = *orderID
	}
	if userID != nil {
		filter["userId"] = *userID
	}
	if normalized := strings.ToLower(strings.TrimSpace(status)); normalized != "" {
		filter["status"] = normalized
	}
	return listAll[Payment](ctx, s.payments, filter, bson.D{{Key: "createdAt", Value: -1}})
}

func (s *Store) GetPayment(ctx context.Context, id primitive.ObjectID) (Payment, error) {
	return getByID[Payment](ctx, s.payments, id)
}

func (s *Store) UpdatePayment(ctx context.Context, id primitive.ObjectID, req UpdatePaymentRequest) (Payment, error) {
	existing, err := s.GetPayment(ctx, id)
	if err != nil {
		return Payment{}, err
	}

	payment, err := s.assemblePayment(ctx, CreatePaymentRequest{
		OrderID:     existing.OrderID.Hex(),
		Method:      req.Method,
		Provider:    req.Provider,
		ProviderRef: req.ProviderRef,
		Status:      req.Status,
	}, &existing)
	if err != nil {
		return Payment{}, err
	}

	if _, err := s.payments.ReplaceOne(ctx, bson.M{"_id": existing.ID}, payment); err != nil {
		return Payment{}, mapMongoError(err)
	}

	if payment.Status == PaymentStatusSucceeded {
		order, err := s.GetOrder(ctx, payment.OrderID)
		if err != nil {
			return Payment{}, err
		}
		order.Status = OrderStatusPaid
		order.UpdatedAt = nowUTC()
		if _, err := s.orders.ReplaceOne(ctx, bson.M{"_id": order.ID}, order); err != nil {
			return Payment{}, mapMongoError(err)
		}
	}

	return payment, nil
}

func (s *Store) DeletePayment(ctx context.Context, id primitive.ObjectID) error {
	if _, err := s.GetPayment(ctx, id); err != nil {
		return err
	}

	_, err := s.payments.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *Store) assemblePayment(ctx context.Context, req CreatePaymentRequest, existing *Payment) (Payment, error) {
	orderID, err := primitive.ObjectIDFromHex(normalizeString(req.OrderID))
	if err != nil {
		return Payment{}, fmt.Errorf("%w: invalid orderId", ErrValidation)
	}

	order, err := s.GetOrder(ctx, orderID)
	if err != nil {
		return Payment{}, err
	}
	if order.Status == OrderStatusCancelled {
		return Payment{}, fmt.Errorf("%w: cannot create a payment for a cancelled order", ErrValidation)
	}

	method := normalizeString(req.Method)
	provider := normalizeString(req.Provider)
	status := normalizeStatus(req.Status, PaymentStatusSucceeded)
	providerRef := normalizeString(req.ProviderRef)

	if existing != nil {
		if normalizeString(req.Method) == "" {
			method = existing.Method
		}
		if normalizeString(req.Provider) == "" {
			provider = existing.Provider
		}
		if normalizeString(req.ProviderRef) == "" {
			providerRef = existing.ProviderRef
		}
		if normalizeString(req.Status) == "" {
			status = existing.Status
		}
	}

	if method == "" || provider == "" {
		return Payment{}, fmt.Errorf("%w: method and provider are required", ErrValidation)
	}
	if !validPaymentStatus(status) {
		return Payment{}, fmt.Errorf("%w: invalid payment status", ErrValidation)
	}

	paymentID := primitive.NewObjectID()
	if existing != nil {
		paymentID = existing.ID
	}
	if providerRef == "" {
		providerRef = generatePaymentRef(order.OrderNumber, paymentID)
	}

	now := nowUTC()
	amount := roundMoney(order.Total)
	currency := order.Currency
	if existing != nil {
		amount = existing.Amount
		currency = existing.Currency
	}

	payment := Payment{
		ID:          paymentID,
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		UserID:      order.UserID,
		Amount:      amount,
		Currency:    currency,
		Method:      method,
		Provider:    provider,
		ProviderRef: providerRef,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if existing != nil {
		payment.CreatedAt = existing.CreatedAt
		payment.UpdatedAt = now
	}

	return payment, nil
}
