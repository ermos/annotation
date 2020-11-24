# Annotation 🎉
> Annotation implementation for Go

Annotation's library allows you to use annotations inside you're golang application, it is fully customizable, you can use you're own parser and build you're own annotation rule. You can use or find some examples in `parser` directory.

## Installation

```bash
go get github.com/ermos/annotation
```

## What does it look like? 🧐

```go
/*
	@Route("POST", "/auth/signin")
	@Desc("sign in to the application")
	@Response([200, 500])
	@Payload("username", string)
	@Payload("password", string)
	@?Payload("2fa", string)
  
  Allows users to get JSON Web Token :)
*/
func (Handler) SignIn() {
	// Logic here
}
```

Annotations can be recognized from others comments thanks to this `@` before each resource. In this example, we have six resource and one (useless?) comment. Each resource is recognize from their keys, for example, `@Route` is a key that contains `("POST", "/auth/signin")`. Data's key logic is decided by you're parser, always for this example, if you looking into `parser/api.go`, in the `(API)._route` method, you can find the regex that decides how you need to write the data.

## Usage

The most important thing to know before starting to use it, **annotations are generated from .go source file**. Like you can see, this can't be used in production, because you can't access to your source file, so who use it correctly? Firstly, we need to implement a build mode, after that, we can simply store annotations in a `json` file when we build the binary, with `go generate` for example, and use it with importing `json` file into our application on startup.

The inconvenient about this solution, you need to pass the `json` file on your production server too. You can solve it with embed you're `json` file into you're binary when you compile it. You have some awesome library to do that like [packr](https://github.com/gobuffalo/packr) or [pkger](https://github.com/markbates/pkger).

## Simple Example

Work in progress..

## Custom Parser

Work in progress..

## Contribution

Work in progress..
