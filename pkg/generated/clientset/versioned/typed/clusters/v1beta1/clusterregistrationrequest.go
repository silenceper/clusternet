/*
Copyright 2021 The Clusternet Authors.

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
// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/clusternet/clusternet/pkg/apis/clusters/v1beta1"
	scheme "github.com/clusternet/clusternet/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ClusterRegistrationRequestsGetter has a method to return a ClusterRegistrationRequestInterface.
// A group's client should implement this interface.
type ClusterRegistrationRequestsGetter interface {
	ClusterRegistrationRequests() ClusterRegistrationRequestInterface
}

// ClusterRegistrationRequestInterface has methods to work with ClusterRegistrationRequest resources.
type ClusterRegistrationRequestInterface interface {
	Create(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.CreateOptions) (*v1beta1.ClusterRegistrationRequest, error)
	Update(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.UpdateOptions) (*v1beta1.ClusterRegistrationRequest, error)
	UpdateStatus(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.UpdateOptions) (*v1beta1.ClusterRegistrationRequest, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.ClusterRegistrationRequest, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.ClusterRegistrationRequestList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ClusterRegistrationRequest, err error)
	ClusterRegistrationRequestExpansion
}

// clusterRegistrationRequests implements ClusterRegistrationRequestInterface
type clusterRegistrationRequests struct {
	client rest.Interface
}

// newClusterRegistrationRequests returns a ClusterRegistrationRequests
func newClusterRegistrationRequests(c *ClustersV1beta1Client) *clusterRegistrationRequests {
	return &clusterRegistrationRequests{
		client: c.RESTClient(),
	}
}

// Get takes name of the clusterRegistrationRequest, and returns the corresponding clusterRegistrationRequest object, and an error if there is any.
func (c *clusterRegistrationRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.ClusterRegistrationRequest, err error) {
	result = &v1beta1.ClusterRegistrationRequest{}
	err = c.client.Get().
		Resource("clusterregistrationrequests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ClusterRegistrationRequests that match those selectors.
func (c *clusterRegistrationRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.ClusterRegistrationRequestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.ClusterRegistrationRequestList{}
	err = c.client.Get().
		Resource("clusterregistrationrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested clusterRegistrationRequests.
func (c *clusterRegistrationRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("clusterregistrationrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a clusterRegistrationRequest and creates it.  Returns the server's representation of the clusterRegistrationRequest, and an error, if there is any.
func (c *clusterRegistrationRequests) Create(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.CreateOptions) (result *v1beta1.ClusterRegistrationRequest, err error) {
	result = &v1beta1.ClusterRegistrationRequest{}
	err = c.client.Post().
		Resource("clusterregistrationrequests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterRegistrationRequest).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a clusterRegistrationRequest and updates it. Returns the server's representation of the clusterRegistrationRequest, and an error, if there is any.
func (c *clusterRegistrationRequests) Update(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.UpdateOptions) (result *v1beta1.ClusterRegistrationRequest, err error) {
	result = &v1beta1.ClusterRegistrationRequest{}
	err = c.client.Put().
		Resource("clusterregistrationrequests").
		Name(clusterRegistrationRequest.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterRegistrationRequest).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *clusterRegistrationRequests) UpdateStatus(ctx context.Context, clusterRegistrationRequest *v1beta1.ClusterRegistrationRequest, opts v1.UpdateOptions) (result *v1beta1.ClusterRegistrationRequest, err error) {
	result = &v1beta1.ClusterRegistrationRequest{}
	err = c.client.Put().
		Resource("clusterregistrationrequests").
		Name(clusterRegistrationRequest.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(clusterRegistrationRequest).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the clusterRegistrationRequest and deletes it. Returns an error if one occurs.
func (c *clusterRegistrationRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("clusterregistrationrequests").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *clusterRegistrationRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("clusterregistrationrequests").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched clusterRegistrationRequest.
func (c *clusterRegistrationRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ClusterRegistrationRequest, err error) {
	result = &v1beta1.ClusterRegistrationRequest{}
	err = c.client.Patch(pt).
		Resource("clusterregistrationrequests").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
