package main

import "fmt"

func main() {
	s := NewServer(WithConfigA("123"))
	s.Start()
}

type Server struct {
	configA string
	configB string
	configC string
}

func NewServer(opts ...func(s *Server)) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start() {
	fmt.Println(s.configA, s.configB, s.configC)
}

func WithConfigA(configA string) func(s *Server) {
	return func(s *Server) {
		s.configA = configA
	}
}
