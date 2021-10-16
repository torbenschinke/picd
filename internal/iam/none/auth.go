package none

import (
	"net/http"
)

type Authenticator struct {
}

func (a Authenticator) Authenticate(w http.ResponseWriter, r *http.Request) bool {
	return true
}
