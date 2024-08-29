package entities

import (
	"net/url"
	"strings"
)

type StructRequest struct {
	Header           map[string]string `json:"header" validate:"required"`
	Method           string            `json:"method" validate:"required"`
	RemoteAddr       string            `json:"remoteAddr" validate:"required"`
	URL              *url.URL          `json:"url" validate:"required"`
	RawBody          []byte            `json:"rawBody"`
	headerAttbsLower bool
}

func (c *StructRequest) GetHeader(name string) string {
	if !c.headerAttbsLower {
		for k, v := range c.Header {
			delete(c.Header, k)
			c.Header[strings.ToLower(k)] = v
		}
		c.headerAttbsLower = true
	}

	return c.Header[strings.ToLower(name)]
}
