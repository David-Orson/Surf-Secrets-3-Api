package validation

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/David-Orson/Surf-Secrets-3-Api/store"
)

var s store.Store

type Validation struct {
	Errors []string `json:"errors"`
}

type ValidationFunc func(interface{}) Validation

func Init(st store.Store) {
	s = st
}

func (v *Validation) IsValid() bool {
	if len(v.Errors) == 0 {
		return true
	} else {
		return false
	}
}

func (v *Validation) isDuplicate(field string, table string, column string, value interface{}, count int) {
	if s.IsDuplicate(table, column, value, count) {
		v.Errors = append(v.Errors, field+" already exists")
	}
}

func (v *Validation) isEmail(field string, email string) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(email) {
		v.Errors = append(v.Errors, field+" - Must be a valid email address")
	}
}

func (v *Validation) maxLength(field string, text string, limit int) {
	if len(text) > limit {
		v.Errors = append(v.Errors, field+" - Maximum length is "+strconv.Itoa(limit))
	}
}

func (v *Validation) minLength(field string, text string, limit int) {
	if len(text) < limit {
		v.Errors = append(v.Errors, field+" - Minimum length is "+strconv.Itoa(limit))
	}
}

func (v *Validation) password(text string) {
	v.minLength("Password", text, 6)
	v.maxLength("password", text, 50)
}

func (v *Validation) AddError(text string) {
	v.Errors = append(v.Errors, text)
}

func sanitizeText(text *string) {
	*text = strings.ToValidUTF8(*text, "")
}
