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

// DaprStateStoresClient contains the methods for the DaprStateStores group.
// Don't use this type directly, use NewDaprStateStoresClient() instead.
type DaprStateStoresClient struct {
	ep string
	pl runtime.Pipeline
	rootScope string
}

// NewDaprStateStoresClient creates a new instance of DaprStateStoresClient with the specified values.
func NewDaprStateStoresClient(con *arm.Connection, rootScope string) *DaprStateStoresClient {
	return &DaprStateStoresClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), rootScope: rootScope}
}

// CreateOrUpdate - Creates or updates a DaprStateStore resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprStateStoresClient) CreateOrUpdate(ctx context.Context, daprStateStoreName string, daprStateStoreParameters DaprStateStoreResource, options *DaprStateStoresCreateOrUpdateOptions) (DaprStateStoresCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, daprStateStoreName, daprStateStoreParameters, options)
	if err != nil {
		return DaprStateStoresCreateOrUpdateResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return DaprStateStoresCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return DaprStateStoresCreateOrUpdateResponse{}, client.createOrUpdateHandleError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *DaprStateStoresClient) createOrUpdateCreateRequest(ctx context.Context, daprStateStoreName string, daprStateStoreParameters DaprStateStoreResource, options *DaprStateStoresCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/daprStateStores/{daprStateStoreName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	if daprStateStoreName == "" {
		return nil, errors.New("parameter daprStateStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{daprStateStoreName}", url.PathEscape(daprStateStoreName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(	client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-15-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, daprStateStoreParameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *DaprStateStoresClient) createOrUpdateHandleResponse(resp *http.Response) (DaprStateStoresCreateOrUpdateResponse, error) {
	result := DaprStateStoresCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprStateStoreResource); err != nil {
		return DaprStateStoresCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *DaprStateStoresClient) createOrUpdateHandleError(resp *http.Response) error {
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

// Delete - Deletes an existing daprStateStore resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprStateStoresClient) Delete(ctx context.Context, daprStateStoreName string, options *DaprStateStoresDeleteOptions) (DaprStateStoresDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, daprStateStoreName, options)
	if err != nil {
		return DaprStateStoresDeleteResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return DaprStateStoresDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return DaprStateStoresDeleteResponse{}, client.deleteHandleError(resp)
	}
	return DaprStateStoresDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *DaprStateStoresClient) deleteCreateRequest(ctx context.Context, daprStateStoreName string, options *DaprStateStoresDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/daprStateStores/{daprStateStoreName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	if daprStateStoreName == "" {
		return nil, errors.New("parameter daprStateStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{daprStateStoreName}", url.PathEscape(daprStateStoreName))
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
func (client *DaprStateStoresClient) deleteHandleError(resp *http.Response) error {
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

// Get - Retrieves information about a daprStateStore resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprStateStoresClient) Get(ctx context.Context, daprStateStoreName string, options *DaprStateStoresGetOptions) (DaprStateStoresGetResponse, error) {
	req, err := client.getCreateRequest(ctx, daprStateStoreName, options)
	if err != nil {
		return DaprStateStoresGetResponse{}, err
	}
	resp, err := 	client.pl.Do(req)
	if err != nil {
		return DaprStateStoresGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return DaprStateStoresGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *DaprStateStoresClient) getCreateRequest(ctx context.Context, daprStateStoreName string, options *DaprStateStoresGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/daprStateStores/{daprStateStoreName}"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
	if daprStateStoreName == "" {
		return nil, errors.New("parameter daprStateStoreName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{daprStateStoreName}", url.PathEscape(daprStateStoreName))
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
func (client *DaprStateStoresClient) getHandleResponse(resp *http.Response) (DaprStateStoresGetResponse, error) {
	result := DaprStateStoresGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprStateStoreResource); err != nil {
		return DaprStateStoresGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *DaprStateStoresClient) getHandleError(resp *http.Response) error {
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

// ListByRootScope - Lists information about all daprStateStore resources in the given root scope
// If the operation fails it returns the *ErrorResponse error type.
func (client *DaprStateStoresClient) ListByRootScope(options *DaprStateStoresListByRootScopeOptions) (*DaprStateStoresListByRootScopePager) {
	return &DaprStateStoresListByRootScopePager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByRootScopeCreateRequest(ctx, options)
		},
		advancer: func(ctx context.Context, resp DaprStateStoresListByRootScopeResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.DaprStateStoreList.NextLink)
		},
	}
}

// listByRootScopeCreateRequest creates the ListByRootScope request.
func (client *DaprStateStoresClient) listByRootScopeCreateRequest(ctx context.Context, options *DaprStateStoresListByRootScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Connector/daprStateStores"
	if client.rootScope == "" {
		return nil, errors.New("parameter client.rootScope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", url.PathEscape(client.rootScope))
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

// listByRootScopeHandleResponse handles the ListByRootScope response.
func (client *DaprStateStoresClient) listByRootScopeHandleResponse(resp *http.Response) (DaprStateStoresListByRootScopeResponse, error) {
	result := DaprStateStoresListByRootScopeResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprStateStoreList); err != nil {
		return DaprStateStoresListByRootScopeResponse{}, err
	}
	return result, nil
}

// listByRootScopeHandleError handles the ListByRootScope error response.
func (client *DaprStateStoresClient) listByRootScopeHandleError(resp *http.Response) error {
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

