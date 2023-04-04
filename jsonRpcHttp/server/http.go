package server

import (
	"encoding/json"
	"fmt"
	"github.com/fushiliang321/jsonrpc/common"
	"io"
	"net/http"
)

type (
	Http struct {
		Ip         string
		Port       string
		Server     Server
		BufferSize int
	}
	ErrorResponse struct {
		Code int
		Text string
	}
)

const (
	BufferSize = 1024
)

func NewHttpServer(ip string, port string) *Http {
	return &Http{
		ip,
		port,
		Server{},
		BufferSize,
	}
}

func (p *Http) Start() {
	var (
		mux = http.NewServeMux()
		url = fmt.Sprintf("%s:%s", p.Ip, p.Port)
	)
	mux.HandleFunc("/", p.handleFunc)
	if err := http.ListenAndServe(url, mux); err != nil {
		fmt.Println("json rpc http server start error", err)
	}
}

func (p *Http) Register(s interface{}) {
	if err := p.Server.Register(s); err != nil {
		fmt.Println("json rpc http server register error", err)
	}
}

func (p *Http) SetBuffer(bs int) {
	p.BufferSize = bs
}

func (p *Http) handleFunc(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		data []byte
	)
	w.Header().Set("Content-Type", "application/json")
	if data, err = io.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var (
		res         = p.Server.Handler(data)
		internalErr *common.InternalErr
	)

	switch res.(type) {
	case common.ErrorResponse:
		_res := res.(common.ErrorResponse)
		if _internalErr, ok := _res.Error.Data.(common.InternalErr); ok && _internalErr.Data != nil {
			internalErr = &_internalErr
		}
	case common.ErrorNotifyResponse:
		_res := res.(common.ErrorNotifyResponse)
		if _internalErr, ok := _res.Error.Data.(common.InternalErr); ok && _internalErr.Data != nil {
			internalErr = &_internalErr
		}
	}

	if internalErr != nil {
		if errData, ok := internalErr.Data.(ErrorResponse); ok {
			w.WriteHeader(errData.Code)
		}
	}
	marshal, _ := json.Marshal(res)
	w.Write(marshal)
}
