package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	svc Service
}

func newHandler(svc Service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) routes(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login/init", h.HandleInit)
		r.Get("/user", h.HandleGetUserInfo)
		r.Get("/roles", h.HandleGetRoles)

		r.Put("/login/finalize", h.HandleFinalize)
		r.Get("/organizations", h.HandleGetOrgs)
		r.Get("/warehouses", h.HandleGetWarehouses)

		r.Post("/logout", h.HandleLogout)
	})
}

func (h *handler) HandleInit(w http.ResponseWriter, r *http.Request) {
	var input struct {
		User     string `json:"userName"`
		Pass     string `json:"password"`
		ClientID int    `json:"clientId"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	res, err := h.svc.LoginStepOne(r.Context(), input.User, input.Pass, input.ClientID)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *handler) HandleGetRoles(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) > 7 {
		token = token[7:]
	}

	cID, _ := strconv.Atoi(r.URL.Query().Get("client"))

	// Panggil repo GetRoles yang sudah kita buat tadi
	res, err := h.svc.GetRoles(r.Context(), token, cID)

	if err != nil {
		http.Error(w, "Failed to fetch roles", 500)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *handler) HandleGetOrgs(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")[7:]
	cID, _ := strconv.Atoi(r.URL.Query().Get("client"))
	rID, _ := strconv.Atoi(r.URL.Query().Get("role"))

	res, _ := h.svc.GetOrgs(r.Context(), token, cID, rID)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) HandleGetWarehouses(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")[7:]
	cID, _ := strconv.Atoi(r.URL.Query().Get("client"))
	rID, _ := strconv.Atoi(r.URL.Query().Get("role"))
	oID, _ := strconv.Atoi(r.URL.Query().Get("organization"))

	res, _ := h.svc.GetWarehouses(r.Context(), token, cID, rID, oID)
	json.NewEncoder(w).Encode(res)
}

func (h *handler) HandleFinalize(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token == "" {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var req FinalizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.svc.Finalize(r.Context(), token, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (h *handler) HandleGetUserInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := ""
	if len(authHeader) > 7 {
		token = authHeader[7:]
	}

	if token == "" {
		log.Println("[AUTH ERROR] Request ditolak: Token kosong atau tidak disertakan")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	userId, _ := strconv.Atoi(r.URL.Query().Get("userId"))

	// Panggil service untuk ambil data user
	user, err := h.svc.GetUserById(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil token dari header
	authHeader := r.Header.Get("Authorization")
	token := ""
	if len(authHeader) > 7 {
		token = authHeader[7:] // Menghapus tulisan "Bearer "
	}

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token is required"})
		return
	}

	// 2. Panggil Service
	err := h.svc.Logout(r.Context(), token)
	if err != nil {
		// Kita log errornya tapi tetap kasih response sukses ke user
		// supaya frontend tetap bersih (karena token di frontend pasti dibuang)
		fmt.Printf("Logout error: %v\n", err)
	}

	// 3. Response Sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}
