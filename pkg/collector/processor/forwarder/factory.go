// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package forwarder

import (
	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/confengine"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/define"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/collector/processor"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils/logger"
)

func init() {
	processor.Register(define.ProcessorForwarder, NewFactory)
}

func NewFactory(conf map[string]interface{}, customized []processor.SubConfigProcessor) (processor.Processor, error) {
	return newFactory(conf, customized)
}

func newFactory(conf map[string]interface{}, customized []processor.SubConfigProcessor) (*forwarder, error) {
	clients := confengine.NewTierConfig()
	c := &Config{}
	if err := mapstructure.Decode(conf, c); err != nil {
		return nil, err
	}
	clients.SetGlobal(NewClient(*c))

	for _, custom := range customized {
		cfg := &Config{}
		if err := mapstructure.Decode(custom.Config.Config, cfg); err != nil {
			logger.Errorf("failed to decode config: %v", err)
			continue
		}
		clients.Set(custom.Token, custom.Type, custom.ID, *cfg)
	}

	return &forwarder{
		CommonProcessor: processor.NewCommonProcessor(conf, customized),
		clients:         clients,
	}, nil
}

type forwarder struct {
	processor.CommonProcessor
	clients *confengine.TierConfig // type: *Client
}

func (p forwarder) Name() string {
	return define.ProcessorForwarder
}

func (p forwarder) Clean() {
	for _, obj := range p.clients.All() {
		if err := obj.(*Client).Stop(); err != nil {
			logger.Errorf("failed to stop client, err: %v", err)
		}
	}
}

func (p forwarder) IsDerived() bool {
	return false
}

func (p forwarder) IsPreCheck() bool {
	return false
}

func (p forwarder) Process(record *define.Record) (*define.Record, error) {
	var err error
	client := p.clients.GetByToken(record.Token.Original).(*Client)
	switch record.RecordType {
	case define.RecordTraces:
		pdTraces := record.Data.(ptrace.Traces)
		err = client.ForwardTraces(pdTraces)
	}

	if err != nil {
		return nil, err
	}
	return nil, define.ErrEndOfPipeline
}
