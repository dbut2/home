package url

import (
	"errors"
	"golang.org/x/net/publicsuffix"
	"strconv"
	"strings"
)

type URL struct {
	Protocol   string
	Username   string
	Password   string
	Subdomain  string
	Domain     string
	Suffix     string
	Port       string
	Directory  string
	Filename   string
	Parameters []Parameter
}

type Parameter struct {
	Key string
	Value string
}

func splitAt(s, substr string, i int) (found bool, a, b string) {
	found = i > -1
	if found {
		a = s[:i]
		b = s[i+len(substr):]
	} else {
		a = s
		b = s
	}
	return found, a, b
}

func split(s, substr string) (found bool, a, b string) {
	i := strings.Index(s, substr)
	return splitAt(s, substr, i)
}

func splitLast(s, substr string) (found bool, a, b string) {
	i := strings.LastIndex(s, substr)
	return splitAt(s, substr, i)
}

func New(s string) (url URL, err error) {
	found, protocol, s := split(s, "://")
	if found {
		url.Protocol = protocol
	}

	found, login, s := split(s, "@")
	if found {
		found, username, password := split(login, ":")
		url.Username = username
		if found {
			url.Password = password
		}
	}

	found, s, route := split(s, "/")
	if found {
		found, route, paramStr := split(route, "?")
		if found {
			params := strings.Split(paramStr, "&")
			for _, param := range params {
				found, key, val := split(param, "=")
				parameter := Parameter{}
				parameter.Key = key
				if found {
					parameter.Value = val
				}
			}
		}

		found, directory, filename := splitLast(route, "/")
		if found {
			url.Directory = directory + "/"
		}
		url.Filename = filename
	}

	found, s, port := split(s, ":")
	if found {
		if p, err := strconv.Atoi(port); err != nil || p < 0 || p > 65535 {
			return URL{}, errors.New("bad port number")
		}
		url.Port = port
	}

	suffix, _ := publicsuffix.PublicSuffix(s)
	i := strings.Index(s, suffix)
	if i == 0 {
		url.Domain = s
	} else {
		url.Suffix = "." + suffix
		s = s[:i-1]

		found, subdomain, domain := splitLast(s, ".")
		if found {
			url.Subdomain = subdomain
		}
		url.Domain = domain
	}

	return url, nil
}

func (u URL) String() string {
	url := ""

	if u.Protocol != "" {
		url += u.Protocol + "://"
	} else {
		url += "http://"
	}

	if u.Username != "" {
		url += u.Username
		if u.Password != "" {
			url += ":" + u.Password
		}
		url += "@"
	}

	if u.Subdomain != "" {
		url += u.Subdomain + "."
	}

	if u.Domain != "" {
		url += u.Domain
	}

	if u.Suffix != "" {
		suffix := u.Suffix
		suffix = strings.TrimPrefix(suffix, ".")
		url += "." + suffix
	}

	if u.Port != "" {
		port := u.Port
		port = strings.TrimPrefix(port, ":")
		url += ":" + port
	}

	route := ""

	if u.Directory != "" {
		directory := u.Directory
		directory = strings.TrimPrefix(directory, "/")
		directory = strings.TrimSuffix(directory, "/")
		route += directory + "/"
	}

	if u.Filename != "" {
		filename := u.Filename
		filename = strings.TrimPrefix(filename, "/")
		route += filename
	}

	if len(u.Parameters) > 0 {
		parameters := "?"
		for _, parameter := range u.Parameters {
			parameters += parameter.Key
			if parameter.Value != "" {
				parameters += "=" + parameter.Value
			}
			parameters += "&"
		}
		parameters = strings.TrimSuffix(parameters, "&")
		route += parameters
	}

	if route != "" {
		url += "/" + route
	}

	return url
}
