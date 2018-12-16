package steam

import (
	"time"

	"github.com/battlesrv/go-gsstat"
)

// GetRules ..
func GetRules(addr string, timeout time.Duration) (interface{}, error) {
	var req gsstat.RequestUDP

	if err := req.NewUDPConnection(addr, timeout); err != nil {
		return nil, err
	}
	defer req.Close()

	reqMsg1 := [11]byte{0xFF, 0xFF, 0xFF, 0xFF}
	copy(reqMsg1[4:], []byte("V"))  // header
	copy(reqMsg1[5:], []byte("-1")) // payload
	copy(reqMsg1[7:], []byte{0xFF, 0xFF, 0xFF, 0xFF})

	// []byte("\xFF\xFF\xFF\xFFV-1\xFF\xFF\xFF\xFF")
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
	copy(reqMsg2[4:], []byte("V")) // header
	copy(reqMsg2[5:], challenge)

	// []byte(fmt.Sprintf("\xFF\xFF\xFF\xFFV%s", challenge))
	if err := req.Send(reqMsg2[:]); err != nil {
		return nil, err
	}

	if err := req.ReadFrom(); err != nil {
		return nil, err
	}

	if !checkHeader(req.Buf, 'E') {
		return nil, errUnknownHeader
	}

	if req.Buf[0] == 'þ' {
		req.SetOffset(19)
	} else {
		req.SetOffset(7)
	}

	var rules Rules
	for {
		name := req.String()
		value := req.String()

		if name != "" {
			rules.Rules = append(rules.Rules, Rule{
				Name:  name,
				Value: value,
			})
		} else {
			break
		}

	}

	return rules, nil
}
