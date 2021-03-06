/*
 * Pipeline API
 *
 * Pipeline is a feature rich application platform, built for containers on top of Kubernetes to automate the DevOps experience, continuous application development and the lifecycle of deployments. 
 *
 * API version: latest
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package pipeline
// CreateSecretRequest struct for CreateSecretRequest
type CreateSecretRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tags []string `json:"tags,omitempty"`
	Version int32 `json:"version,omitempty"`
	Values map[string]interface{} `json:"values"`
}
