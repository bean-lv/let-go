package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccesscontrolExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlMaxAge           = "Access-Control-Max-Age"

	HeaderOrigin                = "Origin"
	AccessControlRequestMethod  = "Access-Control-Request-Method"
	AccessControlRequestHeaders = "Access-Control-Request-Headers"
)

type CORS interface {
	PrepareCors(resp http.ResponseWriter, req *http.Request)
}

type myCors struct {
	allowAllOrigins  bool
	allowOrigins     []string
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge           time.Duration
}

func New() CORS {
	c := &myCors{}

	c.allowAllOrigins = true
	c.allowOrigins = []string{"http://localhost:4200"}
	c.allowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE", "PATCH"}
	c.allowHeaders = []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-requested-with", "Token"}
	c.exposeHeaders = []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"}
	c.allowCredentials = true
	c.maxAge = time.Hour

	return c
}

func (c *myCors) PrepareCors(resp http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get(HeaderOrigin)
	requestMethod := req.Header.Get(AccessControlRequestMethod)
	requestHeaders := req.Header.Get(AccessControlRequestHeaders)
	if req.Method == "OPTIONS" &&
		(requestMethod != "" || requestHeaders != "") {
		headers := c.preflightRequest(origin, requestMethod, requestHeaders)
		for k, v := range headers {
			resp.Header().Set(k, v)
		}
		resp.WriteHeader(http.StatusOK)
		return
	}
	headers := c.getHeaders()
	for k, v := range headers {
		resp.Header().Set(k, v)
	}
}

func (c *myCors) preflightRequest(origin, requestMethod, requestHeaders string) (headers map[string]string) {
	headers = make(map[string]string)

	if !c.allowAllOrigins && !c.isOriginAllowed(origin) {
		return
	}
	if c.allowAllOrigins {
		headers[AccessControlAllowOrigin] = "*"
	} else {
		headers[AccessControlAllowOrigin] = origin
	}

	for _, method := range c.allowMethods {
		if method == requestMethod {
			headers[AccessControlAllowMethods] = strings.Join(c.allowMethods, ",")
			break
		}
	}

	var allowedHeaders []string
	for _, rHeader := range strings.Split(requestHeaders, ",") {
		rHeader = strings.TrimSpace(rHeader)
	lookupLoop:
		for _, ah := range c.allowHeaders {
			if strings.ToLower(rHeader) == strings.ToLower(ah) {
				allowedHeaders = append(allowedHeaders, rHeader)
				break lookupLoop
			}
		}
	}
	if len(allowedHeaders) > 0 {
		headers[AccessControlAllowHeaders] = strings.Join(allowedHeaders, ",")
	}

	if len(c.exposeHeaders) > 0 {
		headers[AccesscontrolExposeHeaders] = strings.Join(c.exposeHeaders, ",")
	}

	headers[AccessControlAllowCredentials] = strconv.FormatBool(c.allowCredentials)

	if c.maxAge > time.Duration(0) {
		headers[AccessControlMaxAge] = strconv.FormatInt(int64(c.maxAge/time.Second), 10)
	}

	return
}

func (c *myCors) isOriginAllowed(origin string) bool {
	for _, o := range c.allowOrigins {
		if o == origin {
			return true
		}
	}
	return false
}

func (c *myCors) getHeaders() (headers map[string]string) {
	headers = make(map[string]string)

	if c.allowAllOrigins {
		headers[AccessControlAllowOrigin] = "*"
	} else if len(c.allowOrigins) > 0 {
		headers[AccessControlAllowOrigin] = strings.Join(c.allowOrigins, ",")
	}

	if len(c.allowMethods) > 0 {
		headers[AccessControlAllowMethods] = strings.Join(c.allowMethods, ",")
	}

	if len(c.allowHeaders) > 0 {
		headers[AccessControlAllowHeaders] = strings.Join(c.allowHeaders, ",")
	}

	if len(c.exposeHeaders) > 0 {
		headers[AccesscontrolExposeHeaders] = strings.Join(c.exposeHeaders, ",")
	}

	headers[AccessControlAllowCredentials] = strconv.FormatBool(c.allowCredentials)

	if c.maxAge > time.Duration(0) {
		headers[AccessControlMaxAge] = strconv.FormatInt(int64(c.maxAge/time.Second), 10)
	}

	return
}
