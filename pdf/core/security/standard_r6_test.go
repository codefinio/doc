package security

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func BenchmarkAlg2b(b *testing.B) {

	r := rand.New(rand.NewSource(1234567))
	const n = 20
	pass := make([]byte, n)
	r.Read(pass)
	data := make([]byte, n+8+48)
	r.Read(data)
	user := make([]byte, 48)
	r.Read(user)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = alg2b(data, pass, user)
	}
}

func TestStdHandlerR6(t *testing.T) {
	var cases = []struct {
		Name      string
		EncMeta   bool
		UserPass  string
		OwnerPass string
	}{
		{
			Name: "simple", EncMeta: true,
			UserPass: "user", OwnerPass: "owner",
		},
		{
			Name: "utf8", EncMeta: false,
			UserPass: "æøå-u", OwnerPass: "æøå-o",
		},
		{
			Name: "long", EncMeta: true,
			UserPass:  strings.Repeat("user", 80),
			OwnerPass: strings.Repeat("owner", 80),
		},
	}

	const (
		perms = 0x12345678
	)

	for _, R := range []int{5, 6} {
		R := R
		t.Run(fmt.Sprintf("R=%d", R), func(t *testing.T) {
			for _, c := range cases {
				c := c
				t.Run(c.Name, func(t *testing.T) {
					sh := stdHandlerR6{}
					d := &StdEncryptDict{
						R: R, P: perms,
						EncryptMetadata: c.EncMeta,
					}

					ekey, err := sh.GenerateParams(d, []byte(c.OwnerPass), []byte(c.UserPass))
					if err != nil {
						t.Fatal("Failed to encrypt:", err)
					}

					key, uperm, err := sh.alg2a(d, []byte(c.UserPass))
					if err != nil || uperm != perms {
						t.Error("Failed to authenticate user pass:", err)
					} else if !bytes.Equal(ekey, key) {
						t.Error("wrong encryption key")
					}

					key, uperm, err = sh.alg2a(d, []byte(c.OwnerPass))
					if err != nil || uperm != PermOwner {
						t.Error("Failed to authenticate owner pass:", err, uperm)
					} else if !bytes.Equal(ekey, key) {
						t.Error("wrong encryption key")
					}

					d.P = PermOwner

					key, uperm, err = sh.alg2a(d, []byte(c.UserPass))
					if R == 5 {

						if err != nil || uperm != PermOwner {
							t.Error("Failed to authenticate user pass:", err)
						}
					} else {

						if err == nil || uperm == PermOwner {
							t.Error("was able to elevate permissions with R=6")
						}
					}
				})
			}
		})
	}
}
