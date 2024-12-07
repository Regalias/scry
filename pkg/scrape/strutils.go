package scrape

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func ToAscii(str string) (string, error) {
	result, _, err := transform.String(transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn))), str)
	if err != nil {
		return "", err
	}
	return result, nil
}

func ParsePrice(priceString string) (int64, error) {
	price_parts := strings.Split(priceString, ".")
	if len(price_parts) != 2 {
		return 0, fmt.Errorf("multiple '.' chars found. got: '%s'", priceString)
	}
	// Truncate off the cents to 2 digits
	cents := price_parts[1]
	if len(cents) > 2 {
		cents = cents[:2]
	}
	price_int, err := strconv.ParseInt(price_parts[0]+cents, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed parsing int64 for '%s': %w", priceString, err)
	}

	return price_int, nil
}
