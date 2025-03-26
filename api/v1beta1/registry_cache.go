package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RegistryCache struct {
	// Upstream is the remote registry host to cache.
	Upstream string `json:"upstream"`
	// RemoteURL is the remote registry URL. The format must be `<scheme><host>[:<port>]` where
	// `<scheme>` is `https://` or `http://` and `<host>[:<port>]` corresponds to the Upstream
	//
	// If defined, the value is set as `proxy.remoteurl` in the registry [configuration](https://github.com/distribution/distribution/blob/main/docs/content/recipes/mirror.md#configure-the-cache)
	// and in containerd configuration as `server` field in [hosts.toml](https://github.com/containerd/containerd/blob/main/docs/hosts.md#server-field) file.
	// +optional
	RemoteURL *string `json:"remoteURL,omitempty"`
	// Volume contains settings for the registry cache volume.
	// +optional
	Volume *Volume `json:"volume,omitempty"`
	// GarbageCollection contains settings for the garbage collection of content from the cache.
	// Defaults to enabled garbage collection.
	// +optional
	GarbageCollection *GarbageCollection `json:"garbageCollection,omitempty"`
	// SecretReferenceName is the name of the reference for the Secret containing the upstream registry credentials.
	// +optional
	SecretReferenceName *string `json:"secretReferenceName,omitempty"`
	// Proxy contains settings for a proxy used in the registry cache.
	// +optional
	Proxy *Proxy `json:"proxy,omitempty"`

	// HTTP contains settings for the HTTP server that hosts the registry cache.
	HTTP *HTTP `json:"http,omitempty"`
}

// Volume contains settings for the registry cache volume.
type Volume struct {
	// Size is the size of the registry cache volume.
	// Defaults to 10Gi.
	// This field is immutable.
	// +optional
	// +default="10Gi"
	Size *resource.Quantity `json:"size,omitempty"`
	// StorageClassName is the name of the StorageClass used by the registry cache volume.
	// This field is immutable.
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`
}

// GarbageCollection contains settings for the garbage collection of content from the cache.
type GarbageCollection struct {
	// TTL is the time to live of a blob in the cache.
	// Set to 0s to disable the garbage collection.
	// Defaults to 168h (7 days).
	// +default="168h"
	TTL metav1.Duration `json:"ttl"`
}

// Proxy contains settings for a proxy used in the registry cache.
type Proxy struct {
	// HTTPProxy field represents the proxy server for HTTP connections which is used by the registry cache.
	// +optional
	HTTPProxy *string `json:"httpProxy,omitempty"`
	// HTTPSProxy field represents the proxy server for HTTPS connections which is used by the registry cache.
	// +optional
	HTTPSProxy *string `json:"httpsProxy,omitempty"`
}

// HTTP contains settings for the HTTP server that hosts the registry cache.
type HTTP struct {
	// TLS indicates whether TLS is enabled for the HTTP server of the registry cache.
	// Defaults to true.
	TLS bool `json:"tls,omitempty"`
}
