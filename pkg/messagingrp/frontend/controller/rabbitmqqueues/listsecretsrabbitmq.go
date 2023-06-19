/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rabbitmqqueues

import (
	"context"
	"net/http"

	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
	ctrl "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	"github.com/project-radius/radius/pkg/armrpc/rest"
	"github.com/project-radius/radius/pkg/linkrp/frontend/deployment"
	"github.com/project-radius/radius/pkg/linkrp/renderers"
	msg_dm "github.com/project-radius/radius/pkg/messagingrp/datamodel"
	msg_conv "github.com/project-radius/radius/pkg/messagingrp/datamodel/converter"
	frontend_ctrl "github.com/project-radius/radius/pkg/messagingrp/frontend/controller"
)

var _ ctrl.Controller = (*ListSecretsRabbitMQQueue)(nil)

// ListSecretsRabbitMQQueue is the controller implementation to list secrets for the to access the connected rabbitMQ resource resource id passed in the request body.
type ListSecretsRabbitMQQueue struct {
	ctrl.Operation[*msg_dm.RabbitMQQueue, msg_dm.RabbitMQQueue]
	dp deployment.DeploymentProcessor
}

// NewListSecretsRabbitMQQueue creates a new instance of ListSecretsRabbitMQQueue.
func NewListSecretsRabbitMQQueue(opts frontend_ctrl.Options) (ctrl.Controller, error) {
	return &ListSecretsRabbitMQQueue{
		Operation: ctrl.NewOperation(opts.Options,
			ctrl.ResourceOptions[msg_dm.RabbitMQQueue]{
				RequestConverter:  msg_conv.RabbitMQQueueDataModelFromVersioned,
				ResponseConverter: msg_conv.RabbitMQQueueDataModelToVersioned,
			}),
		dp: opts.DeployProcessor,
	}, nil
}

// Run returns secrets values for the specified RabbitMQQueue resource
func (ctrl *ListSecretsRabbitMQQueue) Run(ctx context.Context, w http.ResponseWriter, req *http.Request) (rest.Response, error) {
	sCtx := v1.ARMRequestContextFromContext(ctx)

	// Request route for listsecrets has name of the operation as suffix which should be removed to get the resource id.
	// route id format: subscriptions/<subscription_id>/resourceGroups/<resource_group>/providers/Applications.Messaging/rabbitMQQueues/<resource_name>/listsecrets
	parsedResourceID := sCtx.ResourceID.Truncate()
	resource, _, err := ctrl.GetResource(ctx, parsedResourceID)
	if err != nil {
		return nil, err
	}

	if resource == nil {
		return rest.NewNotFoundResponse(sCtx.ResourceID), nil
	}

	secrets, err := ctrl.dp.FetchSecrets(ctx, deployment.ResourceData{ID: sCtx.ResourceID, Resource: resource, OutputResources: resource.Properties.Status.OutputResources, ComputedValues: resource.ComputedValues, SecretValues: resource.SecretValues})
	if err != nil {
		return nil, err
	}

	msgSecrets := msg_dm.RabbitMQSecrets{}
	if connectionString, ok := secrets[renderers.ConnectionStringValue].(string); ok {
		msgSecrets.ConnectionString = connectionString
	}

	versioned, _ := msg_conv.RabbitMQSecretsDataModelToVersioned(&msgSecrets, sCtx.APIVersion)
	return rest.NewOKResponse(versioned), nil
}
