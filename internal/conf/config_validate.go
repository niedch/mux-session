package conf

import "errors"

func validateConfig(conf *Config) error {
	for _, project := range conf.Project {
		if err := validateProjectConfig(project); err != nil {
			return err
		}
	}

	return validateProjectConfig(conf.Default)
}

func validateProjectConfig(project ProjectConfig) error {
	if err := validatePrimaryMarker(project); err != nil {
		return err
	}

	if err := validatePanelConfig(project); err != nil {
		return err
	}

	return nil
}

func validatePrimaryMarker(project ProjectConfig) error {
	primaryCount := 0
	for _, window := range project.WindowConfig {
		// Check if this window is marked as primary
		if window.Primary != nil && *window.Primary {
			primaryCount++
		}

		// Validate that only one window is marked as primary
		if primaryCount > 1 {
			return errors.New("only one window can be marked as primary in project configuration")
		}
	}

	return nil
}

func validatePanelConfig(project ProjectConfig) error {
	for _, window := range project.WindowConfig {
		for _, panel := range window.PanelConfig {
			if panel.PanelDirection != "v" && panel.PanelDirection != "h" {
				return errors.New("panel_direction must be 'v' or 'h'")
			}
		}
	}

	return nil
}
