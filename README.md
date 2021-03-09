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

**Update** : Go 1.16 include a new embed system. Now, you can directly import your json file into your binary with
``go:embed`` directive, see more information [here](https://golang.org/pkg/embed/).

## Simple Example

Work in progress..

## Custom Parser

We want to build a simple cron parser for a dynamic cron system.
In this example, we use [robfig/cron](https://github.com/robfig/cron) package.

First, we will create our parser package and include in a go file.

The principle is simple, we have an array of structure and we will populate it with data get into
annotation.

So, your package need a ``structure`` to receive data and a ``parser's function`` for populate it.

You can design your structure like what you want, this is the goal of this package.

Your function need to contains two parameters.
In first parameter, she needs a ``reflect.Value``, this is your structure array,
and a ``annotation.Result`` that contains all found annotations.

The result of your function need to be an error, if you return an error,
your program can catch it from ``annotation.Fetch`'s method.

See the result :
```go
type Cron struct {
    Quartz string
    TZ     string
    Method string
}

func ToCron(rv reflect.Value, ar annotation.Result) (err error) {
	...
	return nil
}
```

After that, we need to use ``annotation.Result`` for getting data and parse it,
in firstly we can group each annotation for each method like that :

```go
func ToCron(rv reflect.Value, ar annotation.Result) (err error) {
    mapper := make(map[string]annotation.Result)
    
    for _, item := range ar {
        if mapper[item.Method] == nil {
            mapper[item.Method] = annotation.Result{}
        }
        mapper[item.Method] = append(mapper[item.Method], item)
    }
    
    return nil
}
```

We rebuild ``annotation.Result``'s parameter to a map with string key where the key is the
method name.

Next, we can loop into your new map for insert each annotation information to his method
and finally append structure into the array of structure :
```go
func ToCron(rv reflect.Value, ar annotation.Result) (err error) {
    ...
    
    for _, list := range mapper {
        var c Cron
        for _, item := range list {
            if c.Method == "" {
                c.Method = item.Method
            }
            err = c.insert(item.Key, item.Data)
            if err != nil {
                return
            }
        }
        rv.Set(reflect.Append(rv, reflect.ValueOf(a)))
    }
    
    return nil
}
```

You have probably seen the ``c.insert``'s function, this is simply a switch case based on
the annotation key that allows to use the right parsing process, see :

```go
func (c *Cron) insert(key string, data string) error {
	switch strings.ToLower(key) {
	    case "cron":
	    	return a._cron(data)
	    case "tz":
	    	return a._tz(data)
	}
	return nil
}

func (c *Cron) _cron(data string) error {
	c.Quartz = data
    return nil
}
```

All it's good, we can now use it into `annotation`'s package :

```go
func main() {
	...
	var c []cronParser.Cron
	
	err := annotation.Fetch("./internal/cron", &c, cronParser.ToCron)
	if err != nil {
		log.Fatal(err)
	}
	
	err = annotation.Save(c, "cron.json")
	if err != nil {
		log.Fatal(err)
	}
	...
}
```

Example of a controller :

```go
/*
    @cron("3 30 * * * *")
    @tz("Europe/Paris")
 */
func (Handler) Method(...) ... {
	...
}
```
