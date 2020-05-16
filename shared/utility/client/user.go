//go:generate gobox tools/easymock

package client

import (
	"html/template"
	"strings"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ignition/shared/utility/coalesce"
)

type UserClientData struct {
	Host           string
	OriginalHost   string
	Ip             string
	OriginalIp     string
	Schema         string
	OriginalSchema string
}

func (d UserClientData) SitePath() template.URL {
	return template.URL(d.Schema + d.Host)
}

type UserClient interface {
	GetUserClientData(context ctx.Context) UserClientData
}

func GetUserClient() UserClient {
	return userClient{}
}

type userClient struct{}

func (c userClient) GetUserClientData(context ctx.Context) UserClientData {
	type userClientDataContext struct{}
	return context.PersistData(userClientDataContext{}, func() interface{} {
		return c.buildUserClientData(context)
	}).(UserClientData)
}

func (_ userClient) buildUserClientData(context ctx.Context) UserClientData {
	req := context.Request()

	Schema := coalesce.Strings(func() string {
		return req.Header.Get("X-Forwarded-Proto")
	}, func() string {
		return "http"
	})

	Host := coalesce.Strings(func() string {
		return req.Header.Get("Host")
	}, func() string {
		return req.Header.Get("X-Host")
	}, func() string {
		return req.Host
	})

	Ip := coalesce.Strings(func() string {
		return req.Header.Get("X-Real-Ip")
	}, func() string {
		return req.Header.Get("X-Forwarded-For")
	}, func() string {
		return req.RemoteAddr
	})

	return UserClientData{
		Host:           cleanHostAddress(Host),
		OriginalHost:   Host,
		Ip:             cleanIpAddress(Ip),
		OriginalIp:     Ip,
		Schema:         Schema + "://",
		OriginalSchema: Schema,
	}
}

func cleanHostAddress(hostAdress string) string {
	host := strings.Split(hostAdress, ":")
	if len(host) == 1 || host[1] != "80" {
		return hostAdress
	}
	return host[0]
}

func cleanIpAddress(ipAddress string) string {
	address := strings.Split(ipAddress, ":")
	return address[0]
}
