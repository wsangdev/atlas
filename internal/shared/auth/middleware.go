package auth

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	HeaderAdminKey  = "X-Admin-Key"
	HeaderDeviceID  = "X-Device-Id"
	HeaderDeviceKey = "X-Device-Key"
	contextDeviceID = "auth_device_id"
)

// DeviceKeyVerifier valida la API key de un dispositivo. Tracking/fleet no
// conocen esta interfaz: el adaptador se inyecta desde app.go.
type DeviceKeyVerifier interface {
	VerifyDeviceKey(deviceID, apiKey string) (bool, error)
}

// AdminKey protege las rutas de usuarios/panel con la key de entorno.
func AdminKey(adminKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		provided := c.GetHeader(HeaderAdminKey)
		if adminKey == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(adminKey)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "api key de admin invalida"})
			return
		}
		c.Next()
	}
}

// DeviceKey autentica trackers y deja el device_id en el contexto.
func DeviceKey(verifier DeviceKeyVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader(HeaderDeviceID)
		apiKey := c.GetHeader(HeaderDeviceKey)
		if deviceID == "" || apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "headers X-Device-Id y X-Device-Key son requeridos",
			})
			return
		}

		ok, err := verifier.VerifyDeviceKey(deviceID, apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "no se pudo validar el dispositivo"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "dispositivo no autorizado"})
			return
		}

		c.Set(contextDeviceID, deviceID)
		c.Next()
	}
}

// DeviceIDFromContext devuelve el device autenticado por el middleware.
func DeviceIDFromContext(c *gin.Context) (string, bool) {
	value, exists := c.Get(contextDeviceID)
	if !exists {
		return "", false
	}
	id, ok := value.(string)
	return id, ok
}
