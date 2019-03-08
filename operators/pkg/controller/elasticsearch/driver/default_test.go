// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package driver

import (
	"encoding/json"
	"testing"

	"github.com/elastic/k8s-operators/operators/pkg/controller/common/reconciler"
	esclient "github.com/elastic/k8s-operators/operators/pkg/controller/elasticsearch/client"
	"github.com/elastic/k8s-operators/operators/pkg/controller/elasticsearch/observer"
	"github.com/elastic/k8s-operators/operators/pkg/controller/elasticsearch/reconcile"
	"github.com/elastic/k8s-operators/operators/pkg/utils/k8s"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	k8sreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	ClusterStateSample = `
{
    "cluster_name": "elasticsearch-sample",
    "compressed_size_in_bytes": 10281,
    "cluster_uuid": "fW1CurdKQpa-vsEYgTwkvg",
    "version": 28,
    "state_uuid": "0_7Tkm3ERdeB5eOqEgdOcA",
    "master_node": "EizpW8QWRty_T1nJpr-dNQ",
    "nodes": {
        "EizpW8QWRty_T1nJpr-dNQ": {
            "name": "elasticsearch-sample-es-fnsgkkdl85",
            "ephemeral_id": "hd8VlWVdTlyCriXKDW-5kg",
            "transport_address": "172.17.0.10:9300",
            "attributes": {
                "xpack.installed": "true"
            }
        },
        "NRqCLTmhTLuSxzlWcTae3A": {
            "name": "elasticsearch-sample-es-79gc6p57rs",
            "ephemeral_id": "VHAy3TOxTby3fNaPpMgfkg",
            "transport_address": "172.17.0.9:9300",
            "attributes": {
                "xpack.installed": "true"
            }
        },
        "q--ANfDnTKW2WS9pEBuLWQ": {
            "name": "elasticsearch-sample-es-jfpqbt2s4q",
            "ephemeral_id": "USglep8YTW-4vZ9M7PyRqA",
            "transport_address": "172.17.0.7:9300",
            "attributes": {
                "xpack.installed": "true"
            }
        }
    },
    "routing_table": {
        "indices": {
            "shakespeare": {
                "shards": {
                    "0": [
                        {
                            "state": "STARTED",
                            "primary": true,
                            "node": "q--ANfDnTKW2WS9pEBuLWQ",
                            "relocating_node": null,
                            "shard": 0,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "TtAx_PMwRCmanPR7XddWmg"
                            }
                        },
                        {
                            "state": "STARTED",
                            "primary": false,
                            "node": "EizpW8QWRty_T1nJpr-dNQ",
                            "relocating_node": null,
                            "shard": 0,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "QddiDZTHTuStDTIKSOIk5A"
                            }
                        }
                    ],
                    "1": [
                        {
                            "state": "STARTED",
                            "primary": true,
                            "node": "NRqCLTmhTLuSxzlWcTae3A",
                            "relocating_node": null,
                            "shard": 1,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "IzFuExmARziQWcX8RlaZdg"
                            }
                        },
                        {
                            "state": "STARTED",
                            "primary": false,
                            "node": "EizpW8QWRty_T1nJpr-dNQ",
                            "relocating_node": null,
                            "shard": 1,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "XqIv4y1rQf6aL5C63Xsbhg"
                            }
                        }
                    ],
                    "2": [
                        {
                            "state": "STARTED",
                            "primary": false,
                            "node": "q--ANfDnTKW2WS9pEBuLWQ",
                            "relocating_node": null,
                            "shard": 2,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "XCAywOULRf66CR2xugkIpg"
                            }
                        },
                        {
                            "state": "STARTED",
                            "primary": true,
                            "node": "EizpW8QWRty_T1nJpr-dNQ",
                            "relocating_node": null,
                            "shard": 2,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "yNuj-Rw7QkC74opnoRQIqQ"
                            }
                        }
                    ],
                    "3": [
                        {
                            "state": "STARTED",
                            "primary": true,
                            "node": "q--ANfDnTKW2WS9pEBuLWQ",
                            "relocating_node": null,
                            "shard": 3,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "foOkK0oWTAaFTg-M41sMgQ"
                            }
                        },
                        {
                            "state": "STARTED",
                            "primary": false,
                            "node": "NRqCLTmhTLuSxzlWcTae3A",
                            "relocating_node": null,
                            "shard": 3,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "MdjjvB9KTfu4gs_skXDyXg"
                            }
                        }
                    ],
                    "4": [
                        {
                            "state": "STARTED",
                            "primary": false,
                            "node": "q--ANfDnTKW2WS9pEBuLWQ",
                            "relocating_node": null,
                            "shard": 4,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "exBumbxRT6KY7LVmGOSIZA"
                            }
                        },
                        {
                            "state": "STARTED",
                            "primary": true,
                            "node": "NRqCLTmhTLuSxzlWcTae3A",
                            "relocating_node": null,
                            "shard": 4,
                            "index": "shakespeare",
                            "allocation_id": {
                                "id": "pUhEb1k5TC24EKD-OjS7Iw"
                            }
                        }
                    ]
                }
            }
        }
    }
}
`
)

func newPod(name, namespace string) corev1.Pod {
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return pod
}

func Test_defaultDriver_attemptPodsDeletion(t *testing.T) {
	var clusterState esclient.ClusterState
	b := []byte(ClusterStateSample)
	err := json.Unmarshal(b, &clusterState)
	if err != nil {
		t.Error(err)
	}
	pod1 := newPod("elasticsearch-sample-es-79gc6p57rs", "default")
	pod2 := newPod("elasticsearch-sample-es-fnsgkkdl85", "default")
	pod3 := newPod("elasticsearch-sample-es-jfpqbt2s4q", "default")
	pod4 := newPod("elasticsearch-sample-es-nope", "default")

	expectedResult1 := reconcile.Results{}
	expectedResult1.WithResult(defaultRequeue).WithResult(defaultRequeue)

	expectedEmptyResult := reconcile.Results{}
	expectedEmptyResult.WithResult(k8sreconcile.Result{})

	nn := types.NamespacedName{
		Namespace: "default",
		Name:      "elasticsearch-sample",
	}

	type fields struct {
		Options Options
	}

	type args struct {
		ToDelete       []corev1.Pod
		reconcileState *reconcile.State
		resourcesState *reconcile.ResourcesState
		observedState  observer.State
		results        *reconcile.Results
		esClient       *esclient.Client
		namespacedName types.NamespacedName
	}

	type want struct {
		results              *reconcile.Results
		fulfilledExpectation bool
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    want
	}{
		{
			name: "Do not delete a pod with migrating data",
			args: args{
				namespacedName: nn,
				ToDelete:       []corev1.Pod{pod1, pod2},
				resourcesState: &reconcile.ResourcesState{
					CurrentPods: []corev1.Pod{pod1, pod2, pod3},
				},
				observedState: observer.State{
					ClusterState: &clusterState,
				},
				reconcileState: &reconcile.State{},
				results:        &reconcile.Results{},
			},
			fields: fields{
				Options: Options{
					PodsExpectations: reconciler.NewExpectations(),
				},
			},
			wantErr: false,
			want: want{
				results:              &expectedResult1,
				fulfilledExpectation: true, // pod deletion is delayed, do not expect anything
			},
		},
		{
			name: "Delete a pod with no data",
			args: args{
				namespacedName: nn,
				ToDelete:       []corev1.Pod{pod4},
				resourcesState: &reconcile.ResourcesState{
					CurrentPods: []corev1.Pod{pod1, pod2, pod3, pod4},
				},
				observedState: observer.State{
					ClusterState: &clusterState,
				},
				reconcileState: &reconcile.State{},
				results:        &reconcile.Results{},
			},
			fields: fields{
				Options: Options{
					PodsExpectations: reconciler.NewExpectations(),
					Client:           k8s.WrapClient(fake.NewFakeClient()),
				},
			},
			wantErr: false,
			want: want{
				results:              &expectedEmptyResult,
				fulfilledExpectation: false, // pod4 is expected to be deleted
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDriver{
				Options: tt.fields.Options,
			}
			if err := d.attemptPodsDeletion(
				tt.args.ToDelete, tt.args.reconcileState, tt.args.resourcesState,
				tt.args.observedState, tt.args.results, tt.args.esClient, tt.args.namespacedName); (err != nil) != tt.wantErr {
				t.Errorf("defaultDriver.attemptPodsDeletion() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.EqualValues(t, tt.want.results, tt.args.results)
			assert.EqualValues(t, tt.want.fulfilledExpectation, tt.fields.Options.PodsExpectations.Fulfilled(nn))
		})
	}
}