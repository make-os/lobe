module gitlab.com/makeos/mosdef

go 1.13

replace (
	github.com/bitfield/script => github.com/ncodes/script v0.14.1
	github.com/cbroglie/mustache => github.com/ncodes/mustache v1.0.2-0.20200429192435-945fed20e1e2
	github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.3.5-0.20190526074819-1df300866540
	github.com/golangci/golangci-lint => github.com/golangci/golangci-lint v1.18.0
	github.com/tendermint/tendermint => github.com/ncodes/tendermint v0.32.7-0.20200119162731-39690ff2d37e
)

require (
	github.com/AlecAivazis/survey/v2 v2.0.7
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496
	github.com/bitfield/script v0.14.1
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/c-bata/go-prompt v0.2.3
	github.com/cbroglie/mustache v1.0.1
	github.com/cenkalti/backoff/v4 v4.0.2
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgraph-io/badger v1.6.1
	github.com/dgryski/go-farm v0.0.0-20191112170834-c2139c5d712b // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/ellcrys/go-prompt v1.2.1
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fatih/color v1.7.0
	github.com/fatih/structs v1.1.0
	github.com/gen2brain/beeep v0.0.0-20200420150314-13046a26d502
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/go-stack/stack v1.8.0
	github.com/gobuffalo/packr v1.30.1
	github.com/gogo/protobuf v1.3.1
	github.com/gohugoio/hugo v0.69.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.4.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/rpc v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hokaccha/go-prettyjson v0.0.0-20190818114111-108c894c2c0e // indirect
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/imdario/mergo v0.3.9
	github.com/imroc/req v0.3.0
	github.com/ipfs/go-bitswap v0.2.15
	github.com/ipfs/go-block-format v0.0.2
	github.com/ipfs/go-cid v0.0.6
	github.com/ipfs/go-ds-badger v0.2.3
	github.com/ipfs/go-graphsync v0.0.5 // indirect
	github.com/ipfs/go-ipfs-blockstore v1.0.0
	github.com/ipfs/go-ipfs-exchange-interface v0.0.1
	github.com/ipfs/go-ipns v0.0.2
	github.com/ipfs/go-log/v2 v2.1.1 // indirect
	github.com/ipfs/go-todocounter v0.0.2 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/lestrrat-go/file-rotatelogs v2.2.0+incompatible
	github.com/lestrrat-go/strftime v0.0.0-20190725011945-5c849dd2c51d // indirect
	github.com/libp2p/go-libp2p v0.8.3
	github.com/libp2p/go-libp2p-core v0.5.6
	github.com/libp2p/go-libp2p-kad-dht v0.8.1
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-libp2p-routing v0.1.0
	github.com/libp2p/go-msgio v0.0.4
	github.com/libp2p/go-sockaddr v0.1.0 // indirect
	github.com/manifoldco/promptui v0.7.0 // indirect
	github.com/mattn/go-tty v0.0.0-20190424173100-523744f04859 // indirect
	github.com/mingrammer/commonregex v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/multiformats/go-multihash v0.0.13
	github.com/ncodes/go-prettyjson v0.0.0-20180528130907-d229c224a219
	github.com/ncodes/termbox-go v0.0.0-20200418040025-38ba6e5628f1
	github.com/neurosnap/sentences v1.0.6 // indirect
	github.com/olebedev/emitter v0.0.0-20190110104742-e8d1457e6aee
	github.com/olekukonko/tablewriter v0.0.4
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
	github.com/pkg/errors v0.9.1
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942 // indirect
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robertkrimen/otto v0.0.0-20180617131154-15f95af6e78d
	github.com/shopspring/decimal v0.0.0-20190905144223-a36b5d85f337
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.1
	github.com/stretchr/objx v0.2.0
	github.com/stumble/gorocksdb v0.0.3 // indirect
	github.com/tebeka/strftime v0.1.3 // indirect
	github.com/tendermint/iavl v0.12.4
	github.com/tendermint/tendermint v0.32.6
	github.com/tendermint/tm-db v0.2.0
	github.com/thoas/go-funk v0.4.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	github.com/vmihailenco/msgpack/v4 v4.3.11
	github.com/ztrue/tracerr v0.3.0 // indirect
	go.dedis.ch/kyber/v3 v3.0.11
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	golang.org/x/net v0.0.0-20200528225125-3c3fba18258b // indirect
	golang.org/x/sys v0.0.0-20200523222454-059865788121 // indirect
	golang.org/x/tools v0.0.0-20191216052735-49a3e744a425
	gonum.org/v1/gonum v0.7.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	gopkg.in/jdkato/prose.v2 v2.0.0-20190814032740-822d591a158c
	gopkg.in/neurosnap/sentences.v1 v1.0.6 // indirect
	gopkg.in/oleiade/lane.v1 v1.0.0
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.3.0
)
