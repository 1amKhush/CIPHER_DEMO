package transportlayer

import (
	"fmt"
	"net"
)
//to connect the client to provider
func ConnectToProvider(
	address string,
) (net.Conn, error) {

	conn, err := net.Dial(
		"tcp",
		address,
	)

	if err != nil {
		return nil, err
	}

	fmt.Println(
		"connected to provider",
	)

	return conn, nil
}
