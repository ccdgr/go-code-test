package concurrent

import (
	"testing"
	"time"
)

func TestQueryAll(t *testing.T) {
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	res := QueryAll(keys, 3*time.Second)
	for k, v := range res {
		t.Logf("key: %s, res: %v", k, v)
	}
}
