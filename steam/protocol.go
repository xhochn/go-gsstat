package steam

// SourceA2SInfo - https://developer.valvesoftware.com/wiki/Server_queries#A2S_INFO
type SourceA2SInfo struct {
	Protocol    byte   `json:"protocol"`
	Name        string `json:"name"`
	Map         string `json:"map"`
	Folder      string `json:"folder"`
	Game        string `json:"game"`
	ID          uint16 `json:"id"`
	Players     byte   `json:"players"`
	MaxPlayers  byte   `json:"max_players"`
	Bots        byte   `json:"bots"`
	ServerType  string `json:"server_type"`
	Environment string `json:"environment"`
	Visibility  string `json:"visibility"`
	VAC         bool   `json:"vac"`
	TheShip     struct {
		Mode      string `json:"mode,omitempty"`
		Witnesses byte   `json:"witnesses,omitempty"`
		Duration  byte   `json:"duration,omitempty"`
	} `json:"the_ship"`
	Version string `json:"version"`
	EDF     struct {
		Port     uint16 `json:"port,omitempty"`
		SteamID  uint64 `json:"steam_id,omitempty"`
		SourceTV struct {
			Port uint16 `json:"port,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"source_tv"`
		Keywords string `json:"keywords,omitempty"`
		GameID   uint64 `json:"game_id,omitempty"`
	} `json:"edf"`
}

// ObsoleteGoldSourceA2SInfo - https://developer.valvesoftware.com/wiki/Server_queries#Obsolete_GoldSource_Response
type ObsoleteGoldSourceA2SInfo struct {
	Address     string `json:"address"`
	Name        string `json:"name"`
	Map         string `json:"map"`
	Folder      string `json:"folder"`
	Game        string `json:"game"`
	Players     byte   `json:"players"`
	MaxPlayers  byte   `json:"max_players"`
	Protocol    byte   `json:"protocol"`
	ServerType  string `json:"server_type,omitempty"`
	Environment string `json:"environment,omitempty"`
	Visibility  string `json:"visibility"`
	Mod         string `json:"mod"`
	IsMod       struct {
		Link         string `json:"link,omitempty"`
		DownloadLink string `json:"download_link,omitempty"`
		Version      uint64 `json:"version,omitempty"`
		Size         uint64 `json:"size,omitempty"`
		Type         string `json:"type,omitempty"`
		DLL          string `json:"dll,omitempty"`
	} `json:"is_mod"`
	VAC  bool `json:"vac"`
	Bots byte `json:"bots"`
}

// PlayersInfo - https://developer.valvesoftware.com/wiki/Server_queries#A2S_PLAYER
type PlayersInfo struct {
	Count   byte      `json:"count"`
	Players []Players `json:"players"`
}

// Players - https://developer.valvesoftware.com/wiki/Server_queries#A2S_RULES
type Players struct {
	Index    byte    `json:"index"`
	Name     string  `json:"name"`
	Score    uint32  `json:"score"`
	Duration float32 `json:"duration"`
}

// Rules ..
type Rules struct {
	Rules []Rule `json:"rules"`
}

// Rule ..
type Rule struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func checkHeader(buf []byte, lit byte) bool {
	header := 4

	if buf[0] == 'þ' {
		header = 16
	}

	return buf[header] == lit
}

func getChallenge(buf []byte) ([]byte, error) {
	var challenge []byte
	challenge = append(challenge, buf[5])
	challenge = append(challenge, buf[6])
	challenge = append(challenge, buf[7])
	challenge = append(challenge, buf[8])

	return challenge, nil
}

func serverType(b byte, old bool) string {
	switch b {
	case 'd':
		return "dedicated server"
	case 'l':
		return "non-dedicated server"
	case 'p':
		if old {
			return "HLTV server"
		}
		return "SourceTV relay (proxy)"
	default:
		return "unknown server type"
	}
}

func environment(b byte) string {
	switch b {
	case 'l':
		return "Linux"
	case 'w':
		return "Windows"
	case 'm' | 'o':
		return "MacOS"
	default:
		return "unknown environment"
	}
}

func visibility(b byte) string {
	switch b {
	case 0:
		return "public"
	case 1:
		return "private"
	default:
		return "unknown type visibility"
	}
}
