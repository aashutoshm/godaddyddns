package ddns

import (
	"bytes"
	"context"
	json "encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	prodEndpoint       = "https://api.godaddy.com"
	oteEndpoint        = "https://api.ote-godaddy.com"
	externalIPEndpoint = "http://ipv4.icanhazip.com"
)

// DDNS is DDNS client object
type DDNS struct {
	*http.Client
	domain   string
	endpoint string
}

// New returns new DDNS structure
func New(domain string) *DDNS {
	endpoint := prodEndpoint
	if viper.GetBool("ote") {
		endpoint = oteEndpoint
	}
	log.Printf("using godaddy api endpoint %s", endpoint)
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
		domain:   domain,
		endpoint: endpoint,
	}
}

func (d *DDNS) query(ctx context.Context, method, uri string, payload io.Reader) ([]byte, error) {
	log.Printf("querying %s %s", method, uri)
	req, err := http.NewRequest(method, d.endpoint, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "sso-key "+viper.GetString("access_key")+":"+viper.GetString("secret_key"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("StatusCode: %s", http.StatusText(resp.StatusCode))

	return ioutil.ReadAll(resp.Body)
}

// Refresh refreshes subdomain dns record
func (d *DDNS) Refresh(ctx context.Context, subdomain string) error {
	ip, err := d.GetExternalIP(ctx)
	if err != nil {
		return err
	}

	log.Printf("Updating %s.%s of type A to %s", subdomain, d.domain, ip)

	req := []*RecordRequest{NewRecordRequest(subdomain, ip)}
	uri := "/v1/domains/" + d.domain + "/records/A/" + subdomain

	payload := &bytes.Buffer{}
	encoder := json.NewEncoder(payload)
	err = encoder.Encode(req)
	if err != nil {
		return err
	}
	res, err := d.query(ctx, "PUT", uri, payload)
	if err != nil {
		return err
	}
	log.Printf("response: %s", string(res))
	return nil
}

// GetCurrentIP returns IP assigned to subdomain
func (d *DDNS) GetCurrentIP(ctx context.Context, subdomain string) ([]byte, error) {
	uri := "/v1/domains/" + d.domain + "/records/A/"
	return d.query(ctx, "GET", uri, nil)
}

// GetExternalIP retrieves external IP
func (d *DDNS) GetExternalIP(ctx context.Context) (string, error) {
	req, err := http.NewRequest("GET", externalIPEndpoint, nil)
	if err != nil {
		return "", err
	}
	resp, err := d.Do(req.WithContext(ctx))
	if err != nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
		return "", err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// RecordRequest represents request json structure for records
// PUT https://api.godaddy.com/v1/domains/{domain}/records/A/{subdomain}
type RecordRequest struct {
	Data string // ip
	TTL  int
}

// NewRecordRequest returns new RecordRequest structure
func NewRecordRequest(subdomain, ip string) *RecordRequest {
	return &RecordRequest{
		Data: ip,
		TTL:  1800,
	}
}
