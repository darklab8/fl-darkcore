package registry

import (
	"encoding/json"
	"net/http"

	"github.com/darklab8/fl-darkcore/darkcore/core_types"
	"github.com/darklab8/fl-darkcore/darkcore/settings/logus"
)

type ErrorMessage struct {
	// Refactor to html friendly page.
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

func NewErrorMsg(err error) string {
	result, err := json.Marshal(&ErrorMessage{
		Msg:  err.Error(),
		Type: "serializing_error",
	})
	logus.Log.CheckError(err, "failed to marshal error")
	return string(result)
}

type Registion struct {
	endpoints []*Endpoint
}

func NewRegister() *Registion {
	r := &Registion{}
	return r
}

func (r *Registion) Register(endpoint *Endpoint) {
	r.endpoints = append(r.endpoints, endpoint)
}

func (r *Registion) Foreach(callback func(*Endpoint)) {
	for _, endpoint := range r.endpoints {
		callback(endpoint)
	}
}

type Endpoint struct {
	Url     core_types.Url
	Handler func(w http.ResponseWriter, r *http.Request)
}
