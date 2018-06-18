This package is the command-line version of [Jargon](https://github.com/clipperhouse/jargon).

### Command line

```bash
go install github.com/clipperhouse/jargon/cmd/jargon
```

To display usage, simply type:

```bash
jargon
```

Use `-f` to lemmatize a file and pipe to stdout:

```bash
jargon -f file.txt
```

If you’re dealing with large files, you might wish to pipe the results into another file

```bash
jargon -f file.txt > result.txt
```

Use `-s` to lemmatize a string and pipe to stdout

```bash
jargon -s "Here is a string with Ruby and SQL"
```

Use `-u` to fetch a URL and lemmatize, and pipe to stdout

```bash
jargon -u https://en.wikipedia.org/wiki/Programming_language
```
