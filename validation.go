package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// different tests for different types of validation
var isAlpha = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
var isAlphaNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
var isEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString
var isUUID = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[8|9|a|A|b|B][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$").MatchString

// the validPair struct is used to store boh the tag and the value
// of a field within the User Struct
type validPair struct {
	tag   string
	value string
}

type validationError struct {
	s string
}

func (e validationError) Error() string {
	return e.s
}

// load the field names and tags into the valid map[string]string
func parseValid(u *User) map[string]validPair {
	valid := make(map[string]validPair)
	t := reflect.TypeOf(*u)
	v := reflect.ValueOf(*u)
	val := reflect.Indirect(reflect.ValueOf(u))
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i).Name
		fieldName, found := t.FieldByName(field)
		if !found {
			continue
		}
		temp := fmt.Sprint(v.Field(i).Interface())
		valid[field] = validPair{
			tag:   fieldName.Tag.Get("valid"),
			value: temp,
		}
	}
	return valid
}

// ValidateUser will validate the User struct to ensure it is correct
func ValidateUser(u *User) (bool, error) {
	v := parseValid(u)
	result, err := validate(v)
	return result, err
}

// used to validate the different fields
func validate(v map[string]validPair) (bool, error) {
	valid := true
	var err validationError
	for field := range v {
		switch {
		case strings.Contains(v[field].tag, "req") && v[field].value == "":
			err.s += fmt.Sprintf("%s: value is required; ", field)
			valid = false
		case strings.Contains(v[field].tag, "alpha") && !isAlpha(v[field].value):
			err.s += fmt.Sprintf("%s: value can only contain letters; ", field)
			valid = false
		case strings.Contains(v[field].tag, "alph-num") && !isAlphaNum(v[field].value):
			err.s += fmt.Sprintf("%s: value can only contain alphanumeric characters; ", field)
			valid = false
		case strings.Contains(v[field].tag, "email") && !isEmail(v[field].value):
			err.s += fmt.Sprintf("%s: value must be a valid email;", field)
			valid = false
		case strings.Contains(v[field].tag, "uuid") && !isUUID(v[field].value):
			err.s += fmt.Sprintf("%s: value must be a valid UUID;", field)
			//valid = false
		}
	}
	return valid, err
}
