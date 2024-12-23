/*
Copyright 2021 The Dapr Authors
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

package common

import (
	"context"

	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/actor/config"
)

const (
	// AppAPITokenEnvVar is the environment variable for app api token.
	AppAPITokenEnvVar = "APP_API_TOKEN"  /* #nosec */
	APITokenKey       = "dapr-api-token" /* #nosec */
)

// Service represents Dapr callback service.
//
//nolint:interfacebloat
type Service interface {
	// AddHealthCheckHandler sets a health check handler, name: http (router) and grpc (invalid).
	AddHealthCheckHandler(name string, fn HealthCheckHandler) error
	// AddServiceInvocationHandler appends provided service invocation handler with its name to the service.
	AddServiceInvocationHandler(name string, fn ServiceInvocationHandler) error
	// AddTopicEventHandler appends provided event handler with its topic and optional metadata to the service.
	// Note, retries are only considered when there is an error. Lack of error is considered as a success.
	AddTopicEventHandler(sub *Subscription, fn TopicEventHandler) error
	// AddTopicEventSubscriber appends the provided subscriber with its topic and optional metadata to the service.
	// Note, retries are only considered when there is an error. Lack of error is considered as a success.
	AddTopicEventSubscriber(sub *Subscription, subscriber TopicEventSubscriber) error
	// AddBindingInvocationHandler appends provided binding invocation handler with its name to the service.
	AddBindingInvocationHandler(name string, fn BindingInvocationHandler) error
	// RegisterActorImplFactory Register a new actor to actor runtime of go sdk
	// Deprecated: use RegisterActorImplFactoryContext instead.
	RegisterActorImplFactory(f actor.Factory, opts ...config.Option)
	// RegisterActorImplFactoryContext Register a new actor to actor runtime of go sdk
	RegisterActorImplFactoryContext(f actor.FactoryContext, opts ...config.Option)
	// AddJobEventHandler appends a provided job event handler with its name to
	// the service.
	AddJobEventHandler(name string, fn JobEventHandler) error
	// Start starts service.
	Start() error
	// Stop stops the previously started service.
	Stop() error
	// Gracefully stops the previous started service.
	GracefulStop() error
}

type (
	ServiceInvocationHandler func(ctx context.Context, in *InvocationEvent) (out *Content, err error)
	TopicEventHandler        func(ctx context.Context, e *TopicEvent) (retry bool, err error)
	BindingInvocationHandler func(ctx context.Context, in *BindingEvent) (out []byte, err error)
	JobEventHandler          func(ctx context.Context, in *JobEvent) error
	HealthCheckHandler       func(context.Context) error
)

type TopicEventSubscriber interface {
	Handle(ctx context.Context, e *TopicEvent) (retry bool, err error)
}

// Handle converts TopicEventHandler into an adapter that implements
// TopicEventSubscriber.
func (h TopicEventHandler) Handle(ctx context.Context, e *TopicEvent) (retry bool, err error) {
	return h(ctx, e)
}
