package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config is the top-level application configuration.
type Config struct {
	Theme    string         `toml:"theme"`
	Keys     KeyBindings    `toml:"keys"`
	Behavior BehaviorConfig `toml:"behavior"`
}

// BehaviorConfig controls configurable behaviors.
type BehaviorConfig struct {
	// What Enter does on a file: "edit" (default) or "preview"
	EnterAction string `toml:"enter_action"`
	// What Space does on a file: "preview" (default) or "edit"
	SpaceAction string `toml:"space_action"`
}

// KeyBindings defines all configurable key bindings.
// Each field is a string or list of strings representing key combos.
type KeyBindings struct {
	Quit        StringOrList `toml:"quit"`
	TogglePanel StringOrList `toml:"toggle_panel"`
	SwapPanels  StringOrList `toml:"swap_panels"`
	Copy        StringOrList `toml:"copy"`
	Move        StringOrList `toml:"move"`
	Mkdir       StringOrList `toml:"mkdir"`
	Delete      StringOrList `toml:"delete"`
	Rename      StringOrList `toml:"rename"`
	View        StringOrList `toml:"view"`
	Edit        StringOrList `toml:"edit"`

	// Navigation
	Up       StringOrList `toml:"up"`
	Down     StringOrList `toml:"down"`
	PageUp   StringOrList `toml:"page_up"`
	PageDown StringOrList `toml:"page_down"`
	Home     StringOrList `toml:"home"`
	End      StringOrList `toml:"end"`
	GoBack   StringOrList `toml:"go_back"`

	// Selection
	ToggleSelect StringOrList `toml:"toggle_select"`
	SelectUp     StringOrList `toml:"select_up"`
	SelectDown   StringOrList `toml:"select_down"`

	// Search
	QuickSearch StringOrList `toml:"quick_search"`

	// Go to path
	GoTo       StringOrList `toml:"goto"`
	FuzzyFind  StringOrList `toml:"fuzzy_find"`
	Bookmarks  StringOrList `toml:"bookmarks"`
	Help       StringOrList `toml:"help"`
}

// StringOrList can unmarshal from either a single string or a list of strings.
type StringOrList []string

func (s *StringOrList) UnmarshalTOML(data any) error {
	switch v := data.(type) {
	case string:
		*s = []string{v}
	case []any:
		for _, item := range v {
			if str, ok := item.(string); ok {
				*s = append(*s, str)
			}
		}
	}
	return nil
}

// Default returns a config with all defaults.
func Default() Config {
	return Config{
		Theme: "",
		Behavior: BehaviorConfig{
			EnterAction: "edit",
			SpaceAction: "preview",
		},
		Keys: DefaultKeyBindings(),
	}
}

// DefaultKeyBindings returns the default key bindings.
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Quit:        StringOrList{"f10", "ctrl+c"},
		TogglePanel: StringOrList{"tab"},
		SwapPanels:  StringOrList{"ctrl+u"},
		Copy:        StringOrList{"f5"},
		Move:        StringOrList{"f6"},
		Mkdir:       StringOrList{"f7"},
		Delete:      StringOrList{"f8"},
		Rename:      StringOrList{"shift+f6"},
		View:        StringOrList{"f3"},
		Edit:        StringOrList{"f4"},

		Up:       StringOrList{"up", "k"},
		Down:     StringOrList{"down", "j"},
		PageUp:   StringOrList{"pgup"},
		PageDown: StringOrList{"pgdown"},
		Home:     StringOrList{"home"},
		End:      StringOrList{"end"},
		GoBack:   StringOrList{"backspace"},

		ToggleSelect: StringOrList{"ctrl+t"},
		SelectUp:     StringOrList{"shift+up"},
		SelectDown:   StringOrList{"shift+down"},

		QuickSearch: StringOrList{"ctrl+s"},

		GoTo:      StringOrList{"ctrl+g"},
		FuzzyFind: StringOrList{"f9", "ctrl+p"},
		Bookmarks: StringOrList{"f2", "ctrl+b"},
		Help:      StringOrList{"f1"},
	}
}

// Load reads config from ~/.config/mdc/config.toml, merging with defaults.
func Load() Config {
	cfg := Default()

	configPath := filepath.Join(configDirPath(), "config.toml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg
	}

	var fileCfg Config
	if err := toml.Unmarshal(data, &fileCfg); err != nil {
		return cfg
	}

	// Merge: only override non-zero values
	if fileCfg.Theme != "" {
		cfg.Theme = fileCfg.Theme
	}
	if fileCfg.Behavior.EnterAction != "" {
		cfg.Behavior.EnterAction = fileCfg.Behavior.EnterAction
	}
	if fileCfg.Behavior.SpaceAction != "" {
		cfg.Behavior.SpaceAction = fileCfg.Behavior.SpaceAction
	}

	mergeKeys(&cfg.Keys, &fileCfg.Keys)

	return cfg
}

func mergeKeys(dst, src *KeyBindings) {
	mergeKey(&dst.Quit, src.Quit)
	mergeKey(&dst.TogglePanel, src.TogglePanel)
	mergeKey(&dst.SwapPanels, src.SwapPanels)
	mergeKey(&dst.Copy, src.Copy)
	mergeKey(&dst.Move, src.Move)
	mergeKey(&dst.Mkdir, src.Mkdir)
	mergeKey(&dst.Delete, src.Delete)
	mergeKey(&dst.Rename, src.Rename)
	mergeKey(&dst.View, src.View)
	mergeKey(&dst.Edit, src.Edit)
	mergeKey(&dst.Up, src.Up)
	mergeKey(&dst.Down, src.Down)
	mergeKey(&dst.PageUp, src.PageUp)
	mergeKey(&dst.PageDown, src.PageDown)
	mergeKey(&dst.Home, src.Home)
	mergeKey(&dst.End, src.End)
	mergeKey(&dst.GoBack, src.GoBack)
	mergeKey(&dst.ToggleSelect, src.ToggleSelect)
	mergeKey(&dst.SelectUp, src.SelectUp)
	mergeKey(&dst.SelectDown, src.SelectDown)
	mergeKey(&dst.QuickSearch, src.QuickSearch)
	mergeKey(&dst.GoTo, src.GoTo)
	mergeKey(&dst.FuzzyFind, src.FuzzyFind)
	mergeKey(&dst.Bookmarks, src.Bookmarks)
	mergeKey(&dst.Help, src.Help)
}

func mergeKey(dst *StringOrList, src StringOrList) {
	if len(src) > 0 {
		*dst = src
	}
}

// ConfigDir returns the mdc config directory path.
func ConfigDir() string {
	return configDirPath()
}

// configDirPath returns ~/.config/mdc, respecting XDG_CONFIG_HOME.
func configDirPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "mdc")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "mdc")
}
