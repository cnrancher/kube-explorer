package server

import (
	"github.com/rancher/apiserver/pkg/types"
)

type deleteOptionStore struct {
	types.Store
}

func (s *deleteOptionStore) Delete(apiOp *types.APIRequest, schema *types.APISchema, id string) (types.APIObject, error) {
	query := apiOp.Request.URL.Query()
	query.Add("propagationPolicy", "Background")
	apiOp.Request.URL.RawQuery = query.Encode()
	return s.Store.Delete(apiOp, schema, id)
}
