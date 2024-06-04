package utils

import (
	"fmt"
	"testing"
)

func TestGetDB(t *testing.T) {
	db := GetDB()
	fmt.Println(db.Statement)
}
