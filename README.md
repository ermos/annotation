# Annotation 🎉
> Annotation implementation for Go

Annotation's library allows you to use annotations inside you're golang application, it is fully customizable, you can use you're own parser and build you're own annotation rule. You can use or find some examples in `parser` directory.

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

Annotations can be recognize from other comments thanks to this `@` before each ressource. In this example, we have six ressource and one (useless?) comment.
Each ressource is recognize from their keys, for example, `@Route` is a key that contains `("POST", "/auth/signin")`. Data's key logic is decided by you're parser, always for this example, if you looking into `parser/api.go`, in the `(API)._route` method, you can find the regex that decide how you need to write the data.

## Usage

Work in progress..

## Custom Parser

Work in progress..

## Contribution

Work in progress..
