package auth

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetBaseUrl() string
	GetInitialToken(ctx context.Context, user, pass string) (*IdempiereAuthResponse, error)
	GetRoles(ctx context.Context, token string, clientId int) (map[string]any, error)
	GetOrgs(ctx context.Context, token string, clientId, roleId int) (map[string]any, error)
	GetWarehouses(ctx context.Context, token string, clientId, roleId, orgId int) (map[string]any, error)
	UpdateContext(ctx context.Context, token string, data FinalizeRequest) (*IdempiereAuthResponse, error)
	GetUserInfo(ctx context.Context, token string) (*UserInfoResponse, error)
	GetUserById(ctx context.Context, userId int) (*UserInfoResponse, error)
	Logout(ctx context.Context, token string) error
}

type authRepo struct {
	baseUrl string
	client  *http.Client
	db      *sqlx.DB
}

func NewRepository(baseUrl string, db *sqlx.DB) Repository {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &authRepo{
		baseUrl: baseUrl,
		client:  &http.Client{Transport: tr, Timeout: 30 * time.Second},
		db:      db,
	}
}

func (r *authRepo) GetBaseUrl() string {
	return r.baseUrl
}

func (r *authRepo) GetInitialToken(ctx context.Context, user, pass string) (*IdempiereAuthResponse, error) {
	body, _ := json.Marshal(map[string]string{"userName": user, "password": pass})
	req, _ := http.NewRequestWithContext(ctx, "POST", r.baseUrl+"/auth/tokens", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth failed")
	}

	defer resp.Body.Close()

	var result IdempiereAuthResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

func (r *authRepo) UpdateContext(ctx context.Context, token string, data FinalizeRequest) (*IdempiereAuthResponse, error) {

	body, _ := json.Marshal(data)

	req, _ := http.NewRequestWithContext(
		ctx,
		"PUT",
		r.baseUrl+"/auth/tokens",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed update context: %s", string(b))
	}

	var result IdempiereAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *authRepo) GetUserInfo(ctx context.Context, username string) (*UserInfoResponse, error) {
	var user UserInfoResponse

	log.Printf("username=%s", username)

	// Query langsung ke table ad_user berdasarkan column 'value' (username)
	// Kita ambil AD_User_ID sebagai UserID dan Name sebagai UserName
	query := `
        SELECT ad_user_id as userId, name as userName 
        FROM ad_user 
        WHERE name = $1 AND isactive = 'Y' 
        LIMIT 1
    `

	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, fmt.Errorf("GetUserInfo failed: %w", err)
	}
	return &user, nil
}

func (r *authRepo) GetUserById(ctx context.Context, userId int) (*UserInfoResponse, error) {
	var user UserInfoResponse

	// Query langsung ke table ad_user berdasarkan column 'value' (username)
	// Kita ambil AD_User_ID sebagai UserID dan Name sebagai UserName
	query := `
        SELECT ad_user_id as userId, name as userName 
        FROM ad_user 
        WHERE ad_user_id = $1 AND isactive = 'Y' 
        LIMIT 1
    `

	err := r.db.GetContext(ctx, &user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("user not found or database error")
	}
	log.Printf("Berhasil ambil data user: %+v", user)

	return &user, nil
}

func (r *authRepo) GetRoles(ctx context.Context, token string, clientId int) (map[string]any, error) {
	url := fmt.Sprintf("%s/auth/roles?client=%d", r.baseUrl, clientId)
	return r.fetchRaw(ctx, token, url)
}

func (r *authRepo) GetOrgs(ctx context.Context, token string, clientId, roleId int) (map[string]any, error) {
	url := fmt.Sprintf("%s/auth/organizations?client=%d&role=%d", r.baseUrl, clientId, roleId)
	return r.fetchRaw(ctx, token, url)
}

func (r *authRepo) GetWarehouses(ctx context.Context, token string, clientId, roleId, orgId int) (map[string]any, error) {
	url := fmt.Sprintf("%s/auth/warehouses?client=%d&role=%d&organization=%d", r.baseUrl, clientId, roleId, orgId)
	return r.fetchRaw(ctx, token, url)
}

// Implementasi di authRepo
func (r *authRepo) Logout(ctx context.Context, token string) error {
	// Endpoint standar iDempiere untuk invalidate token
	url := fmt.Sprintf("%s/auth/tokens", r.baseUrl)

	req, _ := http.NewRequestWithContext(ctx, "POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// iDempiere biasanya return 204 No Content atau 200 OK
	if resp.StatusCode > 299 {
		return fmt.Errorf("idempiere returned status: %d", resp.StatusCode)
	}
	return nil
}

// Helper baru: mengembalikan map mentah dari iDempiere
func (r *authRepo) fetchRaw(ctx context.Context, token, url string) (map[string]any, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := r.client.Do(req)
	if err != nil {
		log.Printf("Err Repo: fetchRaw->resp %+v", err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("Success Repo: fetchRaw->resp %+v", resp)

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
