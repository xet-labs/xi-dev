// validate/validate.go
package lib

import "regexp"

type validateLib struct {
	uname *regexp.Regexp
	uid   *regexp.Regexp
	slug  *regexp.Regexp
}

var Validate = &validateLib{
	uname: regexp.MustCompile(`^@[a-zA-Z0-9_]{3,33}$`),
	uid:   regexp.MustCompile(`^[0-9]{1,20}$`),
	slug:  regexp.MustCompile(`^[a-zA-Z0-9_-]{3,64}$`),
}

func (v *validateLib) Uname(s string) bool { return v.uname.MatchString(s) }
func (v *validateLib) UID(s string) bool   { return v.uid.MatchString(s) }
func (v *validateLib) Slug(s string) bool  { return v.slug.MatchString(s) }
