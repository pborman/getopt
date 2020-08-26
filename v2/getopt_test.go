// Copyright 2020 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package getopt

import (
	"testing"
)

func TestMandatory(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []string
		err  string
	}{
		{
			name: "required option present",
			in:   []string{"test", "-r"},
		},
		{
			name: "required option not present",
			in:   []string{"test", "-o"},
			err:  "test: option -r is mandatory",
		},
		{
			name: "no options",
			in:   []string{"test"},
			err:  "test: option -r is mandatory",
		},
	} {
		reset()
		var val bool
		Flag(&val, 'o')
		Flag(&val, 'r').Mandatory()
		parse(tt.in)
		if s := checkError(tt.err); s != "" {
			t.Errorf("%s: %s", tt.name, s)
		}
	}
}

func TestGroup(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   []string
		err  string
	}{
		{
			name: "no args",
			in:   []string{"test"},
			err:  "test: exactly one of the following options must be specified: -A, -B",
		},
		{
			name: "one of each",
			in:   []string{"test", "-A", "-C"},
		},
		{
			name: "Two in group One",
			in:   []string{"test", "-A", "-B"},
			err:  "test: options -A and -B are mutually exclusive",
		},
		{
			name: "Two in group Two",
			in:   []string{"test", "-A", "-D", "-C"},
			err:  "test: options -C and -D are mutually exclusive",
		},
	} {
		reset()
		var val bool
		Flag(&val, 'o')
		Flag(&val, 'A').SetGroup("One")
		Flag(&val, 'B').SetGroup("One")
		Flag(&val, 'C').SetGroup("Two")
		Flag(&val, 'D').SetGroup("Two")
		RequiredGroup("One")
		parse(tt.in)
		if s := checkError(tt.err); s != "" {
			t.Errorf("%s: %s", tt.name, s)
		}
	}
}
