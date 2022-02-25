package kysdkid

import (
	"fmt"
	"testing"
)

func TestIdGenerator_NextId(t *testing.T) {
	id := NewIdGenerator()
	for {
		fmt.Println(id.NextId("rmppp"))
	}
}
