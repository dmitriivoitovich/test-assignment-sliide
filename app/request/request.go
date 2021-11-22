package request

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	forwardedForHeaderName = "X-FORWARDED-FOR"

	countParamName = "count"
	defaultCount   = 5
	maxCount       = 100

	offsetParamName = "offset"
	maxOffset       = 10 * 1000
)

var (
	ErrInvalidParameterValue = errors.New("invalid parameter")
)

type Request struct {
	Count  uint64
	Offset uint64
	UserIP net.IP
}

func NewRequest(count, offset uint64, ip net.IP) *Request {
	return &Request{
		Count:  count,
		Offset: offset,
		UserIP: ip,
	}
}

func (r *Request) Parse(httpRequest *http.Request) error {
	if err := r.parseCount(httpRequest); err != nil {
		return err
	}

	if err := r.parseOffset(httpRequest); err != nil {
		return err
	}

	if err := r.parseUserIP(httpRequest); err != nil {
		return err
	}

	return nil
}

func (r *Request) parseCount(req *http.Request) error {
	var err error

	r.Count, err = r.queryParamUint64(req, countParamName)
	if err != nil || r.Count > maxCount {
		return ErrInvalidParameterValue
	}
	if r.Count == 0 {
		r.Count = defaultCount
	}

	return nil
}

func (r *Request) parseOffset(req *http.Request) error {
	var err error

	r.Offset, err = r.queryParamUint64(req, offsetParamName)
	if err != nil || r.Offset > maxOffset {
		return ErrInvalidParameterValue
	}

	return nil
}

func (r *Request) parseUserIP(req *http.Request) error {
	forwardedForHeader := req.Header.Get(forwardedForHeaderName)

	ips := strings.Split(forwardedForHeader, ",")
	if len(ips) > 0 {
		r.UserIP = net.ParseIP(strings.TrimSpace(ips[0]))
		if r.UserIP != nil {
			return nil
		}
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return err
	}

	r.UserIP = net.ParseIP(ip)

	return nil
}

func (r *Request) queryParamUint64(req *http.Request, key string) (uint64, error) {
	value := strings.TrimSpace(req.URL.Query().Get(key))

	if value != "" {
		res, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return 0, err
		}

		return res, nil
	}

	return 0, nil
}
