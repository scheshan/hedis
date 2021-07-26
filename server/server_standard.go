package server

type StandardServer struct {
	*baseServer
}

func NewStandardServer(config *ServerConfig) *StandardServer {
	srv := &StandardServer{}
	srv.baseServer = newBaseServer(config)

	return srv
}
