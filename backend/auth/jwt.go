package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var rsaPublicKey *rsa.PublicKey

// Public Key ที่ copy มาจาก Keycloak (แก้ได้ถ้ามีใหม่)
var publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArqq5a2FFQx1OqPFkDtVawE+4fgeysR9v/QhKgKSRv4FRbynG0ACfBEDiKhHreBKYKodf21BpfC+kcC9MIlsv1NXWprN2xxIqCm+ywFR2JNn4h7yvaUrELXy79BRJN3xQi0MQSvxm1zpPVDnETrjgp7C/gRbKLD3Nb96f6xaMAZdOPAHyhG2t+qqv2fxf01k5G4nLxXKtSwVQ/hPvfOxKVZvH172xt4qklYXXljYIBQaBhRLR3U/z5sFpLICMMdBzP+gk7oh3+rayX6XHXh7yuZnhL2krW3ctCfwNwHVG8cUyO9Pc3abDwXhKdTaxAd5DFhHzgHsUH+rEZLAT2UIFkwIDAQAB
-----END PUBLIC KEY-----`

func init() {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	var ok bool
	rsaPublicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		panic("key is not of type *rsa.PublicKey")
	}
}

// Middleware ตรวจสอบ JWT Token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return rsaPublicKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			c.JSON(401, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
