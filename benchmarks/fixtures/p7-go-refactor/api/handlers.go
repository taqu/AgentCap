// Package api exposes the inventory store over HTTP.
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"inventory/query"
	"inventory/store"
	"inventory/util"
)

type Server struct {
	Store *store.Store
	Quota func(r *http.Request) bool
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /items/{sku}", s.getItem)
	mux.HandleFunc("POST /items", s.createItem)
	mux.HandleFunc("GET /search", s.search)
	mux.HandleFunc("POST /items/{sku}/adjust", s.adjust)
	return mux
}

func writeError(w http.ResponseWriter, c Code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(HTTPStatus(c))
	json.NewEncoder(w).Encode(map[string]string{"code": string(c), "error": msg})
}

func (s *Server) getItem(w http.ResponseWriter, r *http.Request) {
	if s.Quota != nil && !s.Quota(r) {
		writeError(w, CodeQuota, "quota exceeded")
		return
	}
	it, err := s.Store.Get(r.Context(), r.PathValue("sku"))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, CodeNotFound, "no such item")
		return
	}
	if err != nil {
		writeError(w, CodeInternal, err.Error())
		return
	}
	json.NewEncoder(w).Encode(it)
}

func (s *Server) createItem(w http.ResponseWriter, r *http.Request) {
	var it store.Item
	if err := json.NewDecoder(r.Body).Decode(&it); err != nil || it.SKU == "" {
		writeError(w, CodeBadRequest, "invalid item")
		return
	}
	if _, err := s.Store.Get(r.Context(), it.SKU); err == nil {
		writeError(w, CodeConflict, "sku exists")
		return
	}
	s.Store.Put(it)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	e, err := query.Parse(r.URL.Query().Get("q"))
	if err != nil {
		writeError(w, CodeBadQuery, err.Error())
		return
	}
	var out []store.Item
	for _, it := range s.Store.List() {
		if query.Match(e, it) {
			out = append(out, it)
		}
	}
	json.NewEncoder(w).Encode(out)
}

func (s *Server) adjust(w http.ResponseWriter, r *http.Request) {
	it, err := s.Store.Get(r.Context(), r.PathValue("sku"))
	if err != nil {
		writeError(w, CodeNotFound, "no such item")
		return
	}
	delta, err := strconv.Atoi(r.URL.Query().Get("by"))
	if err != nil {
		writeError(w, CodeBadRequest, "by must be an integer")
		return
	}
	it.Qty = util.Clamp(it.Qty+delta, 0, 1_000_000)
	s.Store.Put(it)
	json.NewEncoder(w).Encode(it)
}
