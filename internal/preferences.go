package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ExportPreferences struct {
	WaterMeterFields []string `yaml:"water_meter_fields"`
}

type Preferences struct {
	Export ExportPreferences `yaml:"export"`
}

// LoadPreferences loads preferences.yaml from disk.
func LoadPreferences(path string) (*Preferences, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var prefs Preferences
	if err := yaml.Unmarshal(data, &prefs); err != nil {
		return nil, err
	}

	return &prefs, nil
}
