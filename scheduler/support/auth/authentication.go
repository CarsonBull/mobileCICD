package auth

import (
	"github.com/CarsonBull/mobileCICD/scheduler/support"
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/lestrrat/go-jwx/jwk"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/jws"
	"gopkg.in/square/go-jose.v2/jwt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Authenticator struct {
	Decoder func(ctx context.Context, r *http.Request) (interface{}, error)
}

type Role int

const (
	User Role = iota
	Admin
)

type AuthenticatedRequest struct {
	AccountId string
	Role Role
	Object interface{}
}

func (auth Authenticator) AuthenticateUser(ctx context.Context, r *http.Request) (interface{}, error) {


	var accountId string
	var role Role
	authHeader := "Authorization"

	bearerToken := r.Header.Get(authHeader)

	reg, _ := regexp.Compile("^Bearer (.*)$")

	tokens := reg.FindStringSubmatch(bearerToken)

	if len(tokens) != 2 {
		return nil, support.ErrUnauthorized
	}

	token := tokens[1]

	if token == "" {
		fmt.Println("No JWT in header")
		return nil, support.ErrUnauthorized
	}

	if os.Getenv("ENV") == "dev" {
		accountId = getAccountIdFromJwt(token)

		if accountId == "" {
			return nil, support.ErrUnauthorized
		}

	} else {

		token := r.Header.Get(authHeader)

		if token == "" {
			return nil, support.ErrUnauthorized
		}

		var ok bool

		accountId, ok = checkIdentity(token)

		if !ok {
			return nil, support.ErrUnauthorized
		}
	}

	attrs, _ := getTokenAttributes(token)
	role = getRolesFromAttributes(attrs)

	request, err := auth.Decoder(ctx, r)

	return AuthenticatedRequest{AccountId:accountId, Role:role, Object:request}, err
}

func checkIdentity(token string) (string, bool) {

	parsedJWT, err := jwt.ParseSigned(token)
	if err != nil {
		log.Printf("failed to parse JWT:%v", err)
	}

	keyId := parsedJWT.Headers[0].KeyID

	// TODO: cache this shit
	// https://github.com/patrickmn/go-cache
	set, err := jwk.Fetch("https://cognito-idp.us-west-2.amazonaws.com/us-west-2_TKCY92xxy/.well-known/jwks.json")
	if err != nil {
		log.Printf("failed to parse JWK: %s", err)
		return "", false
	}

	// If you KNOW you have exactly one key, you can just
	// use set.Keys[0]
	keys := set.LookupKeyID(keyId)
	if len(keys) == 0 {
		log.Printf("failed to lookup key: %s", err)
		return "", false
	}

	key, err := keys[0].Materialize()
	if err != nil {
		log.Printf("failed to generate public key: %s", err)
		return "", false
	}

	allClaims := make(map[string]interface{})
	err = parsedJWT.Claims(key, &allClaims)

	if err != nil {
		log.Printf("Failed :%+v", err)
	}

	ok, err := verifyToken(token, key.(*rsa.PublicKey), allClaims["exp"].(float64))

	if err != nil {
		return "", false
	}

	// TODO: do we want to check for dev here? This setup allows them to have an expired jwt. Don't really like that
	if os.Getenv("ENVIRONMENT") != "dev" && (ok == false  || err != nil) {
		return "", false
	}

	log.Println("Successful authentication")

	return allClaims["sub"].(string), true
}


func verifyToken(token string, publicKey *rsa.PublicKey, exp float64) (bool, error) {

	err := jws.Verify(token, publicKey)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	now := time.Now()

	if now.After(time.Unix(int64(exp), 0)) {
		log.Println("Token expired")
		return false, nil
	}

	return true, nil
}

func getAccountIdFromJwt(token string) string {

	allClaims, err := getTokenAttributes(token)

	if err != nil {
		log.Printf("Error parsing jwt: %s", err)
		return ""
	}

	return allClaims["sub"].(string)
}

func getTokenAttributes(token string) (map[string]interface{}, error) {

	parsedJWT, err := jwt.ParseSigned(token)
	if err != nil {
		log.Printf("failed to parse JWT:%v", err)
	}

	keyId := parsedJWT.Headers[0].KeyID

	// TODO: cache this shit
	// https://github.com/patrickmn/go-cache
	set, err := jwk.Fetch("https://cognito-idp.us-west-2.amazonaws.com/us-west-2_TKCY92xxy/.well-known/jwks.json")
	if err != nil {
		log.Printf("failed to parse JWK: %s. Make sure the userpool id is correct", err)
		return nil, err
	}

	// If you KNOW you have exactly one key, you can just
	// use set.Keys[0]
	keys := set.LookupKeyID(keyId)
	if len(keys) == 0 {
		log.Printf("failed to lookup key")
		return nil, errors.New("Key failure")
	}

	key, err := keys[0].Materialize()
	if err != nil {
		log.Printf("failed to generate public key: %s", err)
		return nil, errors.New("test")
	}

	allClaims := make(map[string]interface{})
	err = parsedJWT.Claims(key, &allClaims)

	return allClaims, nil
}

func getRolesFromAttributes(attributes map[string]interface{}) Role {

	if attributes["cognito:groups"] == nil {
		return User
	}

	for _, val:= range attributes["cognito:groups"].([]interface{}) {
		if val == "Admins" {
			log.Printf("Admin user login")
			return Admin
		}
	}

	return User
}