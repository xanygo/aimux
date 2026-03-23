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

func dumpRequest(s *Service, node *Node, req *http.Request, body []byte) {
	w := rpcDumpWriter()
	if w == nil {
		return
	}
	buf := bp.Get()
	defer bp.Put(buf)

	buf.WriteString("<----- request start -------->\n")
	fmt.Fprintf(buf, "service=%s node=%s\n", s.Name, node.Name)
	fmt.Fprintf(buf, "now=%s remote=%s logid=%s\n\n", time.Now().String(), req.RemoteAddr, xlog.FindLogID(req.Context()))

	fmt.Fprintf(buf, "%s %s %s\n", req.Method, req.URL.RequestURI(), req.Proto)

	// 2. Header（逐行处理）
	sensitiveKeys := map[string]struct{}{
		"authorization": {},
		"cookie":        {},
		"token":         {},
		"password":      {},
	}

	for key, values := range req.Header {
		lk := strings.ToLower(key)

		for _, v := range values {
			if lk == "cookie" {
				var tmp []string
				for _, cv := range strings.Split(v, ";") {
					arr := strings.SplitN(cv, "=", 2)
					if len(arr) == 2 {
						tmp = append(tmp, fmt.Sprintf("%s=%s", arr[0], strings.Repeat("*", len(arr[1]))))
					} else {
						tmp = append(tmp, cv)
					}
				}
				buf.WriteString(strings.Join(tmp, ";"))
				continue
			}
			if _, ok := sensitiveKeys[lk]; ok {
				fmt.Fprintf(buf, "%s: %s\n", key, strings.Repeat("*", len(v)))
			} else {
				fmt.Fprintf(buf, "%s: %s\n", key, v)
			}
		}
	}

	buf.WriteString("\nbody:\n")
	buf.Write(body)
	buf.WriteString("\n\n<----request finished---->\n\n")

	w.Write(buf.Bytes())
}

func dumpResponse(s *Service, node *Node, req *http.Request, resp *http.Response, rd io.Reader) {
	w := rpcDumpWriter()
	if w == nil {
		return
	}
	bf, _ := httputil.DumpResponse(resp, false)

	fmt.Fprint(w, "<----- response start -------->\n")
	fmt.Fprintf(w, "service=%s node=%s endpoint=%s\n", s.Name, node.Name, node.Endpoint)
	fmt.Fprintf(w, "now=%s remote=%s logid=%s\n\n", time.Now().String(), req.RemoteAddr, xlog.FindLogID(req.Context()))

	w.Write(bf)
	fmt.Fprint(w, "\nbody:\n")
	io.Copy(w, rd)
	fmt.Fprint(w, "<----- response finished -------->\n\n")
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
