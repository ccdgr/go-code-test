package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestScanIPs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	ips := []string{"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4", "5.5.5.5", "6.6.6.6"}

	results := ScanIPs(ctx, ips, 3)

	for _, r := range results {
		fmt.Printf("IP: %s, 开放: %v\n", r.IP, r.Found)
	}
}
