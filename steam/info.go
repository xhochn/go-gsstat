package steam

import (
	"fmt"
	"time"

	"github.com/battlesrv/go-gsstat"
)

// GetInfo ..
func GetInfo(addr string, timeout time.Duration) (interface{}, error) {
	var req gsstat.RequestUDP

	if err := req.NewUDPConnection(addr, timeout); err != nil {
		return nil, err
	}
	defer req.Close()

	// https://developer.valvesoftware.com/wiki/Server_queries#Request_Format
	reqMsg := [25]byte{0xFF, 0xFF, 0xFF, 0xFF}
	copy(reqMsg[4:], []byte("T"))                   // header
	copy(reqMsg[5:], []byte("Source Engine Query")) // payload
	copy(reqMsg[25:], []byte{0x00})

	if err := req.Send(reqMsg[:]); err != nil {
		return nil, err
	}

	if err := req.ReadFrom(); err != nil {
		return nil, err
	}

	if checkHeader(req.Buf, 'm') {
		return getObsoleteGoldSourceA2SInfo(&req)
	}

	if checkHeader(req.Buf, 'I') {
		return getSourceA2SInfo(&req)
	}

	return nil, fmt.Errorf("package header is UNKNOWN")
}

func getSourceA2SInfo(raw *gsstat.RequestUDP) (*SourceA2SInfo, error) {
	raw.Read(5)

	var info SourceA2SInfo
	info.Protocol = raw.Read(1)[0]
	info.Name = raw.String()
	info.Map = raw.String()
	info.Folder = raw.String()
	info.Game = raw.String()
	info.ID = gsstat.BytesToUint16(raw.Read(2))
	info.Players = raw.Read(1)[0]
	info.MaxPlayers = raw.Read(1)[0]
	info.Bots = raw.Read(1)[0]
	info.ServerType = serverType(raw.Read(1)[0], false)
	info.Environment = environment(raw.Read(1)[0])
	info.Visibility = visibility(raw.Read(1)[0])
	info.VAC = gsstat.ByteToBool(raw.Read(1)[0])

	if info.ID == 2400 {
		// info.IsTheShip = true
		switch raw.Read(1)[0] {
		case 0:
			info.TheShip.Mode = "Hunt"
		case 1:
			info.TheShip.Mode = "Elimination"
		case 2:
			info.TheShip.Mode = "Duel"
		case 3:
			info.TheShip.Mode = "Deathmatch"
		case 4:
			info.TheShip.Mode = "VIP Team"
		case 5:
			info.TheShip.Mode = "Team Elimination"
		}
		info.TheShip.Witnesses = raw.Read(1)[0]
		info.TheShip.Duration = raw.Read(1)[0]
	}

	info.Version = raw.String()

	if raw.Offset < len(raw.Buf) {
		// info.EDFExists = true
		edf := raw.Read(1)[0]
		if edf&0x80 != 0 {
			info.EDF.Port = gsstat.BytesToUint16(raw.Read(2))
		}
		if edf&0x10 != 0 {
			info.EDF.SteamID = gsstat.BytesToUint64(raw.Read(8))
		}
		if edf&0x40 != 0 {
			info.EDF.SourceTV.Port = gsstat.BytesToUint16(raw.Read(2))
			info.EDF.SourceTV.Name = raw.String()
		}
		if edf&0x20 != 0 {
			info.EDF.Keywords = raw.String()
		}
		if edf&0x01 != 0 {
			info.EDF.GameID = gsstat.BytesToUint64(raw.Read(8))
		}
	}

	return &info, nil
}

func getObsoleteGoldSourceA2SInfo(raw *gsstat.RequestUDP) (*ObsoleteGoldSourceA2SInfo, error) {
	raw.Read(5)

	var info ObsoleteGoldSourceA2SInfo
	info.Address = raw.String()
	info.Name = raw.String()
	info.Map = raw.String()
	info.Folder = raw.String()
	info.Game = raw.String()
	info.Players = raw.Read(1)[0]
	info.MaxPlayers = raw.Read(1)[0]
	info.Protocol = raw.Read(1)[0]
	info.ServerType = serverType(raw.Read(1)[0], false)
	info.Environment = environment(raw.Read(1)[0])
	info.Visibility = visibility(raw.Read(1)[0])

	switch raw.Read(1)[0] {
	case 0:
		info.Mod = "Half-Life"
	case 1:
		info.Mod = "Half-Life mod"
		info.IsMod.Link = raw.String()
		info.IsMod.DownloadLink = raw.String()
		raw.Read(1)
		info.IsMod.Version = gsstat.BytesToUint64(raw.Read(8))
		info.IsMod.Size = gsstat.BytesToUint64(raw.Read(8))
		switch raw.Read(1)[0] {
		case 0:
			info.IsMod.Type = "single and multiplayer mod"
		case 1:
			info.IsMod.Type = "multiplayer only mod"
		}
		switch raw.Read(1)[0] {
		case 0:
			info.IsMod.DLL = "it uses the Half-Life DLL"
		case 1:
			info.IsMod.Type = " it uses its own DLL"
		}
	}

	info.VAC = gsstat.ByteToBool(raw.Read(1)[0])
	info.Bots = raw.Read(1)[0]

	return &info, nil
}
