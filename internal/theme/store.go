// Package theme is the package that manages themes
package theme

type ThemeConfig struct {
	ThemeDefaults   map[string]string            `json:"defaults"`
	TargetOverrides map[string]map[string]string `json:"overrides,omitempty"`
}

type Store struct {
	Cfg ThemeConfig
}

func NewStore(seed ThemeConfig) *Store {
	if seed.ThemeDefaults == nil {
		seed.ThemeDefaults = map[string]string{}
	}
	if seed.TargetOverrides == nil {
		seed.TargetOverrides = map[string]map[string]string{}
	}
	return &Store{Cfg: seed}
}

func (s *Store) Resolve(targetID, fieldKey, fieldDefault string) string {
	if v := s.GetOverride(targetID, fieldKey); v != "" {
		return v
	}
	if v, ok := s.Cfg.ThemeDefaults[fieldKey]; ok {
		return v
	}
	return fieldDefault
}

func (s *Store) GetOverride(targetID, fieldKey string) string {
	if m := s.Cfg.TargetOverrides[targetID]; m != nil {
		return m[fieldKey]
	}
	return ""
}

func (s *Store) HasOverride(targetID, fieldKey string) bool {
	return s.GetOverride(targetID, fieldKey) != ""
}

func (s *Store) HasDefault(fieldKey string) bool {
	_, ok := s.Cfg.ThemeDefaults[fieldKey]
	return ok
}

func (s *Store) SetOverride(targetID, fieldKey, val string) {
	if def, ok := s.Cfg.ThemeDefaults[fieldKey]; ok && def == val {
		if m := s.Cfg.TargetOverrides[targetID]; m != nil {
			delete(m, fieldKey)
			if len(m) == 0 {
				delete(s.Cfg.TargetOverrides, targetID)
			}
		}
		return
	}
	if s.Cfg.TargetOverrides[targetID] == nil {
		s.Cfg.TargetOverrides[targetID] = map[string]string{}
	}
	s.Cfg.TargetOverrides[targetID][fieldKey] = val
}
