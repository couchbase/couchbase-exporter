//  Copyright (c) 2019 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package collectors

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/couchbase/couchbase-exporter/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	rebalanceSuccess = "rebalance_success"
	notFound         = "node not found"
)

var (
	ErrNotFound = fmt.Errorf(notFound)
)

func strToFloatArr(floatsStr string) []float64 {
	floatsStrArr := strings.Split(floatsStr, " ")

	var floatsArr []float64

	for _, f := range floatsStrArr {
		i, err := strconv.ParseFloat(f, 64)
		if err == nil {
			floatsArr = append(floatsArr, i)
		}
	}

	return floatsArr
}

func setGaugeVec(vec prometheus.GaugeVec, stats []float64, labelValues ...string) {
	if len(stats) > 0 {
		vec.WithLabelValues(labelValues...).Set(stats[len(stats)-1])
	}
}

func getClusterBalancedStatus(c util.CbClient) (bool, error) {
	node, err := c.Nodes()
	if err != nil {
		return false, fmt.Errorf("unable to retrieve nodes, %w", err)
	}

	return node.Counters[rebalanceSuccess] > 0 || (node.Balanced && node.RebalanceStatus == "none"), nil
}

func getCurrentNode(c util.CbClient) (string, error) {
	nodes, err := c.Nodes()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve nodes: %w", err)
	}

	for _, node := range nodes.Nodes {
		if node.ThisNode { // "ThisNode" is a boolean value indicating that it is the current node
			return node.Hostname, nil // hostname seems to work? just don't use for single node setups
		}
	}

	return "", ErrNotFound
}

func getPerNodeBucketStats(client util.CbClient, bucketName, nodeName string) map[string]interface{} {
	url := getSpecificNodeBucketStatsURL(client, bucketName, nodeName)

	var bucketStats objects.PerNodeBucketStats
	err := client.Get(url, &bucketStats)

	if err != nil {
		log.Error("unable to GET PerNodeBucketStats %s", err)
	}

	return bucketStats.Op.Samples
}

// /pools/default/buckets/<bucket-name>/nodes/<node-name>/stats.
func getSpecificNodeBucketStatsURL(client util.CbClient, bucket, node string) string {
	servers, err := client.Servers(bucket)
	if err != nil {
		log.Error("unable to retrieve Servers %s", err)
	}

	correctURI := ""

	for _, server := range servers.Servers {
		if server.Hostname == node {
			correctURI = server.Stats["uri"]
		}
	}

	return correctURI
}

func collectPerNodeBucketMetrics(client util.CbClient, node string, refreshTime int, config *objects.CollectorConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	clusterName, err := client.ClusterName()
	if err != nil {
		log.Error("%s", err)
		return
	}

	outerErr := util.Retry(ctx, 10*time.Second, 10, func() (bool, error) {
		rebalanced, err := getClusterBalancedStatus(client)
		if err != nil {
			log.Error("Unable to get rebalance status %s", err)
		}

		if !rebalanced {
			log.Info("Waiting for Rebalance... retrying...")
			return false, err
		}
		go func() {
			metrics := make(map[string]*prometheus.GaugeVec)
			for {
				buckets, err := client.Buckets()
				if err != nil {
					log.Error("Unable to get buckets %s", err)
				}

				for _, bucket := range buckets {
					log.Debug("Collecting per-node bucket stats, node=%s, bucket=%s", node, bucket.Name)

					samples := getPerNodeBucketStats(client, bucket.Name, node)

					for _, value := range config.Metrics {
						if !value.Enabled {
							continue
						}
						if mt, ok := metrics[value.Name]; ok {
							setGaugeVec(*mt, strToFloatArr(fmt.Sprint(samples[value.Name])), bucket.Name, node, clusterName)
						} else {
							mt := value.GetPrometheusGaugeVec(config.Namespace, config.Subsystem)
							metrics[value.Name] = mt
							setGaugeVec(*mt, strToFloatArr(fmt.Sprint(samples[value.Name])), bucket.Name, node, clusterName)
						}
					}
				}
				time.Sleep(time.Second * time.Duration(refreshTime))
			}
		}()
		log.Info("Per Node Bucket Stats Go Thread executed successfully")
		return true, nil
	})
	if outerErr != nil {
		log.Error("Getting Per Node Bucket Stats failed %s", outerErr)
	}
}

func RunPerNodeBucketStatsCollection(client util.CbClient, refreshTime int, config *objects.CollectorConfig) {
	log.Info("Initial Collection of Node Stats")

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	outerErr := util.Retry(ctx, 10*time.Second, 10, func() (bool, error) {
		currNode, err := getCurrentNode(client)
		if err != nil {
			log.Error("could not get current node, will retry. %s", err)
			return false, &util.RetryError{
				N: 0,
				E: err,
			}
		}

		collectPerNodeBucketMetrics(client, currNode, refreshTime, config)

		return true, nil
	})
	if outerErr != nil {
		log.Error("getting default stats failed: %s", outerErr)
	}
}
