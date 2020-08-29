package validator

import "regexp"

// ValidateString validate string using regex
func ValidateString(str string, pattern string) bool {
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}

// ValidateEmail validate email using regex
func ValidateEmail(email string) bool {
	return ValidateString(email, `(?i)^([a-z\-0-9\_]+(\.[a-z\-0-9\_]+)?(\.[a-z\-0-9\_]+)?)@[a-z0-9\-_]+(\.[a-z0-9\-_]+)(\.[a-z0-9\-_]+)?(\.[a-z0-9\-_]+)?$`)
}

// ValidateAccountDisplayName validate firstname and last name
func ValidateAccountDisplayName(name string) bool {
	return ValidateString(name, `(?i)^[^!@#$%^&\*()"'?<>\[\]]{1,50}$`)
}
