package annotation

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type API struct {
	Controller 		string			`json:"controller"`
	Routes 			[]route 		`json:"routes"`
}

type route struct {
	Method 	string	`json:"method"`
	Route 	string	`json:"route"`
}

func TestFetchOk(t *testing.T) {
	var api []API
	err := Fetch("test/fetch", &api, _toApiFake)
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, a := range api {
		if len(a.Routes) != 1 {
			t.Errorf("cannot get %s's route", a.Controller)
		}
	}
}

func TestFetchWrongDirectory(t *testing.T) {
	var api []API
	err := Fetch("notexist/directory", &api, _toApiFake)
	if err == nil {
		t.Errorf("wrong directory don't return error")
	}
}

func TestFetchNotPointer(t *testing.T) {
	var api []API
	err := Fetch("test/fetch", api, _toApiFake)
	if err == nil {
		t.Errorf("not pointer don't return error")
	}
}

func TestFetchNotSlice(t *testing.T) {
	var api API
	err := Fetch("test/fetch", &api, _toApiFake)
	if err == nil {
		t.Errorf("not slice don't return error")
	}
}

func TestSaveOk(t *testing.T) {
	var api []API
	err := Fetch("test/fetch", &api, _toApiFake)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = Save(api, "./api.json")
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestSaveWrongData(t *testing.T) {
	x := map[string]interface{}{
		"foo": make(chan int),
	}
	err := Save(x, "./api.json")
	if err == nil {
		t.Errorf("save wrong data don't return error")
	}
}

func _toApiFake(rv reflect.Value, ar Result) (err error) {
	mapper := make(map[string]Result)
	for _, item := range ar {
		if mapper[item.Method] == nil {
			mapper[item.Method] = Result{}
		}
		mapper[item.Method] = append(mapper[item.Method], item)
	}
	for _, list := range mapper {
		var a API
		for _, item := range list {
			if a.Controller == "" {
				a.Controller = item.Method
			}
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