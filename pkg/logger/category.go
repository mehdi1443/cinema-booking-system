package logger

type Category string
type SubCategory string
type ExtraKey string

// Categories
const (
	General         Category = "General"
	IO              Category = "IO"
	Internal        Category = "Internal"
	Postgres        Category = "Postgres"
	Redis           Category = "Redis"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
	Prometheus      Category = "Prometheus"
	Auth            Category = "Auth"
)

// SubCategories
const (
	// General
	Startup         SubCategory = "Startup"
	ExternalService SubCategory = "ExternalService"

	// Postgres
	Migration SubCategory = "Migration"
	Select    SubCategory = "Select"
	Rollback  SubCategory = "Rollback"
	Update    SubCategory = "Update"
	Delete    SubCategory = "Delete"
	Insert    SubCategory = "Insert"

	// Internal
	Api                 SubCategory = "Api"
	HashPassword        SubCategory = "HashPassword"
	DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"
	FailedToCreateUser  SubCategory = "FailedToCreateUser"

	// Validation
	MobileValidation   SubCategory = "MobileValidation"
	PasswordValidation SubCategory = "PasswordValidation"

	// IO
	RemoveFile SubCategory = "RemoveFile"

	// Auth
	JWT  SubCategory = "JWT"
	RBAC SubCategory = "RBAC"

	// Limiter
	Limiter SubCategory = "Limiter"

	// Captcha
	Captcha SubCategory = "Captcha"
)

// Extra Keys
const (
	// System & Metadata
	AppName    ExtraKey = "AppName"
	LoggerName ExtraKey = "Logger"
	ClientIp   ExtraKey = "ClientIp"
	HostIp     ExtraKey = "HostIp"

	// Request & Response
	Method       ExtraKey = "Method"
	Path         ExtraKey = "Path"
	StatusCode   ExtraKey = "StatusCode"
	BodySize     ExtraKey = "BodySize"
	Latency      ExtraKey = "Latency"
	RequestBody  ExtraKey = "RequestBody"
	ResponseBody ExtraKey = "ResponseBody"

	// User & Business Context
	UserID ExtraKey = "UserId"

	// Error Details
	ErrorMessage ExtraKey = "ErrorMessage"
)
