// Copyright 2025 The Terraform SOPS backend Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transformer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	defBuckets           = []float64{.001, .002, .003, .004, .005, .01, .025, .05, .1, .25, .5, 1, 2.5}
	vaultRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transformer_vault_request_duration_seconds",
			Help:    "Histogram for the durations to the Vault client.",
			Buckets: defBuckets,
		},
		[]string{"request"},
	)
	keyServiceRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transformer_key_service_request_duration_seconds",
			Help:    "Histogram for the durations to the key service.",
			Buckets: defBuckets,
		},
		[]string{"request", "key_type"},
	)
	transformerRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "transformer_service_request_duration_seconds",
			Help:    "Histogram for the durations to the transformer service.",
			Buckets: defBuckets,
		},
		[]string{"request"},
	)
)
