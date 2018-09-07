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
Usage:

jargon accepts piped UTF8 text from Stdin and pipes lemmatized text to Stdout

  Example: echo "I luv Rails" | jargon

Alternatively, use jargon 'standalone' by passing flags for inputs and outputs:

  -f string
    	Input file path
  -o string
    	Output file path
  -s string
    	A (quoted) string to lemmatize
  -u string
    	A URL to fetch and lemmatize

  Example: jargon -f /path/to/original.txt -o /path/to/lemmatized.txt
```
