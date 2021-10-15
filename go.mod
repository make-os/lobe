module github.com/make-os/kit

go 1.13

replace (
	github.com/bitfield/script => github.com/ncodes/script v0.14.1
	github.com/btcsuite/btcutil => github.com/make-os/btcutil v1.0.3-0.20201208011646-272219d09635
	github.com/imdario/mergo => github.com/ncodes/mergo v0.3.10-0.20200627182710-b10b58df675a
	github.com/tendermint/tendermint => github.com/make-os/tendermint v0.34.0-rc4.0.20201212174221-93c4161d9329
)

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/AlekSi/pointer v1.1.0
	github.com/Azure/azure-amqp-common-go v1.1.4 // indirect
	github.com/Azure/go-autorest v11.1.2+incompatible // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210920160938-87db9fbc61c7 // indirect
	github.com/araddon/dateparse v0.0.0-20210207001429-0eec95c9db7e
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/bitfield/script v0.14.1
	github.com/briandowns/spinner v1.11.1
	github.com/btcsuite/btcutil v1.0.2
	github.com/c-bata/go-prompt v0.2.3
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/coreos/go-semver v0.3.0
	github.com/cosmos/iavl v0.15.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dgraph-io/badger/v2 v2.2007.2
	github.com/dgraph-io/ristretto v0.0.4-0.20200906165740-41ebdbffecfd // indirect
	github.com/dgryski/go-farm v0.0.0-20191112170834-c2139c5d712b // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/emirpasic/gods v1.12.0
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fatih/color v1.7.0
	github.com/fatih/structs v1.1.0
	github.com/gen2brain/beeep v0.0.0-20200526185328-e9c15c258e28
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-openapi/strfmt v0.19.11 // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/gohugoio/hugo v0.88.1
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.6
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/rpc v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-ds-badger2 v0.1.0
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/kyokomi/emoji v2.2.4+incompatible // indirect
	github.com/lestrrat-go/file-rotatelogs v2.2.0+incompatible
	github.com/lestrrat-go/strftime v0.0.0-20190725011945-5c849dd2c51d // indirect
	github.com/libp2p/go-libp2p v0.12.0
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/libp2p/go-libp2p-kad-dht v0.11.1
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-sockaddr v0.1.0 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/markbates/inflect v1.0.4 // indirect
	github.com/mattn/go-tty v0.0.0-20190424173100-523744f04859 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/multiformats/go-multihash v0.0.14
	github.com/ncodes/go-prettyjson v0.0.1
	github.com/nicksnyder/go-i18n/v2 v2.1.1 // indirect
	github.com/olebedev/emitter v0.0.0-20190110104742-e8d1457e6aee
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.5.0
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942 // indirect
	github.com/prometheus/common v0.14.0
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robertkrimen/otto v0.0.0-20180617131154-15f95af6e78d
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shopspring/decimal v0.0.0-20190905144223-a36b5d85f337
	github.com/sirupsen/logrus v1.7.0
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/objx v0.2.0
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942
	github.com/tebeka/strftime v0.1.3 // indirect
	github.com/tendermint/tendermint v0.34.0
	github.com/tendermint/tm-db v0.6.3
	github.com/thoas/go-funk v0.4.0
	github.com/tidwall/gjson v1.7.4
	github.com/uber-go/atomic v1.3.2 // indirect
	github.com/uber/jaeger-client-go v2.15.0+incompatible // indirect
	github.com/uber/jaeger-lib v1.5.0 // indirect
	github.com/urfave/cli v1.22.4 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/vmihailenco/msgpack/v4 v4.3.11
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.dedis.ch/kyber/v3 v3.0.11
	go.mongodb.org/mongo-driver v1.7.3 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/net v0.0.0-20211011170408-caeb26a5c8c0 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/yaml.v2 v2.4.0
	pack.ag/amqp v0.11.0 // indirect
)
