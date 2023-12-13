// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package fta

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/config"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/define"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/etl"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/logging"
	template "github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/template/etl"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/transfer/utils"
)

const (
	// IngestTimeField 事件数据采集时间字段
	IngestTimeField = "bk_ingest_time"
	// CleanTimeField 事件数据清洗时间字段
	CleanTimeField = "bk_clean_time"
	// TagsField 事件标签字段
	TagsField = "tags"
	// DefaultEventIDField 默认事件ID字段
	DefaultEventIDField = "__bk_event_id__"
	// EventIDField 事件ID字段
	EventIDField = "event_id"
	// AlertNameField 告警名称字段
	AlertNameField = "alert_name"
	// DefaultPluginIDField 默认插件ID字段
	DefaultPluginIDField = "bk_plugin_id"
	// PluginIDField 插件ID字段
	PluginIDField = "plugin_id"
	// CleanConfigField 清洗配置字段
	CleanConfigField = "fta_clean_configs"
)

// Alert 告警名称匹配规则
type Alert struct {
	Trigger `mapstructure:",squash"`

	Name string `json:"name" mapstructure:"name"`
}

// CleanConfig 清洗配置
type CleanConfig struct {
	Alerts         []*Alert `json:"alert_config" mapstructure:"alert_config"`
	Normalizations []*struct {
		Field string `json:"field" mapstructure:"field"`
		Expr  string `json:"expr" mapstructure:"expr"`
	} `json:"normalization_config" mapstructure:"normalization_config"`
	Trigger `mapstructure:",squash"`
}

// Init 清洗配置初始化
func (c *CleanConfig) Init() error {
	for _, alert := range c.Alerts {
		err := alert.Init()
		if err != nil {
			return errors.WithMessagef(err, "alert init error for config->(%+v)", alert)
		}
	}
	return c.Trigger.Init()
}

// convertToExprMap 将字段提取配置转换为map格式
func convertToExprMap(c interface{}) (map[string]string, error) {
	var fields []struct {
		Field string `json:"field"`
		Expr  string `json:"expr"`
	}

	err := mapstructure.Decode(c, &fields)
	if err != nil {
		return nil, errors.WithMessagef(err, "decode expr config failed")
	}

	fieldExpr := make(map[string]string)
	for _, cfg := range fields {
		if cfg.Expr == "" {
			continue
		}
		fieldExpr[cfg.Field] = cfg.Expr
	}
	return fieldExpr, nil
}

// NewAlertFTAProcessor 创建FTA告警处理器
func NewAlertFTAProcessor(ctx context.Context, name string) (*template.RecordProcessor, error) {
	pipeConfig := config.PipelineConfigFromContext(ctx)
	helper := utils.NewMapHelper(pipeConfig.Option)
	configs, _ := helper.Get(CleanConfigField)

	// 清洗配置
	var cleanConfigs []*CleanConfig
	err := mapstructure.Decode(configs, &cleanConfigs)
	if err != nil {
		logging.Errorf("decode fta clean config failed: %+v", err)
	}
	for _, cleanConfig := range cleanConfigs {
		err := cleanConfig.Init()
		if err != nil {
			logging.Errorf("init clean config failed: %+v", err)
		}
	}

	// 默认告警名称配置
	var defaultAlerts []*Alert
	alertsConfig, ok := helper.Get(config.PipelineConfigOptFTAAlertsKey)
	if ok {
		err := mapstructure.Decode(alertsConfig, &defaultAlerts)
		if err != nil {
			return nil, errors.Errorf("decode fta alerts config failed: %+v", err)
		}
	}
	for _, alert := range defaultAlerts {
		err := alert.Init()
		if err != nil {
			logging.Errorf("init alert config failed: %+v", err)
		}
	}

	// 默认字段表达式配置
	fieldsCfg, _ := helper.Get(config.PipelineConfigOptFTAFieldMappingKey)
	defaultExprMap, err := convertToExprMap(fieldsCfg)
	if err != nil {
		logging.Errorf("convert to expr map failed: %+v", err)
	}

	decoder := etl.NewPayloadDecoder()
	return template.NewRecordProcessorWithDecoderFn(
		name, config.PipelineConfigFromContext(ctx),
		etl.NewFunctionalRecord("", func(from etl.Container, to etl.Container) error {
			data := etl.ContainerToMap(from)

			var alerts []*Alert
			var exprMap map[string]string

			// 判断是否满足匹配规则，如果不满足，则使用默认清洗配置
			var matchedCleanConfig *CleanConfig
			for _, cleanConfig := range cleanConfigs {
				if cleanConfig.IsMatch(data) {
					matchedCleanConfig = cleanConfig
					break
				}
			}
			if matchedCleanConfig != nil {
				alerts = matchedCleanConfig.Alerts
				exprMap, _ = convertToExprMap(matchedCleanConfig.Normalizations)
			}
			if alerts == nil {
				alerts = defaultAlerts
			}
			if exprMap == nil {
				exprMap = defaultExprMap
			}

			// 默认字段处理
			IngestTime, _ := from.Get(IngestTimeField)
			stamp, err := etl.TransformAutoTimeStamp(IngestTime)
			if err != nil {
				return err
			}
			_ = to.Put(IngestTimeField, stamp)

			PluginID, err := from.Get(DefaultPluginIDField)
			if err != nil {
				return err
			}
			_ = to.Put(PluginIDField, PluginID)

			newTimeStamp, _ := etl.TransformAutoTimeStamp(time.Now().UTC())
			_ = to.Put(CleanTimeField, newTimeStamp)

			// 按照配置的字段表达式，提取字段
			rt := config.ResultTableConfigFromContext(ctx)
			err = rt.VisitUserSpecifiedFields(func(config *config.MetaFieldConfig) error {
				// tags后面再处理
				if config.FieldName == TagsField {
					return nil
				}

				// 读取字段表达式
				expr, ok := exprMap[config.FieldName]
				if !ok {
					return nil
				}
				compiledExpr, err := utils.CompileJMESPathCustom(expr)
				if err != nil {
					return errors.WithMessagef(err, "compile expr %s failed", expr)
				}

				// 提取字段
				field, err := compiledExpr.Search(data)
				if err != nil {
					return errors.WithMessagef(err, "search expr %s failed", expr)
				}

				// 字段类型转换
				fieldTypeTransformFn := etl.NewTransformByField(config)
				field, err = fieldTypeTransformFn(field)

				_ = to.Put(config.FieldName, field)
				return nil
			})

			// 告警名称匹配
			for _, alert := range alerts {
				// 对满足匹配规则的数据，设置告警名称
				if alert.IsMatch(data) {
					_ = to.Put(AlertNameField, alert.Name)
					logging.Debugf("fta alert name matched->(%s) data->(%+v)", alert.Name, data)
					break
				}
			}

			// 事件ID处理
			eventID, _ := to.Get(EventIDField)
			if eventID == nil || eventID == "" {
				eventID, err = from.Get(DefaultEventIDField)
				if err != nil {
					return errors.Errorf("get key->(%s) from data failed: %+v, origin data: %+v",
						DefaultEventIDField, err, from)
				}
				logging.Debugf("using event_id->(%s) from origin data: %+v", eventID, from)
				_ = to.Put(EventIDField, eventID)
			}

			// tags处理
			if tagExpr, ok := exprMap[TagsField]; ok {
				// 提取tags字段
				compiledExpr, err := utils.CompileJMESPathCustom(tagExpr)
				if err != nil {
					return errors.WithMessagef(err, "compile expr %s failed", tagExpr)
				}
				tags, err := compiledExpr.Search(data)
				if err != nil {
					return nil
				}

				// 转换为统一格式 [{"key": "a", "value": "b"}]
				var tagsList []map[string]interface{}
				switch t := tags.(type) {
				case map[string]interface{}:
					// 针对 tags 为 {"a": "b"} 格式的转换
					for key, value := range t {
						tagsList = append(tagsList, map[string]interface{}{
							"key":   key,
							"value": value,
						})
					}
				case []interface{}:
					// 针对 tags 为 [{"key": "a", "value": "b"}] 的转换
					for _, item := range t {
						mapItem, ok := item.(map[string]interface{})
						if !ok {
							continue
						}
						key := mapItem["key"]
						value := mapItem["value"]
						if key == nil || value == nil {
							continue
						}
						tagsList = append(tagsList, map[string]interface{}{
							"key":   key,
							"value": value,
						})
					}
				}
				_ = to.Put(TagsField, tagsList)
			}

			return nil
		}), decoder.Decode,
	), nil
}

func init() {
	define.RegisterDataProcessor("fta-alert", func(ctx context.Context, name string) (processor define.DataProcessor, e error) {
		pipeConfig := config.PipelineConfigFromContext(ctx)
		if pipeConfig == nil {
			return nil, errors.Wrapf(define.ErrOperationForbidden, "pipeline config is empty")
		}
		return NewAlertFTAProcessor(ctx, pipeConfig.FormatName(name))
	})
}
