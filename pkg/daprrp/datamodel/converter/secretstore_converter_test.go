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

package converter

import (
	"encoding/json"
	"errors"
	"testing"

	v1 "github.com/radius-project/radius/pkg/armrpc/api/v1"
	"github.com/radius-project/radius/pkg/daprrp/api/v20220315privatepreview"
	"github.com/radius-project/radius/pkg/daprrp/datamodel"
	linkrp_util "github.com/radius-project/radius/pkg/linkrp/api/v20220315privatepreview"
	"github.com/stretchr/testify/require"
)

// Validates type conversion between versioned client side data model and RP data model.
func TestDaprSecretStoreDataModelToVersioned(t *testing.T) {
	testset := []struct {
		dataModelFile string
		apiVersion    string
		apiModelType  any
		err           error
	}{
		{
			"../../api/v20220315privatepreview/testdata/secretstore_manual_resourcedatamodel.json",
			"2022-03-15-privatepreview",
			&v20220315privatepreview.DaprSecretStoreResource{},
			nil,
		},
		{
			"../../api/v20220315privatepreview/testdata/secretstore_manual_resourcedatamodel.json",
			"unsupported",
			nil,
			v1.ErrUnsupportedAPIVersion,
		},
	}

	for _, tc := range testset {
		t.Run(tc.apiVersion, func(t *testing.T) {
			c, err := linkrp_util.LoadTestData(tc.dataModelFile)
			require.NoError(t, err)
			dm := &datamodel.DaprSecretStore{}
			_ = json.Unmarshal(c, dm)
			am, err := SecretStoreDataModelToVersioned(dm, tc.apiVersion)
			if tc.err != nil {
				require.ErrorAs(t, tc.err, &err)
			} else {
				require.NoError(t, err)
				require.IsType(t, tc.apiModelType, am)
			}
		})
	}
}

func TestDaprSecretStoreDataModelFromVersioned(t *testing.T) {
	testset := []struct {
		versionedModelFile string
		apiVersion         string
		err                error
	}{
		{
			"../../api/v20220315privatepreview/testdata/secretstore_manual_resource.json",
			"2022-03-15-privatepreview",
			nil,
		},
		{
			"../../api/v20220315privatepreview/testdata/secretstore_invalidrecipe_resource.json",
			"2022-03-15-privatepreview",
			errors.New("json: cannot unmarshal number into Go struct field DaprSecretStoreProperties.properties.version of type string"),
		},
		{
			"../../api/v20220315privatepreview/testdata/secretstore_invalidvalues_resource.json",
			"2022-03-15-privatepreview",
			&v1.ErrClientRP{Code: "BadRequest", Message: "error(s) found:\n\trecipe details cannot be specified when resourceProvisioning is set to manual\n\tmetadata must be specified when resourceProvisioning is set to manual\n\ttype must be specified when resourceProvisioning is set to manual\n\tversion must be specified when resourceProvisioning is set to manual"},
		},
		{
			"../../api/v20220315privatepreview/testdata/secretstore_invalidvalues_resource.json",
			"unsupported",
			v1.ErrUnsupportedAPIVersion,
		},
	}

	for _, tc := range testset {
		t.Run(tc.apiVersion, func(t *testing.T) {
			c, err := linkrp_util.LoadTestData(tc.versionedModelFile)
			require.NoError(t, err)
			dm, err := SecretStoreDataModelFromVersioned(c, tc.apiVersion)
			if tc.err != nil {
				require.ErrorAs(t, tc.err, &err)
			} else {
				require.NoError(t, err)
				require.IsType(t, tc.apiVersion, dm.InternalMetadata.UpdatedAPIVersion)
			}
		})
	}
}
