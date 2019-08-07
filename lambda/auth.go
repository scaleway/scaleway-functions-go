package lambda

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

// ApplicationClaim represents the claims related to an application
// composed of either NamespaceID or ApplicationID of the linked JWT
type ApplicationClaim struct {
	NamespaceID   string `json:"namespace_id"`
	ApplicationID string `json:"application_id"`
}

// Authenticate incoming request based on multiple factors:
// - 1: Whether the function's privacy has been set to private, if public, just leave this middleware
// - 2: Get the public key injected in this function runtime (done automatically by Scaleway)
// - 3: Check whether a Token has been sent via a specific Headers reserved by Scaleway
// - 4: Parse the incoming JWT with the public key
// - 5: Check the "Application Claims" linked to the JWT
// - 6: Both FunctionID and NamespaceID are injected via environment variables by Scaleway
// ---  so we have to check the authenticity of the incoming token by comparing the claims
func authenticate(req *http.Request) (err error) {
	isPublicFunction := os.Getenv("SCW_PUBLIC")
	if isPublicFunction == "true" {
		return
	}

	publicKey := os.Getenv("SCW_PUBLIC_KEY")

	// Check that encoded key may be parsed back to a valid RSA Private Key
	block, _ := pem.Decode([]byte(publicKey))
	parsedKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return
	}
	if parsedKey == nil {
		err = errors.New("Invalid public key")
		return
	}

	requestToken := req.Header.Get("SCW_FUNCTIONS_TOKEN")
	if requestToken == "" {
		err = errors.New("Authentication token not present in the request's header")
		return
	}

	claims := jwt.MapClaims{}

	_, err = jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (i interface{}, e error) {
		return &parsedKey, nil
	})
	if err != nil {
		return
	}

	marshalledClaims, err := json.Marshal(claims["application_claim"])
	if err != nil {
		return
	}

	parsedClaims := []ApplicationClaim{}
	if err = json.Unmarshal(marshalledClaims, &parsedClaims); err != nil {
		return
	}

	if len(parsedClaims) == 0 {
		err = errors.New("Invalid Claims")
		return
	}
	applicationClaims := parsedClaims[0]

	applicationID := os.Getenv("SCW_APPLICATION_ID")
	namespaceID := os.Getenv("SCW_NAMESPACE_ID")

	if applicationClaims.NamespaceID != namespaceID && applicationClaims.ApplicationID != applicationID {
		err = errors.New("Invalid claims")
	}
	return
}
