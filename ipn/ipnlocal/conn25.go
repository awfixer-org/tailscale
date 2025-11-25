package ipnlocal

import (
	"encoding/json"
	"net/http"

	"tailscale.com/appc"
)

// ok, so ipnlocal imports appc
// so we can't put the peeraip registration and handling that conn25 wants in appc, because we get an import loop
// possible solutions:
//  1. ipnlocal stops importing appc
// 		(how would this work?)
//  2. conn25 doesn't go in appc
// 		(this might really be a part of 1. not sure, Adrian said put it in there so not approaching this as the first thing, how does ipnlocal _have_ an appc without importing appc
// 		anyway (obvz it can't but how does the whole thing work then?))
//  3. put the peerapi registration and handling for conn25 in ipnlocal <- that's what we're doing here

func init() {
	RegisterPeerAPIHandler("/v0/connector/transit-ip/", handleConnectorTransitIP)
}

func handleConnectorTransitIP(h PeerAPIHandler, w http.ResponseWriter, r *http.Request) {
	var req appc.ConnectorTransitIPRequest
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	resp := h.LocalBackend().Conn25.HandleConnectorTransitIPRequest(h.Peer().ID(), req)
	bs, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}
