package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig()
	fmt.Println(config)
}
