# `./combinator 🧩`

Parser combinators for building recursive descent parsers

## `features`

- Zero dependencies
- Full Unicode support
- Line/column tracking for error messages
- Recursive grammars with `Rule` and `Ref`
- Expression parsing with `ChainL1`/`ChainR1`
- Composable: small parsers combine into larger ones

## `install`

```bash
go get github.com/dottermi/x/combinator@latest
```

## `usage`

```go
package main

import (
    "fmt"
    "github.com/dottermi/x/combinator"
)

func main() {
    // Parse a comma-separated list of integers
    parser := combinator.SepBy(
        combinator.Integer(),
        combinator.Char(','),
    )

    result := combinator.Parse(parser, "1,2,3")

    if result.OK {
        fmt.Println(result.Value) // [1 2 3]
    } else {
        fmt.Println(result.Err)
    }
}
```

## `api`

### Primitives

<details>
<summary><code>Char(r rune)</code> - matches a single specific character</summary>

```go
result := combinator.Parse(combinator.Char('a'), "abc")
// result.Value == 'a'
```
</details>

<details>
<summary><code>String(s string)</code> - matches an exact string</summary>

```go
result := combinator.Parse(combinator.String("hello"), "hello world")
// result.Value == "hello"
```
</details>

<details>
<summary><code>Satisfy(pred func(rune) bool)</code> - matches a character that satisfies the predicate</summary>

```go
vowel := combinator.Satisfy(func(r rune) bool {
    return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u'
})
result := combinator.Parse(vowel, "apple")
// result.Value == 'a'
```
</details>

<details>
<summary><code>Any()</code> - matches any single character</summary>

```go
result := combinator.Parse(combinator.Any(), "xyz")
// result.Value == 'x'
```
</details>

<details>
<summary><code>EOF()</code> - matches end of input</summary>

```go
complete := combinator.Seq(combinator.String("end"), combinator.EOF())
result := combinator.Parse(complete, "end")
// result.OK == true

result = combinator.Parse(complete, "end!")
// result.OK == false
```
</details>

<details>
<summary><code>OneOf(chars string)</code> - matches any character in the string</summary>

```go
op := combinator.OneOf("+-*/")
result := combinator.Parse(op, "+")
// result.Value == '+'
```
</details>

<details>
<summary><code>NoneOf(chars string)</code> - matches any character NOT in the string</summary>

```go
notQuote := combinator.NoneOf("\"")
result := combinator.Parse(notQuote, "a")
// result.Value == 'a'
```
</details>

<details>
<summary><code>Range(from, to rune)</code> - matches any character in the inclusive range</summary>

```go
lowercase := combinator.Range('a', 'z')
result := combinator.Parse(lowercase, "m")
// result.Value == 'm'
```
</details>

### Characters

<details>
<summary><code>Digit()</code> - matches a Unicode digit (0-9)</summary>

```go
result := combinator.Parse(combinator.Digit(), "5")
// result.Value == '5'
```
</details>

<details>
<summary><code>Letter()</code> - matches a Unicode letter</summary>

```go
result := combinator.Parse(combinator.Letter(), "A")
// result.Value == 'A'
```
</details>

<details>
<summary><code>Space()</code> - matches a single whitespace character</summary>

```go
result := combinator.Parse(combinator.Space(), " ")
// result.Value == ' '
```
</details>

<details>
<summary><code>Spaces()</code> - matches zero or more whitespace characters</summary>

```go
result := combinator.Parse(combinator.Spaces(), "   hello")
// result.Value == "   "
```
</details>

<details>
<summary><code>Spaces1()</code> - matches one or more whitespace characters</summary>

```go
result := combinator.Parse(combinator.Spaces1(), "  \t\n")
// result.Value == "  \t\n"
```
</details>

<details>
<summary><code>Alpha()</code> - matches an ASCII letter (a-z, A-Z)</summary>

```go
result := combinator.Parse(combinator.Alpha(), "Z")
// result.Value == 'Z'
```
</details>

<details>
<summary><code>AlphaNum()</code> - matches a letter or digit</summary>

```go
result := combinator.Parse(combinator.AlphaNum(), "7")
// result.Value == '7'
```
</details>

<details>
<summary><code>Lower()</code> - matches a lowercase letter</summary>

```go
result := combinator.Parse(combinator.Lower(), "a")
// result.Value == 'a'
```
</details>

<details>
<summary><code>Upper()</code> - matches an uppercase letter</summary>

```go
result := combinator.Parse(combinator.Upper(), "A")
// result.Value == 'A'
```
</details>

<details>
<summary><code>HexDigit()</code> - matches a hexadecimal digit (0-9, a-f, A-F)</summary>

```go
result := combinator.Parse(combinator.HexDigit(), "F")
// result.Value == 'F'
```
</details>

<details>
<summary><code>OctDigit()</code> - matches an octal digit (0-7)</summary>

```go
result := combinator.Parse(combinator.OctDigit(), "7")
// result.Value == '7'
```
</details>

<details>
<summary><code>BinDigit()</code> - matches a binary digit (0 or 1)</summary>

```go
result := combinator.Parse(combinator.BinDigit(), "1")
// result.Value == '1'
```
</details>

<details>
<summary><code>Newline()</code> - matches a newline character</summary>

```go
result := combinator.Parse(combinator.Newline(), "\n")
// result.Value == '\n'
```
</details>

<details>
<summary><code>Tab()</code> - matches a tab character</summary>

```go
result := combinator.Parse(combinator.Tab(), "\t")
// result.Value == '\t'
```
</details>

<details>
<summary><code>CRLF()</code> - matches Windows line ending "\r\n"</summary>

```go
result := combinator.Parse(combinator.CRLF(), "\r\n")
// result.Value == "\r\n"
```
</details>

<details>
<summary><code>EndOfLine()</code> - matches Unix or Windows line endings</summary>

```go
result := combinator.Parse(combinator.EndOfLine(), "\n")
// result.Value matches newline or CRLF
```
</details>

### Combinators

<details>
<summary><code>Seq(parsers ...Parser)</code> - runs parsers in sequence</summary>

```go
ab := combinator.Seq(combinator.Char('a'), combinator.Char('b'))
result := combinator.Parse(ab, "ab")
// result.Value == []any{'a', 'b'}
```
</details>

<details>
<summary><code>Choice(parsers ...Parser)</code> - tries parsers in order, returns first success</summary>

```go
boolean := combinator.Choice(combinator.String("true"), combinator.String("false"))
result := combinator.Parse(boolean, "true")
// result.Value == "true"
```
</details>

<details>
<summary><code>Many(p Parser)</code> - matches zero or more occurrences</summary>

```go
digits := combinator.Many(combinator.Digit())
result := combinator.Parse(digits, "123abc")
// result.Value == []any{'1', '2', '3'}
```
</details>

<details>
<summary><code>Many1(p Parser)</code> - matches one or more occurrences</summary>

```go
digits := combinator.Many1(combinator.Digit())
result := combinator.Parse(digits, "123")
// result.Value == []any{'1', '2', '3'}

result = combinator.Parse(digits, "abc")
// result.OK == false
```
</details>

<details>
<summary><code>Opt(p Parser)</code> - makes a parser optional</summary>

```go
sign := combinator.Opt(combinator.Char('-'))
result := combinator.Parse(sign, "42")
// result.Value == nil

result = combinator.Parse(sign, "-42")
// result.Value == '-'
```
</details>

<details>
<summary><code>And(p1, p2 Parser)</code> - sequences two parsers, returns second result</summary>

```go
skipThenDigit := combinator.And(combinator.Char('_'), combinator.Digit())
result := combinator.Parse(skipThenDigit, "_5")
// result.Value == '5'
```
</details>

<details>
<summary><code>Left(p1, p2 Parser)</code> - sequences two parsers, returns first result</summary>

```go
num := combinator.Left(combinator.Integer(), combinator.Char(';'))
result := combinator.Parse(num, "42;")
// result.Value == int64(42)
```
</details>

<details>
<summary><code>Right(p1, p2 Parser)</code> - sequences two parsers, returns second result</summary>

```go
num := combinator.Right(combinator.Spaces(), combinator.Integer())
result := combinator.Parse(num, "   42")
// result.Value == int64(42)
```
</details>

<details>
<summary><code>Count(n int, p Parser)</code> - matches exactly n occurrences</summary>

```go
hex := combinator.Count(2, combinator.HexDigit())
result := combinator.Parse(hex, "FF")
// result.Value == []any{'F', 'F'}
```
</details>

### Transform

<details>
<summary><code>Map(p Parser, fn func(any) any)</code> - transforms the result</summary>

```go
upper := combinator.Map(combinator.Letter(), func(v any) any {
    return unicode.ToUpper(v.(rune))
})
result := combinator.Parse(upper, "a")
// result.Value == 'A'
```
</details>

<details>
<summary><code>MapErr(p Parser, fn func(error) error)</code> - transforms the error message</summary>

```go
digit := combinator.MapErr(combinator.Digit(), func(err error) error {
    return fmt.Errorf("digit required: %w", err)
})
result := combinator.Parse(digit, "x")
// result.Err contains "digit required: ..."
```
</details>

<details>
<summary><code>Label(p Parser, label string)</code> - adds a descriptive name to errors</summary>

```go
digit := combinator.Label(combinator.Range('0', '9'), "digit")
result := combinator.Parse(digit, "x")
// result.Err contains "expected digit: ..."
```
</details>

<details>
<summary><code>Skip(p Parser)</code> - runs parser but discards result</summary>

```go
skipComment := combinator.Skip(combinator.String("//"))
result := combinator.Parse(skipComment, "//")
// result.Value == nil
```
</details>

<details>
<summary><code>SkipMany(p Parser)</code> - matches zero or more, discards all</summary>

```go
skipSpaces := combinator.SkipMany(combinator.Space())
result := combinator.Parse(skipSpaces, "   ")
// result.Value == nil
```
</details>

<details>
<summary><code>SkipMany1(p Parser)</code> - matches one or more, discards all</summary>

```go
skipSpaces := combinator.SkipMany1(combinator.Space())
result := combinator.Parse(skipSpaces, "   ")
// result.Value == nil
```
</details>

<details>
<summary><code>Not(p Parser)</code> - inverts result (negative lookahead)</summary>

```go
notKeyword := combinator.Seq(combinator.Not(combinator.String("if")), combinator.Ident())
result := combinator.Parse(notKeyword, "iffy")
// result.OK == false (starts with "if")

result = combinator.Parse(notKeyword, "foo")
// result.OK == true
```
</details>

<details>
<summary><code>LookAhead(p Parser)</code> - tests without consuming input</summary>

```go
peek := combinator.LookAhead(combinator.Digit())
result := combinator.Parse(peek, "5abc")
// result.Value == '5', but position unchanged
```
</details>

<details>
<summary><code>Lazy(p *Parser)</code> - defers parser evaluation for mutual recursion</summary>

```go
var p combinator.Parser
p = combinator.Choice(
    combinator.Char('a'),
    combinator.Seq(combinator.Char('('), combinator.Lazy(&p), combinator.Char(')')),
)
result := combinator.Parse(p, "(a)")
// result.OK == true
```
</details>

<details>
<summary><code>Ref(r *Rule)</code> - creates parser from Rule for recursive grammars</summary>

```go
var expr combinator.Rule
expr = func() combinator.Parser {
    return combinator.Choice(
        combinator.Integer(),
        combinator.Parens(combinator.Ref(&expr)),
    )
}
result := combinator.Parse(combinator.Ref(&expr), "((42))")
// result.Value == int64(42)
```
</details>

### Delimiters

<details>
<summary><code>Between(open, close, p Parser)</code> - matches content between delimiters</summary>

```go
quoted := combinator.Between(
    combinator.Char('"'),
    combinator.Char('"'),
    combinator.Many(combinator.NoneOf("\"")),
)
result := combinator.Parse(quoted, "\"hello\"")
// result.Value == []any{'h', 'e', 'l', 'l', 'o'}
```
</details>

<details>
<summary><code>Parens(p Parser)</code> - matches content in parentheses</summary>

```go
expr := combinator.Parens(combinator.Integer())
result := combinator.Parse(expr, "(42)")
// result.Value == int64(42)
```
</details>

<details>
<summary><code>Braces(p Parser)</code> - matches content in curly braces</summary>

```go
block := combinator.Braces(combinator.Ident())
result := combinator.Parse(block, "{foo}")
// result.Value == "foo"
```
</details>

<details>
<summary><code>Brackets(p Parser)</code> - matches content in square brackets</summary>

```go
index := combinator.Brackets(combinator.Integer())
result := combinator.Parse(index, "[0]")
// result.Value == int64(0)
```
</details>

<details>
<summary><code>Angles(p Parser)</code> - matches content in angle brackets</summary>

```go
tag := combinator.Angles(combinator.Ident())
result := combinator.Parse(tag, "<html>")
// result.Value == "html"
```
</details>

<details>
<summary><code>SepBy(p, sep Parser)</code> - matches zero or more separated by delimiter</summary>

```go
items := combinator.SepBy(combinator.Integer(), combinator.Char(','))
result := combinator.Parse(items, "1,2,3")
// result.Value == []any{int64(1), int64(2), int64(3)}
```
</details>

<details>
<summary><code>SepBy1(p, sep Parser)</code> - matches one or more separated by delimiter</summary>

```go
items := combinator.SepBy1(combinator.Ident(), combinator.Char(','))
result := combinator.Parse(items, "a,b,c")
// result.Value == []any{"a", "b", "c"}
```
</details>

<details>
<summary><code>EndBy(p, end Parser)</code> - matches zero or more each followed by terminator</summary>

```go
statements := combinator.EndBy(combinator.Ident(), combinator.Char(';'))
result := combinator.Parse(statements, "a;b;c;")
// result.Value == []any{"a", "b", "c"}
```
</details>

<details>
<summary><code>EndBy1(p, end Parser)</code> - matches one or more each followed by terminator</summary>

```go
statements := combinator.EndBy1(combinator.Ident(), combinator.Char(';'))
result := combinator.Parse(statements, "foo;bar;")
// result.Value == []any{"foo", "bar"}
```
</details>

### Tokens

<details>
<summary><code>Ident()</code> - matches a programming language identifier</summary>

```go
result := combinator.Parse(combinator.Ident(), "myVar123")
// result.Value == "myVar123"
```
</details>

<details>
<summary><code>Keyword(kw string)</code> - matches keyword with word boundary</summary>

```go
result := combinator.Parse(combinator.Keyword("if"), "if (x)")
// result.Value == "if"

result = combinator.Parse(combinator.Keyword("if"), "iffy")
// result.OK == false (no word boundary)
```
</details>

<details>
<summary><code>Integer()</code> - matches an optionally negative integer</summary>

```go
result := combinator.Parse(combinator.Integer(), "-42")
// result.Value == int64(-42)
```
</details>

<details>
<summary><code>Float()</code> - matches a decimal number</summary>

```go
result := combinator.Parse(combinator.Float(), "-3.14")
// result.Value == float64(-3.14)
```
</details>

<details>
<summary><code>StringLit()</code> - matches a double-quoted string with escapes</summary>

```go
result := combinator.Parse(combinator.StringLit(), `"hello \"world\""`)
// result.Value == `hello "world"`
```
</details>

<details>
<summary><code>CharLit()</code> - matches a single-quoted character literal</summary>

```go
result := combinator.Parse(combinator.CharLit(), `'a'`)
// result.Value == 'a'
```
</details>

### Expressions

<details>
<summary><code>Lexeme(p Parser)</code> - parses and consumes trailing whitespace</summary>

```go
num := combinator.Lexeme(combinator.Integer())
two := combinator.Seq(num, num)
result := combinator.Parse(two, "42   17")
// result.Value == []any{int64(42), int64(17)}
```
</details>

<details>
<summary><code>Symbol(s string)</code> - matches string and consumes trailing whitespace</summary>

```go
result := combinator.Parse(combinator.Symbol("if"), "if   ")
// result.Value == "if"
```
</details>

<details>
<summary><code>Token()</code> - matches identifier and consumes trailing whitespace</summary>

```go
result := combinator.Parse(combinator.Token(), "variableName   ")
// result.Value == "variableName"
```
</details>

<details>
<summary><code>IntToken()</code> - matches integer and consumes trailing whitespace</summary>

```go
result := combinator.Parse(combinator.IntToken(), "123   ")
// result.Value == int64(123)
```
</details>

<details>
<summary><code>FloatToken()</code> - matches float and consumes trailing whitespace</summary>

```go
result := combinator.Parse(combinator.FloatToken(), "3.14   ")
// result.Value == float64(3.14)
```
</details>

<details>
<summary><code>StringToken()</code> - matches string literal and consumes trailing whitespace</summary>

```go
result := combinator.Parse(combinator.StringToken(), `"hello"   `)
// result.Value == "hello"
```
</details>

<details>
<summary><code>ChainL1(p, op Parser)</code> - parses left-associative binary expressions</summary>

```go
addOp := combinator.Map(combinator.Char('+'), func(_ any) any {
    return func(a, b any) any {
        return a.(int64) + b.(int64)
    }
})
expr := combinator.ChainL1(combinator.Integer(), addOp)
result := combinator.Parse(expr, "1+2+3")
// result.Value == int64(6), computed as ((1+2)+3)
```
</details>

<details>
<summary><code>ChainR1(p, op Parser)</code> - parses right-associative binary expressions</summary>

```go
powOp := combinator.Map(combinator.Char('^'), func(_ any) any {
    return func(a, b any) any {
        return math.Pow(float64(a.(int64)), float64(b.(int64)))
    }
})
expr := combinator.ChainR1(combinator.Integer(), powOp)
result := combinator.Parse(expr, "2^3^2")
// result.Value == 512.0, computed as 2^(3^2)
```
</details>

## `license`

MIT
