# Secure P2P File Chunking and Transmission

This application is a secure Peer-to-Peer (P2P) file processing and networking system. It splits files into small, encrypted chunks, hosts them persistently on Seeder nodes, and utilizes a central Tracker server to connect peers directly. Leechers query the Tracker for available files, connect directly to Seeders, and securely verify and reconstruct files locally.

## Features

- **P2P Tracker Architecture**: A central tracker server (`tracker.go`) manages peer discovery. It connects Seeders holding chunks with Leechers requesting downloads. All file transfers occur directly peer-to-peer!
- **Dynamic Chunking**: Files are intelligently split into chunks based on total file size.
- **AES-GCM Encryption**: Each chunk is independently encrypted using AES-GCM with uniquely generated 32-byte keys and 16-byte nonces.
- **Merkle Tree Integrity**: A Merkle Tree is generated across all encrypted chunks. Before any decryption occurs on the Leecher, the entire downloaded file tree is verified to prevent CPU waste on tampered or corrupted network packets.
- **Persistent Storage Model**: The chunks are persistently stored to disk by the Seeder, allowing for efficient multi-peer seeding without redundant encryption overhead.
- **Modular Network Transport**: Built atop a flexible `Transport` interface (TCP), easily allowing integration of QUIC or custom protocols in the future.

## Architecture & Roles

The system is comprised of three roles:

1. **Tracker Server (`-mode=tracker`)**
   Runs a lightweight signaling server that maintains a registry of files and the addresses of the Seeders hosting them. It answers queries from Leechers on where to download specific files.

2. **Seeder (`-mode=seeder`)**
   Takes an input file, chunks it, encrypts it, generates a Merkle tree, and saves the `.bin` chunks and metadata to disk. It then connects to the Tracker to announce its availability (`Announce`) and starts a direct TCP listener to serve chunks to incoming Leechers.

3. **Leecher (`-mode=leecher`)**
   Connects to the Tracker to request a specific file (`Query`). The Tracker responds with the address of an active Seeder. The Leecher then connects directly to that Seeder to download the metadata, keys, and chunks. It reconstructs and verifies the Merkle tree in memory, and if verification passes, decrypts the file directly to disk.

## Usage & Testing

First, compile the application:

```bash
go build -o chunk_app .
```

### 1. Start the Tracker Server
Launch the central tracker on port `8080`.
```bash
./chunk_app -mode=tracker -port=8080
```

### 2. Start a Seeder (Host a File)
In a new terminal window, run a Seeder. Tell it the tracker address, your IP address (so the tracker knows where you are), the file to host, and what port it should listen on (`9090`).
```bash
./chunk_app -mode=seeder -tracker=127.0.0.1:8080 -myip=127.0.0.1 -port=9090 -file=demo.txt
```

### 3. Start a Leecher (Download a File)
In a third terminal window, run a Leecher. Ask the tracker for the location of the file, then download it.
```bash
./chunk_app -mode=leecher -tracker=127.0.0.1:8080 -file=demo.txt
```

### Interactive Mode
If you prefer not to use command-line flags, you can simply run the application without them. The application will interactively prompt you in the terminal for the required information, such as the mode, tracker IP address, port, and file paths. This is especially useful when running the application across different laptops on a local network!
```bash
./chunk_app
```

Upon successful download, the reconstructed and verified file will be saved securely in the `received_files/` directory!
