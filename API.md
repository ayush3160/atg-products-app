# Sample ATG API

Base URL:

```bash
export BASE_URL="http://localhost:8080"
```

The API returns raw JSON documents for single resources and JSON arrays for list endpoints.

## Domain Model

### User

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string | MongoDB ObjectID |
| `firstName` | string | Required |
| `lastName` | string | Required |
| `email` | string | Required, unique |
| `phone` | string | Optional |
| `status` | string | `active` or `disabled` |
| `createdAt` | string | RFC3339 |
| `updatedAt` | string | RFC3339 |

### Address

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string | MongoDB ObjectID |
| `userId` | string | Required relationship to `users` |
| `label` | string | Example: `home`, `billing` |
| `line1` | string | Required |
| `line2` | string | Optional |
| `city` | string | Required |
| `state` | string | Required |
| `postalCode` | string | Required |
| `country` | string | Required |
| `isDefaultShipping` | boolean | Relationship hint |
| `isDefaultBilling` | boolean | Relationship hint |
| `createdAt` | string | RFC3339 |
| `updatedAt` | string | RFC3339 |

### Product

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string | MongoDB ObjectID |
| `sku` | string | Required, unique |
| `name` | string | Required |
| `description` | string | Optional |
| `price` | number | Required |
| `currency` | string | Defaults to `USD` |
| `stock` | integer | Required |
| `active` | boolean | Product availability |
| `createdAt` | string | RFC3339 |
| `updatedAt` | string | RFC3339 |

### Order

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string | MongoDB ObjectID |
| `orderNumber` | string | Human-friendly identifier |
| `userId` | string | Required relationship to `users` |
| `userSnapshot` | object | Embedded customer summary |
| `shippingAddressId` | string | Required relationship to `addresses` |
| `shippingAddressSnapshot` | object | Embedded shipping summary |
| `billingAddressId` | string | Required relationship to `addresses` |
| `billingAddressSnapshot` | object | Embedded billing summary |
| `items` | array | Product snapshots and quantities |
| `currency` | string | Derived from products |
| `subtotal` | number | Sum of item totals |
| `tax` | number | Order tax |
| `shippingFee` | number | Shipping charge |
| `discount` | number | Discount amount |
| `total` | number | Final order total |
| `status` | string | `draft`, `pending_payment`, `paid`, `cancelled`, `fulfilled`, `refunded` |
| `createdAt` | string | RFC3339 |
| `updatedAt` | string | RFC3339 |

### Payment

| Field | Type | Notes |
| --- | --- | --- |
| `id` | string | MongoDB ObjectID |
| `orderId` | string | Required relationship to `orders` |
| `orderNumber` | string | Snapshot of the order number |
| `userId` | string | Derived from the order |
| `amount` | number | Mirrors the order total on create |
| `currency` | string | Mirrors the order currency |
| `method` | string | Required, example: `card` |
| `provider` | string | Required |
| `providerRef` | string | Optional, auto-generated when omitted |
| `status` | string | `pending`, `succeeded`, `failed`, `refunded` |
| `createdAt` | string | RFC3339 |
| `updatedAt` | string | RFC3339 |

## Flows

1. Create a user.
2. Create one or more addresses for that user.
3. Create products.
4. Create an order that references the user, addresses, and products.
5. Create a payment for the order.
6. Update the records if needed.
7. Delete in reverse order.

## Curl Examples

### Health

```bash
curl "$BASE_URL/healthz"
curl "$BASE_URL/readyz"
```

### Users

Create:

```bash
curl -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Alicia",
    "lastName": "Stone",
    "email": "alicia.stone@example.com",
    "phone": "+2348012345678"
  }'
```

List:

```bash
curl "$BASE_URL/users"
```

Get:

```bash
curl "$BASE_URL/users/$USER_ID"
```

Update:

```bash
curl -X PUT "$BASE_URL/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Alicia",
    "lastName": "Stone-Updated",
    "email": "alicia.stone@example.com",
    "phone": "+2348999999999",
    "status": "active"
  }'
```

Delete:

```bash
curl -X DELETE "$BASE_URL/users/$USER_ID"
```

### Addresses

Create through the root collection:

```bash
curl -X POST "$BASE_URL/addresses" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "'"$USER_ID"'",
    "label": "home",
    "line1": "12 Sample Street",
    "line2": "Suite 4",
    "city": "Lagos",
    "state": "Lagos",
    "postalCode": "100001",
    "country": "Nigeria",
    "isDefaultShipping": true,
    "isDefaultBilling": false
  }'
```

Create through the nested user route:

```bash
curl -X POST "$BASE_URL/users/$USER_ID/addresses" \
  -H "Content-Type: application/json" \
  -d '{
    "label": "billing",
    "line1": "18 Billing Avenue",
    "line2": "Floor 2",
    "city": "Lagos",
    "state": "Lagos",
    "postalCode": "100002",
    "country": "Nigeria",
    "isDefaultShipping": false,
    "isDefaultBilling": true
  }'
```

List:

```bash
curl "$BASE_URL/addresses"
curl "$BASE_URL/addresses?userId=$USER_ID"
curl "$BASE_URL/users/$USER_ID/addresses"
```

Get:

```bash
curl "$BASE_URL/addresses/$ADDRESS_ID"
```

Update:

```bash
curl -X PUT "$BASE_URL/addresses/$ADDRESS_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "'"$USER_ID"'",
    "label": "home-updated",
    "line1": "12 Sample Street",
    "line2": "Suite 8",
    "city": "Lagos",
    "state": "Lagos",
    "postalCode": "100001",
    "country": "Nigeria",
    "isDefaultShipping": true,
    "isDefaultBilling": false
  }'
```

Delete:

```bash
curl -X DELETE "$BASE_URL/addresses/$ADDRESS_ID"
```

### Products

Create:

```bash
curl -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "SKU-001",
    "name": "Sample Keyboard",
    "description": "Mechanical keyboard with sample layout",
    "price": 79.99,
    "currency": "USD",
    "stock": 25,
    "active": true
  }'
```

List:

```bash
curl "$BASE_URL/products"
curl "$BASE_URL/products?active=true"
```

Get:

```bash
curl "$BASE_URL/products/$PRODUCT_ID"
```

Update:

```bash
curl -X PUT "$BASE_URL/products/$PRODUCT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "SKU-001",
    "name": "Sample Keyboard Pro",
    "description": "Updated mechanical keyboard",
    "price": 84.99,
    "currency": "USD",
    "stock": 20,
    "active": true
  }'
```

Delete:

```bash
curl -X DELETE "$BASE_URL/products/$PRODUCT_ID"
```

### Orders

Create:

```bash
curl -X POST "$BASE_URL/orders" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "'"$USER_ID"'",
    "shippingAddressId": "'"$SHIPPING_ADDRESS_ID"'",
    "billingAddressId": "'"$BILLING_ADDRESS_ID"'",
    "items": [
      { "productId": "'"$PRODUCT_ONE_ID"'", "quantity": 2 },
      { "productId": "'"$PRODUCT_TWO_ID"'", "quantity": 1 }
    ],
    "tax": 8.25,
    "shippingFee": 5.0,
    "discount": 3.0,
    "status": "pending_payment"
  }'
```

List:

```bash
curl "$BASE_URL/orders"
curl "$BASE_URL/orders?userId=$USER_ID"
curl "$BASE_URL/orders?userId=$USER_ID&status=pending_payment"
curl "$BASE_URL/users/$USER_ID/orders"
```

Get:

```bash
curl "$BASE_URL/orders/$ORDER_ID"
```

Update:

```bash
curl -X PUT "$BASE_URL/orders/$ORDER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "'"$USER_ID"'",
    "shippingAddressId": "'"$SHIPPING_ADDRESS_ID"'",
    "billingAddressId": "'"$BILLING_ADDRESS_ID"'",
    "items": [
      { "productId": "'"$PRODUCT_ONE_ID"'", "quantity": 2 },
      { "productId": "'"$PRODUCT_TWO_ID"'", "quantity": 1 }
    ],
    "tax": 8.25,
    "shippingFee": 5.0,
    "discount": 3.0,
    "status": "fulfilled"
  }'
```

Delete:

```bash
curl -X DELETE "$BASE_URL/orders/$ORDER_ID"
```

### Payments

Create through the root collection:

```bash
curl -X POST "$BASE_URL/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": "'"$ORDER_ID"'",
    "method": "card",
    "provider": "sandbox-payments",
    "providerRef": "tx-001",
    "status": "succeeded"
  }'
```

Create through the nested order route:

```bash
curl -X POST "$BASE_URL/orders/$ORDER_ID/payments" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "card",
    "provider": "sandbox-payments",
    "providerRef": "tx-001",
    "status": "succeeded"
  }'
```

List:

```bash
curl "$BASE_URL/payments"
curl "$BASE_URL/payments?orderId=$ORDER_ID"
curl "$BASE_URL/payments?userId=$USER_ID"
curl "$BASE_URL/payments?orderId=$ORDER_ID&status=succeeded"
curl "$BASE_URL/orders/$ORDER_ID/payments"
```

Get:

```bash
curl "$BASE_URL/payments/$PAYMENT_ID"
```

Update:

```bash
curl -X PUT "$BASE_URL/payments/$PAYMENT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "card",
    "provider": "sandbox-payments",
    "providerRef": "tx-001-refunded",
    "status": "refunded"
  }'
```

Delete:

```bash
curl -X DELETE "$BASE_URL/payments/$PAYMENT_ID"
```

## Notes

- Orders embed snapshots of the user, shipping address, billing address, and products at the time they are created.
- Payments derive their amount from the order total.
- Deleting a user cascades through their addresses, orders, and payments.
- Deleting an order cascades through its payments.
