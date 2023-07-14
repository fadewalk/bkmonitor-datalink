// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package gse

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/outputs"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/libgse/gse"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/libgse/monitoring"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/libgse/monitoring/report/bkpipe"
	"github.com/TencentBlueKing/bkmonitor-datalink/pkg/libgse/storage"
)

const maxSyncAgentInfoTimeout = 10 // unit: second

var (
	metricGseTaskPublishTotal  = "gse_publish_total"  // 按任务计算发送次数
	metricGseTaskPublishFailed = "gse_publish_failed" // 按任务计算发送失败次数
)

var (
	metricGseAgentInfoFailed = monitoring.NewInt("gse_agent_info_failed") // 获取Agent失败次数
	metricGseSendTotal       = monitoring.NewInt("gse_send_total")        // 发送给gse client的事件数

	metricGsePublishReceived = monitoring.NewInt("gse_publish_received") // publish：接收事件数
	metricGsePublishDropped  = monitoring.NewInt("gse_publish_dropped")  // publish：丢失的事件数（缺少dataid）
	metricGsePublishTotal    = monitoring.NewInt("gse_publish_total")    // publish：发送事件数
	metricGsePublishFailed   = monitoring.NewInt("gse_publish_failed")   // publish：发送失败数

	metricGseReportReceived  = monitoring.NewInt("gse_report_received")   // report：接收事件数
	metricGseReportSendTotal = monitoring.NewInt("gse_report_send_total") // report：发送事件数
	metricGseReportFailed    = monitoring.NewInt("gse_report_failed")     // report：发送失败数
)

func init() {
	outputs.RegisterType("gse", MakeGSE)
}

var MarshalFunc = json.Marshal

// Output : gse output, for libbeat output
type Output struct {
	cli *gse.GseClient
	aif *AgentInfoFetcher
}

// agentInfoLoader 全局 agentInfo 加载器
type agentInfoLoader struct {
	once  sync.Once
	fetch func() (gse.AgentInfo, error)
}

// Init 初始化 fetch 实现 只初始一次
func (ail *agentInfoLoader) Init(fn func() (gse.AgentInfo, error)) {
	ail.once.Do(func() {
		ail.fetch = fn
	})
}

// Fetch 获取 AgentInfo 如果 fetch 实例为空则返回空对象
func (ail *agentInfoLoader) Fetch() (gse.AgentInfo, error) {
	if ail.fetch != nil {
		return ail.fetch()
	}
	return gse.AgentInfo{}, errors.New("no agentInfo fetcher found")
}

var ail = &agentInfoLoader{}

func GetAgentInfo() (gse.AgentInfo, error) {
	return ail.Fetch()
}

// MakeGSE create a gse client
func MakeGSE(im outputs.IndexManager, beat beat.Info, stats outputs.Observer, cfg *common.Config) (outputs.Group, error) {
	c := defaultConfig
	err := cfg.Unpack(&c)
	if err != nil {
		logp.Err("unpack config error, %v", err)
		return outputs.Fail(err)
	}
	logp.Info("gse config: %+v", c)

	// create gse client
	cli, err := gse.NewGseClient(cfg)
	if err != nil {
		return outputs.Fail(err)
	}
	fetcher := NewAgentInfoFetcher(c, cli)
	ail.Init(fetcher.Fetch)
	output := &Output{
		cli: cli,
		aif: fetcher,
	}

	// start gse client
	err = output.cli.Start()
	if err != nil {
		logp.Err("init output failed, %v", err)
		return outputs.Fail(err)
	}
	logp.Info("start gse output")

	// wait to get agent info
	agentInfo, err := output.cli.GetAgentInfo()
	count := maxSyncAgentInfoTimeout
	for {
		if count <= 0 {
			return outputs.Fail(fmt.Errorf("get agent info timeout"))
		}
		if !agentInfo.IsEmpty() {
			break
		}
		count--
		// sleep 1s, then continue to get agent info
		time.Sleep(1 * time.Second)
		agentInfo, err = output.cli.GetAgentInfo()
	}

	// init bkmonitoring
	agentInfo, _ = output.aif.Fetch()
	bkpipe.InitSender(output, agentInfo)

	return outputs.Success(int(c.EventBufferMax), 0, output)
}

// MakeGSEWithoutCheckConn create a gse client without check connection
func MakeGSEWithoutCheckConn(im outputs.IndexManager, beat beat.Info, stats outputs.Observer, cfg *common.Config) (outputs.Group, error) {
	c := defaultConfig
	err := cfg.Unpack(&c)
	if err != nil {
		logp.Err("unpack config error, %v", err)
		return outputs.Fail(err)
	}
	logp.Info("gse config: %+v", c)

	// create gse client
	cli, err := gse.NewGseClient(cfg)
	if err != nil {
		return outputs.Fail(err)
	}
	fetcher := NewAgentInfoFetcher(c, cli)
	ail.Init(fetcher.Fetch)
	output := &Output{
		cli: cli,
		aif: fetcher,
	}

	go func() {
		// start gse client
		output.cli.StartWithoutCheckConn()
		logp.Info("start gse output")

		// wait to get agent info
		var agentInfo gse.AgentInfo
		for {
			agentInfo, _ = output.cli.GetAgentInfo()
			if !agentInfo.IsEmpty() {
				break
			}
			// sleep 1s, then continue to get agent info
			time.Sleep(1 * time.Second)
		}

		// init bkmonitoring
		agentInfo, _ = output.aif.Fetch()
		bkpipe.InitSender(output, agentInfo)
	}()

	return outputs.Success(int(c.EventBufferMax), 0, output)
}

// Publish implement output interface
func (c *Output) Publish(batch publisher.Batch) error {
	events := batch.Events()
	for i := range events {
		if events[i].Content.Fields == nil {
			metricGsePublishDropped.Add(1)
			continue
		}
		metricGsePublishReceived.Add(1)
		err := c.PublishEvent(&events[i])
		if err != nil {
			logp.Err("publish event failed: %v", err)
			metricGsePublishFailed.Add(1)
		} else {
			metricGsePublishTotal.Add(1)
		}
	}

	batch.ACK()
	return nil
}

// String returns the name of the output client
func (c *Output) String() string {
	return "gse"
}

// PublishEvent implement output interface
// data is event, must contain 'dataid' filed
// data will attach agent info, see publishEventAttachInfo
func (c *Output) PublishEvent(event *publisher.Event) error {
	// get dataid from event
	content := event.Content
	data := content.Fields
	val, err := data.GetValue("dataid")
	if err != nil {
		logp.Err("event lost dataid field, %v", err)
		return fmt.Errorf("event lost dataid")
	}

	dataid := c.getdataid(val)
	if dataid <= 0 {
		return fmt.Errorf("dataid %d <= 0", dataid)
	}

	if content.Meta != nil {
		data.Put("@meta", content.Meta)
	}

	if err := c.publishEventAttachInfo(dataid, data); err != nil {
		return err
	}

	return nil
}

// Close : close gse out put
func (c *Output) Close() error {
	logp.Err("gse output close")
	c.cli.Close()
	return nil
}

// publishEventAttachInfo attach agentinfo and gseindex
// will add bizid, cloudid, ip, gseindex
func (c *Output) publishEventAttachInfo(dataid int32, data common.MapStr) error {
	// 是否兼容原采集器输出
	isStandardFormat := true
	if _, ok := data["_time_"]; ok {
		isStandardFormat = false
	}

	// add gseindex
	var gseIndex uint64
	if ok, _ := data.HasKey("gseindex"); !ok {
		gseIndex = getGseIndex(dataid)
	}

	// add bizid, cloudid, ip
	info, _ := c.aif.Fetch()
	if info.IsEmpty() {
		metricGseAgentInfoFailed.Add(1)
		return fmt.Errorf("agent info is empty")
	}

	if isStandardFormat {
		data["bizid"] = info.Bizid
		// 按需补充 避免覆盖拨测 bk_biz_id 字段
		if _, ok := data["bk_biz_id"]; !ok {
			data["bk_biz_id"] = info.BKBizID
		}
		data["cloudid"] = info.Cloudid
		data["ip"] = info.IP
		data["hostname"] = info.Hostname
		data["gseindex"] = gseIndex
		data["bk_agent_id"] = info.BKAgentID
		data["bk_host_id"] = info.HostID
	} else {
		data["_bizid_"] = info.Bizid
		data["_cloudid_"] = info.Cloudid
		data["_server_"] = info.IP
		data["_gseindex_"] = gseIndex
		data["_agentid_"] = info.BKAgentID
		data["_hostid_"] = info.HostID
		data.Delete("dataid")
		data.Delete("time")
	}

	return c.reportCommonData(dataid, data)
}

func getGseIndex(dataid int32) uint64 {
	index := uint64(0)
	gseIndexKey := fmt.Sprintf("gseindex_%s", String(dataid))
	if indexStr, err := storage.Get(gseIndexKey); nil == err {
		if index, err = strconv.ParseUint(indexStr, 10, 64); nil != err {
			logp.Err("fail to get gseindex %v", err)
			index = 0
		}
	}
	index++
	storage.Set(gseIndexKey, fmt.Sprintf("%v", index), 0)
	return index
}

// Report implement interface for bkmonitor
func (c *Output) Report(dataid int32, data common.MapStr) error {
	if dataid <= 0 {
		return fmt.Errorf("dataid %d <= 0", dataid)
	}
	metricGseReportReceived.Add(1)
	err := c.reportCommonData(dataid, data)
	if err != nil {
		metricGseReportFailed.Add(1)
		return err
	}
	metricGseReportSendTotal.Add(1)
	return nil
}

// ReportRaw implement interface for monitor
// send op raw data, without attach anything
func (c *Output) ReportRaw(dataid int32, data interface{}) error {
	if dataid <= 0 {
		return fmt.Errorf("dataid %d <= 0", dataid)
	}

	buf, err := MarshalFunc(data)
	if err != nil {
		logp.Err("convert to json failed: %v", err)
		return err
	}

	logp.Debug("gse", "report data to %d", dataid)
	// report op data

	msg := gse.NewGseOpMsg(buf, dataid, 0, 0, 0)
	logp.Debug("gse", "report data : %s", string(buf))
	// TODO compatible op data bug fixed after agent D48
	// send every op data with new connection
	c.cli.SendWithNewConnection(msg)
	// c.cli.Send(msg)

	return nil
}

var sendHook func(float64)

func RegisterSendHook(f func(float64)) { sendHook = f }

// reportCommonData send common data
func (c *Output) reportCommonData(dataid int32, data common.MapStr) error {
	// change data to json format
	buf, err := MarshalFunc(data)
	if err != nil {
		monitoring.NewIntWithDataID(int(dataid), metricGseTaskPublishFailed).Add(1)
		logp.Err("json marshal failed, content: %+v, err: %+v", data, err)
		return err
	}
	if sendHook != nil {
		sendHook(float64(len(buf)))
	}

	// new dynamic msg
	msg := gse.NewGseDynamicMsg(buf, dataid, 0, 0)

	// send data
	c.cli.Send(msg)

	// 发包计数
	metricGseSendTotal.Add(1)
	monitoring.NewIntWithDataID(int(dataid), metricGseTaskPublishTotal).Add(1)

	return nil
}

func (c *Output) getdataid(dataID interface{}) int32 {
	switch dataID.(type) {
	case int, int8, int16, int32, int64:
		return int32(reflect.ValueOf(dataID).Int())
	case uint, uint8, uint16, uint32, uint64:
		return int32(reflect.ValueOf(dataID).Uint())
	case string:
		dataid, err := strconv.ParseInt(dataID.(string), 10, 32)
		if err != nil {
			logp.Err("can not get dataid, %s", dataID.(string))
			return -1
		}
		return int32(dataid)
	default:
		logp.Err("unexpected type %T for the dataid ", dataID)
		return 0
	}
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
