package download

import (
	"fmt"
	"testing"
)

func TestGenSlice(t *testing.T) {
	slc := genSlice(10)
	fmt.Println(slc)
}
