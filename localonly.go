package localonly

import (
	"net"

	"github.com/miekg/dns"
)

type LocalOnlyResponseWriter struct {
	dns.ResponseWriter
	NetMask int
}

// WriteMsg implements the dns.ResponseWriter interface.
func (l *LocalOnlyResponseWriter) WriteMsg(res *dns.Msg) error {
	if res.Rcode != dns.RcodeSuccess {
		return l.ResponseWriter.WriteMsg(res)
	}

	if res.Question[0].Qtype == dns.TypeAXFR || res.Question[0].Qtype == dns.TypeIXFR {
		return l.ResponseWriter.WriteMsg(res)
	}
	host, _, err := net.SplitHostPort(l.RemoteAddr().String())
	if err == nil {
		clientIP := net.ParseIP(host)
		if clientIP != nil {
			network := net.IPNet{
				IP:   clientIP,
				Mask: net.CIDRMask(l.NetMask, 32),
			}
			res.Answer = filterAnswer(res.Answer, network)
		}
	}

	return l.ResponseWriter.WriteMsg(res)
}

func filterAnswer(in []dns.RR, network net.IPNet) []dns.RR {
	var res []dns.RR
	for _, r := range in {
		switch r.Header().Rrtype {
		case dns.TypeA:
			ar := r.(*dns.A)
			if network.Contains(ar.A) {
				log.Debug("keep ", r)
				res = append(res, r)
			}
		default:
			res = append(res, r)
		}
	}
	return res
}

// Write implements the dns.ResponseWriter interface.
func (r *LocalOnlyResponseWriter) Write(buf []byte) (int, error) {
	// Should we pack and unpack here to fiddle with the packet... Not likely.
	log.Warning("LocalOnly called with Write: not shuffling records")
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}
