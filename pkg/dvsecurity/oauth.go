// package dvsecurity provides server security, including sessions, login, jwt token
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvsecurity

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
)

const (
	security = "security"
)

var SecretKey = "6876jgj6876876hkh8989899"

func securityEndPointHandler(ctx *dvcontext.RequestContext) bool {
	token := GetToken()
	refreshToken := GetRefreshToken()
	result := map[string]interface{}{
		"access_token":       token,
		"expires_in":         18000,
		"refresh_expires_in": 1800,
		"refresh_token":      refreshToken,
		"token_type":         "bearer",
		"not-before-policy":  0,
		"session_state":      "467a8f63-3c4c-40a6-8e6c-d5388b59c7a0",
		"scope":              "profile email",
	}
	ctx.ExtraAsDvObject.Set("OAUTH", result)
	return true
}

func GetToken() string {
	claims := map[string]string{
		"jti": "8b5e5323-7503-409d-8831-db0dbdb66b69",
		// exp: 1590427779,
		// nbf: 0,
		// iat: 1590409779,
		"iss": "microcore",
		"aud": "account",
		"sub": "6402d452-3004-4d87-9a0a-9b5783fbac97",
		"typ": "Bearer",
		"azp": "frontend",
		// auth_time: 0,
		"session_state": "d11bf3d4-381b-4fe5-99f1-05f0338ab800",
		"acr":           "1",
		"scope":         "email profile",
		// email_verified: false,
		"name":               "Tenant Admin",
		"preferred_username": "admin@gmail.com",
		"given_name":         "Tenant",
		"family_name":        "Admin",
		"tenant-id":          "4556c255-be35-47fb-aba1-8d888999db40",
		"email":              "admin@gmail.com",
	}
	return GenerateHs256Jwt(claims, SecretKey, 36000)
}

func GetRefreshToken() string {
	claims := map[string]string{
		"jti": "8b5e5323-7503-409d-8831-db0dbdb66b69",
		// exp: 1590427779,
		// nbf: 0,
		// iat: 1590409779,
		"iss": "microcore",
		"aud": "account",
		"sub": "6402d452-3004-4d87-9a0a-9b5783fbac97",
		"typ": "Bearer",
		"azp": "frontend",
		// auth_time: 0,
		"session_state": "d11bf3d4-381b-4fe5-99f1-05f0338ab800",
		"acr":           "1",
		"scope":         "email profile",
		// email_verified: false,
		"name":               "Tenant Admin",
		"preferred_username": "admin@gmail.com",
		"given_name":         "Tenant",
		"family_name":        "Admin",
		"tenant-id":          "4556c255-be35-47fb-aba1-8d888999db40",
		"email":              "admin@gmail.com",
	}
	return GenerateHs256Jwt(claims, SecretKey, 36000)
}

func registerSecurityAction() bool {
	return dvmodules.RegisterActionProcessor(security, securityEndPointHandler, false)
}

var registered = registerSecurityAction()
