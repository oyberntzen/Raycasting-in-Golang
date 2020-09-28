package networking

import (
	"bytes"
	"encoding/gob"
	"io"
	"net"

	"github.com/oyberntzen/Raycasting-in-Golang/game"
)

/*
Packet Protocol Structure

Client -> Server
- Player Information
	PacketID Uint8
	Username String
- Input
	PacketID Uint8

	Up       Uint16
	Down     Uint16
	Left     Uint16
	Right    Uint16
	MouseX   Uint16
	MouseY   Uint16

Server -> Client
- Server Information
	PacketID Uint8
	PlayerID Uint8
	Cells    [][]Uint8
	Sprites  []Sprite
- Snapshot
	PacketID     Uint8
	ThisPlayer   Player
	OtherPlayers []Player

Server <--> Client
- Event
	PacketID Uint8
	Event    Event
*/

//Packet is the struct converted to first to get the PacketID
type Packet struct {
	ID   PacketID
	Data []byte
}

/*
Client -> Server
*/

//PlayerInfo contains information about the player
type PlayerInfo struct {
	Username string
}

//Input contains information about input done by a player
type Input struct {
	TimeStamp             float32
	Up, Down, Left, Right float32
	MouseX, MouseY        int16
	Jump, Shoot           bool
}

/*
Server -> Client
*/

//ServerInfo contains information about the server
type ServerInfo struct {
	ThisPlayer game.Player
	Cells      [][]uint8
	Sprites    []game.Sprite
}

//Snapshot contains information about every player
type Snapshot struct {
	ThisPlayer   game.Player
	OtherPlayers []game.Player
}

/*
Server <--> Client
*/

//Event contains an event, can be sent from both client and server
type Event struct {
	Event EventID
}

/*
Protocol struct
*/

//Protocol is built on top of gob to make communication easier
type Protocol struct {
	enc *gob.Encoder
	dec *gob.Decoder

	buffer *bytes.Buffer
	bufEnc *gob.Encoder
	bufDec *gob.Decoder
}

//CreateProtocol creates a new protocol
func CreateProtocol(conn net.Conn) Protocol {
	gob.Register(PlayerInfo{})
	gob.Register(Input{})
	gob.Register(ServerInfo{})
	gob.Register(Snapshot{})
	gob.Register(Event{})

	buffer := &bytes.Buffer{}

	return Protocol{gob.NewEncoder(conn), gob.NewDecoder(conn), buffer, gob.NewEncoder(buffer), gob.NewDecoder(buffer)}
}

//Send sends the packet
func (prot *Protocol) Send(data interface{}, id PacketID) error {
	packet := Packet{}
	packet.ID = id

	prot.bufEnc.Encode(data)
	packet.Data = prot.buffer.Bytes()
	prot.buffer.Read(packet.Data)

	return prot.enc.Encode(packet)
}

//Recieve recieves a packet. PacketID is set to NilPacket if packet is invalid or no packet was found
func (prot *Protocol) Recieve() (PacketID, []byte, error) {
	var packet Packet
	err := prot.dec.Decode(&packet)
	if err == io.EOF {
		return NilPacket, nil, nil
	}
	if err != nil {
		return NilPacket, nil, err
	}
	return packet.ID, packet.Data, nil
}

//DecodePlayerInfo decodes []byte sent from server or client to PlayerInfo
func (prot *Protocol) DecodePlayerInfo(data []byte) PlayerInfo {
	prot.buffer.Write(data)
	var playerInfo PlayerInfo
	prot.bufDec.Decode(&playerInfo)
	return playerInfo
}

//DecodeInput decodes []byte sent from server or client to Input
func (prot *Protocol) DecodeInput(data []byte) Input {
	prot.buffer.Write(data)
	var input Input
	prot.bufDec.Decode(&input)
	return input
}

//DecodeServerInfo decodes []byte sent from server or client to ServerInfo
func (prot *Protocol) DecodeServerInfo(data []byte) ServerInfo {
	prot.buffer.Write(data)
	var serverInfo ServerInfo
	prot.bufDec.Decode(&serverInfo)
	return serverInfo
}

//DecodeSnapshot decodes []byte sent from server or client to Snapshot
func (prot *Protocol) DecodeSnapshot(data []byte) Snapshot {
	prot.buffer.Write(data)
	var snapshot Snapshot
	prot.bufDec.Decode(&snapshot)
	return snapshot
}

//DecodeEvent decodes []byte sent from server or client to Event
func (prot *Protocol) DecodeEvent(data []byte) Event {
	prot.buffer.Write(data)
	var event Event
	prot.bufDec.Decode(&event)
	return event
}

/*
Extra structs
*/

//PacketID is the id of a packet
type PacketID uint8

//NilPacket is PacketID for packet that is failed to be decoded
var NilPacket PacketID = 0

//PlayerInfoPacket is PacketID for PlayerInfo
var PlayerInfoPacket PacketID = 1

//InputPacket is PacketID for Input
var InputPacket PacketID = 2

//ServerInfoPacket is PacketID for ServerInfo
var ServerInfoPacket PacketID = 3

//SnapshotPacket is PacketID for Snapshot
var SnapshotPacket PacketID = 4

//EventPacket is PacketID for Event
var EventPacket PacketID = 5

//EventID is used to send events from and to the server
type EventID uint8

//JumpEvent is an event for jumping, sent from client to server
var JumpEvent EventID = 0

//ShootEvent is an event for when a player shoots, sent from client to server
var ShootEvent EventID = 1

//ShotEvent is an event for when a player gets shot, sent from server to the client that is shot
var ShotEvent EventID = 2

/*
Extra conversion functions
*/

/*//ServerInfoFromEnviroment converts enviroment to ServerInfo
func ServerInfoFromEnviroment(env game.Enviroment) ServerInfo {
	serverInfo := ServerInfo{}

	serverInfo.Cells = make([][]uint8, env.Cellsy)
	for y, row := range env.Cells {
		for _, cell := range row {
			serverInfo.Cells[y] = append(serverInfo.Cells[y], uint8(cell))
		}
	}

	serverInfo.Sprites = make([]Sprite, len(env.Sprites))
	for i, sprite := range env.Sprites {
		serverInfo.Sprites[i] = Sprite{sprite.PosX, sprite.PosY, 0, uint8(sprite.Texture)}
	}

	return serverInfo
}

//EnviromentFromServerInfo converts ServerInfo to enviroment
func EnviromentFromServerInfo(serverInfo ServerInfo, dir string) game.Enviroment {
	cells := make([][]int, len(serverInfo.Cells))
	for y, row := range serverInfo.Cells {
		for _, cell := range row {
			cells[y] = append(cells[y], int(cell))
		}
	}

	sprites := make([]game.Sprite, len(serverInfo.Sprites))
	for i, sprite := range serverInfo.Sprites {
		sprites[i] = game.NewSprite(sprite.X, sprite.Y, int(sprite.Texture), 0, false, 0)
	}

	level := game.Level{Cells: cells, Sprites: sprites}
	env := game.Enviroment{}
	env.Init(level, dir)

	return env
}
*/
