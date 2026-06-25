package auth

import (
	module "github.com/eonebyte/go-teus/wire"
	"github.com/go-chi/chi/v5"
)

const ModuleName = "auth"

// Module wires repository → service → handler and satisfies module.Module.
type Module struct {
	handler *handler
}

// -- Factory pattern: orderFactory builds this module. --

type orderFactory struct{}

// NewFactory returns the factory for the order module.
// Register it in the global registry at program start.
func NewFactory() module.Factory {
	return &orderFactory{}
}

func (f *orderFactory) Create(deps module.Dependencies) module.Module {
	// repo := NewRepository("https://api-erp-dev.adyawinsa.com/idempiere/api/v1", deps.DB)
	repo := NewRepository("https://192.168.3.40:8443/api/v1", deps.DB)
	svc := NewService(repo)
	h := newHandler(svc)
	return &Module{handler: h}
}

// -- module.Module implementation --

func (m *Module) Name() string { return ModuleName }

func (m *Module) RegisterRoutes(r chi.Router) {
	m.handler.routes(r)
}
