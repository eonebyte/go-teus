package auth

import (
	"context"
	"log"
)

type Service interface {
	LoginStepOne(ctx context.Context, user, pass string, clientId int) (*LoginSessionResponse, error)
	Finalize(ctx context.Context, token string, req FinalizeRequest) (*IdempiereAuthResponse, error)
	GetUserInfo(ctx context.Context, token string) (*UserInfoResponse, error)
	GetUserById(ctx context.Context, userId int) (*UserInfoResponse, error)
	GetRoles(ctx context.Context, token string, clientId int) (map[string]any, error)
	GetOrgs(ctx context.Context, token string, clientId, roleId int) (map[string]any, error)
	GetWarehouses(ctx context.Context, token string, clientId, roleId, orgId int) (map[string]any, error)
	Logout(ctx context.Context, token string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) LoginStepOne(ctx context.Context, user, pass string, clientId int) (*LoginSessionResponse, error) {

	// 1. Dapatkan Token Awal
	auth, err := s.repo.GetInitialToken(ctx, user, pass)
	if err != nil {
		log.Printf("LoginStepOne->GetInitialToken: %v", err)
		return nil, err
	}

	log.Printf("data user: %+v", user)
	// 2. Gunakan token tersebut untuk mendapatkan detail User (ID & Name)
	userInfo, err := s.repo.GetUserInfo(ctx, user)
	if err != nil {
		log.Printf("LoginStepOne->GetUserInfo: %+v", err)
		log.Printf("LoginStepOne->GetUserInfo: %+v", err)
		// Jika gagal ambil info DB, kita bisa gagalkan atau lanjut dengan info kosong
		return nil, err
	}

	// 3. Ambil Roles (Opsional, jika ingin langsung dikirim di step 1)
	roles, err := s.repo.GetRoles(ctx, auth.Token, clientId)
	if err != nil {
		log.Printf("LoginStepOne->GetRoles: %+v", err)
		return nil, err
	}

	// 4. Bungkus semua data untuk Frontend
	var uid int
	var uname string
	if userInfo != nil {
		uid = userInfo.UserID
		uname = userInfo.UserName
	}

	return &LoginSessionResponse{
		Token:        auth.Token,
		RefreshToken: auth.RefreshToken,
		UserID:       uid,
		UserName:     uname, // Tambahkan ini di struct response kamu
		Roles:        roles,
	}, nil
}

func (s *service) Finalize(ctx context.Context, token string, req FinalizeRequest) (*IdempiereAuthResponse, error) {

	resp, err := s.repo.UpdateContext(ctx, token, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Pastikan return type-nya *UserInfoResponse, bukan map
func (s *service) GetUserInfo(ctx context.Context, token string) (*UserInfoResponse, error) {
	// Panggil repo yang mengembalikan (*UserInfoResponse, error)
	user, err := s.repo.GetUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetUserById(ctx context.Context, userId int) (*UserInfoResponse, error) {
	// Panggil repo yang mengembalikan (*UserInfoResponse, error)
	user, err := s.repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetRoles(ctx context.Context, token string, clientId int) (map[string]any, error) {
	// Panggil repo yang mengembalikan (*UserInfoResponse, error)
	roles, err := s.repo.GetRoles(ctx, token, clientId)
	if err != nil {
		log.Printf("Service->GetRoles: %v", err)
		return nil, err
	}

	return roles, nil
}

func (s *service) GetOrgs(ctx context.Context, token string, clientId, roleId int) (map[string]any, error) {
	// Panggil repo yang mengembalikan (*UserInfoResponse, error)
	roles, err := s.repo.GetOrgs(ctx, token, clientId, roleId)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *service) GetWarehouses(ctx context.Context, token string, clientId, roleId, orgId int) (map[string]any, error) {
	// Panggil repo yang mengembalikan (*UserInfoResponse, error)
	roles, err := s.repo.GetWarehouses(ctx, token, clientId, roleId, orgId)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	// Kamu bisa tambah logic di sini, misal: Log user siapa yang logout
	return s.repo.Logout(ctx, token)
}
