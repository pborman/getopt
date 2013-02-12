package getopt

type stringValue string

func (s *stringValue) Set(value string, opt Option) error {
	*s = stringValue(value)
	return nil
}

func (s *stringValue) String() string {
	return string(*s)
}

func String(name rune, value string, helpvalue ...string) *string {
	return CommandLine.String(name, value, helpvalue...)
}

func (s *Set) String(name rune, value string, helpvalue ...string) *string {
	p := value
	s.StringVarLong(&p, "", name, helpvalue...)
	return &p
}

func StringLong(name string, short rune, value string, helpvalue ...string) *string {
	return CommandLine.StringLong(name, short, value, helpvalue...)
}

func (s *Set) StringLong(name string, short rune, value string, helpvalue ...string) *string {
	p := value
	s.StringVarLong(&p, name, short, helpvalue...)
	return &p
}

func StringVar(p *string, name rune, helpvalue ...string) Option {
	return CommandLine.StringVar(p, name, helpvalue...)
}

func (s *Set) StringVar(p *string, name rune, helpvalue ...string) Option {
	return s.VarLong((*stringValue)(p), "", name, helpvalue...)
}

func StringVarLong(p *string, name string, short rune, helpvalue ...string) Option {
	return CommandLine.StringVarLong(p, name, short, helpvalue...)
}

func (s *Set) StringVarLong(p *string, name string, short rune, helpvalue ...string) Option {
	return s.VarLong((*stringValue)(p), name, short, helpvalue...)
}
