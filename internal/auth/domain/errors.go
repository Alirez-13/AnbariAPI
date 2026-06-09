// # SINGLE REASON: Define auth domain errors.
package domain

import "errors"

var ErrNotFound = errors.New("auth record not found")
