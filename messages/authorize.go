package messages

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/sebdah/recharged-shared/rpc"
	"github.com/sebdah/recharged-shared/types"
)

type AuthorizeReq struct {
	messageType string        `json:"-" type:"string"`
	IdTag       types.IdToken `json:"idTag"`
}

type AuthorizeConf struct {
	IdTagInfo *types.IdTagInfo `json:"idTagInfo"`
	// PriceScheme, Not yet implemented
}

func NewAuthorizeReq(payload string) (req *AuthorizeReq, err error) {
	req = new(AuthorizeReq)
	req.messageType = "Authorize"

	decoder := json.NewDecoder(strings.NewReader(payload))
	err = decoder.Decode(&req)
	if err != nil {
		log.Printf("Unable to parse payload: %s", err.Error())
		err = rpc.NewFormationViolation()
		return
	}

	return
}

func NewAuthorizeConf() (conf *AuthorizeConf) {
	conf = new(AuthorizeConf)
	return
}

// String representation
func (this *AuthorizeReq) String() (str string) {
	js, _ := json.Marshal(this)
	str = string(js)
	return
}

// String representation
func (this *AuthorizeConf) String() (str string) {
	js, _ := json.Marshal(this)
	str = string(js)
	return
}
