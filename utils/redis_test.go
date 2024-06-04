package utils

import (
	"fmt"
	"testing"
)

func TestGetRediscli(t *testing.T) {
	res := GetRediscli()
	fmt.Println(res)
}
