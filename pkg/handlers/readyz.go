package handlers

import (
	"fmt"
	"net/http"

	httputil "github.com/couchbase/couchbase-exporter/pkg/http/util"
	"github.com/couchbase/couchbase-exporter/pkg/util"
)

// Readyz handler responds with a 200 when it gets a 200 back.
func Readyz(client util.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Client.Get(client.URL("pools/default"))
		if err != nil {
			httputil.RespondErr(w, r, fmt.Errorf("error occurred performing Get: %w", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			httputil.Respond(w, r, "Cluster available... ", http.StatusOK)
			return
		}

		httputil.Respond(w, r, "Cluster unavailable... ", resp.StatusCode)
	}
}
