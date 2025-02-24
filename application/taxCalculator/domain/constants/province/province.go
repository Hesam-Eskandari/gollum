package province

import (
	"errors"
	"strings"
)

var ErrInvalidProvince = errors.New("invalid province")

type Province string

const (
	BC      Province = "british-columbia"
	Federal Province = "federal"
)

func FromString(pr string) (Province, error) {
	pr = strings.ToLower(pr)
	switch pr {
	case "british-columbia", "bc", "britishColumbia":
		return BC, nil
	case "federal", "fed", "canada":
		return Federal, nil
	default:
		return Federal, ErrInvalidProvince
	}
}
