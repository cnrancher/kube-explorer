package cluster

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rancher/apiserver/pkg/types"
	"github.com/rancher/steve/pkg/podimpersonation"
	"github.com/rancher/steve/pkg/resources/cluster"
	"github.com/rancher/steve/pkg/server"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Register(_ context.Context, server *server.Server, displayName string) error {
	cg := server.ClientFactory
	shell := &shell{
		cg:           cg,
		namespace:    shellPodNS,
		impersonator: podimpersonation.New("shell", cg, time.Hour, getShellPodImage),
	}

	clusterSchema := server.BaseSchemas.LookupSchema("management.cattle.io.cluster")
	if clusterSchema == nil {
		return errors.New("failed to find management.cattle.io.cluster in base schema")
	}
	if clusterSchema.LinkHandlers == nil {
		clusterSchema.LinkHandlers = make(map[string]http.Handler)
	}
	clusterSchema.LinkHandlers["shell"] = shell
	clusterSchema.Store = func() types.Store {
		return &displaynameWrapper{Store: clusterSchema.Store, displayName: displayName}
	}()
	return nil
}

type displaynameWrapper struct {
	types.Store
	displayName string
}

func (s *displaynameWrapper) ByID(apiOp *types.APIRequest, schema *types.APISchema, id string) (types.APIObject, error) {
	obj, err := s.Store.ByID(apiOp, schema, id)
	if err != nil {
		return obj, err
	}
	if obj.ID != "local" {
		return obj, nil
	}
	if c, ok := obj.Object.(*cluster.Cluster); ok {
		c.Spec.DisplayName = getDisplayNameWithContext(s.displayName)
	}
	return obj, nil
}

func (s *displaynameWrapper) List(apiOp *types.APIRequest, schema *types.APISchema) (types.APIObjectList, error) {
	rtn, err := s.Store.List(apiOp, schema)
	if err != nil {
		return rtn, err
	}
	for _, obj := range rtn.Objects {
		if obj.ID != "local" {
			continue
		}
		if c, ok := obj.Object.(*cluster.Cluster); ok {
			c.Spec.DisplayName = getDisplayNameWithContext(s.displayName)
		}
	}
	return rtn, nil
}

func getDisplayNameWithContext(CurrentKubeContext string) string {
	if CurrentKubeContext != "" {
		return fmt.Sprintf("%s Cluster", cases.Title(language.English).String(CurrentKubeContext))
	}
	return "Local Cluster"
}
