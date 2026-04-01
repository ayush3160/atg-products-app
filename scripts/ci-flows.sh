#!/bin/sh
set -eu

BASE_URL="${BASE_URL:-http://api:8080}"

log() {
	printf '\n==> %s\n' "$1"
}

request() {
	method="$1"
	path="$2"
	body="${3:-}"
	if [ -n "$body" ]; then
		curl -fsS -X "$method" "$BASE_URL$path" \
			-H "Content-Type: application/json" \
			--data "$body"
	else
		curl -fsS -X "$method" "$BASE_URL$path"
	fi
}

log "Create user"
USER_JSON=$(request POST /users "$(jq -n \
	--arg firstName "Alicia" \
	--arg lastName "Stone" \
	--arg email "alicia.stone@example.com" \
	--arg phone "+2348012345678" \
	'{firstName:$firstName,lastName:$lastName,email:$email,phone:$phone}')")
printf '%s\n' "$USER_JSON"
USER_ID=$(printf '%s' "$USER_JSON" | jq -r '.id')

log "Get user"
request GET "/users/$USER_ID" | jq .

log "Update user"
USER_JSON=$(request PUT "/users/$USER_ID" "$(jq -n \
	--arg firstName "Alicia" \
	--arg lastName "Stone-Updated" \
	--arg email "alicia.stone@example.com" \
	--arg phone "+2348999999999" \
	'{firstName:$firstName,lastName:$lastName,email:$email,phone:$phone,status:"active"}')")
printf '%s\n' "$USER_JSON"

log "Create shipping address"
SHIPPING_JSON=$(request POST /addresses "$(jq -n \
	--arg userId "$USER_ID" \
	--arg label "home" \
	--arg line1 "12 Sample Street" \
	--arg line2 "Suite 4" \
	--arg city "Lagos" \
	--arg state "Lagos" \
	--arg postalCode "100001" \
	--arg country "Nigeria" \
	'{userId:$userId,label:$label,line1:$line1,line2:$line2,city:$city,state:$state,postalCode:$postalCode,country:$country,isDefaultShipping:true,isDefaultBilling:false}')")
printf '%s\n' "$SHIPPING_JSON"
SHIPPING_ID=$(printf '%s' "$SHIPPING_JSON" | jq -r '.id')

log "Create billing address"
BILLING_JSON=$(request POST "/users/$USER_ID/addresses" "$(jq -n \
	--arg label "billing" \
	--arg line1 "18 Billing Avenue" \
	--arg line2 "Floor 2" \
	--arg city "Lagos" \
	--arg state "Lagos" \
	--arg postalCode "100002" \
	--arg country "Nigeria" \
	'{label:$label,line1:$line1,line2:$line2,city:$city,state:$state,postalCode:$postalCode,country:$country,isDefaultShipping:false,isDefaultBilling:true}')")
printf '%s\n' "$BILLING_JSON"
BILLING_ID=$(printf '%s' "$BILLING_JSON" | jq -r '.id')

log "Get address"
request GET "/addresses/$SHIPPING_ID" | jq .

log "Update address"
SHIPPING_JSON=$(request PUT "/addresses/$SHIPPING_ID" "$(jq -n \
	--arg userId "$USER_ID" \
	--arg label "home-updated" \
	--arg line1 "12 Sample Street" \
	--arg line2 "Suite 8" \
	--arg city "Lagos" \
	--arg state "Lagos" \
	--arg postalCode "100001" \
	--arg country "Nigeria" \
	'{userId:$userId,label:$label,line1:$line1,line2:$line2,city:$city,state:$state,postalCode:$postalCode,country:$country,isDefaultShipping:true,isDefaultBilling:false}')")
printf '%s\n' "$SHIPPING_JSON"

log "List addresses"
request GET "/addresses?userId=$USER_ID" | jq .

log "Create product one"
PRODUCT_ONE_JSON=$(request POST /products "$(jq -n \
	--arg sku "SKU-001" \
	--arg name "Sample Keyboard" \
	--arg description "Mechanical keyboard with sample layout" \
	--argjson price 79.99 \
	--arg currency "USD" \
	--argjson stock 25 \
	'{sku:$sku,name:$name,description:$description,price:$price,currency:$currency,stock:$stock,active:true}')")
printf '%s\n' "$PRODUCT_ONE_JSON"
PRODUCT_ONE_ID=$(printf '%s' "$PRODUCT_ONE_JSON" | jq -r '.id')

log "Create product two"
PRODUCT_TWO_JSON=$(request POST /products "$(jq -n \
	--arg sku "SKU-002" \
	--arg name "Sample Mouse" \
	--arg description "Wireless mouse with sample layout" \
	--argjson price 29.50 \
	--arg currency "USD" \
	--argjson stock 50 \
	'{sku:$sku,name:$name,description:$description,price:$price,currency:$currency,stock:$stock,active:true}')")
printf '%s\n' "$PRODUCT_TWO_JSON"
PRODUCT_TWO_ID=$(printf '%s' "$PRODUCT_TWO_JSON" | jq -r '.id')

log "Get product"
request GET "/products/$PRODUCT_ONE_ID" | jq .

log "Update product"
PRODUCT_ONE_JSON=$(request PUT "/products/$PRODUCT_ONE_ID" "$(jq -n \
	--arg sku "SKU-001" \
	--arg name "Sample Keyboard Pro" \
	--arg description "Updated mechanical keyboard" \
	--argjson price 84.99 \
	--arg currency "USD" \
	--argjson stock 20 \
	'{sku:$sku,name:$name,description:$description,price:$price,currency:$currency,stock:$stock,active:true}')")
printf '%s\n' "$PRODUCT_ONE_JSON"

log "List products"
request GET "/products?active=true" | jq .

log "Create order"
ORDER_JSON=$(request POST /orders "$(jq -n \
	--arg userId "$USER_ID" \
	--arg shippingAddressId "$SHIPPING_ID" \
	--arg billingAddressId "$BILLING_ID" \
	--arg productOneId "$PRODUCT_ONE_ID" \
	--arg productTwoId "$PRODUCT_TWO_ID" \
	'{
		userId:$userId,
		shippingAddressId:$shippingAddressId,
		billingAddressId:$billingAddressId,
		items:[
			{productId:$productOneId,quantity:2},
			{productId:$productTwoId,quantity:1}
		],
		tax:8.25,
		shippingFee:5.00,
		discount:3.00,
		status:"pending_payment"
	}')")
printf '%s\n' "$ORDER_JSON"
ORDER_ID=$(printf '%s' "$ORDER_JSON" | jq -r '.id')

log "Get order"
request GET "/orders/$ORDER_ID" | jq .

log "List orders"
request GET "/orders?userId=$USER_ID&status=pending_payment" | jq .

log "Create payment"
PAYMENT_JSON=$(request POST "/orders/$ORDER_ID/payments" "$(jq -n \
	--arg method "card" \
	--arg provider "sandbox-payments" \
	--arg providerRef "tx-001" \
	'{method:$method,provider:$provider,providerRef:$providerRef,status:"succeeded"}')")
printf '%s\n' "$PAYMENT_JSON"
PAYMENT_ID=$(printf '%s' "$PAYMENT_JSON" | jq -r '.id')

log "Get payment"
request GET "/payments/$PAYMENT_ID" | jq .

log "List payments"
request GET "/payments?orderId=$ORDER_ID&status=succeeded" | jq .

log "Update order"
request PUT "/orders/$ORDER_ID" "$(jq -n \
	--arg userId "$USER_ID" \
	--arg shippingAddressId "$SHIPPING_ID" \
	--arg billingAddressId "$BILLING_ID" \
	--arg productOneId "$PRODUCT_ONE_ID" \
	--arg productTwoId "$PRODUCT_TWO_ID" \
	'{
		userId:$userId,
		shippingAddressId:$shippingAddressId,
		billingAddressId:$billingAddressId,
		items:[
			{productId:$productOneId,quantity:2},
			{productId:$productTwoId,quantity:1}
		],
		tax:8.25,
		shippingFee:5.00,
		discount:3.00,
		status:"fulfilled"
	}')" | jq .

log "Update payment"
request PUT "/payments/$PAYMENT_ID" "$(jq -n \
	--arg method "card" \
	--arg provider "sandbox-payments" \
	--arg providerRef "tx-001-refunded" \
	'{method:$method,provider:$provider,providerRef:$providerRef,status:"refunded"}')" | jq .

log "List users"
request GET /users | jq .

log "Delete payment"
request DELETE "/payments/$PAYMENT_ID" | jq .

log "Delete order"
request DELETE "/orders/$ORDER_ID" | jq .

log "Delete products"
request DELETE "/products/$PRODUCT_ONE_ID" | jq .
request DELETE "/products/$PRODUCT_TWO_ID" | jq .

log "Delete addresses"
request DELETE "/addresses/$SHIPPING_ID" | jq .
request DELETE "/addresses/$BILLING_ID" | jq .

log "Delete user"
request DELETE "/users/$USER_ID" | jq .

log "Smoke flow complete"
