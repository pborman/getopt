package getopt

import (
	"fmt"
	"strings"
)

type boolValue bool

func (b *boolValue) Set(value string, opt Option) error {
	switch strings.ToLower(value) {
	case "", "1", "true", "on", "t":
		*b = true
	case "0", "false", "off", "f":
		*b = false
	default:
		return fmt.Errorf("invalid value for bool %s: %q", opt.Name(), value)
	}
	return nil
}

func (b *boolValue) String() string {
	if *b {
		return "true"
	}
	return "false"
}

// Bool declares a Flag option, returning a pointer to the bool value to check.
func Bool(name rune, helpvalue ...string) *bool {
	return CommandLine.Bool(name, helpvalue...)
}

// Bool declares a Flag option, returning a pointer to the bool value to check.
func (s *Set) Bool(name rune, helpvalue ...string) *bool {
	var p bool
	s.BoolVarLong(&p, "", name, helpvalue...)
	return &p
}

// BoolLong declares a Flag option, returning a pointer to the bool value to
// check.  If short is not 0 then it is the short equivalent of name.
func BoolLong(name string, short rune, helpvalue ...string) *bool {
	return CommandLine.BoolLong(name, short, helpvalue...)
}

// BoolLong declares a Flag option, returning a pointer to the bool value to
// check.  If short is not 0 then it is the short equivalent of name.
func (s *Set) BoolLong(name string, short rune, helpvalue ...string) *bool {
	var p bool
	s.BoolVarLong(&p, name, short, helpvalue...)
	return &p
}

// BoolVar declares a Flag option.  The argument p points to a bool variable in
// which to store the value of the option.
func BoolVar(p *bool, name rune, helpvalue ...string) Option {
	return CommandLine.BoolVar(p, name, helpvalue...)
}

// BoolVar declares a Flag option.  The argument p points to a bool variable in
// which to store the value of the option.
func (s *Set) BoolVar(p *bool, name rune, helpvalue ...string) Option {
	return s.BoolVarLong(p, "", name, helpvalue...)
}

// BoolVarLong declares a Flag option.  The argument p points to a bool variable in
// which to store the value of the option.
// If short is not 0 then it is the short equivalent of name.
func BoolVarLong(p *bool, name string, short rune, helpvalue ...string) Option {
	return CommandLine.BoolVarLong(p, name, short, helpvalue...)
}

// BoolVarLong declares a Flag option.  The argument p points to a bool variable in
// which to store the value of the option.
// If short is not 0 then it is the short equivalent of name.
func (s *Set) BoolVarLong(p *bool, name string, short rune, helpvalue ...string) Option {
	return s.VarLong((*boolValue)(p), name, short, helpvalue...).SetFlag()
}
