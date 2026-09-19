package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeVerifier struct {
	ok bool
}

func (f fakeVerifier) VerifyDeviceKey(deviceID, apiKey string) (bool, error) {
	return f.ok, nil
}

func routerCon(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/privado", middleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return router
}

func doRequest(router *gin.Engine, headers map[string]string) int {
	req := httptest.NewRequest(http.MethodGet, "/privado", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Code
}

func TestAdminKey(t *testing.T) {
	router := routerCon(AdminKey("secreto-admin"))

	if code := doRequest(router, nil); code != http.StatusUnauthorized {
		t.Fatalf("sin key: code = %d, queria 401", code)
	}
	if code := doRequest(router, map[string]string{HeaderAdminKey: "otra"}); code != http.StatusUnauthorized {
		t.Fatalf("key mala: code = %d, queria 401", code)
	}
	if code := doRequest(router, map[string]string{HeaderAdminKey: "secreto-admin"}); code != http.StatusOK {
		t.Fatalf("key buena: code = %d, queria 200", code)
	}
}

func TestDeviceKey(t *testing.T) {
	router := routerCon(DeviceKey(fakeVerifier{ok: true}))

	if code := doRequest(router, nil); code != http.StatusUnauthorized {
		t.Fatalf("sin headers: code = %d, queria 401", code)
	}
	if code := doRequest(router, map[string]string{HeaderDeviceID: "d1"}); code != http.StatusUnauthorized {
		t.Fatalf("falta key: code = %d, queria 401", code)
	}
	headers := map[string]string{HeaderDeviceID: "d1", HeaderDeviceKey: "k1"}
	if code := doRequest(router, headers); code != http.StatusOK {
		t.Fatalf("headers completos: code = %d, queria 200", code)
	}

	routerMalo := routerCon(DeviceKey(fakeVerifier{ok: false}))
	if code := doRequest(routerMalo, headers); code != http.StatusUnauthorized {
		t.Fatalf("key invalida: code = %d, queria 401", code)
	}
}

func TestHashYVerify(t *testing.T) {
	key, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("generar: %v", err)
	}
	if len(key) != 64 {
		t.Fatalf("len(key) = %d, queria 64", len(key))
	}

	hash := HashAPIKey(key)
	if hash == key {
		t.Fatal("el hash no debe ser igual a la key")
	}
	if !VerifyAPIKey(hash, key) {
		t.Fatal("la key correcta debe verificar")
	}
	if VerifyAPIKey(hash, "otra-key") {
		t.Fatal("una key distinta no debe verificar")
	}
}
