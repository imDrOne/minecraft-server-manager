package nodes

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s1 := "123"
	s2 := "456"
	h := md5.New()
	h.Write([]byte(s1 + s2))
	fmt.Printf("%x", h.Sum(nil))
}
