package app

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"

	OrderStatusDraft         = "draft"
	OrderStatusPending       = "pending_payment"
	OrderStatusPaid          = "paid"
	OrderStatusCancelled     = "cancelled"
	OrderStatusFulfilled     = "fulfilled"
	OrderStatusRefunded      = "refunded"

	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Email     string             `bson:"email" json:"email"`
	Phone     string             `bson:"phone" json:"phone"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Address struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"userId" json:"userId"`
	Label            string             `bson:"label" json:"label"`
	Line1            string             `bson:"line1" json:"line1"`
	Line2            string             `bson:"line2" json:"line2"`
	City             string             `bson:"city" json:"city"`
	State            string             `bson:"state" json:"state"`
	PostalCode       string             `bson:"postalCode" json:"postalCode"`
	Country          string             `bson:"country" json:"country"`
	IsDefaultShipping bool               `bson:"isDefaultShipping" json:"isDefaultShipping"`
	IsDefaultBilling  bool               `bson:"isDefaultBilling" json:"isDefaultBilling"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SKU         string             `bson:"sku" json:"sku"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Currency    string             `bson:"currency" json:"currency"`
	Stock       int                `bson:"stock" json:"stock"`
	Active      bool               `bson:"active" json:"active"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type UserSnapshot struct {
	ID    primitive.ObjectID `bson:"id" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Phone string             `bson:"phone" json:"phone"`
}

type AddressSnapshot struct {
	ID         primitive.ObjectID `bson:"id" json:"id"`
	Label      string             `bson:"label" json:"label"`
	Line1      string             `bson:"line1" json:"line1"`
	Line2      string             `bson:"line2" json:"line2"`
	City       string             `bson:"city" json:"city"`
	State      string             `bson:"state" json:"state"`
	PostalCode  string             `bson:"postalCode" json:"postalCode"`
	Country    string             `bson:"country" json:"country"`
}

type OrderItem struct {
	ProductID primitive.ObjectID `bson:"productId" json:"productId"`
	SKU       string             `bson:"sku" json:"sku"`
	Name      string             `bson:"name" json:"name"`
	UnitPrice float64            `bson:"unitPrice" json:"unitPrice"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	LineTotal float64            `bson:"lineTotal" json:"lineTotal"`
}

type Order struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderNumber            string             `bson:"orderNumber" json:"orderNumber"`
	UserID                 primitive.ObjectID `bson:"userId" json:"userId"`
	UserSnapshot           UserSnapshot       `bson:"userSnapshot" json:"userSnapshot"`
	ShippingAddressID      primitive.ObjectID `bson:"shippingAddressId" json:"shippingAddressId"`
	ShippingAddressSnapshot AddressSnapshot    `bson:"shippingAddressSnapshot" json:"shippingAddressSnapshot"`
	BillingAddressID       primitive.ObjectID `bson:"billingAddressId" json:"billingAddressId"`
	BillingAddressSnapshot AddressSnapshot    `bson:"billingAddressSnapshot" json:"billingAddressSnapshot"`
	Items                  []OrderItem        `bson:"items" json:"items"`
	Currency               string             `bson:"currency" json:"currency"`
	Subtotal               float64            `bson:"subtotal" json:"subtotal"`
	Tax                    float64            `bson:"tax" json:"tax"`
	ShippingFee            float64            `bson:"shippingFee" json:"shippingFee"`
	Discount               float64            `bson:"discount" json:"discount"`
	Total                  float64            `bson:"total" json:"total"`
	Status                 string             `bson:"status" json:"status"`
	CreatedAt              time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt              time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Payment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID     primitive.ObjectID `bson:"orderId" json:"orderId"`
	OrderNumber string             `bson:"orderNumber" json:"orderNumber"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Amount      float64            `bson:"amount" json:"amount"`
	Currency    string             `bson:"currency" json:"currency"`
	Method      string             `bson:"method" json:"method"`
	Provider    string             `bson:"provider" json:"provider"`
	ProviderRef string             `bson:"providerRef" json:"providerRef"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}

type CreateAddressRequest struct {
	UserID           string `json:"userId"`
	Label            string `json:"label"`
	Line1            string `json:"line1"`
	Line2            string `json:"line2"`
	City             string `json:"city"`
	State            string `json:"state"`
	PostalCode       string `json:"postalCode"`
	Country          string `json:"country"`
	IsDefaultShipping bool   `json:"isDefaultShipping"`
	IsDefaultBilling  bool   `json:"isDefaultBilling"`
}

type UpdateAddressRequest struct {
	UserID           string `json:"userId"`
	Label            string `json:"label"`
	Line1            string `json:"line1"`
	Line2            string `json:"line2"`
	City             string `json:"city"`
	State            string `json:"state"`
	PostalCode       string `json:"postalCode"`
	Country          string `json:"country"`
	IsDefaultShipping bool   `json:"isDefaultShipping"`
	IsDefaultBilling  bool   `json:"isDefaultBilling"`
}

type CreateProductRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Stock       int     `json:"stock"`
	Active      *bool   `json:"active"`
}

type UpdateProductRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Stock       int     `json:"stock"`
	Active      *bool   `json:"active"`
}

type OrderItemRequest struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type CreateOrderRequest struct {
	UserID            string             `json:"userId"`
	ShippingAddressID string             `json:"shippingAddressId"`
	BillingAddressID  string             `json:"billingAddressId"`
	Items             []OrderItemRequest `json:"items"`
	Tax               float64            `json:"tax"`
	ShippingFee       float64            `json:"shippingFee"`
	Discount          float64            `json:"discount"`
	Status            string             `json:"status"`
}

type UpdateOrderRequest struct {
	UserID            string             `json:"userId"`
	ShippingAddressID string             `json:"shippingAddressId"`
	BillingAddressID  string             `json:"billingAddressId"`
	Items             []OrderItemRequest `json:"items"`
	Tax               float64            `json:"tax"`
	ShippingFee       float64            `json:"shippingFee"`
	Discount          float64            `json:"discount"`
	Status            string             `json:"status"`
}

type CreatePaymentRequest struct {
	OrderID     string `json:"orderId"`
	Method      string `json:"method"`
	Provider    string `json:"provider"`
	ProviderRef string `json:"providerRef"`
	Status      string `json:"status"`
}

type UpdatePaymentRequest struct {
	Method      string `json:"method"`
	Provider    string `json:"provider"`
	ProviderRef string `json:"providerRef"`
	Status      string `json:"status"`
}
