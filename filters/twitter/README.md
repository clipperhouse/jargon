This package provides a [jargon](https://github.com/clipperhouse/jargon) filter to identify Twitter-style hanshtags and handles, and combining them into a single token.

By default, the jargon tokenizer sees `@` and `#` as separate tokens. This filter looks for those tokens, followed by a token which meets Twitter's rules for legal handles and hastags. The result of two such tokens is a single token, i.e.:

"@" + "somename" → "@somename"
"#" + "sometag" → "#sometag"

### Usage

```go
tokens := jargon.Tokenize(reader)
twittered := tokens.Filter(twitter.Hashtags, twitter.Handles)
```