package codegen

import (
	"github.com/tinhtran24/gqlgen/codegen/config"
)

func (b *builder) buildTypes() (map[string]*config.TypeReference, error) {
	ret := map[string]*config.TypeReference{}

	for _, ref := range b.Binder.References {
		for ref != nil {
			ret[ref.UniquenessKey()] = ref

			ref = ref.Elem()
		}
	}
	return ret, nil
}
