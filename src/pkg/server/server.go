package server

var (
	Upstream *Server
)

type Server struct {
	id string
}

func (s *Server) ID() string {
	return s.id
}
