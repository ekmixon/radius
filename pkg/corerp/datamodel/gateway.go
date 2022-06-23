// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package datamodel

import (
	v1 "github.com/project-radius/radius/pkg/armrpc/api/v1"
)

// Gateway represents Gateway resource.
type Gateway struct {
	v1.TrackedResource

	// SystemData is the systemdata which includes creation/modified dates.
	SystemData v1.SystemData `json:"systemData,omitempty"`
	// Properties is the properties of the resource.
	Properties GatewayProperties `json:"properties"`

	// InternalMetadata is the internal metadata which is used for conversion.
	v1.InternalMetadata
}

// ResourceTypeName returns the qualified name of the resource
func (g Gateway) ResourceTypeName() string {
	return "Applications.Core/gateways"
}

// GatewayProperties represents the properties of Gateway.
type GatewayProperties struct {
	v1.BasicResourceProperties
	ProvisioningState v1.ProvisioningState       `json:"provisioningState,omitempty"`
	Application       string                     `json:"application,omitempty"`
	Internal          bool                       `json:"internal,omitempty"`
	Hostname          *GatewayPropertiesHostname `json:"hostname,omitempty"`
	Routes            []GatewayRoute             `json:"routes,omitempty"`
}

// GatewayRoute represents the route attached to Gateway.
type GatewayRoute struct {
	Destination   string `json:"destination,omitempty"`
	Path          string `json:"path,omitempty"`
	ReplacePrefix string `json:"replacePrefix,omitempty"`
}

// GatewayPropertiesHostname - Declare hostname information for the Gateway.
type GatewayPropertiesHostname struct {
	FullyQualifiedHostname string `json:"fullyQualifiedHostname,omitempty"`
	Prefix                 string `json:"prefix,omitempty"`
}
