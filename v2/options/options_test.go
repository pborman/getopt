package options

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pborman/getopt/v2"
)

type theOptions struct {
	Name    string        `getopt:"--name=NAME      name of the widget"`
	Count   int           `getopt:"--count -c=COUNT number of widgets"`
	Verbose bool          `getopt:"-v               be verbose"`
	N       int           `getopt:"-n=NUMBER        set n to NUMBER"`
	Timeout time.Duration `getopt:"--timeout        duration of run"`
	Lazy    string
}

var myOptions = theOptions{
	Count: 42,
}

// This is the help we expect from theOptions.  If you change theOptions then
// you must change this string.  Note that getopt.HelpColumn must be set to 25.
var theHelp = `
Usage: program [-v] [-c COUNT] [--lazy value] [-n NUMBER] [--name NAME] [--timeout value] [parameters ...]
 -c, --count=COUNT    number of widgets
     --lazy=value     unspecified
 -n NUMBER            set n to NUMBER
     --name=NAME      name of the widget
     --timeout=value  duration of run
 -v                   be verbose
`[1:]

func TestLookup(t *testing.T) {
	opt := &struct {
		Option string `getopt:"--option -o"`
	}{
		Option: "value",
	}
	if o := Lookup(opt, "option"); o.(string) != "value" {
		t.Errorf("--option returned value %q, want %q", o, "value")
	}
	if o := Lookup(opt, "o"); o.(string) != "value" {
		t.Errorf("-o returned value %q, want %q", o, "value")
	}
}

func TestHelp(t *testing.T) {
	getopt.HelpColumn = 25
	opts, s := RegisterNew(&myOptions)
	if dopts, ok := opts.(*theOptions); !ok {
		t.Errorf("RegisterNew returned type %T, want %T", dopts, opts)
	}
	var buf bytes.Buffer
	s.SetProgram("program")
	s.PrintUsage(&buf)
	if help := buf.String(); help != theHelp {
		t.Errorf("Got help:\n%s\nWant:\n%s", help, theHelp)
	}
}

func TestRegisterSet(t *testing.T) {
	opts := &struct {
		Name string `getopt:"--the_name"`
	}{
		Name: "bob",
	}
	s := getopt.New()
	RegisterSet(opts, s)
	s.VisitAll(func(o getopt.Option) {
		if o.Name() != "--the_name" {
			t.Errorf("unexpected option found: %q", o.Name())
			return
		}
		if v := o.String(); v != "bob" {
			t.Errorf("%s=%q, want %q", o.Name(), v, "bob")
		}
	})
	s.Parse([]string{"", "--the_name", "fred"})
	s.VisitAll(func(o getopt.Option) {
		if o.Name() != "--the_name" {
			t.Errorf("unexpected option found: %q", o.Name())
			return
		}
		if v := o.String(); v != "fred" {
			t.Errorf("%s=%q, want %q", o.Name(), v, "fred")
		}
	})
}

func TestParseTag(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		tag  *optTag
		err  string
	}{
		{
			name: "nothing",
		},
		{
			name: "dash",
			in:   "-",
		},
		{
			name: "dash-dash",
			in:   "--",
		},
		{
			name: "long arg",
			in:   "--option",
			tag: &optTag{
				long: "option",
			},
		},
		{
			name: "short arg",
			in:   "-o",
			tag: &optTag{
				short: 'o',
			},
		},
		{
			name: "long help",
			in:   "--option this is an option",
			tag: &optTag{
				long: "option",
				help: "this is an option",
			},
		},
		{
			name: "long help1",
			in:   "--option -- this is an option",
			tag: &optTag{
				long: "option",
				help: "this is an option",
			},
		},
		{
			name: "long help2",
			in:   "--option - this is an option",
			tag: &optTag{
				long: "option",
				help: "this is an option",
			},
		},
		{
			name: "long help3",
			in:   "--option -- -this is an option",
			tag: &optTag{
				long: "option",
				help: "-this is an option",
			},
		},
		{
			name: "long and short arg",
			in:   "--option -o",
			tag: &optTag{
				long:  "option",
				short: 'o',
			},
		},
		{
			name: "short and long arg",
			in:   "-o --option",
			tag: &optTag{
				long:  "option",
				short: 'o',
			},
		},
		{
			name: "long arg with param",
			in:   "--option=PARAM",
			tag: &optTag{
				long:  "option",
				param: "PARAM",
			},
		},
		{
			name: "short arg with param",
			in:   "-o=PARAM",
			tag: &optTag{
				short: 'o',
				param: "PARAM",
			},
		},
		{
			name: "everything",
			in:   "--option=PARAM -o -- - this is help",
			tag: &optTag{
				long:  "option",
				short: 'o',
				param: "PARAM",
				help:  "- this is help",
			},
		},
		{
			name: "two longs",
			in:   "--option1 --option2",
			err:  "tag has too many long names",
		},
		{
			name: "two shorts",
			in:   "-a -b",
			err:  "tag has too many short names",
		},
		{
			name: "two parms",
			in:   "--option=PARAM1 -o=PARAM2",
			err:  "tag has multiple parameter names",
		},
		{
			name: "missing option",
			in:   "no option",
			err:  "tag missing option name",
		},
		{
			name: "long param only",
			in:   "--=PARAM",
			err:  "tag missing option name",
		},
		{
			name: "short param only",
			in:   "-=PARAM",
			err:  "tag missing option name",
		},
		{
			name: "two many dashes",
			in:   "---option",
			err:  "tag must not start with ---",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tag, err := parseTag(tt.in)
			switch {
			case err == nil && tt.err != "":
				t.Fatalf("did not get expected error %v", tt.err)
			case err != nil && tt.err == "":
				t.Fatalf("unexpected error %v", err)
			case err == nil:
			case !strings.Contains(err.Error(), tt.err):
				t.Fatalf("got error %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(tag, tt.tag) {
				t.Errorf("got %v, want %v", tag, tt.tag)
			}
		})
	}
}
