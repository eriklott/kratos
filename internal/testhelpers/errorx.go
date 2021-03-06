package testhelpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/ory/x/logrusx"

	"github.com/ory/viper"

	"github.com/ory/herodot"

	"github.com/ory/kratos/driver/configuration"
	"github.com/ory/kratos/selfservice/errorx"
	"github.com/ory/kratos/x"
)

func NewErrorTestServer(t *testing.T, reg interface{ errorx.PersistenceProvider }) *httptest.Server {
	logger := logrusx.New("", "", logrusx.ForceLevel(logrus.TraceLevel))
	writer := herodot.NewJSONWriter(logger)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		e, err := reg.SelfServiceErrorPersister().Read(r.Context(), x.ParseUUID(r.URL.Query().Get("error")))
		require.NoError(t, err)
		t.Logf("Found error in NewErrorTestServer: %s", e.Errors)
		writer.Write(w, r, e.Errors)
	}))
	t.Cleanup(ts.Close)
	viper.Set(configuration.ViperKeyURLsError, ts.URL)
	return ts
}

func NewRedirTS(t *testing.T, body string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(body) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(ts.Close)
	viper.Set(configuration.ViperKeyURLsDefaultReturnTo, ts.URL)
	return ts
}
