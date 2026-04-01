package app

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Store) CreateOrder(ctx context.Context, req CreateOrderRequest) (Order, error) {
	order, err := s.assembleOrder(ctx, req, nil)
	if err != nil {
		return Order{}, err
	}

	if _, err := s.orders.InsertOne(ctx, order); err != nil {
		return Order{}, mapMongoError(err)
	}

	return order, nil
}

func (s *Store) ListOrders(ctx context.Context, userID *primitive.ObjectID, status string) ([]Order, error) {
	filter := bson.M{}
	if userID != nil {
		filter["userId"] = *userID
	}
	if normalized := strings.ToLower(strings.TrimSpace(status)); normalized != "" {
		filter["status"] = normalized
	}
	return listAll[Order](ctx, s.orders, filter, bson.D{{Key: "createdAt", Value: -1}})
}

func (s *Store) GetOrder(ctx context.Context, id primitive.ObjectID) (Order, error) {
	return getByID[Order](ctx, s.orders, id)
}

func (s *Store) UpdateOrder(ctx context.Context, id primitive.ObjectID, req UpdateOrderRequest) (Order, error) {
	existing, err := s.GetOrder(ctx, id)
	if err != nil {
		return Order{}, err
	}

	order, err := s.assembleOrder(ctx, CreateOrderRequest{
		UserID:            req.UserID,
		ShippingAddressID: req.ShippingAddressID,
		BillingAddressID:  req.BillingAddressID,
		Items:             req.Items,
		Tax:               req.Tax,
		ShippingFee:       req.ShippingFee,
		Discount:          req.Discount,
		Status:            req.Status,
	}, &existing)
	if err != nil {
		return Order{}, err
	}

	if _, err := s.orders.ReplaceOne(ctx, bson.M{"_id": existing.ID}, order); err != nil {
		return Order{}, mapMongoError(err)
	}

	return order, nil
}

func (s *Store) DeleteOrder(ctx context.Context, id primitive.ObjectID) error {
	if _, err := s.GetOrder(ctx, id); err != nil {
		return err
	}

	if _, err := s.payments.DeleteMany(ctx, bson.M{"orderId": id}); err != nil {
		return err
	}

	_, err := s.orders.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *Store) assembleOrder(ctx context.Context, req CreateOrderRequest, existing *Order) (Order, error) {
	userID, err := primitive.ObjectIDFromHex(normalizeString(req.UserID))
	if err != nil {
		return Order{}, fmt.Errorf("%w: invalid userId", ErrValidation)
	}

	shippingAddressID, err := primitive.ObjectIDFromHex(normalizeString(req.ShippingAddressID))
	if err != nil {
		return Order{}, fmt.Errorf("%w: invalid shippingAddressId", ErrValidation)
	}

	billingAddressID, err := primitive.ObjectIDFromHex(normalizeString(req.BillingAddressID))
	if err != nil {
		return Order{}, fmt.Errorf("%w: invalid billingAddressId", ErrValidation)
	}

	if len(req.Items) == 0 {
		return Order{}, fmt.Errorf("%w: at least one item is required", ErrValidation)
	}
	if req.Tax < 0 || req.ShippingFee < 0 || req.Discount < 0 {
		return Order{}, fmt.Errorf("%w: monetary values cannot be negative", ErrValidation)
	}

	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return Order{}, err
	}

	shippingAddress, err := s.GetAddress(ctx, shippingAddressID)
	if err != nil {
		return Order{}, err
	}
	if shippingAddress.UserID != userID {
		return Order{}, fmt.Errorf("%w: shipping address does not belong to the user", ErrValidation)
	}

	billingAddress, err := s.GetAddress(ctx, billingAddressID)
	if err != nil {
		return Order{}, err
	}
	if billingAddress.UserID != userID {
		return Order{}, fmt.Errorf("%w: billing address does not belong to the user", ErrValidation)
	}

	items := make([]OrderItem, 0, len(req.Items))
	var subtotal float64
	currency := ""

	for _, itemReq := range req.Items {
		productID, err := primitive.ObjectIDFromHex(normalizeString(itemReq.ProductID))
		if err != nil {
			return Order{}, fmt.Errorf("%w: invalid productId", ErrValidation)
		}
		if itemReq.Quantity <= 0 {
			return Order{}, fmt.Errorf("%w: item quantity must be greater than zero", ErrValidation)
		}

		product, err := s.GetProduct(ctx, productID)
		if err != nil {
			return Order{}, err
		}
		if !product.Active {
			return Order{}, fmt.Errorf("%w: product %s is not active", ErrValidation, product.SKU)
		}
		if product.Stock < itemReq.Quantity {
			return Order{}, fmt.Errorf("%w: product %s does not have enough stock", ErrValidation, product.SKU)
		}

		itemCurrency := normalizeCurrency(product.Currency)
		if currency == "" {
			currency = itemCurrency
		} else if currency != itemCurrency {
			return Order{}, fmt.Errorf("%w: all order items must use the same currency", ErrValidation)
		}

		unitPrice := roundMoney(product.Price)
		lineTotal := roundMoney(unitPrice * float64(itemReq.Quantity))
		subtotal = roundMoney(subtotal + lineTotal)

		items = append(items, OrderItem{
			ProductID: product.ID,
			SKU:       product.SKU,
			Name:      product.Name,
			UnitPrice: unitPrice,
			Quantity:  itemReq.Quantity,
			LineTotal: lineTotal,
		})
	}

	if currency == "" {
		currency = "USD"
	}

	total := roundMoney(subtotal + req.Tax + req.ShippingFee - req.Discount)
	if total < 0 {
		return Order{}, fmt.Errorf("%w: order total cannot be negative", ErrValidation)
	}

	status := normalizeStatus(req.Status, OrderStatusPending)
	if existing != nil && normalizeString(req.Status) == "" {
		status = existing.Status
	}
	if !validOrderStatus(status) {
		return Order{}, fmt.Errorf("%w: invalid order status", ErrValidation)
	}

	now := nowUTC()
	orderID := primitive.NewObjectID()
	order := Order{
		ID:                      orderID,
		OrderNumber:             generateOrderNumber(orderID),
		UserID:                  user.ID,
		UserSnapshot:            snapshotUser(user),
		ShippingAddressID:       shippingAddress.ID,
		ShippingAddressSnapshot: snapshotAddress(shippingAddress),
		BillingAddressID:        billingAddress.ID,
		BillingAddressSnapshot:  snapshotAddress(billingAddress),
		Items:                   items,
		Currency:                currency,
		Subtotal:                roundMoney(subtotal),
		Tax:                     roundMoney(req.Tax),
		ShippingFee:             roundMoney(req.ShippingFee),
		Discount:                roundMoney(req.Discount),
		Total:                   total,
		Status:                  status,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	if existing != nil {
		order.ID = existing.ID
		order.OrderNumber = existing.OrderNumber
		order.CreatedAt = existing.CreatedAt
		order.UpdatedAt = now
	}

	return order, nil
}
