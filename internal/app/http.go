package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("%w: request body is required", ErrValidation)
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("%w: request body is required", ErrValidation)
		}
		return fmt.Errorf("%w: invalid JSON body: %v", ErrValidation, err)
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("%w: request body must contain a single JSON object", ErrValidation)
		}
		return fmt.Errorf("%w: invalid JSON body: %v", ErrValidation, err)
	}

	return nil
}

func parsePathObjectID(r *http.Request, key string) (primitive.ObjectID, error) {
	raw := strings.TrimSpace(r.PathValue(key))
	if raw == "" {
		return primitive.NilObjectID, fmt.Errorf("%w: %s is required", ErrValidation, key)
	}

	id, err := primitive.ObjectIDFromHex(raw)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("%w: invalid %s", ErrValidation, key)
	}
	return id, nil
}

func parseQueryObjectID(r *http.Request, key string) (*primitive.ObjectID, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return nil, nil
	}

	id, err := primitive.ObjectIDFromHex(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid %s", ErrValidation, key)
	}
	return &id, nil
}

func parseQueryBool(r *http.Request, key string) (*bool, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid %s", ErrValidation, key)
	}
	return &value, nil
}

func parseQueryString(r *http.Request, key string) string {
	return strings.TrimSpace(r.URL.Query().Get(key))
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrConflict):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
