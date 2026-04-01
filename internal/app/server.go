package app

import "net/http"

type apiServer struct {
	store *Store
}

func NewServer(store *Store) http.Handler {
	server := &apiServer{store: store}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", server.handleRoot)
	mux.HandleFunc("GET /healthz", server.handleHealth)
	mux.HandleFunc("GET /readyz", server.handleReady)

	mux.HandleFunc("POST /users", server.handleCreateUser)
	mux.HandleFunc("GET /users", server.handleListUsers)
	mux.HandleFunc("GET /users/{userID}", server.handleGetUser)
	mux.HandleFunc("PUT /users/{userID}", server.handleUpdateUser)
	mux.HandleFunc("DELETE /users/{userID}", server.handleDeleteUser)
	mux.HandleFunc("POST /users/{userID}/addresses", server.handleCreateUserAddress)
	mux.HandleFunc("GET /users/{userID}/addresses", server.handleListUserAddresses)
	mux.HandleFunc("GET /users/{userID}/orders", server.handleListUserOrders)

	mux.HandleFunc("POST /addresses", server.handleCreateAddress)
	mux.HandleFunc("GET /addresses", server.handleListAddresses)
	mux.HandleFunc("GET /addresses/{addressID}", server.handleGetAddress)
	mux.HandleFunc("PUT /addresses/{addressID}", server.handleUpdateAddress)
	mux.HandleFunc("DELETE /addresses/{addressID}", server.handleDeleteAddress)

	mux.HandleFunc("POST /products", server.handleCreateProduct)
	mux.HandleFunc("GET /products", server.handleListProducts)
	mux.HandleFunc("GET /products/{productID}", server.handleGetProduct)
	mux.HandleFunc("PUT /products/{productID}", server.handleUpdateProduct)
	mux.HandleFunc("DELETE /products/{productID}", server.handleDeleteProduct)

	mux.HandleFunc("POST /orders", server.handleCreateOrder)
	mux.HandleFunc("GET /orders", server.handleListOrders)
	mux.HandleFunc("GET /orders/{orderID}", server.handleGetOrder)
	mux.HandleFunc("PUT /orders/{orderID}", server.handleUpdateOrder)
	mux.HandleFunc("DELETE /orders/{orderID}", server.handleDeleteOrder)
	mux.HandleFunc("POST /orders/{orderID}/payments", server.handleCreateOrderPayment)
	mux.HandleFunc("GET /orders/{orderID}/payments", server.handleListOrderPayments)

	mux.HandleFunc("POST /payments", server.handleCreatePayment)
	mux.HandleFunc("GET /payments", server.handleListPayments)
	mux.HandleFunc("GET /payments/{paymentID}", server.handleGetPayment)
	mux.HandleFunc("PUT /payments/{paymentID}", server.handleUpdatePayment)
	mux.HandleFunc("DELETE /payments/{paymentID}", server.handleDeletePayment)

	return mux
}

func (s *apiServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "sample-atg-app",
		"version": "v1",
		"endpoints": []string{
			"GET /healthz",
			"POST /users",
			"POST /addresses",
			"POST /products",
			"POST /orders",
			"POST /payments",
		},
	})
}

func (s *apiServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *apiServer) handleReady(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *apiServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	user, err := s.store.CreateUser(r.Context(), req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (s *apiServer) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.ListUsers(r.Context())
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (s *apiServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "userID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	user, err := s.store.GetUser(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (s *apiServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "userID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var req UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	user, err := s.store.UpdateUser(r.Context(), id, req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (s *apiServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "userID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if err := s.store.DeleteUser(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id.Hex()})
}

func (s *apiServer) handleCreateAddress(w http.ResponseWriter, r *http.Request) {
	var req CreateAddressRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	address, err := s.store.CreateAddress(r.Context(), req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, address)
}

func (s *apiServer) handleCreateUserAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := parsePathObjectID(r, "userID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var req CreateAddressRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}
	req.UserID = userID.Hex()

	address, err := s.store.CreateAddress(r.Context(), req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, address)
}

func (s *apiServer) handleListAddresses(w http.ResponseWriter, r *http.Request) {
	userID, err := parseQueryObjectID(r, "userId")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	addresses, err := s.store.ListAddresses(r.Context(), userID)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, addresses)
}

func (s *apiServer) handleListUserAddresses(w http.ResponseWriter, r *http.Request) {
	userID, err := parsePathObjectID(r, "userID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	addresses, err := s.store.ListAddresses(r.Context(), &userID)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, addresses)
}

func (s *apiServer) handleGetAddress(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "addressID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	address, err := s.store.GetAddress(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, address)
}

func (s *apiServer) handleUpdateAddress(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "addressID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var req UpdateAddressRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	address, err := s.store.UpdateAddress(r.Context(), id, req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, address)
}

func (s *apiServer) handleDeleteAddress(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "addressID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if err := s.store.DeleteAddress(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id.Hex()})
}

func (s *apiServer) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	product, err := s.store.CreateProduct(r.Context(), req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, product)
}

func (s *apiServer) handleListProducts(w http.ResponseWriter, r *http.Request) {
	active, err := parseQueryBool(r, "active")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	products, err := s.store.ListProducts(r.Context(), active)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

func (s *apiServer) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "productID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	product, err := s.store.GetProduct(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (s *apiServer) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "productID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var req UpdateProductRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	product, err := s.store.UpdateProduct(r.Context(), id, req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (s *apiServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "productID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if err := s.store.DeleteProduct(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id.Hex()})
}

func (s *apiServer) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	order, err := s.store.CreateOrder(r.Context(), req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, order)
}

func (s *apiServer) handleListOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := parseQueryObjectID(r, "userId")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	status := parseQueryString(r, "status")
	orders, err := s.store.ListOrders(r.Context(), userID, status)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

func (s *apiServer) handleListUserOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := parsePathObjectID(r, "userID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	orders, err := s.store.ListOrders(r.Context(), &userID, "")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

func (s *apiServer) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "orderID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	order, err := s.store.GetOrder(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (s *apiServer) handleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "orderID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var req UpdateOrderRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	order, err := s.store.UpdateOrder(r.Context(), id, req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, order)
}

func (s *apiServer) handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "orderID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if err := s.store.DeleteOrder(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id.Hex()})
}

func (s *apiServer) handleCreateOrderPayment(w http.ResponseWriter, r *http.Request) {
	orderID, err := parsePathObjectID(r, "orderID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var req CreatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}
	req.OrderID = orderID.Hex()

	payment, err := s.store.CreatePayment(r.Context(), req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, payment)
}

func (s *apiServer) handleListOrderPayments(w http.ResponseWriter, r *http.Request) {
	orderID, err := parsePathObjectID(r, "orderID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	payments, err := s.store.ListPayments(r.Context(), &orderID, nil, "")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, payments)
}

func (s *apiServer) handleCreatePayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	payment, err := s.store.CreatePayment(r.Context(), req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, payment)
}

func (s *apiServer) handleListPayments(w http.ResponseWriter, r *http.Request) {
	orderID, err := parseQueryObjectID(r, "orderId")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	userID, err := parseQueryObjectID(r, "userId")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	status := parseQueryString(r, "status")
	payments, err := s.store.ListPayments(r.Context(), orderID, userID, status)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, payments)
}

func (s *apiServer) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "paymentID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	payment, err := s.store.GetPayment(r.Context(), id)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func (s *apiServer) handleUpdatePayment(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "paymentID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	var req UpdatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeStoreError(w, err)
		return
	}

	payment, err := s.store.UpdatePayment(r.Context(), id, req)
	if err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, payment)
}

func (s *apiServer) handleDeletePayment(w http.ResponseWriter, r *http.Request) {
	id, err := parsePathObjectID(r, "paymentID")
	if err != nil {
		writeStoreError(w, err)
		return
	}

	if err := s.store.DeletePayment(r.Context(), id); err != nil {
		writeStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id.Hex()})
}
