package k8s_kubelet

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testMetrics, _ = ioutil.ReadFile("testdata/metrics.txt")

func TestNew(t *testing.T) {
	job := New()

	assert.IsType(t, (*Kubelet)(nil), job)
	assert.Equal(t, defaultURL, job.URL)
	assert.Equal(t, defaultHTTPTimeout, job.Timeout.Duration)
}

func TestKubeProxy_Charts(t *testing.T) { assert.NotNil(t, New().Charts()) }

func TestKubeProxy_Cleanup(t *testing.T) { New().Cleanup() }

func TestKubeProxy_Init(t *testing.T) { assert.True(t, New().Init()) }

func TestKubeProxy_InitNG(t *testing.T) {
	job := New()
	job.URL = ""
	assert.False(t, job.Init())
}

func TestKubeProxy_Check(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testMetrics)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.True(t, job.Check())
}

func TestKubeProxy_CheckNG(t *testing.T) {
	job := New()
	job.URL = "http://127.0.0.1:38001/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestKubeProxy_Collect(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write(testMetrics)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	require.True(t, job.Check())

	//for k, v := range job.Collect() {
	//	fmt.Println(fmt.Sprintf("\"%s\": %d,", k, v))
	//}

	expected := map[string]int64{
		"kubelet_docker_operations_errors_remove_container":            4,
		"rest_client_requests_200":                                     177,
		"apiserver_audit_requests_rejected_total":                      0,
		"apiserver_storage_envelope_transformation_cache_misses_total": 0,
		"apiserver_storage_data_key_generation_bucket_20":              1,
		"kubelet_runtime_operations_container_status":                  90,
		"volume_manager_plugin_kubernetes.io/host-path_state_actual":   15,
		"rest_client_requests_403":                                     2,
		"rest_client_requests_409":                                     1,
		"kubelet_runtime_operations_image_status":                      25,
		"kubelet_runtime_operations_run_podsandbox":                    9,
		"kubelet_docker_operations_stop_container":                     23,
		"volume_manager_plugin_kubernetes.io/configmap_state_desired":  3,
		"volume_manager_plugin_kubernetes.io/host-path_state_desired":  15,
		"rest_client_requests_POST":                                    8,
		"apiserver_storage_data_key_generation_bucket_320":             1,
		"apiserver_storage_data_key_generation_bucket_+Inf":            1,
		"kubelet_runtime_operations_start_container":                   10,
		"kubelet_docker_operations_version":                            472,
		"kubelet_runtime_operations_list_containers":                   2586,
		"kubelet_runtime_operations_status":                            279,
		"kubelet_runtime_operations_errors_remove_container":           4,
		"kubelet_docker_operations_remove_container":                   23,
		"kubelet_docker_operations_inspect_container":                  223,
		"kubelet_pleg_relist_interval_09":                              1016820,
		"kubelet_runtime_operations_remove_container":                  14,
		"volume_manager_plugin_kubernetes.io/configmap_state_actual":   3,
		"kubelet_docker_operations_start_container":                    19,
		"rest_client_requests_GET":                                     37,
		"apiserver_storage_data_key_generation_bucket_80":              1,
		"apiserver_storage_data_key_generation_bucket_640":             1,
		"kubelet_pleg_relist_latency_099":                              31234,
		"kubelet_runtime_operations_exec_sync":                         138,
		"token_count":                                                  0,
		"rest_client_requests_201":                                     43,
		"kubelet_runtime_operations_list_images":                       195,
		"token_fail_count":                                             0,
		"apiserver_storage_data_key_generation_bucket_1280":            1,
		"apiserver_storage_data_key_generation_bucket_20480":           1,
		"kubelet_pleg_relist_interval_099":                             1032022,
		"kubelet_runtime_operations_stop_podsandbox":                   14,
		"rest_client_requests_<error>":                                 8,
		"kubelet_runtime_operations_create_container":                  10,
		"kubelet_docker_operations_create_container":                   19,
		"kubelet_docker_operations_info":                               2,
		"volume_manager_plugin_kubernetes.io/secret_state_actual":      4,
		"kubelet_runtime_operations_errors_container_status":           14,
		"rest_client_requests_PATCH":                                   177,
		"apiserver_storage_data_key_generation_bucket_40":              1,
		"apiserver_storage_data_key_generation_bucket_40960":           1,
		"kubelet_running_pod":                                          9,
		"kubelet_pleg_relist_interval_05":                              1013125,
		"apiserver_storage_data_key_generation_bucket_5120":            1,
		"kubelet_docker_operations_list_containers":                    5157,
		"kubelet_docker_operations_errors_inspect_container":           14,
		"apiserver_storage_data_key_generation_bucket_5":               6,
		"volume_manager_plugin_kubernetes.io/secret_state_desired":     4,
		"kubelet_docker_operations_inspect_image":                      110,
		"apiserver_storage_data_key_generation_failures_total":         0,
		"apiserver_storage_data_key_generation_bucket_2560":            1,
		"kubelet_node_config_error":                                    1,
		"kubelet_running_container":                                    9,
		"kubelet_runtime_operations_version":                           190,
		"apiserver_storage_data_key_generation_bucket_10240":           1,
		"kubelet_pleg_relist_latency_05":                               12741,
		"kubelet_runtime_operations_list_podsandbox":                   2562,
		"kubelet_docker_operations_list_images":                        195,
		"apiserver_storage_data_key_generation_bucket_10":              1,
		"kubelet_pleg_relist_latency_09":                               16211,
		"apiserver_storage_data_key_generation_bucket_160":             1,
		"kubelet_runtime_operations_podsandbox_status":                 77,
	}

	assert.Equal(t, expected, job.Collect())
}

func TestKubeProxy_InvalidData(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("hello and goodbye"))
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}

func TestKubeProxy_404(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}))
	defer ts.Close()

	job := New()
	job.URL = ts.URL + "/metrics"
	require.True(t, job.Init())
	assert.False(t, job.Check())
}
