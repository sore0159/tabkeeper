package main

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"
)

// Amount is assuming a directionality: Positive means
// Eric owes Julie (the most common case).  Negative means
// Julie owes Eric.
// Amount is in US pennies.
type Entry struct {
	Name          string    `json:"name"`
	Amount        int       `json:"amount"`
	DateSubmitted time.Time `json:"dateSubmitted"`
}

func EntryFromPost(r *http.Request) (*Entry, error) {
	name := r.PostFormValue("entry_name")
	amount := r.PostFormValue("entry_amount")
	direction := r.PostFormValue("entry_direction")

	if name == "" {
		return nil, errors.New("ENTRY NAME REQUIRED")
	} else if amount == "" {
		return nil, errors.New("ENTRY AMOUNT REQUIRED")
	}
	aF, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return nil, errors.New("MALFORMED ENTRY AMOUNT FORM VALUE")
	}
	if aF < 0.01 {
		return nil, errors.New("ENTRY AMOUNT REQUIRED TO BE POSITIVE")
	}
	aI := int(math.Floor(aF * 100))

	if direction == "julie_owe" {
		aI *= -1
	} else if direction != "eric_owe" {
		return nil, errors.New("MALFORMED ENTRY DIRECTION FORM VALUE")
	}
	return &Entry{
		DateSubmitted: time.Now(),
		Amount:        aI,
		Name:          name,
	}, nil
}
