package ptti

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// json to struct
func UnMarshaller[T any](w http.ResponseWriter, r *http.Request, input any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(input)
	if err != nil {
		return err
	}

	err = Validate.Struct(input.(T))
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range errs {
				fmt.Printf("❌ Field '%s' failed on '%s' rule",
					fieldErr.Field(), fieldErr.Tag())
			}
			return err
		}
		return err
	}
	return nil
}

// send response
func MarshalValue(w http.ResponseWriter, input any, customMsg ...string) {
	type response struct {
		Message string `bson:"message,omitempty" json:"message,omitempty"`
		Status  string `bson:"status" json:"status"`
		Code    int    `bson:"code,omitempty" json:"code,omitempty"`
		Payload any    `bson:"data,omitempty" json:"data,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := &response{
		Payload: input,
		Status:  "success",
		Message: "success",
		Code:    http.StatusOK,
	}

	if customMsg != nil {
		res.Message = customMsg[0]
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		ErrorHandler(w, err, "Encoding Error")
	}
}
