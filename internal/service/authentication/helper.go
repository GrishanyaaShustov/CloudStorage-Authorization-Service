package authentication

import "google.golang.org/grpc/status"
import "google.golang.org/grpc/codes"

func statusUnimplemented(method string) error {
	return status.Errorf(codes.Unimplemented, "%s is not implemented yet", method)
}
