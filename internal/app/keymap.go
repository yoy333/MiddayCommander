package app

import (
	"github.com/charmbracelet/bubbles/key"

	"github.com/kooler/MiddayCommander/internal/config"
)

// KeyMap defines all global keybindings.
type KeyMap struct {
	Quit        key.Binding
	TogglePanel key.Binding
	SwapPanels  key.Binding
	Copy        key.Binding
	Move        key.Binding
	Mkdir       key.Binding
	Delete      key.Binding
	Rename      key.Binding
	View        key.Binding
	Edit        key.Binding
	GoTo        key.Binding
	FuzzyFind   key.Binding
	Bookmarks   key.Binding
	Help        key.Binding
}

// PanelKeyMap defines per-panel keybindings.
type PanelKeyMap struct {
	Up           key.Binding
	Down         key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	Home         key.Binding
	End          key.Binding
	GoBack       key.Binding
	ToggleSelect key.Binding
	SelectUp     key.Binding
	SelectDown   key.Binding
	QuickSearch  key.Binding
}

// KeyMapFromConfig builds the global keymap from config.
func KeyMapFromConfig(keys config.KeyBindings) KeyMap {
	return KeyMap{
		Quit:        binding(keys.Quit, "quit"),
		TogglePanel: binding(keys.TogglePanel, "switch panel"),
		SwapPanels:  binding(keys.SwapPanels, "swap panels"),
		Copy:        binding(keys.Copy, "copy"),
		Move:        binding(keys.Move, "move"),
		Mkdir:       binding(keys.Mkdir, "mkdir"),
		Delete:      binding(keys.Delete, "delete"),
		Rename:      binding(keys.Rename, "rename"),
		View:        binding(keys.View, "view"),
		Edit:        binding(keys.Edit, "edit"),
		GoTo:        binding(keys.GoTo, "go to"),
		FuzzyFind:   binding(keys.FuzzyFind, "find"),
		Bookmarks:   binding(keys.Bookmarks, "bookmarks"),
		Help:        binding(keys.Help, "help"),
	}
}

// PanelKeyMapFromConfig builds the panel keymap from config.
func PanelKeyMapFromConfig(keys config.KeyBindings) PanelKeyMap {
	return PanelKeyMap{
		Up:           binding(keys.Up, "up"),
		Down:         binding(keys.Down, "down"),
		PageUp:       binding(keys.PageUp, "page up"),
		PageDown:     binding(keys.PageDown, "page down"),
		Home:         binding(keys.Home, "home"),
		End:          binding(keys.End, "end"),
		GoBack:       binding(keys.GoBack, "go back"),
		ToggleSelect: binding(keys.ToggleSelect, "toggle select"),
		SelectUp:     binding(keys.SelectUp, "select up"),
		SelectDown:   binding(keys.SelectDown, "select down"),
		QuickSearch:  binding(keys.QuickSearch, "quick search"),
	}
}

func binding(keys config.StringOrList, help string) key.Binding {
	return key.NewBinding(
		key.WithKeys([]string(keys)...),
		key.WithHelp(keys[0], help),
	)
}
