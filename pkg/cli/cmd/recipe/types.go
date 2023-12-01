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

package recipe

type EnvironmentRecipe struct {
	Name            string `json:"name"`
	ResourceType    string `json:"resourceType"`
	TemplateKind    string `json:"templateKind"`
	TemplatePath    string `json:"templatePath"`
	TemplateVersion string `json:"templateVersion"`
	PlainHTTP       bool   `json:"plainHTTP"`
}

type RecipeParameter struct {
	Name         string      `json:"name,omitempty"`
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Type         string      `json:"type,omitempty"`
	MaxValue     string      `json:"maxValue,omitempty"`
	MinValue     string      `json:"minValue,omitempty"`
}
