package middleware

import (
	"bytes"
	"encoding/json"
	Errors "errors"
	"fmt"
	"io"
	"net/http"
	"product_recommendation/pkg/config"
	"product_recommendation/pkg/constants"
	"product_recommendation/pkg/database"
	"product_recommendation/pkg/errors"
	"product_recommendation/pkg/logger"
	"product_recommendation/pkg/server/routers"
	"product_recommendation/pkg/types"
	"runtime"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func SetDatabase(database database.DBInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constants.FieldDatabase, database)
	}
}

func SetRedis(redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(constants.FieldRedis, redis)
	}
}

// SetLogger logs a gin HTTP request in JSON format, with some additional custom key/values
func SetLogger(c *gin.Context) {
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	body, _ := io.ReadAll(tee)
	c.Request.Body = io.NopCloser(&buf)

	responseBodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = responseBodyLogWriter
	// Start timer
	start := time.Now().UTC()

	// Process Request
	c.Next()

	// Stop timer
	duration := GetDurationInMilliseconds(start)
	status := c.Writer.Status()

	fields := types.H{
		"duration":       fmt.Sprintf("%.2fms", duration),
		"status":         status,
		"referrer":       c.Request.Referer(),
		"request_header": c.Request.Header,
	}
	//get StackTrace
	stackTrace, ok := c.Get(constants.StackTrace)
	if ok {
		fields["err_stack_trace"] = stackTrace
	}
	//Some API handles error for each record when batch writing, so the https code is still 200
	var responseBodyJSON types.H
	err := json.Unmarshal(responseBodyLogWriter.body.Bytes(), &responseBodyJSON)
	if err != nil {
		if err, ok := responseBodyJSON["error"]; ok {
			fields["error"] = err
		} else {
			//responseBodyJSON should be nil if err returned from Unmarshal is not nil
			fields["response_body"] = responseBodyLogWriter.body.String()
		}
	} else {
		fields["response_body"] = responseBodyJSON
	}

	if len(body) > 0 {
		fields["request_body"] = body
	}

	message := fmt.Sprintf("%s %s", c.Request.Method, c.Request.RequestURI)
	if status >= http.StatusBadRequest {
		logger.Error(message, fields)
	} else {
		logger.Info(message, fields)
	}
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if p, ok := r.(runtime.Error); ok {
					err := Errors.New(fmt.Sprintf("%s", p.Error()))
					logger.Error("panic error", types.H{
						"error": err,
					})
					c.JSON(http.StatusInternalServerError, gin.H{"error": err})
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}

func VerifyToken(loginType types.LoginType) gin.HandlerFunc {
	handler := func(c *gin.Context) error {

		var accessToken string
		trySetTokenFromHeader(c, &accessToken)

		if accessToken == "" {
			return errors.ErrorTokenNotFound()
		}
		c.Set("accessToken", accessToken)

		result, err := parseToken(config.GetConfig().Settings.JWTSecret, accessToken, &types.GateUserClaims{})
		if err != nil {
			return err
		}

		castedResult, ok := result.Claims.(*types.GateUserClaims)
		if !ok {
			return errors.ErrorUnknown("claims is not expected struct")
		}

		c.Set(constants.FieldLoginUser, &types.LoginUser{
			UserId:    castedResult.UserId,
			LoginType: castedResult.LoginType,
			Account:   castedResult.Account,
		})

		if castedResult.LoginType != loginType {
			return errors.ErrorTokenNotValid()
		}

		return nil
	}

	return routers.WrapperErr(handler)
}

func parseToken(secret, token string, claims jwt.Claims) (*jwt.Token, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil }
	parser := &jwt.Parser{ValidMethods: []string{"HS256"}}
	result, err := parser.ParseWithClaims(token, claims, keyFunc)

	if result != nil && result.Valid {
		return result, nil
	}
	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors == jwt.ValidationErrorMalformed {
			return result, errors.ErrorTokenMalformed()
		} else if ve.Errors == jwt.ValidationErrorExpired {
			// Token is either expired or not active yet
			return result, errors.ErrorTokenExpired()
		} else if ve.Errors == jwt.ValidationErrorSignatureInvalid {
			return result, errors.ErrorTokenSignatureInvalid()
		}
	}
	return result, errors.ErrorUnknown("unexpected error when parsing token")
}

func trySetTokenFromHeader(c *gin.Context, token *string) {
	if *token != "" {
		return
	}

	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		authorization = c.Request.Header.Get("authorization")
		if authorization == "" {
			return
		}
	}

	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return
	}

	*token = parts[1]
}
