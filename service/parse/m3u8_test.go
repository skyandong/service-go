package parse

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestParse(t *testing.T) {
	_, f, _, _ := runtime.Caller(0)
	filePath := path.Join(path.Dir(f), "../../testdata/test.m3u8")
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	require.NoError(t, err)
	defer file.Close()
	_, err = parse(file)
	require.NoError(t, err)

}
