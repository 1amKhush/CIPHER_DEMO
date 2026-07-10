package packetlayer

type PacketType byte

const (
    ChunkRequest  PacketType = 0x01
    ChunkResponse PacketType = 0x02
    LotteryTicket PacketType = 0x03
    KeyReveal     PacketType = 0x04
    TicketReject  PacketType = 0x05  //not in used

    FileMetadataRequest  PacketType = 0x06
    FileMetadataResponse PacketType = 0x07
)