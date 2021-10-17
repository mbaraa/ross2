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
	PortNumber            string `json:"port_number"`
	DBUser                string `json:"db_user"`
	DBPassword            string `json:"db_password"`
	ClientAddress         string `json:"client_address"`
	GoogleClientID        string `json:"google_client_id"`
	GoogleClientSecret    string `json:"google_client_secret"`
	GoogleCallbackHandler string `json:"google_callback_handler"`
	MachineAddress        string `json:"machine_address"`
}

var instance *Config = nil

// GetInstance returns a singleton instance of type Config
func GetInstance() *Config {
	if instance == nil {
		instance = new(Config).loadConfig()
	}
	return instance
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

// loadConfig loads the configuration from the file `./config.json`
// if the file doesn't exist the program crashes
func (c *Config) loadConfig() *Config {
	confFile, err := os.ReadFile("./config.json")
	if err != nil {
		panic("hello I need my config.json file :)")
	}

	err = json.Unmarshal(confFile, c)
	if err != nil {
		return nil
	}

	if c.MachineAddress == "" {
		c.MachineAddress = c.getMachineIP()
	}

	c.MachineAddress = fmt.Sprintf("http://%s:%s", c.MachineAddress, c.PortNumber)

	return c
}
