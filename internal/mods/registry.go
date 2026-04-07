package mods

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type ModType string

const (
	ModTypeCoreContracts ModType = "core-contracts"
	ModTypeRuleset       ModType = "ruleset"
	ModTypeExpansion     ModType = "expansion"
)

type Dependency struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type Entrypoints struct {
	Data      []string `json:"data,omitempty"`
	Schemas   []string `json:"schemas,omitempty"`
	Templates []string `json:"templates,omitempty"`
	Locales   []string `json:"locales,omitempty"`
}

type Compatibility struct {
	App string `json:"app"`
	API string `json:"api"`
}

type Runtime struct {
	RulesEngine string `json:"rules_engine,omitempty"`
}

type Legal struct {
	License                 string   `json:"license"`
	OGL                     bool     `json:"ogl"`
	Notes                   string   `json:"notes,omitempty"`
	SRDOnly                 bool     `json:"srd_only,omitempty"`
	ProductIdentityIncluded bool     `json:"product_identity_included,omitempty"`
	Notices                 []string `json:"notices,omitempty"`
}

type Manifest struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Version       string        `json:"version"`
	Type          ModType       `json:"type"`
	Ruleset       string        `json:"ruleset,omitempty"`
	Description   string        `json:"description"`
	Dependencies  []Dependency  `json:"dependencies"`
	Entrypoints   Entrypoints   `json:"entrypoints"`
	Compatibility Compatibility `json:"compatibility"`
	Runtime       Runtime       `json:"runtime,omitempty"`
	Legal         Legal         `json:"legal"`
}

type Mod interface {
	Manifest() Manifest
}

type RulesEngine interface {
	RulesetID() string
	Execute(action string, payload map[string]any) (any, error)
}

type RulesetMod interface {
	Mod
	Rules() RulesEngine
}

type Registry struct {
	currentLoadedMods []Mod
	modsByID          map[string]Mod
}

var localeCodePattern = regexp.MustCompile(`^[A-Z]{2}$`)

func NewRegistry() *Registry {
	return &Registry{
		currentLoadedMods: make([]Mod, 0),
		modsByID:          make(map[string]Mod),
	}
}

func (r *Registry) RegisterMod(mod Mod) error {
	m := mod.Manifest()
	if err := m.ValidateBasic(); err != nil {
		return err
	}
	if _, exists := r.modsByID[m.ID]; exists {
		return fmt.Errorf("manifest: duplicate mod id %q", m.ID)
	}
	r.currentLoadedMods = append(r.currentLoadedMods, mod)
	r.modsByID[m.ID] = mod
	return nil
}

func (r *Registry) RegisterStartupMods(mods []Mod) error {
	for i, mod := range mods {
		m := mod.Manifest()
		if err := r.RegisterMod(mod); err != nil {
			return fmt.Errorf("manifest: startup register failed at index %d (id=%q): %w", i, m.ID, err)
		}
	}

	return nil
}

func (r *Registry) Mod(id string) (Mod, bool) {
	m, ok := r.modsByID[id]
	return m, ok
}

func (r *Registry) Manifest(id string) (Manifest, bool) {
	mod, ok := r.modsByID[id]
	if !ok {
		return Manifest{}, false
	}
	return mod.Manifest(), true
}

func (r *Registry) LoadedMods() []Mod {
	out := make([]Mod, len(r.currentLoadedMods))
	copy(out, r.currentLoadedMods)
	return out
}

func (r *Registry) LoadedManifests() []Manifest {
	out := make([]Manifest, 0, len(r.currentLoadedMods))
	for _, mod := range r.currentLoadedMods {
		out = append(out, mod.Manifest())
	}
	return out
}

func (m Manifest) ValidateBasic() error {
	if m.ID == "" {
		return fmt.Errorf("manifest: id is required")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest: version is required")
	}
	if m.Type == "" {
		return fmt.Errorf("manifest: type is required")
	}
	switch m.Type {
	case ModTypeCoreContracts, ModTypeRuleset, ModTypeExpansion:
	default:
		return fmt.Errorf("manifest: unsupported type %q", m.Type)
	}
	if m.Type == ModTypeExpansion && m.Ruleset == "" {
		return fmt.Errorf("manifest: ruleset is required for expansion mods")
	}
	if m.Type == ModTypeRuleset && m.Runtime.RulesEngine == "" {
		return fmt.Errorf("manifest: runtime.rules_engine is required for ruleset mods")
	}
	if m.Type == ModTypeExpansion && m.Runtime.RulesEngine != "" {
		return fmt.Errorf("manifest: expansion mods cannot declare runtime.rules_engine")
	}
	if m.Compatibility.App == "" || m.Compatibility.API == "" {
		return fmt.Errorf("manifest: compatibility.app and compatibility.api are required")
	}
	if m.Legal.License == "" {
		return fmt.Errorf("manifest: legal.license is required")
	}

	for i, dep := range m.Dependencies {
		if dep.ID == "" || dep.Version == "" {
			return fmt.Errorf("manifest: dependency[%d] requires id and version", i)
		}
	}
	if err := validateLocaleEntrypoints(m.Entrypoints.Locales); err != nil {
		return err
	}

	return nil
}

func validateLocaleEntrypoints(paths []string) error {
	for i, p := range paths {
		clean := filepath.ToSlash(filepath.Clean(p))
		parts := strings.Split(clean, "/")

		localeIdx := -1
		for idx, part := range parts {
			if part == "locales" {
				localeIdx = idx
				break
			}
		}

		if localeIdx == -1 || localeIdx+1 >= len(parts) {
			return fmt.Errorf("manifest: locale entrypoint[%d] %q must match locales/<LOCALE>/...", i, p)
		}

		locale := parts[localeIdx+1]
		if !localeCodePattern.MatchString(locale) {
			return fmt.Errorf("manifest: locale entrypoint[%d] %q uses invalid locale %q (expected two-letter uppercase code, e.g. EN/ES/DE)", i, p, locale)
		}
	}

	return nil
}
