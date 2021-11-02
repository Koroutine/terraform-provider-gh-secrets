package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// https://docs.github.com/en/actions/reference/encrypted-secrets#naming-your-secrets
var secretNameRegexp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

func validateSecretNameFunc(v interface{}, keyName string) (we []string, errs []error) {
	name, ok := v.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", keyName)}
	}

	if !secretNameRegexp.MatchString(name) {
		errs = append(errs, errors.New("secret names can only contain alphanumeric characters or underscores and must not start with a number"))
	}

	if strings.HasPrefix(strings.ToUpper(name), "GITHUB_") {
		errs = append(errs, errors.New("secret names must not start with the GITHUB_ prefix"))
	}

	return we, errs
}

// return the pieces of id `left:center:right` as left, center, right
func parseThreePartID(id, left, center, right string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("Unexpected ID format (%q). Expected %s:%s:%s", id, left, center, right)
	}

	return parts[0], parts[1], parts[2], nil
}

// format the strings into an id `a:b:c`
func buildThreePartID(a, b, c string) string {
	return fmt.Sprintf("%s:%s:%s", a, b, c)
}
