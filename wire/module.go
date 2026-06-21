package module

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// Module is the contract every feature module must implement.
// This is the core of the modular monolith: each module is
// self-contained and registers its own routes.
type Module interface {
	// Name returns the unique module identifier (e.g. "order").
	Name() string
	// RegisterRoutes mounts the module's routes onto the given router.
	RegisterRoutes(r chi.Router)
}

// Dependencies holds shared infrastructure passed to every module.
// Add more fields (Redis, Mailer, Logger…) as the project grows.
type Dependencies struct {
	DB *sqlx.DB
}

// Factory builds and wires a Module from shared Dependencies.
// This is the Factory pattern: callers request a module by name
// without knowing its concrete construction logic.
type Factory interface {
	Create(deps Dependencies) Module
}

// Registry holds all registered module factories, keyed by name.
// Modules self-register via Register() at init time.
type Registry struct {
	factories map[string]Factory
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{factories: make(map[string]Factory)}
}

// Register adds a factory under a given module name.
// Panics on duplicate registration (caught at startup, not runtime).
func (r *Registry) Register(name string, f Factory) {
	if _, exists := r.factories[name]; exists {
		panic("module: duplicate registration for " + name)
	}
	r.factories[name] = f
}

// Build instantiates every registered module with the given dependencies.
func (r *Registry) Build(deps Dependencies) []Module {
	modules := make([]Module, 0, len(r.factories))
	for _, f := range r.factories {
		modules = append(modules, f.Create(deps))
	}
	return modules
}
