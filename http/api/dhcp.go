package api

import (
	"github.com/danieldin95/lightstar/network/libvirtn"
	"github.com/danieldin95/lightstar/schema"
	"github.com/gorilla/mux"
	"net/http"
)

type DHCPLease struct {
}

func (le DHCPLease) Router(router *mux.Router) {
	router.HandleFunc("/api/dhcp/lease", le.GET).Methods("GET")
	router.HandleFunc("/api/dhcp/lease/{name}", le.GET).Methods("GET")
}

func (le DHCPLease) GET(w http.ResponseWriter, r *http.Request) {
	name, _ := GetArg(r, "name")
	leases, err := libvirtn.ListLeases(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := make(map[string]schema.DHCPLease, 128)
	for addr, le := range leases {
		data[addr] = schema.DHCPLease{
			Mac:      le.Mac,
			IPAddr:   le.IPAddr,
			Prefix:   le.Prefix,
			Hostname: le.Hostname,
			Type:     le.Type,
		}
	}
	ResponseJson(w, data)
}

func (le DHCPLease) POST(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (le DHCPLease) PUT(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}

func (le DHCPLease) DELETE(w http.ResponseWriter, r *http.Request) {
	ResponseMsg(w, 0, "")
}