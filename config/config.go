package config

import (
	"APIGateWay/constants"
	"APIGateWay/pkg/secure"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	RSA struct {
		PublicPath  string `yaml:"public_path"`
		PrivatePath string `yaml:"private_path"`
	} `yaml:"rsa"`
	System struct {
		MaxGoroutines string `yaml:"max_goroutines"`
	} `yaml:"system"`
	Rabbit struct {
		LogPublisher struct {
			Url       string `yaml:"url"`
			QueueName string `yaml:"queueName"`
		} `yaml:"logPublisher"`
	} `yaml:"Rabbit"`
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Version string `yaml:"version"`
	} `yaml:"server"`
	Postgres struct {
		Host     string `yaml:"host" encrypted:"true"`
		Port     string `yaml:"port" encrypted:"true"`
		User     string `yaml:"user" encrypted:"true"`
		Password string `yaml:"password" encrypted:"true"`
		DBName   string `yaml:"db_name" encrypted:"true"`
		Sslmode  string `yaml:"sslmode" encrypted:"false"`
	} `yaml:"postgres"`
}

func LoadAppConfig() (*Config, error) {
	encryptionManager := flag.Bool("encryption", false, "Encrypt config")
	flag.Parse()
	var cfg *Config
	var err error = nil
	if val, ok := os.LookupEnv("SERVICE_ENV"); ok && val == "prod" {
		cfg, err = LoadConfig("config/config.yml.prod")
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = LoadConfig("config/config.yml")
		if err != nil {
			return nil, err
		}
		fmt.Println(cfg.Rabbit.LogPublisher.QueueName)
	}
	if *encryptionManager {
		fmt.Printf("ENCRYPTING CONFIG!!!")
		if cfg.RSA.PublicPath == "" {
			privateKey, publicKey, err := secure.GenerateKeyPair(4096)
			if err != nil {
				panic(err)
			}
			publicBytes, _ := secure.PublicKeyToBytes(publicKey)
			// #nosec
			_ = os.WriteFile("config/public.key", publicBytes, 0400)
			privateBytes, _ := secure.PrivateKeyToBytes(privateKey)
			// #nosec
			_ = os.WriteFile("config/private.pem", privateBytes, 0400)
			if err != nil {
				panic(err)
			}
			cfg.RSA.PublicPath = "config/public.key"
			cfg.RSA.PrivatePath = "config/private.pem"
		}

		fmt.Println("Approve encryption?")
		char := 'n'
		_, err := fmt.Scanf("%c", &char)
		if err != nil {
			return nil, err
		}
		if char != 'y' {
			os.Exit(0)
		}
		err = EncryptConfig(cfg)

		if err != nil {
			return nil, err
		}

		cfgBytes, err := yaml.Marshal(cfg)
		if err != nil {
			return nil, err
		}

		if val, ok := os.LookupEnv("SERVICE_ENV"); ok && val == "prod" {
			err = os.WriteFile("config/config.yml.prod", cfgBytes, 0644)
			if err != nil {
				panic(err)
			}
		} else {
			err = os.WriteFile("config/config.yml", cfgBytes, 0644)
			if err != nil {
				panic(err)
			}
		}
	}

	err = DecryptConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil

}

func LoadConfig(path string) (*Config, error) {
	// #nosec
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file not found %s %w", path, constants.ErrConfig)
	}
	yamlDecoder := yaml.NewDecoder(file)

	cfg := &Config{}
	err = yamlDecoder.Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("%v %w", err.Error(), constants.ErrConfig)
	}

	return cfg, nil
}

func EncryptConfig(cfg *Config) error {

	rsaPublic, err := os.ReadFile(cfg.RSA.PublicPath)
	if err != nil {
		return fmt.Errorf("read public key: %s %w", err.Error(), constants.ErrConfig)
	}
	rsaPrivate, err := os.ReadFile(cfg.RSA.PrivatePath)
	if err != nil {
		return fmt.Errorf("read private key: %s %w", err.Error(), constants.ErrConfig)
	}
	rsaCypher, err := secure.NewRSACypher(rsaPublic, rsaPrivate)

	err = secure.EncryptStruct(cfg, rsaCypher)
	if err != nil {
		return fmt.Errorf("%v %w", err.Error(), constants.ErrConfig)
	}

	return nil
}

func DecryptConfig(cfg *Config) error {
	rsaPublic, err := os.ReadFile(cfg.RSA.PublicPath)
	if err != nil {
		return fmt.Errorf("read public key: %s %w", err.Error(), constants.ErrConfig)
	}
	rsaPrivate, err := os.ReadFile(cfg.RSA.PrivatePath)
	if err != nil {
		return fmt.Errorf("read private key: %s %w", err.Error(), constants.ErrConfig)
	}
	rsaCypher, err := secure.NewRSACypher(rsaPublic, rsaPrivate)

	err = secure.DecryptStruct(cfg, rsaCypher)
	if err != nil {
		return fmt.Errorf("%v %w", err.Error(), constants.ErrConfig)
	}

	return nil
}
