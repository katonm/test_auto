package service

import (
	"context"

	"github.com/go-kit/kit/log"
)

const (
	valid    = "Balanced"
	notValid = "Not Balanced"
)

type service struct {
	logger log.Logger
}

// Service interface describes a service that adds numbers...
type Service interface {
	Fix(ctx context.Context, strIn string) (strOut string)
	Validate(ctx context.Context, strIn string) (res string)
}

// NewService returns a Service with all of the expected dependencies...
func NewService(logger log.Logger) Service {
	return &service{
		logger: logger,
	}
}

// Fix func implements Service interface...
func (s *service) Fix(ctx context.Context, strIn string) (strOut string) {
	var st []rune
	a := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}

	for _, j := range strIn {
		l := len(st) - 1

		if l < 0 && (j == ')' || j == '}' || j == ']') {
			continue
		}

		if j == '{' || j == '(' || j == '[' {
			st = append(st, j)
			strOut += string(j)

			continue
		}

		if l < 0 {
			continue
		}

		if (st[l] == '(' && j == ')') ||
			(st[l] == '{' && j == '}') ||
			(st[l] == '[' && j == ']') {
			st = st[:l]
			strOut += string(j)

			continue
		}

		strOut += string(a[st[l]])
		st = st[:l]
	}

	for i := len(st); i > 0; i-- {
		l := len(st) - 1
		strOut += string(a[st[l]])
		st = st[:l]
	}

	return strOut
}

// Validate func implements Service interface...
func (s *service) Validate(ctx context.Context, strIn string) (res string) {
	var st []rune

	for _, j := range strIn {
		l := len(st) - 1

		if j == '{' || j == '(' || j == '[' {
			st = append(st, j)

			continue
		}

		if l < 0 {
			return notValid
		}

		if (st[l] == '(' && j != ')') ||
			(st[l] == '{' && j != '}') ||
			(st[l] == '[' && j != ']') {
			return notValid
		}
		st = st[:l]
	}

	if len(st) == 0 {
		return valid
	}

	return notValid
}
