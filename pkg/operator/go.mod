module github.com/TencentBlueKing/bkmonitor-datalink/pkg/operator

go 1.21

require (
	github.com/elastic/beats v7.1.1+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/goware/urlx v0.3.1
	github.com/kylelemons/godebug v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus-operator/prometheus-operator v0.59.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.59.1
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.59.1
	github.com/prometheus/client_golang v1.16.0
	github.com/prometheus/common v0.44.0
	// This version is replaced using replace directive below
	github.com/prometheus/prometheus v1.99.0
	github.com/spf13/cobra v1.4.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/automaxprocs v1.5.2
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.27.3
	k8s.io/apimachinery v0.27.3
	k8s.io/client-go v0.27.3
)

require (
	github.com/Tencent/bk-bcs/bcs-scenarios/kourse v0.0.0-20220914032224-06b1bd3358bc
	github.com/TencentBlueKing/bkmonitor-datalink/pkg/libgse v0.0.0-00010101000000-000000000000
	github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils v0.0.0-00010101000000-000000000000
	github.com/blang/semver/v4 v4.0.0
	github.com/cespare/xxhash/v2 v2.2.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-kit/log v0.2.1
	github.com/valyala/bytebufferpool v1.0.0
)

require (
	cloud.google.com/go/compute v1.22.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/Shopify/sarama v1.27.1 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go v1.44.302 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/containerd/cgroups v1.0.3 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/edsrzf/mmap-go v1.1.0 // indirect
	github.com/efficientgo/tools/core v0.0.0-20220225185207-fe763185946b // indirect
	github.com/elastic/ecs v1.6.0 // indirect
	github.com/elastic/go-lumber v0.1.0 // indirect
	github.com/elastic/go-seccomp-bpf v1.1.0 // indirect
	github.com/elastic/go-structform v0.0.7 // indirect
	github.com/elastic/go-sysinfo v1.4.0 // indirect
	github.com/elastic/go-ucfg v0.7.0 // indirect
	github.com/elastic/go-windows v1.0.0 // indirect
	github.com/elastic/gosigar v0.11.0 // indirect
	github.com/emicklei/go-restful/v3 v3.10.2 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.20.4 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.25.0 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/strfmt v0.21.7 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-openapi/validate v0.22.1 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grafana/regexp v0.0.0-20221122212121-6b5c0a4cb7fd // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/hashstructure v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/natefinch/npipe v0.0.0-20160621034901-c1b8fa8bdcce // indirect
	github.com/nightlyone/lockfile v1.0.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus-community/prom-label-proxy v0.5.0 // indirect
	github.com/prometheus/alertmanager v0.25.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.11.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.mongodb.org/mongo-driver v1.12.0 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.16.0 // indirect
	go.opentelemetry.io/otel/metric v1.16.0 // indirect
	go.opentelemetry.io/otel/sdk v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/goleak v1.2.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/exp v0.0.0-20221031165847-c99f073a8326 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/oauth2 v0.10.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230717213848-3f92550aa753 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230717213848-3f92550aa753 // indirect
	google.golang.org/grpc v1.56.2 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/jcmturner/aescts.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/dnsutils.v1 v1.0.1 // indirect
	gopkg.in/jcmturner/gokrb5.v7 v7.5.0 // indirect
	gopkg.in/jcmturner/rpc.v1 v1.1.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
	k8s.io/apiextensions-apiserver v0.25.0 // indirect
	k8s.io/component-base v0.25.0 // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230525220651-2546d827e515 // indirect
	k8s.io/kubernetes v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/utils v0.0.0-20230711102312-30195339c3c7 // indirect
	sigs.k8s.io/controller-runtime v0.12.3 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.3.0 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	github.com/TencentBlueKing/bkmonitor-datalink/pkg/libgse => ../libgse
	github.com/TencentBlueKing/bkmonitor-datalink/pkg/utils => ../utils
	github.com/elastic/beats v7.1.1+incompatible => github.com/TencentBlueKing/beats v7.1.17-bk+incompatible
	// A replace directive is needed for github.com/prometheus/prometheus to ensure running against the latest version of prometheus.
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v0.40.0
	k8s.io/api => k8s.io/api v0.25.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.25.0
	k8s.io/apiserver => k8s.io/apiserver v0.25.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.25.0
	k8s.io/client-go => k8s.io/client-go v0.25.0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.0
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.0
	k8s.io/code-generator => k8s.io/code-generator v0.25.0
	k8s.io/component-base => k8s.io/component-base v0.25.0
	k8s.io/component-helpers => k8s.io/component-helpers v0.25.0
	k8s.io/controller-manager => k8s.io/controller-manager v0.25.0
	k8s.io/cri-api => k8s.io/cri-api v0.25.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.25.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.25.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.25.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.25.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.25.0
	k8s.io/kubectl => k8s.io/kubectl v0.25.0
	k8s.io/kubelet => k8s.io/kubelet v0.25.0
	k8s.io/kubernetes => k8s.io/kubernetes v1.20.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.25.0
	k8s.io/metrics => k8s.io/metrics v0.25.0
	k8s.io/mount-utils => k8s.io/mount-utils v0.25.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.25.0
)
