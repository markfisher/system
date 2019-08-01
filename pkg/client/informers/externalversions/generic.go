/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package externalversions

import (
	"fmt"

	v1alpha1 "github.com/projectriff/system/pkg/apis/build/v1alpha1"
	knative_v1alpha1 "github.com/projectriff/system/pkg/apis/knative/v1alpha1"
	streaming_v1alpha1 "github.com/projectriff/system/pkg/apis/streaming/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=build.projectriff.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("applications"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Build().V1alpha1().Applications().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("containers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Build().V1alpha1().Containers().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("functions"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Build().V1alpha1().Functions().Informer()}, nil

		// Group=knative.projectriff.io, Version=v1alpha1
	case knative_v1alpha1.SchemeGroupVersion.WithResource("adapters"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Knative().V1alpha1().Adapters().Informer()}, nil
	case knative_v1alpha1.SchemeGroupVersion.WithResource("handlers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Knative().V1alpha1().Handlers().Informer()}, nil

		// Group=streaming.projectriff.io, Version=v1alpha1
	case streaming_v1alpha1.SchemeGroupVersion.WithResource("processors"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Streaming().V1alpha1().Processors().Informer()}, nil
	case streaming_v1alpha1.SchemeGroupVersion.WithResource("streams"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Streaming().V1alpha1().Streams().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
