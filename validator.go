package ptti

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// for now this is an example to create custom validation
// to use add at the end of the properties: SomeDate time.Time `bson:"some_date" json:"some_date" validate:"required,ltnow"`
// the name of the function should be the one declared during RegisterValidation (the first parameter [should be a string])
func ltnow(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}

	today := time.Now().Truncate(24 * time.Hour)
	inputDate := parsedDate.Truncate(24 * time.Hour)

	return inputDate.Before(today)
}

func Init() {
	_ = Validate.RegisterValidation("ltnow", ltnow)
}
