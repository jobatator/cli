package connexion

import (
	"reflect"
	"regexp"
)

// Options -
type Options struct {
	URL      string
	Auth     string
	Host     string
	Port     string
	Username string
	Password string
	Group    string
}

// ParseURL - get a formatted struct object from a URL
func ParseURL(url string) Options {
	var options Options
	exprStr := "^(?P<Auth>(?P<Username>[a-zA-Z_\\-0-9]+)(:(?P<Password>[a-zA-Z_\\-0-9]+))?@)?(?P<Host>[a-zA-Z\\-0-9.]+)(:(?P<Port>[0-9]{1,5}))?(\\/(?P<Group>[0-9a-zA-Z_\\-]+))?$"
	expr, err := regexp.Compile(exprStr)
	if err != nil {
		panic(err)
	}
	res := expr.FindStringSubmatch(url)
	names := expr.SubexpNames()
	for i := range res {
		if i != 0 {
			key := names[i]
			value := res[i]
			optionsReflect := reflect.ValueOf(&options).Elem()
			for field := 0; field < optionsReflect.NumField(); field++ {
				if optionsReflect.Type().Field(field).Name == key {
					optionsReflect.Field(field).SetString(value)
				}
			}
		}
	}
	options.URL = url
	if options.Host == "" {
		options.Host = "127.0.0.1"
	}
	if options.Port == "" {
		options.Port = "8962"
	}
	return options
}
