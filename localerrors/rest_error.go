//Later this will be made into a shared library
package localerrors

import (
	"net/http"

	"github.com/Sora8d/bookstore_utils-go/rest_errors"
)

func NewUnauthorizedError(message string) rest_errors.RestErr {
	return rest_errors.NewRestError("unable to retrieve user information from given access_token", http.StatusUnauthorized, "unauthorized", nil)
}
