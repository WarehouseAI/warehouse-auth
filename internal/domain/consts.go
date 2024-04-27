package domain

type CtxKey string

const (
	HeaderContentType   = "Content-Type"
	HeaderXForwardedFor = "X-Forwarded-For"
	JsonContentType     = "application/json"
	ProtoContentType    = "application/x-protobuf"
	AuthHeader          = "Authorization"
	VersionHeader       = "Coffee-Version"
	VersionDelimiter    = ":"

	S3Endpoint = "storage.yandexcloud.net"

	AccountCtxKey     = CtxKey("account")
	TokenNumberCtxKey = CtxKey("token_number")
)
