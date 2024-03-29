// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"./pkg/apis/fnworkflows/v1alpha1.CustomerImageBuild":       schema_pkg_apis_fnworkflows_v1alpha1_CustomerImageBuild(ref),
		"./pkg/apis/fnworkflows/v1alpha1.CustomerImageBuildSpec":   schema_pkg_apis_fnworkflows_v1alpha1_CustomerImageBuildSpec(ref),
		"./pkg/apis/fnworkflows/v1alpha1.CustomerImageBuildStatus": schema_pkg_apis_fnworkflows_v1alpha1_CustomerImageBuildStatus(ref),
	}
}

func schema_pkg_apis_fnworkflows_v1alpha1_CustomerImageBuild(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "CustomerImageBuild is the Schema for the customerimagebuilds API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/fnworkflows/v1alpha1.CustomerImageBuildSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/fnworkflows/v1alpha1.CustomerImageBuildStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/fnworkflows/v1alpha1.CustomerImageBuildSpec", "./pkg/apis/fnworkflows/v1alpha1.CustomerImageBuildStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_fnworkflows_v1alpha1_CustomerImageBuildSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "CustomerImageBuildSpec defines the desired state of CustomerImageBuild",
				Properties: map[string]spec.Schema{
					"repoURL": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"gitRef": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"imageRepo": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"imageTag": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"retries": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int32",
						},
					},
				},
				Required: []string{"repoURL", "gitRef", "imageRepo", "imageTag"},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_fnworkflows_v1alpha1_CustomerImageBuildStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "CustomerImageBuildStatus defines the observed state of CustomerImageBuild",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
