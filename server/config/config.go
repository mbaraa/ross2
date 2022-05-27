package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

// Config represents the configuration of the server
type Config struct {
	PortNumber      string `json:"port_number"`
	DBUser          string `json:"db_user"`
	DBPassword      string `json:"db_password"`
	DBHost          string `json:"db_host"`
	AllowedClients  string `json:"allowed_clients"`
	MachineAddress  string `json:"machine_address"`
	GoogleClientID  string `json:"google_client_id"`
	UploadDirectory string `json:"upload_directory"`
	Development     bool   `json:"development"`
}

var instance *Config = nil

// GetInstance returns a singleton instance of type Config
func GetInstance() *Config {
	if instance == nil {
		instance = getInstanceUsingOSArgs()
		instance.setMachineIP()
	}
	return instance
}

func getInstanceUsingOSArgs() *Config {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "env":
			return new(Config).loadConfigFromENV()
		case "json":
			return new(Config).loadConfigFromFile()
		}
	}
	return new(Config).loadConfigFromFile() // default config is using JSON
}

func (c *Config) getMachineIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:"+c.PortNumber)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func (c *Config) setMachineIP() {
	if c.MachineAddress == "" {
		c.MachineAddress = c.getMachineIP()
	}

	c.MachineAddress = fmt.Sprintf("http://%s:%s", c.MachineAddress, c.PortNumber)
}

// used with docker
func (c *Config) loadConfigFromENV() *Config {
	return &Config{
		PortNumber:      os.Getenv("PORT_NUMBER"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),
		DBHost:          os.Getenv("DB_HOST"),
		AllowedClients:  os.Getenv("ALLOWED_CLIENTS"),
		MachineAddress:  os.Getenv("MACHINE_ADDRESS"),
		GoogleClientID:  os.Getenv("GOOGLE_CLIENT_ID"),
		UploadDirectory: os.Getenv("UPLOAD_DIRECTORY"),
		Development:     os.Getenv("DEVELOPMENT") == "true",
	}
}

// loadConfigFromFile loads the configuration from the file `./config.json`
// if the file doesn't exist the program crashes
func (c *Config) loadConfigFromFile() *Config {
	confFile, err := os.ReadFile("./config.json")
	if err != nil {
		panic("hello I need my config.json file :)")
	}

	err = json.Unmarshal(confFile, c)
	if err != nil {
		return nil
	}

	return c
}
