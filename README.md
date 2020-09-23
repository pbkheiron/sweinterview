# Prefix and infix calculator

## Running the code
Program is written in go programming langugage.

### Getting go environment
Please follow the instructions https://golang.org/doc/install to install go runtime.
Required version: go1.13+

There is no need to manually install any extra dependencies (go.mod contains list of required deps that will be 
automatically downloaded on first run).

### Console interactive mode
To run the the prefix notation calculator:
```
> go run main.go -n=prefix
> + 1 2
3
> + 1 * 2 3
7
(Ctrl+C to quit)
```

To run the infix notation calculator:
```
> go run main.go (or go run main.go -n=infix)
> 1 + 2
3
> ( 1 + 2 ) * 3
9
(Ctrl+C to quit)
```

### HTTP mode
HTTP uses POST with json for request and responses.  
For the endpoints like the calculator, GET would be a better method (caching, copy-pastable URLs, etc), 
but I didn't like the how the `expr` (expression) query parameter would look like.
Because of the URL encoding and because we use spaces (and maybe +) quite a lot in the input expression,
there would be many %20 and %2B in the query string, making it inconvenient to use by human (e.g. curl over command line).
A nice solution would be to implement more robust tokenizer that can generate tokens even
if the input is not space-separated.

To simplify understanding & running of the sample code (especially for people unfamiliar with go), 
I didn't use any HTTP frameworks.

Request:
```
{
    "expr":"( 1 + 2 ) * 5", 
    "notation":"infix"
}
```

Response:
```
{
    "result": ...
    "err": ... // empty if no error
}
```

Starting the server:
```
go run main.go -p=8080
```

Running sample calculation:
```
curl -X POST localhost:8080/calc -d '{"expr":"( 1 + 2 ) * 5", "notation":"infix"}'
```

## Code structure
```
.
├── calc
│   ├── expr.go         -- expression tree
│   ├── expr_test.go    
│   ├── infix.go        -- infix parser
│   ├── infix_test.go
│   ├── prefix.go       -- prefix parser
│   ├── prefix_test.go
│   └── tokenizer.go    -- tokenizer
├── go.mod
├── go.sum
├── main.go           
├── moretesting
│   └── assert.go       -- testing helpers
└── README.md
```

Evaluating expressions is split into 3 parts: tokenizer, parser and expression tree evaluation.
As I've built expression tree mechanism for infix expressions, I decided to reuse for prefix also 
(solution for prefix expression could be done more easily without expression tree).
 
## Tests
```
expr_test.go: test for expr tree evaluation if the expr tree is correctly build
infix_test.go: tests for infix expression evaluation (expression + expected result)
               tests for infix parser (expression + expected expression tree)
prefix_test.go: tests for prefix expression evaluaion (expression + expected result)
```