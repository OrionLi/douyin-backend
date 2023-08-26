package e

import "google.golang.org/grpc/codes"

const (
	Success                = codes.OK
	Error                  = codes.Unknown
	InvalidParams          = codes.InvalidArgument
	ErrorExistUserNotFound = codes.NotFound
	ErrorExistUser         = codes.AlreadyExists
	ErrorAuthToken         = codes.PermissionDenied
	ErrorAborted           = codes.Aborted
	ErrorNotCompare        = codes.Unauthenticated
)
