package transportlayer

import (
	"fmt"
	"net"
)

//To start the server
func StartServer(
	address string,
) (net.Listener, error) {

	listener, err := net.Listen(
		"tcp",
		address,
	)

	if err != nil {
		return nil, err
	}

	fmt.Println(
		"provider listening on",
		address,
	)

	return listener, nil
}