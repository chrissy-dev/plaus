package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	SiteID     string
	Token      string
	HTTPClient *http.Client
}

func NewClient(baseURL, siteID, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		SiteID:     siteID,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

// Query request/response types for /api/v2/query

type Query struct {
	SiteID     string     `json:"site_id"`
	Metrics    []string   `json:"metrics"`
	DateRange  any        `json:"date_range"`
	Dimensions []string   `json:"dimensions,omitempty"`
	Filters    []any      `json:"filters,omitempty"`
	OrderBy    [][]any    `json:"order_by,omitempty"`
	Pagination *Paginate  `json:"pagination,omitempty"`
}

type Paginate struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type QueryResponse struct {
	Results []ResultRow `json:"results"`
}

type ResultRow struct {
	Dimensions []string  `json:"dimensions"`
	Metrics    []float64 `json:"metrics"`
}

// Typed results for convenience

type Aggregate struct {
	Visitors      int
	Visits        int
	Pageviews     int
	ViewsPerVisit float64
	BounceRate    float64
	VisitDuration int
}

type PageStats struct {
	Page     string
	Visitors int
}

type SourceStats struct {
	Source   string
	Visitors int
}

type TimeSeriesPoint struct {
	Date     string
	Visitors int
}

func (c *Client) query(q Query) (*QueryResponse, error) {
	q.SiteID = c.SiteID

	body, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/api/v2/query", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var result QueryResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetAggregate(dateRange any) (Aggregate, error) {
	resp, err := c.query(Query{
		Metrics:   []string{"visitors", "visits", "pageviews", "views_per_visit", "bounce_rate", "visit_duration"},
		DateRange: dateRange,
	})
	if err != nil {
		return Aggregate{}, err
	}
	if len(resp.Results) == 0 {
		return Aggregate{}, nil
	}
	m := resp.Results[0].Metrics
	return Aggregate{
		Visitors:      int(m[0]),
		Visits:        int(m[1]),
		Pageviews:     int(m[2]),
		ViewsPerVisit: m[3],
		BounceRate:    m[4],
		VisitDuration: int(m[5]),
	}, nil
}

func (c *Client) GetTopPages(dateRange any, limit int) ([]PageStats, error) {
	resp, err := c.query(Query{
		Metrics:    []string{"visitors"},
		DateRange:  dateRange,
		Dimensions: []string{"event:page"},
		OrderBy:    [][]any{{"visitors", "desc"}},
		Pagination: &Paginate{Limit: limit},
	})
	if err != nil {
		return nil, err
	}
	pages := make([]PageStats, len(resp.Results))
	for i, r := range resp.Results {
		pages[i] = PageStats{
			Page:     r.Dimensions[0],
			Visitors: int(r.Metrics[0]),
		}
	}
	return pages, nil
}

func (c *Client) GetTopSources(dateRange any, limit int) ([]SourceStats, error) {
	resp, err := c.query(Query{
		Metrics:    []string{"visitors"},
		DateRange:  dateRange,
		Dimensions: []string{"visit:source"},
		OrderBy:    [][]any{{"visitors", "desc"}},
		Pagination: &Paginate{Limit: limit},
	})
	if err != nil {
		return nil, err
	}
	sources := make([]SourceStats, len(resp.Results))
	for i, r := range resp.Results {
		sources[i] = SourceStats{
			Source:   r.Dimensions[0],
			Visitors: int(r.Metrics[0]),
		}
	}
	return sources, nil
}

func (c *Client) GetTimeSeries(dateRange any, timeDimension string) ([]TimeSeriesPoint, error) {
	resp, err := c.query(Query{
		Metrics:    []string{"visitors"},
		DateRange:  dateRange,
		Dimensions: []string{timeDimension},
	})
	if err != nil {
		return nil, err
	}
	points := make([]TimeSeriesPoint, len(resp.Results))
	for i, r := range resp.Results {
		points[i] = TimeSeriesPoint{
			Date:     r.Dimensions[0],
			Visitors: int(r.Metrics[0]),
		}
	}
	return points, nil
}

func (c *Client) GetRealtimeVisitors() (int, error) {
	req, err := http.NewRequest("GET", c.BaseURL+"/api/v1/stats/realtime/visitors?site_id="+c.SiteID, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var count int
	if err := json.Unmarshal(body, &count); err != nil {
		return 0, err
	}
	return count, nil
}
