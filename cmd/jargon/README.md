This package is the command-line version of [Jargon](https://github.com/clipperhouse/jargon).

### Installation

```bash
go install github.com/clipperhouse/jargon/cmd/jargon
```

(Assumes a [Go installation](https://golang.org/dl/).)

### Usage

To display usage, simply type:

```bash
jargon
```

```
Usage: jargon accepts piped UTF8 text from tools such as cat, curl or echo, via Stdin

  Example: echo "I luv Rails" | jargon

Alternatively, use jargon 'standalone' by passing flags for text sources:

  -f string
    	A file path to lemmatize
  -s string
    	A (quoted) string to lemmatize
  -u string
    	A URL to fetch and lemmatize

  Example: jargon -f /path/to/file.txt

Results are piped to Stdout (regardless of input)
```
