// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"riddhi/states/p2p"
// )

// func main() {

// 	ctx := context.Background()

// 	h, err := p2p.NewHost(
// 		ctx,
// 		p2p.HostOptions{
// 			ListenPort:  4001,
// 			PrivKeyPath: "peer.key",
// 			EnableMDNS:  true,
// 			RelayAddr:   "",
// 		},
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Host started!")

// 	fmt.Println("Peer ID:")
// 	fmt.Println(h.ID())

// 	fmt.Println()

// 	fmt.Println("Listening addresses:")

// 	for _, addr := range h.Addrs() {
// 		fmt.Println(addr)
// 	}

// 	go func() {
// 	for {
// 		time.Sleep(3 * time.Second)

// 		fmt.Println("Connected peers:")

// 		for _, id := range h.Network().Peers() {
// 			fmt.Println(" -", id)
// 		}
// 	}
// }()
// 	select {}
// }

package main

import (
	"context"
	"fmt"
	"io"

	"riddhi/states/p2p"

	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)
type DebugNotify struct{}

func (d *DebugNotify) Listen(network.Network, ma.Multiaddr) {}
func (d *DebugNotify) ListenClose(network.Network, ma.Multiaddr) {}

func (d *DebugNotify) Connected(n network.Network, c network.Conn) {
    fmt.Println("CONNECTED:", c.RemotePeer())
}

func (d *DebugNotify) Disconnected(n network.Network, c network.Conn) {
    fmt.Println("DISCONNECTED:", c.RemotePeer())
}

func (d *DebugNotify) OpenedStream(n network.Network, s network.Stream) {
    fmt.Println("STREAM OPENED:", s.Conn().RemotePeer(), s.Protocol())
}

func (d *DebugNotify) ClosedStream(n network.Network, s network.Stream) {
    fmt.Println("STREAM CLOSED:", s.Conn().RemotePeer(), s.Protocol())
}
func main() {

	ctx := context.Background()

	h, err := p2p.NewHost(
		ctx,
		p2p.HostOptions{
			ListenPort: 4001,
			EnableMDNS: false,
			RelayAddr: "/ip4/10.255.255.254/tcp/9000/ws/p2p/12D3KooWQuGhVvCBAG1zVwn2wXdDCdXbqY3Nh6tYWRQtb7zbDd1F",
			PrivKeyPath: "provider.key",
		//	DisableHolePunch:  false, // Enforce strict relay tunneling
		},
	)

	if err != nil {
		panic(err)
	}

	h.Network().Notify(&DebugNotify{})
	

	p2p.RegisterStreamHandler(
		h,
		func(s network.Stream) {

			fmt.Println("Stream opened from:",
				s.Conn().RemotePeer())
				

			buf := make([]byte, 1024)

			n, _ := s.Read(buf)

			fmt.Println(
				"Received:",
				string(buf[:n]),
			)

			io.WriteString(
				s,
				"Hello Client",
			)

			s.Close()
		},
	)

	fmt.Println("Provider:", h.ID())

	select {}
}