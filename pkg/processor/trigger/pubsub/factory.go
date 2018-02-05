/*
Copyright 2017 The Nuclio Authors.

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

package pubsub

import (
    "github.com/nuclio/nuclio/pkg/errors"
    "github.com/nuclio/nuclio/pkg/functionconfig"
    "github.com/nuclio/nuclio/pkg/processor/runtime"
    "github.com/nuclio/nuclio/pkg/processor/trigger"
    "github.com/nuclio/nuclio/pkg/processor/worker"

    "github.com/nuclio/logger"
)

type factory struct{}

func (f *factory) Create(parentLogger logger.Logger,
    ID string,
    triggerConfiguration *functionconfig.Trigger,
    runtimeConfiguration *runtime.Configuration) (trigger.Trigger, error) {

    // create logger parent
    pubsubLogger := parentLogger.GetChild("pubsub")

    configuration, err := NewConfiguration(ID, triggerConfiguration)
    if err != nil {
        return nil, errors.Wrap(err, "Failed to create configuration")
    }

    // create worker allocator
    workerAllocator, err := worker.WorkerFactorySingleton.CreateFixedPoolWorkerAllocator(pubsubLogger,
        configuration.MaxWorkers,
        runtimeConfiguration)

    if err != nil {
        return nil, errors.Wrap(err, "Failed to create worker allocator")
    }

    // finally, create the trigger
    pubsubTrigger, err := newTrigger(pubsubLogger,
        workerAllocator,
        configuration,
    )

    if err != nil {
        return nil, errors.Wrap(err, "Failed to create pubsub trigger")
    }

    return pubsubTrigger, nil
}

// this called from cmd/processor/app/processor
func init() {
    trigger.RegistrySingleton.Register("pubsub", &factory{})
}
