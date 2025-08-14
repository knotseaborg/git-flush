package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DiffLimit int    `yaml:"DIFF_LIMIT"`
	Model     string `yaml:"MODEL"`
	APIKey    string `yaml:"API_KEY"`
	EndPoint  string `yaml:"END_POINT"`
}

func (c *Config) Edit() error {
	// Determine editor
	var editor string
	if runtime.GOOS == "windows" {
		editor = "notepad"
	} else {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi"
		}
	}

	// Determine config file path
	path, err := getConfigPath()
	if err != nil {
		logger.Error("Failed to determine config path")
		return err
	}

	// Create intermediary dirs if necessary
	// This step is idempotent
	err = os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		logger.Error("Failed to create file path: ", err)
	}

	// Pre-fill the config file
	if err := c.Load(); err != nil {
		if !os.IsNotExist(err) {
			logger.Error("Failed to read config file. Creating new config file ", err)
		}
		if err = c.write(path); err != nil {
			logger.Error("Failed to edit config file")
		}
	}

	// Open config file in editor
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.Error("Error opening editor:", err)
		return err
	}

	fmt.Println("Config edited — hope you brought toilet paper🧻")
	return nil
}

func (c *Config) Load() error {
	// Determine config path
	path, err := getConfigPath()
	if err != nil {
		logger.Error("Failed to determine config path")
		return err
	}
	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		logger.Error("Failed to unmarshal config data: ", err)
	}
	return nil
}

func (c *Config) write(path string) error {
	if data, err := yaml.Marshal(c); err != nil {
		logger.Error("Failed to create default config file")
		logger.Info("Manually create the config file")
		return err
	} else {
		// Save updated config
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
	}
	return nil
}

func getConfigPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			logger.Error("%APPDATA% is not set. Cannot determine suitable path for config file")
		}
		return filepath.Join(appdata, "git-flush", "config.yaml"), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "git-flush", "config.yaml"), nil
	case "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".config", "git-flush", "config.yaml"), nil
	default:
		logger.Error("Incompatible OS detected!")
	}
	return "", errors.New("incompatible OS detected")
}

// DefaultConfig returns a default config
func DefaultConfig() *Config {
	// Initialize with default values
	config := &Config{
		DiffLimit: 5000,
		Model:     "gpt-4.1-nano",
		EndPoint:  "https://api.openai.com/v1/responses",
	}

	return config
}
