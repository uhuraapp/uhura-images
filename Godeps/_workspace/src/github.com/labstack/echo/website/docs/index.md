# Echo

Build simple and performant systems!

---

## Overview

Echo is a fast HTTP router (zero memory allocation) and micro web framework in Go.

## Features

- Fast HTTP router which smartly resolves conflicting routes.
- Extensible middleware/handler, supports:
	- Middleware
		- `func(*echo.Context)`
		- `func(*echo.Context) error`
		- `func(echo.HandlerFunc) echo.HandlerFunc`
		- `func(http.Handler) http.Handler`
		- `http.Handler`
		- `http.HandlerFunc`
		- `func(http.ResponseWriter, *http.Request)`
		- `func(http.ResponseWriter, *http.Request) error`
	- Handler
		- `echo.HandlerFunc`
		- `func(*echo.Context) error`
		- `func(*echo.Context)`
		- `http.Handler`
		- `http.HandlerFunc`
		- `func(http.ResponseWriter, *http.Request)`
		- `func(http.ResponseWriter, *http.Request) error`
- Sub routing with groups.
- Handy encoding/decoding functions.
- Serve static files, including index.
- Centralized HTTP error handling.
- Use a customized function to bind request body to a Go type.
- Register a view render so you can use any HTML template engine.

## Getting Started

### Installation

```sh
$ go get github.com/labstack/echo
```

###[Hello, World!](https://github.com/labstack/echo/tree/master/examples/hello)

Create `server.go` with the following content

```go
package main

import (
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// Handler
func hello(c *echo.Context) {
	c.String(http.StatusOK, "Hello, World!\n")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger)

	// Routes
	e.Get("/", hello)

	// Start server
	e.Run(":4444")
}
```

`echo.New()` returns a new instance of Echo.

`e.Use(mw.Logger)` adds logging middleware to the chain. It logs every request made to the server,
producing output

```sh
2015/04/25 12:15:20 GET / 200 7.544µs
2015/04/25 12:15:26 GET / 200 3.681µs
2015/04/25 12:15:29 GET / 200 5.434µs
```

`e.Get("/", hello)` Registers a GET route for path `/` with hello handler, so
whenever server receives a request at `/`, hello handler is called.

In hello handler `c.String(http.StatusOK, "Hello, World!\n")` sends a text/plain
response to the client with 200 status code.

`e.Run(":4444")` Starts HTTP server at network address `:4444`.

Now start the server using command

```sh
$ go run server.go
```

Browse to [http://localhost:4444](http://localhost:4444) and you should see
Hello, World! on the page.

### Next?
- Browse [examples](https://github.com/labstack/echo/tree/master/examples)
- Head over to [Guide](guide.md)

## Contribute

**Use issues for everything**

- Report issues
- Discuss before sending pull request
- Suggest new features
- Improve/fix documentation

## Credits
- [Vishal Rana](https://github.com/vishr) - Author
- [Nitin Rana](https://github.com/nr17) - Consultant
- [Contributors](https://github.com/labstack/echo/graphs/contributors)

## License

[MIT](https://github.com/labstack/echo/blob/master/LICENSE)
