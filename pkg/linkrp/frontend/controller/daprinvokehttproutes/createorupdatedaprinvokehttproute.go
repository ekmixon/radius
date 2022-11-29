// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package daprinvokehttproutes

import (
	"context"
	"errors"
	"net/http"

	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	ctrl "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	"github.com/project-radius/radius/pkg/armrpc/rest"
	"github.com/project-radius/radius/pkg/linkrp/datamodel"
	"github.com/project-radius/radius/pkg/linkrp/datamodel/converter"
	"github.com/project-radius/radius/pkg/linkrp/frontend/deployment"
	"github.com/project-radius/radius/pkg/rp/outputresource"
	"github.com/project-radius/radius/pkg/ucp/store"
)

var _ ctrl.Controller = (*CreateOrUpdateDaprInvokeHttpRoute)(nil)

// CreateOrUpdateDaprInvokeHttpRoute is the controller implementation to create or update DaprInvokeHttpRoute link resource.
type CreateOrUpdateDaprInvokeHttpRoute struct {
	ctrl.BaseController
}

// NewCreateOrUpdateDaprInvokeHttpRoute creates a new instance of CreateOrUpdateDaprInvokeHttpRoute.
func NewCreateOrUpdateDaprInvokeHttpRoute(opts ctrl.Options) (ctrl.Controller, error) {
	return &CreateOrUpdateDaprInvokeHttpRoute{ctrl.NewBaseController(opts)}, nil
}

// Run executes CreateOrUpdateDaprInvokeHttpRoute operation.
func (daprHttpRoute *CreateOrUpdateDaprInvokeHttpRoute) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (rest.Response, error) {
	serviceCtx := v1.ARMRequestContextFromContext(ctx)
	newResource, err := daprHttpRoute.Validate(ctx, req, serviceCtx.APIVersion)
	if err != nil {
		return nil, err
	}

	old := &datamodel.DaprInvokeHttpRoute{}
	isNewResource := false
	etag, err := daprHttpRoute.GetResource(ctx, serviceCtx.ResourceID.String(), old)
	if err != nil {
		if errors.Is(&store.ErrNotFound{}, err) {
			isNewResource = true
		} else {
			return nil, err
		}
	}

	if req.Method == http.MethodPatch && isNewResource {
		return rest.NewNotFoundResponse(serviceCtx.ResourceID), nil
	}

	err = ctrl.ValidateETag(*serviceCtx, etag)
	if err != nil {
		return rest.NewPreconditionFailedResponse(serviceCtx.ResourceID.String(), err.Error()), nil
	}

	newResource.SystemData = ctrl.UpdateSystemData(old.SystemData, *serviceCtx.SystemData())
	if !isNewResource {
		newResource.CreatedAPIVersion = old.CreatedAPIVersion
		prop := newResource.Properties.BasicResourceProperties
		if !old.Properties.BasicResourceProperties.EqualLinkedResource(&prop) {
			return rest.NewLinkedResourceUpdateErrorResponse(serviceCtx.ResourceID, &old.Properties.BasicResourceProperties, &newResource.Properties.BasicResourceProperties), nil
		}
	}

	rendererOutput, err := daprHttpRoute.DeploymentProcessor().Render(ctx, serviceCtx.ResourceID, newResource)
	if err != nil {
		return nil, err
	}
	deploymentOutput, err := daprHttpRoute.DeploymentProcessor().Deploy(ctx, serviceCtx.ResourceID, rendererOutput)
	if err != nil {
		return nil, err
	}

	newResource.Properties = deploymentOutput.RadiusResource.(*datamodel.DaprInvokeHttpRoute).Properties
	newResource.Properties.BasicResourceProperties.Status.OutputResources = deploymentOutput.Resources
	newResource.ComputedValues = deploymentOutput.ComputedValues
	newResource.SecretValues = deploymentOutput.SecretValues

	if !isNewResource {
		diff := outputresource.GetGCOutputResources(newResource.Properties.Status.OutputResources, old.Properties.Status.OutputResources)
		err = daprHttpRoute.DeploymentProcessor().Delete(ctx, deployment.ResourceData{ID: serviceCtx.ResourceID, Resource: newResource, OutputResources: diff, ComputedValues: newResource.ComputedValues, SecretValues: newResource.SecretValues, RecipeData: newResource.RecipeData})
		if err != nil {
			return nil, err
		}
	}

	savedResource, err := daprHttpRoute.SaveResource(ctx, serviceCtx.ResourceID.String(), newResource, etag)
	if err != nil {
		return nil, err
	}

	versioned, err := converter.DaprInvokeHttpRouteDataModelToVersioned(newResource, serviceCtx.APIVersion)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"ETag": savedResource.ETag}

	return rest.NewOKResponseWithHeaders(versioned, headers), nil
}

// Validate extracts versioned resource from request and validates the properties.
func (daprHttpRoute *CreateOrUpdateDaprInvokeHttpRoute) Validate(ctx context.Context, req *http.Request, apiVersion string) (*datamodel.DaprInvokeHttpRoute, error) {
	serviceCtx := v1.ARMRequestContextFromContext(ctx)
	content, err := ctrl.ReadJSONBody(req)
	if err != nil {
		return nil, err
	}

	dm, err := converter.DaprInvokeHttpRouteDataModelFromVersioned(content, apiVersion)
	if err != nil {
		return nil, err
	}

	dm.ID = serviceCtx.ResourceID.String()
	dm.TrackedResource = ctrl.BuildTrackedResource(ctx)
	dm.Properties.ProvisioningState = v1.ProvisioningStateSucceeded
	dm.TenantID = serviceCtx.HomeTenantID
	dm.CreatedAPIVersion = dm.UpdatedAPIVersion

	return dm, nil
}
