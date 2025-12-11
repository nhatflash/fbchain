package service

import (
	"time"
	"fmt"
)

func GenerateTenantCode() string {
	now := time.Now()
	unixMilli := now.UnixMilli()
	return fmt.Sprintf("TENANT-%d", unixMilli)
}