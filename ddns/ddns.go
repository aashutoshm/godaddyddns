package ddns

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// DDNS is DDNS client object
type DDNS struct {
	*http.Client
	domain string
}

// New returns new DDNS structure
func New(domain string) *DDNS {
	return &DDNS{
		Client: &http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				MaxIdleConns:        60,
				TLSHandshakeTimeout: 30 * time.Second,
				IdleConnTimeout:     30 * time.Second,
				DisableCompression:  true,
			},
			Timeout: 30 * time.Second,
		},
		domain: domain,
	}
}

func (d *DDNS) query(ctx context.Context, method, uri string, payload io.Reader) error {
	req, err = http.NewRequest(method, c.endpoint+uri, payload)
	if err != nil {
		return err
	}
	// TODO: add auth
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		return err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// TODO: parse results
	return nil
}

// Refresh refreshes subdomain dns record
func (d *DDNS) Refresh(ctx context.Context, subdomain string) error {
	req := []*RecordRequest{NewRecordRequest(subdomain)}
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		if err := json.NewEncoder(pw).Encode(req); err != nil {
			return
		}
	}()
	err := d.query(ctx, "PATCH", "/v1/domains/"+d.domain+"/records", pr)
	if err != nil {
		return err
	}
	return nil
}

// RecordRequest represents request json structure for records
// PATCH https://api.godaddy.com/v1/domains/{domain}/records
type RecordRequest struct {
	Data     string
	Name     string
	Port     int
	Priority int
	Protocol string
	Service  string
	TTL      int
	Type     string
	Weight   int
}

// NewRecordRequest returns new RecordRequest structure
func NewRecordRequest(subdomain string) *RecordRequest {
	return &RecordRequest{
		Name: subdomain,
		Type: "A",
	}
}
