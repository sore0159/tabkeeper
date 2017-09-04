package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Tab struct {
	TotalStr   string
	Total      int
	SplitList  []*ListPart
	AssumeUser int
}

type ListPart struct {
	Title   string
	Entries []*Entry
}

func ProcessTab(entries []*Entry, assume int) (*Tab, error) {
	var sum int
	var title string
	var lp *ListPart
	slist := make([]*ListPart, 0)
	for i, _ := range entries {
		e := entries[len(entries)-(i+1)]
		sum += e.Amount
		if e.Repeatable {
			if e.Amount < 0 {
				e.RepAmount = float64(e.Amount) * -0.01
			} else {
				e.RepAmount = float64(e.Amount) * 0.01
			}
		}
		e.AmountStr = FormatAmount(e.Amount)
		e.DateStr = e.DateSubmitted.Format("02")
		newTitle := e.DateSubmitted.Format("Jan 2006")
		if newTitle != title {
			lp = new(ListPart)
			slist = append(slist, lp)
			lp.Title = newTitle
			title = newTitle
		}
		lp.Entries = append(lp.Entries, e)
	}
	return &Tab{
		TotalStr:   FormatAmount(sum),
		Total:      sum,
		SplitList:  slist,
		AssumeUser: assume,
	}, nil
}

// Amount is assuming a directionality: Positive means
// Eric owes Julie (the most common case).  Negative means
// Julie owes Eric.
// Amount is in US pennies.
type Entry struct {
	Description   string    `json:"description"`
	Amount        int       `json:"amount"`
	DateSubmitted time.Time `json:"dateSubmitted"`
	Repeatable    bool      `json:"repeatable"`

	RepAmount float64 `json:"-"`
	AmountStr string  `json:"-"`
	DateStr   string  `json:"-"`
}

func EntryFromPost(r *http.Request) (*Entry, error) {
	name := r.PostFormValue("entry_name")
	amount := r.PostFormValue("entry_amount")
	direction := r.PostFormValue("entry_direction")
	repeatable := r.PostFormValue("entry_repeatable")

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
		Description:   name,
		Repeatable:    repeatable == "on",
	}, nil
}

func FormatAmount(x int) string {
	if x == 0 {
		return "Julie and Eric have an even balance!"
	} else if x > 0 {
		return fmt.Sprintf("Eric owes Julie $%.2f", float64(x)*0.01)
	} else {
		return fmt.Sprintf("Julie owes Eric $%.2f", float64(-x)*0.01)
	}
}
