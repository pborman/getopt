# getopt

Package getopt provides traditional getopt processing for implementing
commands that use traditional command lines.  The standard Go flag package
cannot be used to write a program that parses flags the way ls or ssh does,
for example.

Getopt supports functionality found in both the standard BSD getopt
as well as (one of the many versions of) the GNU getopt_long.  Being a Go
package, this package makes common usage easy, but still enables more
controlled usage if needed.

## Typical Usage

```
	// Declare the flags to be used
	helpFlag := getopt.Bool('?', "display help")
	cmdFlag := getopt.StringLong("command", 'c', "", "the command)

	func main() {
		// Parse the program arguments
		getopt.Parse()
		// Get the remaining positional parameters
		args := getopt.Args()
		...
	}
```

## Option Formats

Support is provided for both short (-f) and long (--flag) options.  A single
option may have both a short and a long name.  Each option may be a flag
or a value.  A value takes an argument.

Declaring no long names causes this package to process arguments like the
traditional BSD getopt.

Short flags may be combined into a single parameter.  For example, "-a -b -c"
may also be expressed "-abc".  Long flags must stand on their own "--alpha
--beta"

Values require an argument.  For short options the argument may either be
immediately following the short name or as the next argument.  Only one short
value may be combined with short flags in a single argument; the short value
must be after all short flags.  For example, if f is a flag and v is a
value, then:

```
 -vvalue    (sets v to "value")
 -v value   (sets v to "value")
 -fvvalue   (sets f, and sets v to "value")
 -fv value  (sets f, and sets v to "value")
 -vf value  (set v to "f" and value is the first parameter)
```

For the long value option val:

```
 --val value (sets val to "value")
 --val=value (sets val to "value")
 --valvalue  (invalid option "valvalue")
```

### Optional Values

Values with an optional value only set the value if the value is part of
the same argument.  In any event, the option count is increased and the
option is marked as seen.

```
 -v -f          (sets v and f as being seen)
 -vvalue -f     (sets v to "value" and sets f)
 --val -f       (sets v and f as being seen)
 --val=value -f (sets v to "value" and sets f)
```

There is no convience function defined for making the value optional.
The SetOptional method must be called on the actual Option.

```
 v := String("val", 'v', "", "the optional v")
 Lookup("v").SetOptional()

 var s string
 StringVar(&s, "val", 'v', "the optional v).SetOptional()
```

### Parsing

Parsing continues until the first non-option or "--" is encountered.

The short name "-" can be used, but it either is specified as "-" or as
part of a group of options, for example "-f-".  If there are no long
options specified then "--f" could also be used.  If "-" is not declared
as an option then the single "-" will also terminate the option processing
but unlike "--", the "-" will be part of the remaining arguments.

Normally the parsing is performed by calling the Parse function.  If it is
important to see the order of the options then the Getopt function should
be used.  The standard Parse function does the equivalent of:

```
func Parse() {
	if err := getopt.Getopt(os.Args, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		s.PrintUsage(os.Stderr)
		os.Exit(1)
	}
}
```

When calling Getopt it is the responsibility of the caller to print
any errors.

Normally the default option set, CommandLine, is used.  Other option sets
may be created with New.

After parsing, the sets Args will contain the non-option arguments.  If an
error is encountered then Args will begin with argument that caused the
error.

It is valid to call a set's Parse a second time to ammend the current set of
flags or values.  As an example:

```
 var a = getopt.Bool('a', "", "The a flag")
 var b = getopt.Bool('b', "", "The a flag")
 var cmd = ""

 var opts = getopt.CommandLine

 opts.Parse(os.Args)
 if opts.NArgs() > 0 {
     cmd = opts.Arg(0)
     opts.Parse(opts.Args())
 }
 ```

If called with set to { "prog", "-a", "cmd", "-b", "arg" } then both
and and b would be set, cmd would be set to "cmd", and opts.Args()
would return { "arg" }.

Unless an option type explicitly prohibits it, an option may appear more
than once in the arguments.  The last value provided to the option is
the value.

## SYNTAX

For each option type there are an unfortunately large number of ways, 8,
to initialize the option.  This number is derived from three attributes:

  1.  Short or Long name
  2.  Normal vs Var
  3.  Command Line vs Option Set

The first two variations provide 4 signature:

```
 Option(name rune, [value type,]  helpvalue... string)
 OptionLong(name string, short rune, [value type,]  helpvalue... string)
 OptionVar(p *type, name rune, helpvalue... string)
 OptionVarLong(p *type, name string, short rune, helpvalue... string)
```

Foo can actually be expressed in terms of FooLong:

```
 func Foo(name rune, value type, helpvalue... string) *type {
     return FooLong("", name, value, helpvalue...)
 }
```

Normally Foo is used, unless long options are needed.  Setting short to 0
creates only a long option.

The difference bentween Foo and FooVar is that you pass a pointer, p, to
the location of the value to FooVar.  The default value is simply *p.
The initial value of *p is the defaut value of the option.

Foo is actually a wrapper around FooVar:

```
 func Foo(name rune, value type, helpvalue... string) *type {
     p := value
     FooVar(&p, name, helpvalue... string)
     return &p
 }
```


The third variation provides a top-level function and a method on a Set:

```
 func Option(...)
 func (s *Set) Option(...)
```

The top-level function is simply:

```
 func Option(...) *type {
     return CommandLine.Option(...) {
 }
```

To simplfy documentation, typically only the main top-level function is
fully documented.  The others will have documentation when there is
something special about them.

## VALUEHELP

All non-flag options are created with a "valuehelp" as the last parameter.
Valuehelp should be 0, 1, or 2 strings.  The first string, if provided,
is the usage message for the option.  If the second string, if provided,
is the name to use for the value when displaying the usage.  If not
provided the term "value" is assumed.

The usage message for the option created with

```
 StringLong("option", 'o', "defval", "a string of letters")
```

is

```
 -o, -option=value
```

And for

```
 StringLong("option", 'o', "defval", "a string of letters", "string")
```

is

```
 -o, -option=string
```
