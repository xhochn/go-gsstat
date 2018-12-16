package minecraft

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"strconv"
	"time"

	"github.com/battlesrv/go-gsstat"
)

// FullStats ..
type FullStats struct {
	Hostname   string   `json:"hostname"`
	GameType   string   `json:"game_type"`
	GameID     string   `json:"game_id"`
	Version    string   `json:"version"`
	Plugins    string   `json:"plugins"`
	Map        string   `json:"map"`
	Numplayers uint8    `json:"numplayers"`
	Maxplayers uint8    `json:"maxplayers"`
	Hostport   uint16   `json:"hostport"`
	Hostip     string   `json:"hostip"`
	Players    []string `json:"players"`
}

// https://wiki.vg/Query#Generating_a_Session_ID
func genSessionID() ([]byte, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	sessionID := make([]byte, 4)

	if _, err := rand.Read(sessionID[0:]); err != nil {
		return nil, err
	}

	// The session ID is used identify your requests. The following examples use
	// session ID = 1 (encoded as 00 00 00 01 on the hex dumps). Only the lower
	// 4-bits on each byte of the session ID should be used as Minecraft does
	// not process the higher 4-bits on each byte. To convert any 4-byte
	// session ID to a valid Minecraft session ID, simply mask the bits with
	// sessionId & 0x0F0F0F0F.
	for i := range sessionID {
		sessionID[i] = sessionID[i] & 0x0F
	}

	return sessionID[:], nil
}

// https://wiki.vg/Query#Request
func reqChallengeToken(req *gsstat.RequestUDP, sessionID []byte) ([]byte, error) {
	reqMsg := [11]byte{0xFE, 0xFD, 0x09}
	copy(reqMsg[3:], sessionID[:])

	if err := req.Send(reqMsg[:]); err != nil {
		return nil, err
	}

	if err := req.ReadFrom(); err != nil {
		return nil, err
	}

	// TODO: check type packet and SessionPK in buf

	// we recived challenge token as string (Null-terminated string), then convert to Int64
	tInt64, err := strconv.ParseInt(string(req.Buf[5:len(req.Buf)-1]), 10, 64)
	if err != nil {
		return nil, err
	}

	tBuf := &bytes.Buffer{}

	if err = binary.Write(tBuf, binary.BigEndian, tInt64); err != nil {
		return nil, err
	}

	challengeToken := make([]byte, 4)
	tBytes := tBuf.Bytes()
	copy(challengeToken[0:], tBytes[len(tBytes)-4:])
	return challengeToken, nil
}
