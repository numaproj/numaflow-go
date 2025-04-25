package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type MaybeCat struct {
	// if allowData == true, allow data to pass through; otherwise drop it
	allowData bool
}

func (e *MaybeCat) Map(ctx context.Context, keys []string, d mapper.Datum) mapper.Messages {
	if e.allowData {
		return mapper.MessagesBuilder().Append(mapper.NewMessage(d.Value()).WithKeys(keys))
	} else {
		return mapper.MessagesBuilder().Append(mapper.MessageToDrop())
	}
}

func main() {
	config := readConfig()
	fmt.Printf("successfully read in ConfigMap; allowData=%t\n", config.AllowData)
	err := mapper.NewServer(&MaybeCat{allowData: config.AllowData}).Start(context.Background())
	if err != nil {
		log.Panic("Failed to start map function server: ", err)
	}
}

func readConfig() Config {
	filePath := "/etc/config/config.json" // Path to the JSON config file

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
		return Config{}
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Printf("Error unmarshalling config JSON: %v\n", err)
		return Config{}
	}

	return config
}

type Config struct {
	AllowData bool `json:"allowData"`

	AnotherField string `json:"anotherField"`
}
