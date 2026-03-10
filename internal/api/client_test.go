package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetAggregate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/api/v2/query" {
			t.Errorf("path = %s, want /api/v2/query", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer test-token")
		}

		var q Query
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &q)

		if q.SiteID != "example.com" {
			t.Errorf("site_id = %q, want %q", q.SiteID, "example.com")
		}

		resp := QueryResponse{
			Results: []ResultRow{
				{Metrics: []float64{100, 150, 300, 2.0, 45.5, 120}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "example.com", "test-token")
	agg, err := c.GetAggregate("30d")
	if err != nil {
		t.Fatalf("GetAggregate: %v", err)
	}

	if agg.Visitors != 100 {
		t.Errorf("Visitors = %d, want 100", agg.Visitors)
	}
	if agg.Visits != 150 {
		t.Errorf("Visits = %d, want 150", agg.Visits)
	}
	if agg.Pageviews != 300 {
		t.Errorf("Pageviews = %d, want 300", agg.Pageviews)
	}
	if agg.BounceRate != 45.5 {
		t.Errorf("BounceRate = %f, want 45.5", agg.BounceRate)
	}
	if agg.VisitDuration != 120 {
		t.Errorf("VisitDuration = %d, want 120", agg.VisitDuration)
	}
}

func TestGetTopPages(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := QueryResponse{
			Results: []ResultRow{
				{Dimensions: []string{"/"}, Metrics: []float64{50}},
				{Dimensions: []string{"/blog"}, Metrics: []float64{30}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "example.com", "test-token")
	pages, err := c.GetTopPages("30d", 10)
	if err != nil {
		t.Fatalf("GetTopPages: %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("len(pages) = %d, want 2", len(pages))
	}
	if pages[0].Page != "/" || pages[0].Visitors != 50 {
		t.Errorf("pages[0] = %+v, want {/ 50}", pages[0])
	}
	if pages[1].Page != "/blog" || pages[1].Visitors != 30 {
		t.Errorf("pages[1] = %+v, want {/blog 30}", pages[1])
	}
}

func TestGetTopSources(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := QueryResponse{
			Results: []ResultRow{
				{Dimensions: []string{"Google"}, Metrics: []float64{80}},
				{Dimensions: []string{""}, Metrics: []float64{40}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "example.com", "test-token")
	sources, err := c.GetTopSources("30d", 10)
	if err != nil {
		t.Fatalf("GetTopSources: %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("len(sources) = %d, want 2", len(sources))
	}
	if sources[0].Source != "Google" || sources[0].Visitors != 80 {
		t.Errorf("sources[0] = %+v, want {Google 80}", sources[0])
	}
	if sources[1].Source != "" || sources[1].Visitors != 40 {
		t.Errorf("sources[1] = %+v, want { 40}", sources[1])
	}
}

func TestAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "example.com", "bad-token")
	_, err := c.GetAggregate("30d")
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

func TestGetRealtimeVisitors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if !strings.Contains(r.URL.Path, "/api/v1/stats/realtime/visitors") {
			t.Errorf("path = %s, want /api/v1/stats/realtime/visitors", r.URL.Path)
		}
		if got := r.URL.Query().Get("site_id"); got != "example.com" {
			t.Errorf("site_id = %q, want %q", got, "example.com")
		}
		w.Write([]byte("42"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "example.com", "test-token")
	count, err := c.GetRealtimeVisitors()
	if err != nil {
		t.Fatalf("GetRealtimeVisitors: %v", err)
	}
	if count != 42 {
		t.Errorf("count = %d, want 42", count)
	}
}

func TestGetTimeSeriesHourly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := QueryResponse{
			Results: []ResultRow{
				{Dimensions: []string{"2026-03-10 09:00:00"}, Metrics: []float64{25}},
				{Dimensions: []string{"2026-03-10 10:00:00"}, Metrics: []float64{40}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "example.com", "test-token")
	ts, err := c.GetTimeSeries("day", "time:hour")
	if err != nil {
		t.Fatalf("GetTimeSeries: %v", err)
	}
	if len(ts) != 2 {
		t.Fatalf("len = %d, want 2", len(ts))
	}
	if ts[0].Date != "2026-03-10 09:00:00" {
		t.Errorf("ts[0].Date = %q", ts[0].Date)
	}
	if ts[1].Visitors != 40 {
		t.Errorf("ts[1].Visitors = %d, want 40", ts[1].Visitors)
	}
}

func TestEmptyResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(QueryResponse{Results: []ResultRow{}})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "example.com", "test-token")
	agg, err := c.GetAggregate("30d")
	if err != nil {
		t.Fatalf("GetAggregate: %v", err)
	}
	if agg.Visitors != 0 {
		t.Errorf("Visitors = %d, want 0", agg.Visitors)
	}
}
