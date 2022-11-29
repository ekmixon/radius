// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package rediscaches

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	ctrl "github.com/project-radius/radius/pkg/armrpc/frontend/controller"
	radiustesting "github.com/project-radius/radius/pkg/corerp/testing"
	"github.com/project-radius/radius/pkg/linkrp/api/v20220315privatepreview"
	"github.com/project-radius/radius/pkg/linkrp/datamodel"
	"github.com/project-radius/radius/pkg/linkrp/frontend/deployment"
	"github.com/project-radius/radius/pkg/linkrp/renderers"
	"github.com/project-radius/radius/pkg/resourcekinds"
	"github.com/project-radius/radius/pkg/resourcemodel"
	"github.com/project-radius/radius/pkg/rp"
	"github.com/project-radius/radius/pkg/rp/outputresource"
	"github.com/project-radius/radius/pkg/ucp/store"
	"github.com/stretchr/testify/require"
)

func getDeploymentProcessorOutputs(buildComputedValueReferences bool) (renderers.RendererOutput, deployment.DeploymentOutput) {
	var computedValues map[string]renderers.ComputedValueReference
	var portValue interface{}
	if buildComputedValueReferences {
		computedValues = map[string]renderers.ComputedValueReference{
			renderers.Host: {
				LocalID:     outputresource.LocalIDAzureRedis,
				JSONPointer: "/properties/hostName",
			},
			renderers.Port: {
				LocalID:     outputresource.LocalIDAzureRedis,
				JSONPointer: "/properties/sslPort",
			},
		}

		portValue = "10255"
	} else {
		portValue = float64(10255)

		computedValues = map[string]renderers.ComputedValueReference{
			renderers.Host: {
				Value: "myrediscache.redis.cache.windows.net",
			},
			renderers.Port: {
				Value: portValue,
			},
		}
	}

	rendererOutput := renderers.RendererOutput{
		Resources: []outputresource.OutputResource{
			{
				LocalID: outputresource.LocalIDAzureRedis,
				ResourceType: resourcemodel.ResourceType{
					Type:     resourcekinds.AzureRedis,
					Provider: resourcemodel.ProviderAzure,
				},
				Identity: resourcemodel.ResourceIdentity{},
			},
		},
		SecretValues: map[string]rp.SecretValueReference{
			renderers.ConnectionStringValue: {Value: "test-connection-string"},
			renderers.PasswordStringHolder:  {Value: "testpassword"},
		},
		ComputedValues: computedValues,
	}

	deploymentOutput := deployment.DeploymentOutput{
		Resources: []outputresource.OutputResource{
			{
				LocalID: outputresource.LocalIDAzureRedis,
				ResourceType: resourcemodel.ResourceType{
					Type:     resourcekinds.AzureRedis,
					Provider: resourcemodel.ProviderAzure,
				},
			},
		},
		ComputedValues: map[string]interface{}{
			renderers.Host: "myrediscache.redis.cache.windows.net",
			renderers.Port: portValue,
		},
	}

	return rendererOutput, deploymentOutput
}

func TestCreateOrUpdateRedisCache_20220315PrivatePreview(t *testing.T) {
	mctrl := gomock.NewController(t)
	defer mctrl.Finish()

	mStorageClient := store.NewMockStorageClient(mctrl)
	mDeploymentProcessor := deployment.NewMockDeploymentProcessor(mctrl)
	rendererOutput, deploymentOutput := getDeploymentProcessorOutputs(false)
	ctx := context.Background()

	createNewResourceTestCases := []struct {
		desc               string
		headerKey          string
		headerValue        string
		resourceETag       string
		expectedStatusCode int
		shouldFail         bool
		azureResource      bool
	}{
		{"create-new-resource-no-if-match", "If-Match", "", "", http.StatusOK, false, false},
		{"create-new-resource-*-if-match", "If-Match", "*", "", http.StatusPreconditionFailed, true, false},
		{"create-new-resource-etag-if-match", "If-Match", "random-etag", "", http.StatusPreconditionFailed, true, false},
		{"create-new-resource-*-if-none-match", "If-None-Match", "*", "", http.StatusOK, false, true},
	}

	for _, testcase := range createNewResourceTestCases {
		t.Run(testcase.desc, func(t *testing.T) {
			if testcase.azureResource {
				rendererOutput, deploymentOutput = getDeploymentProcessorOutputs(true)
			}

			input, dataModel, expectedOutput := getTestModelsForGetAndListApis20220315privatepreview()
			w := httptest.NewRecorder()
			req, _ := radiustesting.GetARMTestHTTPRequest(ctx, http.MethodGet, testHeaderfile, input)
			req.Header.Set(testcase.headerKey, testcase.headerValue)
			ctx := radiustesting.ARMTestContextFromRequest(req)

			mStorageClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
					return nil, &store.ErrNotFound{}
				})

			expectedOutput.SystemData.CreatedAt = expectedOutput.SystemData.LastModifiedAt
			expectedOutput.SystemData.CreatedBy = expectedOutput.SystemData.LastModifiedBy
			expectedOutput.SystemData.CreatedByType = expectedOutput.SystemData.LastModifiedByType

			if !testcase.shouldFail {
				deploymentOutput.RadiusResource = dataModel
				deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Host = deploymentOutput.ComputedValues[renderers.Host].(string)

				port := deploymentOutput.ComputedValues[renderers.Port]
				if port != nil {
					switch p := port.(type) {
					case float64:
						deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Port = int32(p)
					case int32:
						deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Port = p
					case string:
						converted, _ := strconv.Atoi(p)
						deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Port = int32(converted)
					default:
						panic("unhandled type for the property portx")
					}
				}
				mDeploymentProcessor.EXPECT().Render(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(rendererOutput, nil)
				mDeploymentProcessor.EXPECT().Deploy(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(deploymentOutput, nil)

				mStorageClient.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, obj *store.Object, opts ...store.SaveOptions) error {
						// First time created objects should have the same lastModifiedAt and createdAt
						dataModel.SystemData.CreatedAt = dataModel.SystemData.LastModifiedAt
						obj.ETag = "new-resource-etag"
						obj.Data = dataModel
						return nil
					})
			}

			opts := ctrl.Options{
				StorageClient: mStorageClient,
				GetDeploymentProcessor: func() deployment.DeploymentProcessor {
					return mDeploymentProcessor
				},
			}

			ctl, err := NewCreateOrUpdateRedisCache(opts)
			require.NoError(t, err)
			resp, err := ctl.Run(ctx, w, req)
			require.NoError(t, err)
			_ = resp.Apply(ctx, w, req)
			require.Equal(t, testcase.expectedStatusCode, w.Result().StatusCode)
			if !testcase.shouldFail {
				actualOutput := &v20220315privatepreview.RedisCacheResource{}
				_ = json.Unmarshal(w.Body.Bytes(), actualOutput)
				require.Equal(t, expectedOutput, actualOutput)

				require.Equal(t, "new-resource-etag", w.Header().Get("ETag"))
			}
		})
	}

	updateExistingResourceTestCases := []struct {
		desc               string
		headerKey          string
		headerValue        string
		inputFile          string
		resourceETag       string
		expectedStatusCode int
		shouldFail         bool
	}{
		{"update-resource-no-if-match", "If-Match", "", "", "resource-etag", http.StatusOK, false},
		{"update-resource-with-diff-app", "If-Match", "", "20220315privatepreview_input_diff_app.json", "resource-etag", http.StatusBadRequest, true},
		{"update-resource-*-if-match", "If-Match", "*", "", "resource-etag", http.StatusOK, false},
		{"update-resource-matching-if-match", "If-Match", "matching-etag", "", "matching-etag", http.StatusOK, false},
		{"update-resource-not-matching-if-match", "If-Match", "not-matching-etag", "", "another-etag", http.StatusPreconditionFailed, true},
		{"update-resource-*-if-none-match", "If-None-Match", "*", "", "another-etag", http.StatusPreconditionFailed, true},
	}

	for _, testcase := range updateExistingResourceTestCases {
		t.Run(testcase.desc, func(t *testing.T) {
			input, dataModel, expectedOutput := getTestModelsForGetAndListApis20220315privatepreview()
			if testcase.inputFile != "" {
				input = &v20220315privatepreview.RedisCacheResource{}
				_ = json.Unmarshal(radiustesting.ReadFixture(testcase.inputFile), input)
			}
			w := httptest.NewRecorder()
			req, _ := radiustesting.GetARMTestHTTPRequest(ctx, http.MethodGet, testHeaderfile, input)
			req.Header.Set(testcase.headerKey, testcase.headerValue)
			ctx := radiustesting.ARMTestContextFromRequest(req)

			mStorageClient.
				EXPECT().
				Get(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, id string, _ ...store.GetOptions) (*store.Object, error) {
					return &store.Object{
						Metadata: store.Metadata{ID: id, ETag: testcase.resourceETag},
						Data:     dataModel,
					}, nil
				})

			if !testcase.shouldFail {
				deploymentOutput.RadiusResource = dataModel
				deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Host = deploymentOutput.ComputedValues[renderers.Host].(string)

				port := deploymentOutput.ComputedValues[renderers.Port]
				if port != nil {
					switch p := port.(type) {
					case float64:
						deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Port = int32(p)
					case int32:
						deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Port = p
					case string:
						converted, _ := strconv.Atoi(p)
						deploymentOutput.RadiusResource.(*datamodel.RedisCache).Properties.Port = int32(converted)
					default:
						panic("unhandled type for the property portx")
					}
				}
				mDeploymentProcessor.EXPECT().Render(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(rendererOutput, nil)
				mDeploymentProcessor.EXPECT().Deploy(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(deploymentOutput, nil)
				mDeploymentProcessor.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(1).Return(nil)

				mStorageClient.
					EXPECT().
					Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, obj *store.Object, opts ...store.SaveOptions) error {
						obj.ETag = "updated-resource-etag"
						obj.Data = dataModel
						return nil
					})
			}

			opts := ctrl.Options{
				StorageClient: mStorageClient,
				GetDeploymentProcessor: func() deployment.DeploymentProcessor {
					return mDeploymentProcessor
				},
			}

			ctl, err := NewCreateOrUpdateRedisCache(opts)
			require.NoError(t, err)
			resp, err := ctl.Run(ctx, w, req)
			_ = resp.Apply(ctx, w, req)
			require.NoError(t, err)
			require.Equal(t, testcase.expectedStatusCode, w.Result().StatusCode)

			if !testcase.shouldFail {
				actualOutput := &v20220315privatepreview.RedisCacheResource{}
				_ = json.Unmarshal(w.Body.Bytes(), actualOutput)
				require.Equal(t, expectedOutput, actualOutput)

				require.Equal(t, "updated-resource-etag", w.Header().Get("ETag"))
			}
		})
	}
}
