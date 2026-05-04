package ast

import "github.com/shotowon/kica/internal/tokens"

type NTMapping struct {
	Name tokens.Token
	Ty   Type
}
