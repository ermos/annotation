package parser

import (
"errors"
"github.com/ermos/annotation"
"reflect"
"regexp"
"strconv"
"strings"
)

type API struct {
	Controller 		string			`json:"controller"`
	Routes 			[]route 		`json:"routes"`
	Authorization	[]string		`json:"authorization,omitempty"`
	Middleware		middleware		`json:"middlewares,omitempty"`
	Response 		[]int			`json:"response,omitempty"`
	Version 		string			`json:"version,omitempty"`
	Description 	string			`json:"description,omitempty"`
	Validate 		validate		`json:"validate,omitempty"`
}

type route struct {
	Method 	string	`json:"method"`
	Route 	string	`json:"route"`
}

type middleware struct {
	Before 	[]string 	`json:"before,omitempty"`
	After 	[]string	`json:"after,omitempty"`
}

type validate struct {
	Params 	[]param		`json:"params,omitempty"`
	Payload	[]payload	`json:"payload,omitempty"`
	Queries	[]query		`json:"queries,omitempty"`
}

type param struct {
	Key 	string	`json:"key,omitempty"`
	Type 	string	`json:"type,omitempty"`
}

type payload struct {
	Key 		string	`json:"key,omitempty"`
	Type 		string	`json:"type,omitempty"`
	Description	string	`json:"description,omitempty"`
	Nullable	bool	`json:"nullable,omitempty"`
}

type query struct {
	Key 		string	`json:"key,omitempty"`
	Type 		string	`json:"type,omitempty"`
	Nullable	bool	`json:"nullable,omitempty"`
}

func ToAPI(rv reflect.Value, ar annotation.Result) (err error) {
	mapper := make(map[string]annotation.Result)
	for _, item := range ar {
		if mapper[item.Method] == nil {
			mapper[item.Method] = annotation.Result{}
		}
		mapper[item.Method] = append(mapper[item.Method], item)
	}
	for _, list := range mapper {
		var a API
		for _, item := range list {
			err = a.insert(item.Key, item.Data)
			if err != nil {
				return
			}
		}
		rv.Set(reflect.Append(rv, reflect.ValueOf(a)))
	}
	return nil
}

func (a *API) insert(key string, data string) error {
	switch strings.ToLower(key) {
	case "route":
		return a._route(data)
	case "auth":
		return a._auth(data)
	case "middlewarebefore":
		return a._middleware(data, true)
	case "middlewareafter":
		return a._middleware(data, false)
	case "desc":
		return a._description(data)
	case "param":
		return a._parameter(data)
	case "payload":
		return a._payload(data, false)
	case "?payload":
		return a._payload(data, true)
	case "query":
		return a._query(data, false)
	case "?query":
		return a._query(data, true)
	case "response":
		return a._response(data)
	case "version":
		return a._version(data)
	}
	return nil
}

func (a *API) _route(data string) error {
	rgx := regexp.MustCompile(`^(?:\"|\')(.*)(?:\"|\')(?:,)(?:\x09+| +|)(?:\"|\')(.*)(?:\"|\')$`)
	if !rgx.Match([]byte(data)) {
		return errors.New("incorrect route pattern")
	}
	res := rgx.FindStringSubmatch(data)
	a.Routes = append(a.Routes, route{
		Method: res[1],
		Route: res[2],
	})
	return nil
}

func (a *API) _auth(data string) error {
	rgx := regexp.MustCompile(`(?m)(?:"|')(.*?)(?:"|')`)
	for _, role := range rgx.FindAllString(data, -1) {
		a.Authorization = append(a.Authorization, role)
	}
	return nil
}

func (a *API) _middleware(data string, before bool) error {
	rgx := regexp.MustCompile(`(?m)(?:"|')(.*?)(?:"|')`)
	for _, mw := range rgx.FindAllString(data, -1) {
		if before {
			a.Middleware.Before = append(a.Middleware.Before, mw)
		} else {
			a.Middleware.After = append(a.Middleware.After, mw)
		}
	}
	return nil
}

func (a *API) _response(data string) error {
	rgx := regexp.MustCompile(`(?m)([0-9]+)`)
	for _, response := range rgx.FindAllString(data, -1) {
		responseInt, err := strconv.Atoi(response)
		if err != nil {
			return errors.New("responses value need to be an int")
		}
		a.Response = append(a.Response, responseInt)
	}
	return nil
}

func (a *API) _description(data string) error {
	data = strings.TrimPrefix(data, "\"")
	data = strings.TrimSuffix(data, "\"")
	a.Description = data
	return nil
}

func (a *API) _parameter(data string) error {
	rgx := regexp.MustCompile(`^(?:\"|\')(.*)(?:\"|\'),(?:\x09+| +|)(.*)$`)
	if !rgx.Match([]byte(data)) {
		return errors.New("incorrect document parameter pattern : " + data)
	}
	res := rgx.FindStringSubmatch(data)
	a.Validate.Params = append(a.Validate.Params, param{
		Key: res[1],
		Type: res[2],
	})
	return nil
}

func (a *API) _payload(data string, nullable bool) error {
	rgx := regexp.MustCompile(`^(?:\"|\')(.*)(?:\"|\'),(?:\x09+| +|)(.*)$`)
	if !rgx.Match([]byte(data)) {
		return errors.New("incorrect document field pattern : " + data)
	}
	res := rgx.FindStringSubmatch(data)
	a.Validate.Payload = append(a.Validate.Payload, payload{
		Key: res[1],
		Type: res[2],
		Nullable: nullable,
	})
	return nil
}

func (a *API) _query(data string, nullable bool) error {
	rgx := regexp.MustCompile(`^(?:\"|\')(.*)(?:\"|\'),(?:\x09+| +|)(.*)$`)
	if !rgx.Match([]byte(data)) {
		return errors.New("incorrect document field pattern : " + data)
	}
	res := rgx.FindStringSubmatch(data)
	a.Validate.Queries = append(a.Validate.Queries, query{
		Key: res[1],
		Type: res[2],
		Nullable: nullable,
	})
	return nil
}

func (a *API) _version(data string) error {
	rgx := regexp.MustCompile(`^("|'|)[0-9a-zA-Z.]+("|'|)$`)
	if !rgx.Match([]byte(data)) {
		return errors.New("incorrect version name pattern : " + data)
	}
	data = strings.TrimPrefix(data, "\"")
	data = strings.TrimSuffix(data, "\"")
	a.Version = data
	return nil
}
