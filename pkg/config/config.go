package config

import (
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

type SDConfig struct {
	Region 			string `yaml:"region"`
	AccessKey 		string `yaml:"access_key,omitempty"`
	SecretKey 		string `yaml:"secret_key,omitempty"`
	RefreshInterval model.Duration `yaml:"refresh_interval,omitempty"`
	Port            int `yaml:"port"`
	Filters         []*Filter `yaml:"filters,omitempty"`
}

// DefaultSDConfig is the default EC2 SD configuration.
var DefaultSDConfig = SDConfig{
	Port:            80,
	RefreshInterval: model.Duration(60 * time.Second),
}
// Filter is the configuration for filtering EC2 instances.
type Filter struct {
	Name  string   `yaml:"name"`
	Value string `yaml:"value"`
}

// Config is the top-level configuration for Prometheus's config files.
type Config struct {
	EcsSDConfig SDConfig `yaml:"ecs_sd_config"`
	// original is the input from which the config was parsed.
	original string
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}
	// If the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	cfg.original = s
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, errors.Wrapf(err, "parsing YAML file %s", filename)
	}
	return cfg, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *SDConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultSDConfig
	type plain SDConfig
	err := unmarshal((*plain)(c))
	if err != nil {
		return err
	}
	if c.Region == "" {
		c.Region = "cn-beijing"
	}

	if c.AccessKey == "" {
		c.AccessKey = os.Getenv("ALIYUN_ACCESS_KEY_ID")
	}

	if c.SecretKey == "" {
		c.SecretKey = os.Getenv("ALIYUN_SECRET_ACCESS_KEY")
	}

	if c.AccessKey == "" || c.SecretKey == "" {
		return errors.New("aliyun access_key/secret_key not config")
	}

	for _, f := range c.Filters {
		if len(f.Value) == 0 {
			return errors.New("ECS SD configuration filter values cannot be empty")
		}
	}
	return nil
}