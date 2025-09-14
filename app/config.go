package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	data map[string]string
}

func LoadConfig(path string) *Config {
	cfg := &Config{data: make(map[string]string)}

	// 1️⃣ Load .env (اگر وجود داشت)
	_ = godotenv.Load() // اگر فایل .env نبود، نادیده گرفته میشه

	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			cfg.data[parts[0]] = parts[1]
		}
	}

	// 2️⃣ Load YAML اگر مسیر داده شده بود
	if path != "" {
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			log.Printf("⚠️ Failed to read YAML file: %v", err)
		} else {
			yamlData := make(map[string]interface{})
			if err := yaml.Unmarshal(yamlFile, &yamlData); err != nil {
				log.Printf("⚠️ Failed to parse YAML: %v", err)
			} else {
				flattenYAML("", yamlData, cfg.data)
			}
		}
	}

	return cfg
}

// flattenYAML نقشه‌های تو در تو را به key=value مسطح تبدیل می‌کند
func flattenYAML(prefix string, input map[string]interface{}, out map[string]string) {
	for k, v := range input {
		key := k
		if prefix != "" {
			key = prefix + "_" + strings.ToUpper(k)
		} else {
			key = strings.ToUpper(k)
		}

		switch val := v.(type) {
		case map[string]interface{}:
			flattenYAML(key, val, out)
		case string:
			if _, exists := out[key]; !exists { // ENV اولویت داره
				out[key] = val
			}
		case int, float64, bool:
			if _, exists := out[key]; !exists {
				out[key] = fmt.Sprintf("%v", val)
			}
		}
	}
}

// دسترسی به مقدار
func (c *Config) Get(key string) string {
	return c.data[strings.ToUpper(key)]
}

// دسترسی به int
func (c *Config) GetInt(key string) int {
	var v int
	fmt.Sscanf(c.Get(key), "%d", &v)
	return v
}
