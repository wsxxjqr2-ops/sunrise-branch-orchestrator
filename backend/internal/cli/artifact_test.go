package cli

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestArtifactReadFetchesRawSourceFromAnotherSession(t *testing.T) {
	cfg := setConfigEnv(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/sessions/architect-1/artifact-files/rfc/design.md" || r.URL.Query().Get("source") != "1" {
			t.Fatalf("request = %s %s, want raw artifact path", r.Method, r.URL.String())
		}
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		_, _ = io.WriteString(w, "# architecture\n")
	}))
	defer srv.Close()
	writeRunFileFor(t, cfg, srv)

	var out bytes.Buffer
	c := &commandContext{deps: Deps{Out: &out, HTTPClient: &http.Client{}, ProcessAlive: func(int) bool { return true }}}
	if err := c.readArtifact(t.Context(), artifactReadOptions{session: "architect-1", path: "rfc/design.md"}); err != nil {
		t.Fatalf("readArtifact: %v", err)
	}
	if got := out.String(); got != "# architecture\n" {
		t.Fatalf("artifact output = %q, want raw RFC", got)
	}
}

func TestArtifactReadRejectsEscapingSourcePath(t *testing.T) {
	c := &commandContext{deps: Deps{Out: io.Discard}}
	err := c.readArtifact(t.Context(), artifactReadOptions{session: "architect-1", path: "../secret.md"})
	if err == nil || !strings.Contains(err.Error(), "stay inside") {
		t.Fatalf("error = %v, want path boundary validation", err)
	}
}
