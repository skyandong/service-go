package downloadm3u8

import (
	"fmt"
	"testing"
)

func TestGenSlice(t *testing.T) {
	slc := genSlice(10)
	fmt.Println(slc)
}
