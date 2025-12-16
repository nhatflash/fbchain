package graph

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/nhatflash/fbchain/scalar"
)

func MarshalDate(t scalar.CustomDate) graphql.Marshaler {
	return t
}

func UnmarshalDate(v any) (scalar.CustomDate, error) {
	var t scalar.CustomDate
	err := t.UnmarshalGQL(v)
	return t, err
}