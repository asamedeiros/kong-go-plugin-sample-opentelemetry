package entities

import (
	"net/url"
	"strings"
)

type PermissionRequest struct {
	Header           map[string]string `json:"header" validate:"required"`
	Method           string            `json:"method" validate:"required"`
	RemoteAddr       string            `json:"remoteAddr" validate:"required"`
	URL              *url.URL          `json:"url" validate:"required"`
	RawBody          []byte            `json:"rawBody"`
	headerAttbsLower bool
}

func (c *PermissionRequest) GetHeader(name string) string {
	if !c.headerAttbsLower {
		for k, v := range c.Header {
			delete(c.Header, k)
			c.Header[strings.ToLower(k)] = v
		}
		c.headerAttbsLower = true
	}

	return c.Header[strings.ToLower(name)]
}

type Context struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type PermissionResponse struct {
	Provider        string     `json:"provider"`
	Status          int        `json:"status"`
	MaskedKey       string     `json:"masked_key"`
	UserID          string     `json:"user_id"`
	UserScopes      []string   `json:"user_scopes"`
	Context         []*Context `json:"context"`
	AllowedIPs      string     `json:"allowed_ips"`
	CachedTimestamp int64      `json:"timestamp"`
}

type FormattedResource struct {
	Prefix string                   `json:"prefix"`
	Routes []FormattedResourceRoute `json:"routes"`
}

type FormattedResourceRoute struct {
	Path              string              `json:"path"`
	Hosts             []string            `json:"hosts"`
	Headers           map[string][]string `json:"headers"`
	Scopes            []string            `json:"scopes"`
	IsInternalUseOnly bool                `json:"isInternalUseOnly"`
}

type Entity struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BlockList struct {
	Entity     Entity `json:"entity"`
	AllowedIps string `json:"allowed_ips"`
}

type KeyEncrypted struct {
	Hash string `json:"hash"`
	Mask string `json:"mask"`
}
