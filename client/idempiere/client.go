package idempiere

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client menyimpan konfigurasi iDempiere
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient menginisialisasi client dengan konfigurasi standar
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second, // Timeout penting untuk API
		},
	}
}

// CallAPI melakukan request ke iDempiere
func (c *Client) CallAPI(ctx context.Context, method, path string, payload any) ([]byte, error) {
	// Mengambil token dari context
	token, _ := ctx.Value("token").(string)

	endpoint := fmt.Sprintf("%s/%s", c.BaseURL, strings.TrimLeft(path, "/"))

	var body io.Reader
	if payload != nil {
		bodyBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("gagal marshal payload: %w", err)
		}
		body = bytes.NewBuffer(bodyBytes)

		// [Optional] Debugging log (gunakan logger yang lebih baik di produksi)
		// log.Printf("[DEBUG] Request to %s: %s", endpoint, string(bodyBytes))
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	// MEMPERBAIKI BUG: Gunakan instance HTTPClient
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengeksekusi request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response: %w", err)
	}

	// Validasi status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("idempiere error [status %d]: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// DecodeResponse adalah helper generik untuk memproses JSON response
func DecodeResponse[T any](body []byte) (*T, error) {
	var result T
	if len(body) == 0 {
		return &result, nil
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("gagal parse response: %w", err)
	}

	return &result, nil
}
