// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package dbfilter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/collector/semconv/v1.8.0"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/define"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/internal/generator"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/internal/mapstructure"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/internal/testkits"
)

func TestFactory(t *testing.T) {
	content := `
processor:
  - name: "db_filter/common"
    config:
      slow_query:
        destination: "db.is_slow"
        rules:
        - match: "mysql"
          threshold: 1s
        - match: "redis"
          threshold: 2s
`
	psc := testkits.MustLoadProcessorConfigs(content)
	obj, err := NewFactory(psc[0].Config, nil)
	factory := obj.(*dbFilter)
	assert.NoError(t, err)
	assert.Equal(t, psc[0].Config, factory.MainConfig())

	c := &Config{}
	err = mapstructure.Decode(psc[0].Config, c)
	c.Setup()

	assert.NoError(t, err)
	assert.Equal(t, *c, factory.configs.Get("", "", "").(Config))

	assert.Equal(t, define.ProcessorDbFilter, factory.Name())
	assert.False(t, factory.IsDerived())
	assert.False(t, factory.IsPreCheck())

	duration, ok := c.GetSlowQueryConf("mysql")
	assert.True(t, ok)
	assert.Equal(t, time.Second, duration)
}

func TestSlowMySqlQuery(t *testing.T) {
	content := `
processor:
  - name: "db_filter/common"
    config:
      slow_query:
        destination: "db.is_slow"
        rules:
        - match: "mysql"
          threshold: 1s
`
	factory := testkits.MustCreateFactory(content, NewFactory)

	g := generator.NewTracesGenerator(define.TracesOptions{
		GeneratorOptions: define.GeneratorOptions{
			Attributes: map[string]string{"db.system": "mysql"},
		},
		SpanCount: 1,
	})

	t.Run("mysql slow query", func(t *testing.T) {
		data := g.Generate()
		span := data.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0)

		// duration: 2s
		span.SetEndTimestamp(pcommon.Timestamp(3 * time.Second))
		span.SetStartTimestamp(pcommon.Timestamp(1 * time.Second))

		record := &define.Record{
			RecordType: define.RecordTraces,
			Data:       data,
		}
		_, err := factory.Process(record)
		assert.NoError(t, err)

		span = testkits.FirstSpan(record.Data.(ptrace.Traces))
		testkits.AssertAttrsFoundIntVal(t, span.Attributes(), "db.is_slow", 1)
	})

	t.Run("mysql normal query", func(t *testing.T) {
		data := g.Generate()
		span := data.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0)

		// duration: 0.5s
		span.SetEndTimestamp(pcommon.Timestamp(1500 * time.Millisecond))
		span.SetStartTimestamp(pcommon.Timestamp(1 * time.Second))

		record := &define.Record{
			RecordType: define.RecordTraces,
			Data:       data,
		}
		_, err := factory.Process(record)
		assert.NoError(t, err)

		span = testkits.FirstSpan(record.Data.(ptrace.Traces))
		testkits.AssertAttrsFoundIntVal(t, span.Attributes(), "db.is_slow", 0)
	})
}

func TestSlowQueryDefault(t *testing.T) {
	content := `
processor:
  - name: "db_filter/common"
    config:
      slow_query:
        destination: "db.is_slow_or_else_name"
        rules:
         - match: ""
           threshold: 3s
`
	factory := testkits.MustCreateFactory(content, NewFactory)

	g := generator.NewTracesGenerator(define.TracesOptions{
		GeneratorOptions: define.GeneratorOptions{
			Attributes: map[string]string{"db.system": "elasticsearch"},
		},
		SpanCount: 1,
	})

	// 使用兜底规则
	t.Run("elasticsearch slow query(default)", func(t *testing.T) {
		data := g.Generate()
		span := data.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0)

		// duration: 9s
		span.SetEndTimestamp(pcommon.Timestamp(10 * time.Second))
		span.SetStartTimestamp(pcommon.Timestamp(1 * time.Second))

		record := &define.Record{
			RecordType: define.RecordTraces,
			Data:       data,
		}
		_, err := factory.Process(record)
		assert.NoError(t, err)

		span = testkits.FirstSpan(record.Data.(ptrace.Traces))
		testkits.AssertAttrsFoundIntVal(t, span.Attributes(), "db.is_slow_or_else_name", 1)
	})

	t.Run("not db system", func(t *testing.T) {
		data := g.Generate()
		span := data.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0)
		span.Attributes().Remove(semconv.AttributeDBSystem)

		record := &define.Record{
			RecordType: define.RecordTraces,
			Data:       data,
		}
		_, err := factory.Process(record)
		assert.NoError(t, err)

		span = testkits.FirstSpan(record.Data.(ptrace.Traces))
		testkits.AssertAttrsNotFound(t, span.Attributes(), "db.is_slow_or_else_name")
	})
}
