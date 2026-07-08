package p2p
import (
	"context"
	"fmt"
	
	
	   "github.com/libp2p/go-libp2p/core/host"
		"github.com/libp2p/go-libp2p/core/peer"
		"github.com/libp2p/go-libp2p/core/peerstore"
		"github.com/multiformats/go-multiaddr"
)
func ConnectViaRelay(
    ctx context.Context,
    h host.Host,
    relayAddr string,
    providerID peer.ID,
) error {

    // Inside p2p.ConnectViaRelay:
targetMaddrStr := fmt.Sprintf("%s/p2p-circuit/p2p/%s", relayAddr, providerID.String())

    addr, err := multiaddr.NewMultiaddr(targetMaddrStr)
    if err != nil {
       return fmt.Errorf("failed to parse multiaddress: %w", err)
    }

    // Create an AddrInfo struct bound directly to the Provider's ID
    info := &peer.AddrInfo{
		ID:    providerID,
		Addrs: []multiaddr.Multiaddr{addr},
	}
    if err != nil {
        return err
    }

// Remove old direct addresses
h.Peerstore().ClearAddrs(providerID)

// Save the circuit route under the provider's ID
	h.Peerstore().AddAddrs(
		providerID,
		info.Addrs,
		peerstore.PermanentAddrTTL,
	)

// Explicitly dial the provider through the circuit route
	if err := h.Connect(ctx, *info); err != nil {
		return fmt.Errorf("failed connecting to provider via relay: %w", err)
	}

//     info, err := peer.AddrInfoFromP2pAddr(addr)
//     if err != nil {
//         return err
//     }

//     err = h.Connect(ctx, *info)
// if err != nil {
//     return err
// }

// fmt.Println("Connected through relay address")

 return nil
}