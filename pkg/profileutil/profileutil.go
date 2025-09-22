// pkg/profileutil/profileutil.go
package profileutil

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var reCountry = regexp.MustCompile(`^[A-Za-z]{2}$`)

// --- single value normalizers/validators ---

func NormalizeDisplayName(s string) (string, error) {
	n := strings.TrimSpace(s)
	if n == "" {
		return "", fmt.Errorf("name must not be empty")
	}
	if l := len(n); l < 1 || l > 80 {
		return "", fmt.Errorf("name length must be 1..80")
	}
	return n, nil
}

func NormalizeGender(s string) (string, error) {
	g := strings.ToLower(strings.TrimSpace(s))
	if g != "male" && g != "female" {
		return "", fmt.Errorf("gender must be male or female")
	}
	return g, nil
}

func NormalizeGenderFilter(s string) (string, error) {
	gf := strings.ToLower(strings.TrimSpace(s))
	switch gf {
	case "male", "female", "any":
		return gf, nil
	default:
		return "", fmt.Errorf("gender_filter must be one of: male, female, any")
	}
}

func NormalizeCountryCode(s string) (string, error) {
	cc := strings.ToUpper(strings.TrimSpace(s))
	if !reCountry.MatchString(cc) {
		return "", fmt.Errorf("country_code must be 2 letters")
	}
	return cc, nil
}

func ValidateAgePtr(p *int) error {
	if p == nil {
		return nil
	}
	if *p < 13 || *p > 120 {
		return fmt.Errorf("age must be 13..120")
	}
	return nil
}

func ValidateLevelPtr(p *int) error {
	if p == nil {
		return nil
	}
	if *p < 1 || *p > 6 {
		return fmt.Errorf("level must be 1..6")
	}
	return nil
}

func ValidateLevelRange(minP, maxP *int) error {
	if minP != nil {
		if err := ValidateLevelPtr(minP); err != nil {
			return err
		}
	}
	if maxP != nil {
		if err := ValidateLevelPtr(maxP); err != nil {
			return err
		}
	}
	if minP != nil && maxP != nil && *minP > *maxP {
		return fmt.Errorf("min_level cannot be greater than max_level")
	}
	return nil
}

// --- slices ---

func NormalizeCountries(codes []string) []string {
	if len(codes) == 0 {
		return codes
	}
	m := make(map[string]struct{}, len(codes))
	out := make([]string, 0, len(codes))
	for _, v := range codes {
		cc := strings.ToUpper(strings.TrimSpace(v))
		if !reCountry.MatchString(cc) {
			continue
		}
		if _, ok := m[cc]; ok {
			continue
		}
		m[cc] = struct{}{}
		out = append(out, cc)
	}
	sort.Strings(out)
	return out
}

func DedupInts(in []int) []int {
	if len(in) == 0 {
		return in
	}
	m := make(map[int]struct{}, len(in))
	out := make([]int, 0, len(in))
	for _, v := range in {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}
