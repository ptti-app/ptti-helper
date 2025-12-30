package ptti

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type errStruct struct {
	Code    int    `bson:"code,omitempty" json:"code,omitempty"`
	Status  string `bson:"status" json:"status"`
	Message string `bson:"message,omitempty" json:"message,omitempty"`
	Error   string `bson:"error,omitempty" json:"error,omitempty"`
}

func ErrorHandler(w http.ResponseWriter, err error, customMsg ...string) {
	errRes := &errStruct{}
	errRes.Error = err.Error()
	errRes.Status = "Failed to fetch"

	var verss validator.ValidationErrors
	// add more
	switch {
	case err == mongo.ErrNoDocuments:
		errRes.Message = "No documents found!"
		errRes.Code = http.StatusNotFound
	case mongo.IsNetworkError(err):
		errRes.Message = "Network Error!"
		errRes.Code = http.StatusBadRequest
	case mongo.IsDuplicateKeyError(err):
		errRes.Message = "Duplicated Key!"
		errRes.Code = http.StatusBadRequest
	case strings.Contains(strings.ToLower(err.Error()), "unauthorized"):
		errRes.Message = "Not authorized"
		errRes.Code = http.StatusUnauthorized
	case strings.Contains(strings.ToLower(err.Error()), "invalid token"):
		errRes.Message = "Not authorized"
		errRes.Code = http.StatusUnauthorized
	case errors.As(err, &verss):
		errRes.Message = "Bad user input"
		errRes.Code = http.StatusBadRequest
		errRes.Error = "Validation Failure"
	default:
		errRes.Message = "Internal Server Error!"
		errRes.Code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errRes.Code)

	if customMsg != nil {
		Log().Errorw("custom error msg", customMsg, err)
	}

	if err := json.NewEncoder(w).Encode(errRes); err != nil {
		Log().Error("Json encoding failed")
	}
}
