package steam

import (
	"fmt"
	"time"

	"github.com/battlesrv/go-gsstat"
)

var errUnknownHeader = fmt.Errorf("package header is UNKNOWN")

// GetPlayers ..
func GetPlayers(addr string, timeout time.Duration) (interface{}, error) {
	var req gsstat.RequestUDP

	if err := req.NewUDPConnection(addr, timeout); err != nil {
		return nil, err
	}
	defer req.Close()

	// https://developer.valvesoftware.com/wiki/Server_queries#Request_Format_2
	reqMsg1 := [11]byte{0xFF, 0xFF, 0xFF, 0xFF}
	copy(reqMsg1[4:], []byte("U"))  // header
	copy(reqMsg1[5:], []byte("-1")) // payload
	copy(reqMsg1[7:], []byte{0xFF, 0xFF, 0xFF, 0xFF})

	if err := req.Send(reqMsg1[:]); err != nil {
		return nil, err
	}

	if err := req.ReadFrom(); err != nil {
		return nil, err
	}

	if !checkHeader(req.Buf, 'A') {
		return nil, errUnknownHeader
	}

	challenge, err := getChallenge(req.Buf)
	if err != nil {
		return nil, err
	}

	reqMsg2 := [9]byte{0xFF, 0xFF, 0xFF, 0xFF}
	copy(reqMsg2[4:], []byte("U")) // header
	copy(reqMsg2[5:], challenge)

	if err := req.Send(reqMsg2[:]); err != nil {
		return nil, err
	}

	if err := req.ReadFrom(); err != nil {
		return nil, err
	}

	if !checkHeader(req.Buf, 'D') {
		return nil, errUnknownHeader
	}

	if req.Buf[0] == 'þ' {
		req.SetOffset(17)
	} else {
		req.SetOffset(5)
	}

	var players PlayersInfo
	players.Count = req.Read(1)[0]

	for i := 0; i < int(players.Count); i++ {
		players.Players = append(players.Players, Players{
			Index:    req.Read(1)[0],
			Name:     req.String(),
			Score:    gsstat.BytesToUint32(req.Read(4)),
			Duration: gsstat.BytesToFloat32(req.Read(4)),
		})
	}

	return players, nil
}
