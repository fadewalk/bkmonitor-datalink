// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package structured

import (
	"context"
	"testing"
	"time"

	goRedis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/offline-data-archive/metadata"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/offline-data-archive/policy/stores/shard"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/consul"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/log"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/mock"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/unify-query/redis"
	ir "github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/router/influxdb"
)

type m struct {
	shard []*shard.Shard
}

func (m *m) PublishShard(ctx context.Context, channelValue interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m *m) SubscribeShard(ctx context.Context) <-chan *goRedis.Message {
	//TODO implement me
	panic("implement me")
}

func (m *m) GetShardID(ctx context.Context, sd *shard.Shard) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *m) GetAllShards(ctx context.Context) map[string]*shard.Shard {
	//TODO implement me
	panic("implement me")
}

func (m *m) SetShard(ctx context.Context, k string, sd *shard.Shard) error {
	//TODO implement me
	panic("implement me")
}

func (m *m) GetShard(ctx context.Context, k string) (*shard.Shard, error) {
	//TODO implement me
	panic("implement me")
}

func (m *m) GetDistributedLock(ctx context.Context, key, val string, expiration time.Duration) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *m) RenewalLock(ctx context.Context, key string, renewalDuration time.Duration) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *m) GetPolicies(ctx context.Context, clusterName, tagRouter string) (map[string]*metadata.Policy, error) {
	//TODO implement me
	panic("implement me")
}

func (m *m) GetShards(ctx context.Context, clusterName, tagRouter, database string) (map[string]*shard.Shard, error) {
	//TODO implement me
	panic("implement me")
}

func (m *m) GetReadShardsByTimeRange(ctx context.Context, clusterName, tagRouter, database, retentionPolicy string, start int64, end int64) ([]*shard.Shard, error) {
	log.Infof(ctx, "check offline data archive query: %s %s %s %s %d %d", clusterName, tagRouter, database, retentionPolicy, start, end)
	var shards = make([]*shard.Shard, 0, len(m.shard))
	for _, sd := range m.shard {
		// 验证 meta 字段
		if sd.Meta.ClusterName != clusterName {
			continue
		}
		if sd.Meta.Database != database {
			continue
		}
		if sd.Meta.TagRouter != tagRouter {
			continue
		}
		if sd.Meta.RetentionPolicy != retentionPolicy {
			continue
		}
		if sd.Meta.TagRouter != tagRouter {
			continue
		}

		// 判断是否是过期的 shard，只有过期的 shard 才进行查询
		if sd.Spec.Expired.Unix() > time.Now().Unix() {
			continue
		}

		// 通过时间过滤
		if sd.Spec.Start.UnixNano() >= start && end < sd.Spec.End.UnixNano() {
			shards = append(shards, sd)
		}
	}
	return shards, nil
}

func TestQueryToMetric(t *testing.T) {
	ctx := context.Background()

	mock.SetRedisClient(ctx, "")

	testCases := map[string]struct {
		spaceUid      string
		tableID       string
		field         string
		referenceName string

		clusterName     string
		tagsKey         []string
		db              string
		measurement     string
		retentionPolicy string
		storageID       string
		vmRt            string

		tagRouter         string
		expectedStorageID string
		expired           time.Time

		start string
		end   string
	}{
		"offlineDataArchiveQuery": {
			spaceUid: "q_test", tableID: "pushgateway_bkmonitor_unify_query.__default__", field: "q_test", referenceName: "a",
			start: "0", end: "60",

			storageID: "2", clusterName: "cluster_internal", tagsKey: []string{"bk_biz_id"},
			db: "pushgateway_bkmonitor_unify_query", measurement: "unify_query_request_handler_total", retentionPolicy: "",
			tagRouter: "bk_biz_id==2", expired: time.Now().Add(-time.Minute),

			expectedStorageID: consul.OfflineDataArchive,
		},
		"notOfflineDataArchiveQuery": {
			spaceUid: "q_test", tableID: "pushgateway_bkmonitor_unify_query.__default__", field: "q_test", referenceName: "a",
			start: "0", end: "60",

			storageID: "2", clusterName: "cluster_internal", tagsKey: []string{"bk_biz_id"},
			db: "pushgateway_bkmonitor_unify_query", measurement: "unify_query_request_handler_total", retentionPolicy: "",
			tagRouter: "bk_biz_id==2", expired: time.Now().Add(time.Minute),

			expectedStorageID: "2",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mock.SetSpaceAndProxyMockData(
				ctx, tc.spaceUid, &redis.TsDB{
					TableID:         tc.tableID,
					Field:           []string{tc.field},
					MeasurementType: redis.BkSplitMeasurement,
				}, &ir.Proxy{
					MeasurementType: redis.BKTraditionalMeasurement,
					StorageID:       tc.storageID,
					ClusterName:     tc.clusterName,
					TagsKey:         tc.tagsKey,
					Db:              tc.db,
					Measurement:     tc.measurement,
					VmRt:            tc.vmRt,
				},
			)

			mockMd := &m{
				shard: []*shard.Shard{
					{
						Meta: shard.Meta{
							ClusterName:     tc.clusterName,
							Database:        tc.db,
							RetentionPolicy: tc.retentionPolicy,
							TagRouter:       tc.tagRouter,
						},
						Spec: shard.Spec{
							Start:   time.Unix(0, 0),
							End:     time.Unix(6000, 0),
							Expired: tc.expired,
						},
					},
				},
			}
			mock.SetOfflineDataArchiveMetadata(mockMd)

			query := &Query{
				TableID:       tc.tableID,
				FieldName:     tc.field,
				ReferenceName: tc.referenceName,
				Start:         tc.start,
				End:           tc.end,
				Conditions: Conditions{
					FieldList: []ConditionField{
						{
							DimensionName: "bk_biz_id",
							Operator:      "contains",
							Value:         []string{"2"},
						},
					},
				},
			}

			metric, err := query.ToQueryMetric(ctx, tc.spaceUid)
			assert.Nil(t, err)
			if len(metric.QueryList) > 0 {
				assert.Equal(t, tc.expectedStorageID, metric.QueryList[0].StorageID)
			} else {
				panic("query list length is 0")
			}
		})
	}

}
