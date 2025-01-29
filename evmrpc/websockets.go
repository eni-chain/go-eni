package evmrpc

import (
	"net/http"

	"github.com/eni-chain/go-eni/utils/metrics"
)

type wsConnectionHandler struct {
	underlying http.Handler
}

func (h *wsConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics.IncWebsocketConnects()
	h.underlying.ServeHTTP(w, r)
}

func NewWSConnectionHandler(handler http.Handler) http.Handler {
	return &wsConnectionHandler{underlying: handler}
}
