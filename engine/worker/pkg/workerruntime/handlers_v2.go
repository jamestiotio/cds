package workerruntime

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ovh/cds/sdk"
	"github.com/rockbears/log"
)

func V2_runResultHandler(ctx context.Context, wk Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		btes, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, r, sdk.NewError(sdk.ErrWrongRequest, err))
			return
		}
		defer r.Body.Close()

		switch r.Method {
		case http.MethodGet:
			var filter V2FilterRunResult
			if err := sdk.JSONUnmarshal(btes, &filter); err != nil {
				writeError(w, r, sdk.NewError(sdk.ErrWrongRequest, err))
				return
			}
			if filter.Type == "" {
				writeError(w, r, sdk.ErrWrongRequest)
				return
			}
			response, err := wk.V2GetRunResult(ctx, filter)
			if err != nil {
				writeError(w, r, err)
				return
			}
			writeJSON(w, response, http.StatusOK)
		case http.MethodPost:
			var runResultRequest V2RunResultRequest
			if err := sdk.JSONUnmarshal(btes, &runResultRequest); err != nil {
				writeError(w, r, sdk.NewError(sdk.ErrWrongRequest, err))
				return
			}
			response, err := wk.V2AddRunResult(ctx, runResultRequest)
			if err != nil {
				writeError(w, r, err)
				return
			}
			log.Info(ctx, "run result %s created", response.RunResult.ID)
			writeJSON(w, response, http.StatusCreated)
		case http.MethodPut:
			var runResultRequest V2RunResultRequest
			if err := sdk.JSONUnmarshal(btes, &runResultRequest); err != nil {
				writeError(w, r, sdk.NewError(sdk.ErrWrongRequest, err))
				return
			}
			response, err := wk.V2UpdateRunResult(ctx, runResultRequest)
			if err != nil {
				writeError(w, r, err)
				return
			}
			log.Info(ctx, "run result %s updated", response.RunResult.ID)
			writeJSON(w, response, http.StatusOK)
		default:
			writeError(w, r, sdk.ErrNotFound)
		}
	}
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	b, _ := json.Marshal(data)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	writePlainText(w, err.Error(), 500)
}

func writePlainText(w http.ResponseWriter, data string, status int) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(data))
}
