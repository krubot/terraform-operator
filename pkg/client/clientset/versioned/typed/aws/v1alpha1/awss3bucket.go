
package v1alpha1

import (
	v1alpha1 "github.com/krubot/terraform-operator/pkg/apis/terraform/v1alpha1"
	scheme "github.com/krubot/terraform-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AwsS3BucketsGetter has a method to return a AwsS3BucketInterface.
// A group's client should implement this interface.
type AwsS3BucketsGetter interface {
	AwsS3Buckets(namespace string) AwsS3BucketInterface
}

// AwsS3BucketInterface has methods to work with AwsS3Bucket resources.
type AwsS3BucketInterface interface {
	Create(*v1alpha1.AwsS3Bucket) (*v1alpha1.AwsS3Bucket, error)
	Update(*v1alpha1.AwsS3Bucket) (*v1alpha1.AwsS3Bucket, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.AwsS3Bucket, error)
	List(opts v1.ListOptions) (*v1alpha1.AwsS3BucketList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AwsS3Bucket, err error)
	AwsS3BucketExpansion
}

// awsS3Buckets implements AwsS3BucketInterface
type awsS3Buckets struct {
	client rest.Interface
	ns     string
}

// newAwsS3Buckets returns a AwsS3Buckets
func newAwsS3Buckets(c *krubotV1alpha1Client, namespace string) *awsS3Buckets {
	return &awsS3Buckets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the awsS3Bucket, and returns the corresponding awsS3Bucket object, and an error if there is any.
func (c *awsS3Buckets) Get(name string, options v1.GetOptions) (result *v1alpha1.AwsS3Bucket, err error) {
	result = &v1alpha1.AwsS3Bucket{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("awss3buckets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AwsS3Buckets that match those selectors.
func (c *awsS3Buckets) List(opts v1.ListOptions) (result *v1alpha1.AwsS3BucketList, err error) {
	result = &v1alpha1.AwsS3BucketList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("awss3buckets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested awsS3Buckets.
func (c *awsS3Buckets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("awss3buckets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a awsS3Bucket and creates it.  Returns the server's representation of the awsS3Bucket, and an error, if there is any.
func (c *awsS3Buckets) Create(awsS3Bucket *v1alpha1.AwsS3Bucket) (result *v1alpha1.AwsS3Bucket, err error) {
	result = &v1alpha1.AwsS3Bucket{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("awss3buckets").
		Body(awsS3Bucket).
		Do().
		Into(result)
	return
}

// Update takes the representation of a awsS3Bucket and updates it. Returns the server's representation of the awsS3Bucket, and an error, if there is any.
func (c *awsS3Buckets) Update(awsS3Bucket *v1alpha1.AwsS3Bucket) (result *v1alpha1.AwsS3Bucket, err error) {
	result = &v1alpha1.AwsS3Bucket{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("awss3buckets").
		Name(awsS3Bucket.Name).
		Body(awsS3Bucket).
		Do().
		Into(result)
	return
}

// Delete takes name of the awsS3Bucket and deletes it. Returns an error if one occurs.
func (c *awsS3Buckets) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("awss3buckets").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *awsS3Buckets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("awss3buckets").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched awsS3Bucket.
func (c *awsS3Buckets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AwsS3Bucket, err error) {
	result = &v1alpha1.AwsS3Bucket{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("awss3buckets").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
