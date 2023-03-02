// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package planes

import (
	"context"
	"fmt"
	http "net/http"
	"strings"

	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	armrpc_controller "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	armrpc_rest "github.com/project-radius/radius/pkg/armrpc/rest"
	"github.com/project-radius/radius/pkg/middleware"
	"github.com/project-radius/radius/pkg/ucp/datamodel"
	"github.com/project-radius/radius/pkg/ucp/datamodel/converter"
	ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller"
	"github.com/project-radius/radius/pkg/ucp/resources"
	"github.com/project-radius/radius/pkg/ucp/store"
	"github.com/project-radius/radius/pkg/ucp/ucplog"
)

var _ armrpc_controller.Controller = (*ListPlanesByType)(nil)

// ListPlanesByType is the controller implementation to get the list of UCP planes.
type ListPlanesByType struct {
	ctrl.BaseController
}

// NewListPlanes creates a new ListPlanesByType.
func NewListPlanesByType(opts ctrl.Options) (armrpc_controller.Controller, error) {
	return &ListPlanesByType{ctrl.NewBaseController(opts)}, nil
}

func (e *ListPlanesByType) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (armrpc_rest.Response, error) {
	path := middleware.GetRelativePath(e.Options.BasePath, req.URL.Path)
	// The path is /planes/{planeType}
	planeType := strings.Split(path, resources.SegmentSeparator)[2]
	query := store.Query{
		RootScope:    resources.SegmentSeparator + resources.PlanesSegment,
		IsScopeQuery: true,
		ResourceType: planeType,
	}
	logger := ucplog.FromContextOrDiscard(ctx)
	logger.Info(fmt.Sprintf("Listing planes in scope %s/%s", query.RootScope, planeType))
	result, err := e.StorageClient().Query(ctx, query)
	if err != nil {
		return nil, err
	}
	listOfPlanes, err := e.createResponse(ctx, req, result)
	if err != nil {
		return nil, err
	}
	var ok = armrpc_rest.NewOKResponse(&v1.PaginatedList{
		Value: listOfPlanes,
	})
	return ok, nil
}

func (p *ListPlanesByType) createResponse(ctx context.Context, req *http.Request, result *store.ObjectQueryResult) ([]any, error) {
	apiVersion := ctrl.GetAPIVersion(req)
	listOfPlanes := []any{}
	if len(result.Items) > 0 {
		for _, item := range result.Items {
			var plane datamodel.Plane
			err := item.As(&plane)
			if err != nil {
				return nil, err
			}

			versioned, err := converter.PlaneDataModelToVersioned(&plane, apiVersion)
			if err != nil {
				return nil, err
			}

			listOfPlanes = append(listOfPlanes, versioned)
		}
	}
	return listOfPlanes, nil
}
