package receipt

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	Description string `json:"shortDescription"`
	Price       string `json:"price"`
}
// The rule function type that can assign points
// based on business requirements
type RuleFunc func(Receipt) int

// Calculate the total points by calling each rule function
// Note this will not return an error. To see the output of
// each rule run in debug mode. Ideally validation on the receipt
// would have been done in a layer before this is called so no validation
// happens in this function.
func (r Receipt) CalculatePoints() int {
	ruleFns := []RuleFunc{
		countAlphanumericChars,
		totalIsRoundDollar,
		everyTwoItems,
		itemDescription,
		purchasedAt,
		totalIsMultipleOf25,
	}
	return calculatePoints(r, ruleFns...)
}

// Apply all the rules and return the total points
// awarded to the reciept
func calculatePoints(rec Receipt, fns ...RuleFunc) int {
	points := 0
	for _, fn := range fns {
		points += fn(rec)
	}
	return points
}

// One point for every alphanumeric character in the retailer name.
func countAlphanumericChars(rec Receipt) int {
	count := 0
	for _, r := range rec.Retailer {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count++
		}
	}
	slog.Debug("Assigned points in countAlphanumericChars ", "points", count)
	return count
}

// 50 points if the total is a round dollar amount with no cents.
func totalIsRoundDollar(rec Receipt) int {
	total, err := strconv.ParseFloat(rec.Total, 64)
	if err != nil {
		slog.Debug("Assigned points in totalIsRoundDollar ", "points", 0)
		return 0
	}
	if total == math.Trunc(total) {
		slog.Debug("Assigned points in totalIsRoundDollar ", "points", 50)
		return 50
	}
	slog.Debug("Assigned points in totalIsRoundDollar ", "points", 0)
	return 0
}
// 25 points if the total is a multiple of 0.25
func totalIsMultipleOf25(rec Receipt) int {
	total, err := strconv.ParseFloat(rec.Total, 64)
	if err != nil {
		slog.Debug("Assigned points in totalIsMultipleOf25 ", "points", 0)
		return 0
	}

	if math.Mod(total, 0.25) == 0 {
		slog.Debug("Assigned points in totalIsMultipleOf25 ", "points", 25)
		return 25
	}
	slog.Debug("Assigned points in totalIsMultipleOf25 ", "points", 0)
	return 0
}

// 5 points for every two items on the receipt.
func everyTwoItems(rec Receipt) int {
	points := 0
	if len(rec.Items) >= 2 {
		points += (len(rec.Items) / 2) * 5
	}
	slog.Debug("Assigned points in everyTwoItems ", "points", points)
	return points
}

// If the trimmed length of the item description is a multiple of 3,
// multiply the price by 0.2 and round up to the nearest integer.
func itemDescription(rec Receipt) int {
	points := 0
	for _, item := range rec.Items {
		descriptionLength := utf8.RuneCountInString(strings.TrimSpace(item.Description))
		if descriptionLength%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				continue
			}
			points += int(math.Ceil(itemPrice * 0.2))
		}
	}
	slog.Debug("Assigned points in itemDescription ", "points", points)
	return points
}

// 6 points if the day in the purchase date is odd
// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func purchasedAt(rec Receipt) int {
	points := 0
	purchaseDateTime := fmt.Sprintf("%v %v", rec.PurchaseDate, rec.PurchaseTime)
	purchasedAt, err := time.Parse("2006-01-02 15:04", purchaseDateTime)

	if err != nil {
		slog.Debug("Assigned 0 points in purchasedAt exit", "error", err.Error())
		return points
	}

	if purchasedAt.Day()%2 != 0 {
		points += 6
	}

	if purchasedAt.Hour() >= 14 && purchasedAt.Hour() <= 16 {
		points += 10
	}
	slog.Debug("Assigned points in purchasedAt ", "points", points)
	return points
}
