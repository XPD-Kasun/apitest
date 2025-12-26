package user

import "errors"

var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrNoJWTKeyFound = errors.New("no jwt key is found in env")
