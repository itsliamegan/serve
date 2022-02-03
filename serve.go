package serve

import (
	"github.com/itsliamegan/serve/fileserver"
)

func Start(rootDir string, addr string) error {
	server := fileserver.New(rootDir)
	return server.Listen(addr)
}
