

# Incertus

Probabilities expression calculator with support for throwing dice.

It takes a mathematical expression and returns the frequency of each result value.

```
<result-0>: <frequency-0>
<result-1>: <frequency-1>
<result-2>: <frequency-2>
<result-3>: <frequency-3>
...
<result-N>: <frequency-N>
```


## Command line interface

`incertus [options...] [expression|--from=<file>|--from=-]`

Help:
  - `--help`/`-h`: Print a global help.
  - `--help-function`: Print a list of registered functions.
  - `--help-function=<function-name>`: Print description and usage for a registered function.


Input options:
- `expression`: If provided must be the last argument.
- `--from=<file>` / `-f`: To read the expression(s) from a file.
- `--from=-`: To read the expression(s) from the standard input.

Line options:
- `--single-expression` _(default if `--expression-separator` NOT provided)_: Treats the input as one single expression that can encompass multiple lines.
- `--multiple-expressions` _(default if `--expression-separator` provided)_: Threats the input as one expression per line (see _Line split option_)

Line split option:
- `--expression-separator`: When using `--multiple-expressions`, this string is the separator between expressions.
  - Newline by default (`\n`)
  - This value will be validated to avoid conflicts with valid expressions or operators.
  - This validation can be skipped using `--force-expression-separator` instead.
  - _IMPORTANT: Making this value conflict with valid expressions or operators can cause unexpected results._

Behavior policies to apply when dividing by zero:
- `--on-division-zero=[remove|zero|<N>]` and `--on-modulo-zero=[remove|zero|<N>]`
  - `remove` _(default)_: Don't count those undefined results towards the evaluation.
  - `zero`: Count their frequency as if they were `0`.
  - `<N>`: Count their frequency as if they were an arbitrary integer `N`.

Decimal rounding:
- `--round=[no|<N>]`: Decimal positions to round to. Applied after normalization.

Result separator:
- `--results-separator`:
    String to use when evaluating multiple expressions to indicate where the results of one ends and the results of the next begin.
    It does not get included before the first one nor after the last one.
  - When `-o=text`, the default value is `\n---\n`
  - When `-o=json`, the default value is `\n`.
  - _Note: A newline (`\n`) is not included by default, it must be explicitly provided._


Output format:
- `--output=[text|json|csv]` (also `-o`): Determines the output format of the evaluation results.
  - `text` _(default)_: One line per result like: `<result>:\t<frequency>`
  - `json`: Single line JSON object, where the keys are the results and the values are their frequencies.
  - `csv`: One line per result like: `<result>,<frequency>`, with a single header line `value,frequency`.

Output format modifiers:
- `-o-sort=[no|value|freq]` _(default=value)_: Prints the results sorted by value, frequency or not at all.
- `-o-order=[asc|desc]` _(default=asc)_: When `-o-sort` option active, orders the results on ascending or descending order.
- `-o-string-values=[true|yes|false|no]` _(default=false)_: Prints the value as a string when applicable based on `--output`. This usually means surrounding the value number with quotes `"`.
- `-o-normalized=[no|1|100]`: Normalize the frequencies. This can output decimal numbers.
  - `no` means absolute frequencies
  - `1` means normalized
  - `100` _(default)_ means normalized as percentage (normalized then multiplied by 100)


Debug:
- `--timeit`/`-t`: Time execution phases: lexing, parsing, evaluation and formatting.
  - Prints a report at the end.


## Some Examples


### Common operations

These are standard mathematical operations on simple integer numbers.
```
1
1 + 1
2 * 3
4 + 5 * 6
4 * (5 + 6)
(4 - (5 + 6)) * 7
1 + -2
```

Their result always have a single entry with frequency equal to 1: `<result>: 1`

When used this way, it behaves as a standard calculator program.

### Statistical operations

These are some of the core functions provided. A full list of them, their usage and description can be obtained by running with the flag `--help-function`


Aggregation:
```
# Calculate the mean of a frequency set
# Returns a single value with frequency 1
stats.mean(...)

# Calculate the median of a frequency set
# Returns a single value with frequency 1
stats.median(...)
```


Dice simulation:
```
# Throw one 6-sided dice
dice.sum(6)

# Throw 2 6-sided dice and add their results
dice.sum(2, 6)

# Keep the highest between 10 and a 20-sided dice throw
dice.best(10, dice(20))

# Throw 2 20-sided dice and keep the highest between them
dice.best(dice(20), dice(20))

# Throw 3 20-sided dice and keep the highest between them
dice.best(dice(20), dice.best(dice(20), dice(20)))
```


## The expression language

### Supported operators and functions

We present all the supported operators ordered by precedence.

> On operator grammar:  
> Most operators: `<expr> <operator> <expr>`  
> Unary operators: `<operator> <expr>`  
> Functions: `<func>(<expr>, ...)`  


**Numeric operators**:
- Unary `-`
- `*` multiplication, `//` integer division, `%` modulo
- `+` addition, `-` subtraction
- `<<` Bitwise left shift, `>>` bitwise right shift


**Comparison operators**. They output `1` for true and `0` for false, enabling subsequent operations:
- `==` equal, `!=` not equal, `>` greater, `<` smaller, `>=` greater or equal, `<=` smaller or equal


**Boolean operators**. They output `1` for true and `0` for false, enabling subsequent operations:
- `&&` and, `||` or


**Functions**. They all have the same precedence.

The 2 standard addons include the following functions:  
- **Dice simulation**
  - `dice.sum(N)` with one argument: Populates the frequency of every number between `1` and `N` with `1`.
  - `dice.sum(D, N)` with two arguments: The equivalent to throwing `D` `N`-sided dice and adding their results.
  - `dice.best(A, B)`: Keeping the best between `A` and `B`.
  - `dice.worst(D, N)`: Keeping the worst between `A` and `B`.

> Note: `dice.best` and `dice.worst` are typically used with `dice` like:  

- **Basic statistics**
  - `mean(N)`: Generate a single result with frequency `1` like: `<mean>: 1`
  - `median(N)`: Generate a single result with frequency `1` like: `<median>: 1`


> Note: On any of the provided functions, they evaluate to an empty frequency set if they receive an empty one as argument.


### Division and modulo considerations

Both division (`//`) and modulo `%` accept and generate integers.

There is a common case when dividing a number by zero (`0`) that is undefined.

Behavior policies for this cases are explained in the [Command Line Interface](#command-line-interface) section.


### Identifiers

These are mostly used as function names. All the default ones fulfill these requirements and all added ones must do too.

- They are case-sensitive
- They start with a letter (all ASCII letters between a and z)
- The rest of their body may contain letters, and the characters: '.' and '_'
- Numbers are not allowed. _Note: this might change on further revisions._


### Numbers

Only integers are allowed.
At any point during the evaluation, the range of numbers supported are those representable by a 64 bit integer, including negative numbers.


### Grammar

Heavily inspired by [Crafting Interpreters](https://craftinginterpreters.com).
See [Representing code](https://craftinginterpreters.com/representing-code.html) and [Parsing Expressions](https://craftinginterpreters.com/parsing-expressions.html)

Unambiguous, right-recursion grammar for recursive descent:
```
expression -> boolean

boolean ->
  | comparison
  | comparison "&&" boolean
  | comparison "||" boolean

comparison ->
  | bitop
  | bitop "==" comparison
  | bitop "!=" comparison
  | bitop ">=" comparison
  | bitop "<=" comparison
  | bitop ">" comparison
  | bitop "<" comparison

bitop ->
  | term
  | term "<<" bitop
  | term ">>" bitop

term ->
  | factor
  | factor "+" term
  | factor "-" term

factor ->
  | unary
  | unary "*" factor
  | unary "//" factor
  | unary "%" factor

unary ->
  | "-" unary
  | primary


primary ->
  | NUMBER
  | "(" expression ")"
  | IDENTIFIER call

call ->
  | "(" ")"
  | "(" arguments ")"


arguments ->
  | expression
  | expression "," arguments

```

Recursive descent chosen for being simple, robust and not particularly slow.
It also maps really easily from grammar to implementation.



## Future work

### Embedded frequency sets

Include as a valid expression: `{<N>: <F>, ...}`, where `<N>` is a number and `<F>` is its frequency.

During execution this should be treated like a number, because `<N>` gets transformed to `{<N>: 1}` by the parser.

Examples:
```
{}
{1: 23}
{1: 23, 4: 56}
```


### Embedded lists

Include as a valid expression `[<N>, ...]`, where `<N>` is a number. A frequency set can be obtained by counting the absolute frequency of each number.

During execution this should be treated like a number, because `<N>` gets transformed to `{<N>: 1}` by the parser.

Examples:
```
[]              => {}
[1]             => {1: 1}
[1, 2, 3]       => {1: 1, 2: 1, 3: 1}
[2, 1, 2, 2, 2] => {1: 1, 2: 4}
```
