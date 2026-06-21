package auth

// Response dasar dari iDempiere
type IdempiereAuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// Digunakan untuk Role, Org, dan Warehouse
type AuthOption struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Response awal ke Frontend
type LoginSessionResponse struct {
	Token        string         `json:"token"`
	RefreshToken string         `json:"refreshToken"`
	UserID       int            `json:"userId"`
	UserName     string         `json:"userName"`
	Roles        map[string]any `json:"roles"`
}

type UserInfoResponse struct {
	// Tambahkan tag db agar sqlx bisa mapping alias dari SQL
	UserID   int    `db:"userid" json:"userId"`
	UserName string `db:"username" json:"userName"`
}

// Payload untuk PUT Context iDempiere
type FinalizeRequest struct {
	ClientID       int    `json:"clientId"`
	RoleID         int    `json:"roleId"`
	OrganizationID int    `json:"organizationId"`
	WarehouseID    int    `json:"warehouseId"`
	Language       string `json:"language"`
}

type FinalizeResponse struct {
	Token string `json:"token"`
}
