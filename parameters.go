package conditional

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/coredns/policy/plugin/pkg/response"
	"github.com/miekg/dns"
)

type Parameters struct {
	state      *request.Request
	extractors extractorMap
}

type extractorMap map[string]func(state *request.Request) string

func makeExtractors() extractorMap {
	return extractorMap{
		"type": func(state *request.Request) string {
			return state.Type()
		},
		"name": func(state *request.Request) string {
			return state.Name()
		},
		"class": func(state *request.Request) string {
			return state.Class()
		},
		"proto": func(state *request.Request) string {
			return state.Proto()
		},
		"size": func(state *request.Request) string {
			return strconv.Itoa(state.Len())
		},
		"client_ip": func(state *request.Request) string {
			return addrToRFC3986(state.IP())
		},
		"port": func(state *request.Request) string {
			return addrToRFC3986(state.Port())
		},
		"rcode": func(state *request.Request) string {
			rcode := ""
			rr, ok := state.W.(*response.Reader)
			if ok && rr.Msg != nil {
				rcode = dns.RcodeToString[rr.Msg.Rcode]
				if rcode == "" {
					rcode = strconv.Itoa(rr.Msg.Rcode)
				}
			}
			return rcode
		},
		"rsize": func(state *request.Request) string {
			rsize := ""
			rr, ok := state.W.(*response.Reader)
			if ok && rr.Msg != nil {
				rsize = strconv.Itoa(rr.Msg.Len())
			}
			return rsize
		},
		"rflags": func(state *request.Request) string {
			flags := ""
			rr, ok := state.W.(*response.Reader)
			if ok && rr.Msg != nil {
				flags = flagsToString(rr.Msg.MsgHdr)
			}
			return flags
		},
		"id": func(state *request.Request) string {
			return strconv.Itoa(int(state.Req.Id))
		},
		"opcode": func(state *request.Request) string {
			return strconv.Itoa(int(state.Req.Opcode))
		},
		"do": func(state *request.Request) string {
			return boolToString(state.Do())
		},
		"bufsize": func(state *request.Request) string {
			return strconv.Itoa(state.Size())
		},
		"server_ip": func(state *request.Request) string {
			return addrToRFC3986(state.LocalIP())
		},
		"server_port": func(state *request.Request) string {
			return addrToRFC3986(state.LocalPort())
		},
	}
}

func (p Parameters) Get(s string) (interface{}, error) {
	if v, ok := p.Value(s); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unknown variable '%v'", s)
}

func (p Parameters) Value(s string) (string, bool) {
	fn := p.extractors[s]
	if fn == nil {
		return "", false
	}
	return fn(p.state), true
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// flagsToString checks all header flags and returns those
// that are set as a string separated with commas
func flagsToString(h dns.MsgHdr) string {
	flags := make([]string, 7)
	i := 0

	if h.Response {
		flags[i] = "qr"
		i++
	}

	if h.Authoritative {
		flags[i] = "aa"
		i++
	}
	if h.Truncated {
		flags[i] = "tc"
		i++
	}
	if h.RecursionDesired {
		flags[i] = "rd"
		i++
	}
	if h.RecursionAvailable {
		flags[i] = "ra"
		i++
	}
	if h.Zero {
		flags[i] = "z"
		i++
	}
	if h.AuthenticatedData {
		flags[i] = "ad"
		i++
	}
	if h.CheckingDisabled {
		flags[i] = "cd"
		i++
	}
	return strings.Join(flags[:i], ",")
}

// addrToRFC3986 will add brackets to the address if it is an IPv6 address.
func addrToRFC3986(addr string) string {
	if strings.Contains(addr, ":") {
		return "[" + addr + "]"
	}
	return addr
}

// respIP return the first A or AAAA records found in the Answer of the DNS msg
func respIP(r *dns.Msg) net.IP {
	if r == nil {
		return nil
	}

	var ip net.IP
	for _, rr := range r.Answer {
		switch rr := rr.(type) {
		case *dns.A:
			ip = rr.A

		case *dns.AAAA:
			ip = rr.AAAA
		}
		// If there are several responses, currently
		// only return the first one and break.
		if ip != nil {
			break
		}
	}
	return ip
}
