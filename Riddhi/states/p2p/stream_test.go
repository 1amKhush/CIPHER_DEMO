package p2p

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/host"
)

func TestOpenStreamViaRelayAllowsLimitedConn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tempDir := t.TempDir()

	relayHost, err := NewHost(ctx, HostOptions{
		ListenPort:  0,
		PrivKeyPath: filepath.Join(tempDir, "relay.key"),
		EnableMDNS:  false,
	})
	if err != nil {
		t.Fatalf("create relay host: %v", err)
	}

	r, err := relay.New(relayHost)
	if err != nil {
		t.Fatalf("enable relay service: %v", err)
	}
	defer r.Close()

	relayAddr := findLoopbackRelayAddr(t, relayHost.Addrs())
	relayAddr = fmt.Sprintf("%s/p2p/%s", relayAddr, relayHost.ID())

	providerHost, err := NewHost(ctx, HostOptions{
		ListenPort:       0,
		PrivKeyPath:      filepath.Join(tempDir, "provider.key"),
		EnableMDNS:       false,
		RelayAddr:        relayAddr,
	//	DisableHolePunch: true,
	})
	if err != nil {
		t.Fatalf("create provider host: %v", err)
	}

	clientHost, err := NewHost(ctx, HostOptions{
		ListenPort:       0,
		PrivKeyPath:      filepath.Join(tempDir, "client.key"),
		EnableMDNS:       false,
		RelayAddr:        relayAddr,
		//DisableHolePunch: true,
	})
	if err != nil {
		t.Fatalf("create client host: %v", err)
	}

	providerHost.SetStreamHandler(ProtocolID, func(s network.Stream) {
		defer s.Close()
		buf := make([]byte, 64)
		_, _ = io.ReadFull(s, buf)
		_, _ = s.Write([]byte("pong"))
		_ = s.CloseWrite()
	})

	if err := ConnectViaRelay(ctx, clientHost, relayAddr, providerHost.ID()); err != nil {
		t.Fatalf("connect via relay: %v", err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if clientHost.Network().Connectedness(providerHost.ID()) == network.Limited {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	stream, err := OpenStream(ctx, clientHost, providerHost.ID())
	if err != nil {
		t.Fatalf("open relay-backed stream: %v", err)
	}
	defer stream.Close()

	if _, err := stream.Write([]byte("ping")); err != nil {
		t.Fatalf("write to stream: %v", err)
	}
	if err := stream.CloseWrite(); err != nil {
		t.Fatalf("close write side: %v", err)
	}
}

func findLoopbackRelayAddr(t *testing.T, addrs []ma.Multiaddr) string {
	t.Helper()
	for _, addr := range addrs {
		if strings.Contains(addr.String(), "/ip4/127.0.0.1/") {
			return addr.String()
		}
	}
	t.Fatalf("no loopback relay address found in %v", addrs)
	return ""
}


func monitorConnections(
    ctx context.Context,
    t *testing.T,
    h host.Host,
    peerID peer.ID,
) {
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return

        case <-ticker.C:
            conns := h.Network().ConnsToPeer(peerID)

            t.Log("------ Connections ------")

            if len(conns) == 0 {
                t.Log("No connections")
                continue
            }

            for _, c := range conns {
                t.Logf("Remote: %s", c.RemoteMultiaddr())
                t.Logf("Local : %s", c.LocalMultiaddr())
            }
        }
    }
}

func TestDirectTransport(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tempDir := t.TempDir()

	// -------------------------
	// Provider
	// -------------------------
	providerHost, err := NewHost(ctx, HostOptions{
		ListenPort:  0,
		PrivKeyPath: filepath.Join(tempDir, "provider.key"),
		EnableMDNS:  false,
	})
	if err != nil {
		t.Fatalf("create provider host: %v", err)
	}
	defer providerHost.Close()

	providerHost.SetStreamHandler(ProtocolID, func(s network.Stream) {
		defer s.Close()

	buf := make([]byte, 5)

_, err := io.ReadFull(s, buf)
if err != nil {
    t.Fatalf("provider read: %v", err)
}

if string(buf) != "hello" {
    t.Fatalf("expected hello, got %q", string(buf))
}

		_, err = s.Write([]byte("world"))
		if err != nil {
			t.Errorf("provider write: %v", err)
		}
	})

	// -------------------------
	// Client
	// -------------------------
	clientHost, err := NewHost(ctx, HostOptions{
		ListenPort:  0,
		PrivKeyPath: filepath.Join(tempDir, "client.key"),
		EnableMDNS:  false,
	})
	if err != nil {
		t.Fatalf("create client host: %v", err)
	}
	defer clientHost.Close()

	// -------------------------
	// Direct connect
	// -------------------------
	providerInfo := peer.AddrInfo{
		ID:    providerHost.ID(),
		Addrs: providerHost.Addrs(),
	}

	if err := clientHost.Connect(ctx, providerInfo); err != nil {
		t.Fatalf("client connect: %v", err)
	}

	// -------------------------
	// Open stream
	// -------------------------
	stream, err := OpenStream(
		ctx,
		clientHost,
		providerHost.ID(),
	)
	if err != nil {
		t.Fatalf("open stream: %v", err)
	}
	defer stream.Close()

	// -------------------------
	// Send
	// -------------------------
	_, err = stream.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := stream.CloseWrite(); err != nil {
		t.Fatalf("close write: %v", err)
	}

	// -------------------------
	// Receive
	// -------------------------
reply := make([]byte, 5)

_, error := io.ReadFull(stream, reply)
if err != nil {
    t.Fatal(error)
}

if string(reply) != "world" {
    t.Fatalf("expected world got %q", string(reply))
}
}



func TestHolePunchTransport(t *testing.T) {
	t.Skip("Requires two NATed peers and a public relay")
}

func TestTransportUpgrade(t *testing.T) {
	t.Skip("Requires successful DCUtR hole punching")
}