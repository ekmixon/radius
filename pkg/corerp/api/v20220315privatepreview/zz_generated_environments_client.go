//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20220315privatepreview

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// EnvironmentsClient contains the methods for the Environments group.
// Don't use this type directly, use NewEnvironmentsClient() instead.
type EnvironmentsClient struct {
	ep string
	pl runtime.Pipeline
	rootScope string
}

// NewEnvironmentsClient creates a new instance of EnvironmentsClient with the specified values.
func NewEnvironmentsClient(con *arm.Connection, rootScope string) *EnvironmentsClient {
	return &EnvironmentsClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), rootScope: rootScope}
}

// CreateOrUpdate - Create or update an Environment.
// If the operation fails it returns the *ErrorResponse error type.
func (client *EnvironmentsClient) CreateOrUpdate(ctx context.Context, environmentName string, environmentResource EnvironmentResource, options *EnvironmentsCreateOrUpdateOptions) (EnvironmentsCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, environmentName, environmentResource, options)
	if err != nil {
		return EnvironmentsCreateOrUpdateResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return EnvironmentsCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return EnvironmentsCreateOrUpdateResponse{}, client.createOrUpdateHandleError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *EnvironmentsClient) createOrUpdateCreateRequest(ctx context.Context, environmentName string, environmentResource EnvironmentResource, options *EnvironmentsCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/environments/{environmentName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, environmentResource)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *EnvironmentsClient) createOrUpdateHandleResponse(resp *http.Response) (EnvironmentsCreateOrUpdateResponse, error) {
	result := EnvironmentsCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnvironmentResource); err != nil {
		return EnvironmentsCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *EnvironmentsClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Delete - Delete an Environment.
// If the operation fails it returns the *ErrorResponse error type.
func (client *EnvironmentsClient) Delete(ctx context.Context, environmentName string, options *EnvironmentsDeleteOptions) (EnvironmentsDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, environmentName, options)
	if err != nil {
		return EnvironmentsDeleteResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return EnvironmentsDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return EnvironmentsDeleteResponse{}, client.deleteHandleError(resp)
	}
	return EnvironmentsDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *EnvironmentsClient) deleteCreateRequest(ctx context.Context, environmentName string, options *EnvironmentsDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/environments/{environmentName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *EnvironmentsClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Gets the properties of an Environment.
// If the operation fails it returns the *ErrorResponse error type.
func (client *EnvironmentsClient) Get(ctx context.Context, environmentName string, options *EnvironmentsGetOptions) (EnvironmentsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, environmentName, options)
	if err != nil {
		return EnvironmentsGetResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return EnvironmentsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return EnvironmentsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *EnvironmentsClient) getCreateRequest(ctx context.Context, environmentName string, options *EnvironmentsGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/environments/{environmentName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *EnvironmentsClient) getHandleResponse(resp *http.Response) (EnvironmentsGetResponse, error) {
	result := EnvironmentsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnvironmentResource); err != nil {
		return EnvironmentsGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *EnvironmentsClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListByScope - List all environments in a scope.
// If the operation fails it returns the *ErrorResponse error type.
func (client *EnvironmentsClient) ListByScope(options *EnvironmentsListByScopeOptions) (*EnvironmentsListByScopePager) {
	return &EnvironmentsListByScopePager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByScopeCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp EnvironmentsListByScopeResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.EnvironmentResourceList.NextLink)
		},
	}
}

// listByScopeCreateRequest creates the ListByScope request.
func (client *EnvironmentsClient) listByScopeCreateRequest(ctx context.Context, options *EnvironmentsListByScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/environments"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByScopeHandleResponse handles the ListByScope response.
func (client *EnvironmentsClient) listByScopeHandleResponse(resp *http.Response) (EnvironmentsListByScopeResponse, error) {
	result := EnvironmentsListByScopeResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnvironmentResourceList); err != nil {
		return EnvironmentsListByScopeResponse{}, err
	}
	return result, nil
}

// listByScopeHandleError handles the ListByScope error response.
func (client *EnvironmentsClient) listByScopeHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Update - Update the properties of an existing Environment.
// If the operation fails it returns the *ErrorResponse error type.
func (client *EnvironmentsClient) Update(ctx context.Context, environmentName string, environmentResource EnvironmentResource, options *EnvironmentsUpdateOptions) (EnvironmentsUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, environmentName, environmentResource, options)
	if err != nil {
		return EnvironmentsUpdateResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return EnvironmentsUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return EnvironmentsUpdateResponse{}, client.updateHandleError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *EnvironmentsClient) updateCreateRequest(ctx context.Context, environmentName string, environmentResource EnvironmentResource, options *EnvironmentsUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Core/environments/{environmentName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if environmentName == "" {
		return nil, errors.New("parameter environmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{environmentName}", url.PathEscape(environmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, environmentResource)
}

// updateHandleResponse handles the Update response.
func (client *EnvironmentsClient) updateHandleResponse(resp *http.Response) (EnvironmentsUpdateResponse, error) {
	result := EnvironmentsUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.EnvironmentResource); err != nil {
		return EnvironmentsUpdateResponse{}, err
	}
	return result, nil
}

// updateHandleError handles the Update error response.
func (client *EnvironmentsClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
		errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

