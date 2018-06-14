/*
Copyright 2018 The Kubernetes Authors.

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

package neg

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/ingress-gce/pkg/metrics"
)

const (
	negControllerSubsystem = "neg_controller"
	syncLatencyKey         = "neg_sync_duration_seconds"
	lastSyncTimestampKey   = "sync_timestamp"

	resultSuccess = "success"
	resultError   = "error"

	attachSync = syncType("attach")
	detachSync = syncType("detach")
)

type syncType string

var (
	syncMetricsLabels = []string{
		"key",    // The key to uniquely identify the NEG syncer.
		"type",   // Type of the NEG sync
		"result", // Result of the sync.
	}

	syncLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metrics.GLBC_NAMESPACE,
			Subsystem: negControllerSubsystem,
			Name:      syncLatencyKey,
			Help:      "Sync latency of a NEG syncer",
		},
		syncMetricsLabels,
	)

	lastSyncTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metrics.GLBC_NAMESPACE,
			Subsystem: negControllerSubsystem,
			Name:      lastSyncTimestampKey,
			Help:      "The timestamp of the last execution of NEG controller sync loop.",
		},
		[]string{},
	)
)

var register sync.Once

func registerMetrics() {
	register.Do(func() {
		prometheus.MustRegister(syncLatency)
		prometheus.MustRegister(lastSyncTimestamp)
	})
}

// observeNegSync publish collected metrics for the sync of NEG
func observeNegSync(negName string, syncType syncType, err error, start time.Time) {
	result := resultSuccess
	if err != nil {
		result = resultError
	}
	syncLatency.WithLabelValues(negName, string(syncType), result).Observe(time.Since(start).Seconds())
}
