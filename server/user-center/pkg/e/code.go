package e

import "google.golang.org/grpc/codes"

const (
	Error                      = codes.Unknown
	InvalidParams              = codes.InvalidArgument
	ErrorExistUser             = codes.AlreadyExists
	ErrorExistUserNotFound     = codes.NotFound
	ErrorNotCompare            = codes.Unauthenticated
	ErrorAuthToken             = codes.PermissionDenied
	ErrorAuthCheckTokenTimeout = codes.DeadlineExceeded
)
