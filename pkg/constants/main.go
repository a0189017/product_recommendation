package constants

import (
	"os"
	"product_recommendation/pkg/types"
	"syscall"
)

// current environment
var (
	CurrentVersion = "1.0.0"
)

// Field collection
const (
	FieldSygnaRequestId = "x-sygna-requestId"
	FieldAwsRequestId   = "x-amzn-requestId"
	FieldBridgeApiKey   = "x-api-key"
	FieldDatabase       = "database"
	FieldSession        = "session"
	FieldLoginUser      = "login_user"
	FieldRedis          = "redis"
)

// mysql table name
const (
	ProductRecommendationTableName = "product_recommendation"
)

var SignalsToShutdown = []os.Signal{syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM, os.Interrupt}

const StackTrace = "stack_trace"

type LoginUserContextKeyType string

const LoginUserContextKey = LoginUserContextKeyType("login_user")

// db action
const (
	RecordCreate = "create"
	RecordUpdate = "update"
	RecordDelete = "delete"
)

const (
	LoginTypeOtp   types.LoginType = "otp"
	LoginTypeToken types.LoginType = "token"
)

// redis key
const (
	ProductRecommendationKey = "product_recommendation"
)
