package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"inventory/store"
)

func server() (*Server, http.Handler) {
	st := store.New()
	st.Put(store.Item{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25})
	st.Put(store.Item{SKU: "N-1", Name: `12" nail`, Qty: 3, Price: 0.05})
	s := &Server{Store: st}
	return s, s.Routes()
}

func do(h http.Handler, method, url, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, url, strings.NewReader(body)))
	return rec
}

func TestGetItem(t *testing.T) {
	_, h := server()
	if rec := do(h, "GET", "/items/B-1", ""); rec.Code != 200 || !strings.Contains(rec.Body.String(), "bolt") {
		t.Fatalf("get: %d %s", rec.Code, rec.Body)
	}
	if rec := do(h, "GET", "/items/X", ""); rec.Code != 404 {
		t.Fatalf("missing: %d", rec.Code)
	}
}

func TestCreateConflict(t *testing.T) {
	_, h := server()
	if rec := do(h, "POST", "/items", `{"SKU":"B-1"}`); rec.Code != http.StatusConflict {
		t.Fatalf("conflict: got %d want 409", rec.Code)
	}
	if rec := do(h, "POST", "/items", `{"SKU":"C-1","Name":"cog"}`); rec.Code != http.StatusCreated {
		t.Fatalf("create: %d", rec.Code)
	}
	if rec := do(h, "POST", "/items", `{`); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad body: %d", rec.Code)
	}
}

func TestSearch(t *testing.T) {
	_, h := server()
	rec := do(h, "GET", `/search?q=name+%3D+%2212%5C%22+nail%22`, "")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "N-1") {
		t.Fatalf("search escaped: %d %s", rec.Code, rec.Body)
	}
	if rec := do(h, "GET", "/search?q=qty+%3E", ""); rec.Code != 400 {
		t.Fatalf("bad query: %d", rec.Code)
	}
}

func TestAdjustClamps(t *testing.T) {
	_, h := server()
	rec := do(h, "POST", "/items/N-1/adjust?by=-10", "")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"Qty":0`) {
		t.Fatalf("adjust: %d %s", rec.Code, rec.Body)
	}
}
