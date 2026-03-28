//  Copyright(C) 2026 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2026-03-23

package apigate

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/xanygo/anygo/ds/xsync"
	"github.com/xanygo/anygo/xlog"
)

var bp = xsync.NewBytesBufferPool(1024 * 100)

var sensitiveKeys = map[string]struct{}{
	"authorization": {},
	"cookie":        {},
	"token":         {},
	"password":      {},
}

func dumpRequest(s *Service, node *Node, req *http.Request, body []byte, mod *Model) {
	w := rpcDumpWriter()
	if w == nil {
		return
	}
	buf := bp.Get()
	defer bp.Put(buf)

	var modelID string
	if mod != nil {
		modelID = mod.ID
	}
	buf.WriteString("<----- request start -------->\n")
	fmt.Fprintf(buf, "service=%s node=%s model=%s\n", s.Name, node.Name, modelID)
	fmt.Fprintf(buf, "now=%s remote=%s logid=%s\n\n", time.Now().String(), req.RemoteAddr, xlog.FindLogID(req.Context()))

	fmt.Fprintf(buf, "%s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)

	for key, values := range req.Header {
		lk := strings.ToLower(key)

		for _, v := range values {
			if lk == "cookie" {
				var tmp []string
				for _, cv := range strings.Split(v, ";") {
					arr := strings.SplitN(cv, "=", 2)
					if len(arr) == 2 {
						tmp = append(tmp, fmt.Sprintf("%s=*%d", arr[0], len(arr[1])))
					} else {
						tmp = append(tmp, cv)
					}
				}
				buf.WriteString(strings.Join(tmp, ";"))
				continue
			}
			if _, ok := sensitiveKeys[lk]; ok {
				fmt.Fprintf(buf, "%s: *%d\n", key, len(v))
			} else {
				fmt.Fprintf(buf, "%s: %s\n", key, v)
			}
		}
	}

	fmt.Fprintf(buf, "\nbody (len=%d):\n", len(body))
	buf.Write(body)
	buf.WriteString("\n\n<----request finished---->\n\n")

	w.Write(buf.Bytes())
}

func dumpError(s *Service, node *Node, req *http.Request, typ string, err error) {
	w := rpcDumpWriter()
	if w == nil {
		return
	}
	bf := bp.Get()
	defer bp.Put(bf)

	fmt.Fprintf(bf, "<----- response %s error -------->\n", typ)
	fmt.Fprintf(bf, "service=%s node=%s endpoint=%s\n", s.Name, node.Name, node.Endpoint)
	fmt.Fprintf(bf, "now=%s remote=%s logid=%s\n\n", time.Now().String(), req.RemoteAddr, xlog.FindLogID(req.Context()))
	fmt.Fprintf(bf, "%s\n\n", err.Error())

	w.Write(bf.Bytes())
}

func dumpResponse(s *Service, node *Node, req *http.Request, resp *http.Response, rd io.Reader) {
	w := rpcDumpWriter()
	if w == nil {
		return
	}

	fmt.Fprint(w, "<----- response start -------->\n")
	fmt.Fprintf(w, "service=%s node=%s endpoint=%s\n", s.Name, node.Name, node.Endpoint)
	fmt.Fprintf(w, "now=%s remote=%s logid=%s\n\n", time.Now().String(), req.RemoteAddr, xlog.FindLogID(req.Context()))

	bf, _ := httputil.DumpResponse(resp, false)
	w.Write(bf)
	fmt.Fprint(w, "\nbody:\n")
	n, err := io.Copy(w, rd)
	fmt.Fprintf(w, "\n<----- response finished (n=%d, err=%v) -------->\n\n", n, err)
}

var dumpWriterStore = &xsync.OnceInit[io.Writer]{
	New: func() io.Writer {
		return nil
	},
}

func rpcDumpWriter() io.Writer {
	return dumpWriterStore.Load()
}

func SetDumpWriter(w io.Writer) {
	dumpWriterStore.Store(w)
}
