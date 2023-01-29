// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------
package planes

import (
	"context"
	"errors"
	"fmt"
	http "net/http"

	"github.com/go-logr/logr"
	armrpc_controller "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	armrpc_rest "github.com/project-radius/radius/pkg/armrpc/rest"
	"github.com/project-radius/radius/pkg/middleware"
	"github.com/project-radius/radius/pkg/ucp/datamodel"
	ctrl "github.com/project-radius/radius/pkg/ucp/frontend/controller"
	"github.com/project-radius/radius/pkg/ucp/resources"
	"github.com/project-radius/radius/pkg/ucp/store"
	"go.opentelemetry.io/otel"
)

var _ armrpc_controller.Controller = (*DeletePlane)(nil)

// DeletePlane is the controller implementation to delete a UCP Plane.
type DeletePlane struct {
	ctrl.BaseController
}

// NewDeletePlane creates a new DeletePlane.
func NewDeletePlane(opts ctrl.Options) (armrpc_controller.Controller, error) {
	return &DeletePlane{ctrl.NewBaseController(opts)}, nil
}

func (p *DeletePlane) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (armrpc_rest.Response, error) {
	tr := otel.Tracer("planes")
	ctx, span := tr.Start(ctx, "deletePlane")
	defer span.End()
	req = req.WithContext(ctx)

	path := middleware.GetRelativePath(p.Options.BasePath, req.URL.Path)
	//spanAttrKey := attribute.Key(middleware.UCP_REQ_URI)
	//span.SetAttributes(spanAttrKey.String(path))

	resourceId, err := resources.ParseScope(path)
	if err != nil {
		span.RecordError(err)
		return armrpc_rest.NewBadRequestResponse(err.Error()), nil
	}
	existingPlane := datamodel.Plane{}
	etag, err := p.GetResource(ctx, resourceId.String(), &existingPlane)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, &store.ErrNotFound{}) {
			restResponse := armrpc_rest.NewNoContentResponse()
			return restResponse, nil
		}
		return nil, err
	}

	err = p.DeleteResource(ctx, resourceId.String(), etag)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	logger := logr.FromContextOrDiscard(ctx)
	restResponse := armrpc_rest.NewOKResponse(nil)
	logger.Info(fmt.Sprintf("Successfully deleted plane %s", resourceId))
	return restResponse, nil
}
