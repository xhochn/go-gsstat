package minecraft

import (
	"time"

	"github.com/battlesrv/go-gsstat"
)

// GetStats ..
func GetStats(addr string, timeout time.Duration) (*FullStats, error) {
	var req gsstat.RequestUDP

	if err := req.NewUDPConnection(addr, timeout); err != nil {
		return nil, err
	}
	defer req.Close()

	sessionID, err := genSessionID()
	if err != nil {
		return nil, err
	}

	challengeToken, err := reqChallengeToken(&req, sessionID)
	if err != nil {
		return nil, err
	}

	// https://wiki.vg/Query#Full_stat
	reqMsg := [15]byte{0xFE, 0xFD, 0x00}
	copy(reqMsg[3:], sessionID[:])
	copy(reqMsg[7:], challengeToken[:])

	if err := req.Send(reqMsg[:]); err != nil {
		return nil, err
	}

	if err := req.ReadFrom(); err != nil {
		return nil, err
	}

	// https://wiki.vg/Query#Response_3
	req.Read(16) // skip first 16 bytes
	var fullStat FullStats
	stop := false

	// https://wiki.vg/Query#K.2C_V_section
	for {
		switch req.String() {
		case "hostname":
			fullStat.Hostname = req.String()
		case "gametype":
			fullStat.GameType = req.String()
		case "game_id":
			fullStat.GameID = req.String()
		case "version":
			fullStat.Version = req.String()
		case "plugins":
			fullStat.Plugins = req.String()
		case "map":
			fullStat.Map = req.String()
		case "numplayers":
			fullStat.Numplayers = gsstat.StringToUint8(req.String())
		case "maxplayers":
			fullStat.Maxplayers = gsstat.StringToUint8(req.String())
		case "hostport":
			fullStat.Hostport = gsstat.StringToUint16(req.String())
		case "hostip":
			fullStat.Hostip = req.String()
		default:
			req.Read(10) // skip 10 bytes

			for {
				name := req.String()

				if len(name) > 0 {
					fullStat.Players = append(fullStat.Players, name)
				}

				if len(name) == 0 {
					stop = true
					break
				}
			}
		}
		if stop {
			break
		}
	}

	return &fullStat, nil
}
