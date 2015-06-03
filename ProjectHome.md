The getopt package provides BSD and POSIX compatible getopt and getopt\_long functions.  It is as easy to use as the standard flag package but also allow special processing of each option as they are parsed.

Parsing is extremely close to BSD and getopt\_long with the POSIXLY\_CORRECT option.

While there are several other option parsing packages available, this package was written to be able to support the parsing of the ssh command line, as well as being very close to traditional command line parsing in UNIX.  Also, it was fun to write.