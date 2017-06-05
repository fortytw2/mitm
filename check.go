package mitm

import "net/http"

type mitmCtxKey struct{}

// Check returns true if the connection has been marked as MITM-ed
func Check(r *http.Request) bool {
	v, ok := r.Context().Value(mitmCtxKey{}).(bool)
	if !ok {
		return false
	}
	return v
}
