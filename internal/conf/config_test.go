package conf

import (
	"testing"
)

func TestGetProjectConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		dir      string
		expected ProjectConfig
	}{
		{
			name: "returns default when no project found",
			config: &Config{
				Default: ProjectConfig{
					Name: stringPtr("default-project"),
				},
				Project: []ProjectConfig{},
			},
			dir: "project",
			expected: ProjectConfig{
				Name: stringPtr("default-project"),
			},
		},
		{
			name: "returns project when found",
			config: &Config{
				Default: ProjectConfig{
					Name: stringPtr("default-project"),
				},
				Project: []ProjectConfig{
					{
						Name: stringPtr("test-project"),
					},
				},
			},
			dir: "test-project",
			expected: ProjectConfig{
				Name: stringPtr("test-project"),
			},
		},
		{
			name: "returns default when project exists but dir doesn't match",
			config: &Config{
				Default: ProjectConfig{
					Name: stringPtr("default-project"),
				},
				Project: []ProjectConfig{
					{
						Name: stringPtr("test-project"),
					},
				},
			},
			dir: "another-project",
			expected: ProjectConfig{
				Name: stringPtr("default-project"),
			},
		},
		{
			name: "returns project with env vars when found",
			config: &Config{
				Default: ProjectConfig{
					Name: stringPtr("default-project"),
				},
				Project: []ProjectConfig{
					{
						Name: stringPtr("test-project"),
						Env: map[string]string{
							"TEST_VAR": "test_value",
						},
					},
				},
			},
			dir: "test-project",
			expected: ProjectConfig{
				Name: stringPtr("test-project"),
				Env: map[string]string{
					"TEST_VAR": "test_value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetProjectConfig(tt.dir)

			if result.Name == nil && tt.expected.Name != nil {
				t.Errorf("Expected Name to be %s, but got nil", *tt.expected.Name)
			} else if result.Name != nil && tt.expected.Name == nil {
				t.Errorf("Expected Name to be nil, but got %s", *result.Name)
			} else if result.Name != nil && tt.expected.Name != nil && *result.Name != *tt.expected.Name {
				t.Errorf("Expected Name to be %s, but got %s", *tt.expected.Name, *result.Name)
			}

			if len(result.Env) != len(tt.expected.Env) {
				t.Errorf("Expected %d env vars, but got %d", len(tt.expected.Env), len(result.Env))
			}
			for key, value := range tt.expected.Env {
				if result.Env[key] != value {
					t.Errorf("Expected env var %s to be %s, but got %s", key, value, result.Env[key])
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
