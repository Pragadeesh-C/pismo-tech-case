package service

import "errors"

var ErrAccountAlreadyExists = errors.New("account already exists")
var ErrDBDown = errors.New("db down")
var ErrDocNumEmpty = errors.New("document number is empty")
var ErrInvalidAccountID = errors.New("invalid account id")
var ErrAccountNotFound = errors.New("account not found")
