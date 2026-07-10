package p2p

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// RegisterStreamHandler registers the handler for our protocol.
func RegisterStreamHandler(
	h host.Host,
	handler network.StreamHandler,
) {
	h.SetStreamHandler(
		ProtocolID,
		handler,
	)
}

// OpenStream opens a new stream to another peer.
func OpenStream(
	ctx context.Context,
	h host.Host,
	peerID peer.ID,
) (network.Stream, error) {

	fmt.Println("Opening stream...")
	fmt.Println(peerID)
	fmt.Println(ProtocolID)

	stream, err := h.NewStream(
		network.WithAllowLimitedConn(ctx, "relay-stream"), //very imp step
		peerID,
		ProtocolID,
	)
	fmt.Println("Returned from NewStream")

	if err != nil {
		return nil, fmt.Errorf(
			"failed to open stream: %w",
			err,
		)
	}

	return stream, nil
}
