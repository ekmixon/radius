// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package daprstatestorev1alpha3

import (
	"github.com/project-radius/radius/pkg/azure/radclient"
	"github.com/project-radius/radius/pkg/handlers"
	"github.com/project-radius/radius/pkg/radrp/outputresource"
	"github.com/project-radius/radius/pkg/renderers"
	"github.com/project-radius/radius/pkg/resourcekinds"
)

func GetDaprStateStoreAzureStorage(resource renderers.RendererResource) ([]outputresource.OutputResource, error) {
	properties := radclient.DaprStateStoreAzureTableStorageResourceProperties{}
	err := resource.ConvertDefinition(&properties)
	if err != nil {
		return nil, err
	}
	resourceKind := resourcekinds.DaprStateStoreAzureStorage
	localID := outputresource.LocalIDDaprStateStoreAzureStorage

	if properties.Resource == nil || *properties.Resource == "" {
		return nil, renderers.ErrResourceMissingForResource
	}
	accountID, err := renderers.ValidateResourceID(*properties.Resource, StorageAccountResourceType, "Storage Account")
	if err != nil {
		return nil, err
	}

	// generate data we can use to connect to a Storage Account
	outputResource := outputresource.OutputResource{
		LocalID:      localID,
		ResourceKind: resourceKind,
		Resource: map[string]string{
			handlers.KubernetesNameKey:       resource.ResourceName,
			handlers.KubernetesNamespaceKey:  resource.ApplicationName,
			handlers.KubernetesAPIVersionKey: "dapr.io/v1alpha1",
			handlers.KubernetesKindKey:       "Component",

			handlers.StorageAccountIDKey:   accountID.ID,
			handlers.StorageAccountNameKey: accountID.Types[0].Name,
		},
	}
	return []outputresource.OutputResource{outputResource}, nil
}
