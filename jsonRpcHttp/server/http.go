package server

import (
	"fmt"
	"io"
	"net/http"
)

const (
	BufferSize = 1024
)

type Http struct {
	Ip         string
	Port       string
	Server     Server
	BufferSize int
}

func NewHttpServer(ip string, port string) *Http {
	return &Http{
		ip,
		port,
		Server{},
		BufferSize,
	}
}

func (p *Http) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handleFunc)
	var url = fmt.Sprintf("%s:%s", p.Ip, p.Port)
	err := http.ListenAndServe(url, mux)
	if err != nil {
		fmt.Println("json rpc http server start error", err)
	}
}

func (p *Http) Register(s interface{}) {
	err := p.Server.Register(s)
	if err != nil {
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
	res := p.Server.Handler(data)
	w.Write(res)
}
