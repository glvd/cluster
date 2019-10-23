# IPFS Cluster Changelog

### v0.11.0 - 2019-09-13

#### Summary

IPFS Cluster v0.11.0 is the biggest release in the project's history. Its main
feature is the introduction of the new CRDT "consensus" component. Leveraging
Pubsub, Bitswap and the DHT and using CRDTs, cluster peers can track the
global pinset without needing to be online or worrying about the rest of the
peers as it happens with the original Raft approach.

The CRDT component brings a lots of features around it, like RPC
authorization, which effectively lets cluster peers run in clusters where only
a trusted subset of nodes can access peer endpoints and made modifications to
the pinsets.

We have additionally taken lots of steps to improve configuration management
of peers, separating the peer identity from the rest of the configuration and
allowing to use remote configurations fetched from an HTTP url (which may well
be the local IPFS gateway). This allows cluster administrators to provide
the configurations needed for any peers to join a cluster as followers.

The CRDT arrival incorporates a large number of improvements in peerset
management, bootstrapping, connection management and auto-recovery of peers
after network disconnections. We have improved the peer monitoring system,
added support for efficient Pin-Update-based pinning, reworked timeout control
for pinning and fixed a number of annoying bugs.

This release is mostly backwards compatible with the previous one and
clusters should keep working with the same configurations, but users should
have a look to the sections below and read the updated documentation, as a
number of changes have been introduced to support both consensus components.

Consensus selection happens during initialization of the configuration (see
configuration changes below). Migration of the pinset is necessary by doing
`state export` (with Raft configured), followed by `state import` (with CRDT
configured). Note that all peers should be configured with the same consensus
type.


#### List of changes

##### Features


* crdt: introduce crdt-based consensus component | [glvd/cluster#685](https://github.com/glvd/cluster/issues/685) | [glvd/cluster#804](https://github.com/glvd/cluster/issues/804) | [glvd/cluster#787](https://github.com/glvd/cluster/issues/787) | [glvd/cluster#798](https://github.com/glvd/cluster/issues/798) | [glvd/cluster#805](https://github.com/glvd/cluster/issues/805) | [glvd/cluster#811](https://github.com/glvd/cluster/issues/811) | [glvd/cluster#816](https://github.com/glvd/cluster/issues/816) | [glvd/cluster#820](https://github.com/glvd/cluster/issues/820) | [glvd/cluster#856](https://github.com/glvd/cluster/issues/856) | [glvd/cluster#857](https://github.com/glvd/cluster/issues/857) | [glvd/cluster#834](https://github.com/glvd/cluster/issues/834) | [glvd/cluster#856](https://github.com/glvd/cluster/issues/856) | [glvd/cluster#867](https://github.com/glvd/cluster/issues/867) | [glvd/cluster#874](https://github.com/glvd/cluster/issues/874) | [glvd/cluster#885](https://github.com/glvd/cluster/issues/885) | [glvd/cluster#899](https://github.com/glvd/cluster/issues/899) | [glvd/cluster#906](https://github.com/glvd/cluster/issues/906) | [glvd/cluster#918](https://github.com/glvd/cluster/issues/918)
* configs: separate identity and configuration | [glvd/cluster#760](https://github.com/glvd/cluster/issues/760) | [glvd/cluster#766](https://github.com/glvd/cluster/issues/766) | [glvd/cluster#780](https://github.com/glvd/cluster/issues/780)
* configs: support running with a remote `service.json` (http) | [glvd/cluster#868](https://github.com/glvd/cluster/issues/868)
* configs: support a `follower_mode` option | [glvd/cluster#803](https://github.com/glvd/cluster/issues/803) | [glvd/cluster#864](https://github.com/glvd/cluster/issues/864)
* service/configs: do not load API components if no config present | [glvd/cluster#452](https://github.com/glvd/cluster/issues/452) | [glvd/cluster#836](https://github.com/glvd/cluster/issues/836)
* service: add `ipfs-cluster-service init --peers` flag to initialize with given peers | [glvd/cluster#835](https://github.com/glvd/cluster/issues/835) | [glvd/cluster#839](https://github.com/glvd/cluster/issues/839) | [glvd/cluster#870](https://github.com/glvd/cluster/issues/870)
* cluster: RPC auth: block rpc endpoints for non trusted peers | [glvd/cluster#775](https://github.com/glvd/cluster/issues/775) | [glvd/cluster#710](https://github.com/glvd/cluster/issues/710) | [glvd/cluster#666](https://github.com/glvd/cluster/issues/666) | [glvd/cluster#773](https://github.com/glvd/cluster/issues/773) | [glvd/cluster#905](https://github.com/glvd/cluster/issues/905)
* cluster: introduce connection manager | [glvd/cluster#791](https://github.com/glvd/cluster/issues/791)
* cluster: support new `PinUpdate` option for new pins | [glvd/cluster#869](https://github.com/glvd/cluster/issues/869) | [glvd/cluster#732](https://github.com/glvd/cluster/issues/732)
* cluster: trigger `Recover` automatically on a configurable interval | [glvd/cluster#831](https://github.com/glvd/cluster/issues/831) | [glvd/cluster#887](https://github.com/glvd/cluster/issues/887)
* cluster: enable mDNS discovery for peers | [glvd/cluster#882](https://github.com/glvd/cluster/issues/882) | [glvd/cluster#900](https://github.com/glvd/cluster/issues/900)
* IPFS Proxy: Support `pin/update` | [glvd/cluster#732](https://github.com/glvd/cluster/issues/732) | [glvd/cluster#768](https://github.com/glvd/cluster/issues/768) | [glvd/cluster#887](https://github.com/glvd/cluster/issues/887)
* monitor: Accrual failure detection. Leaderless re-pinning | [glvd/cluster#413](https://github.com/glvd/cluster/issues/413) | [glvd/cluster#713](https://github.com/glvd/cluster/issues/713) | [glvd/cluster#714](https://github.com/glvd/cluster/issues/714) | [glvd/cluster#812](https://github.com/glvd/cluster/issues/812) | [glvd/cluster#813](https://github.com/glvd/cluster/issues/813) | [glvd/cluster#814](https://github.com/glvd/cluster/issues/814) | [glvd/cluster#815](https://github.com/glvd/cluster/issues/815)
* datastore: Expose badger configuration | [glvd/cluster#771](https://github.com/glvd/cluster/issues/771) | [glvd/cluster#776](https://github.com/glvd/cluster/issues/776)
* IPFSConnector: pin timeout start counting from last received block | [glvd/cluster#497](https://github.com/glvd/cluster/issues/497) | [glvd/cluster#738](https://github.com/glvd/cluster/issues/738)
* IPFSConnector: remove pin method options | [glvd/cluster#875](https://github.com/glvd/cluster/issues/875)
* IPFSConnector: `unpin_disable` removes the ability to unpin anything from ipfs (experimental) | [glvd/cluster#793](https://github.com/glvd/cluster/issues/793) | [glvd/cluster#832](https://github.com/glvd/cluster/issues/832)
* REST API Client: Load-balancing Go client | [glvd/cluster#448](https://github.com/glvd/cluster/issues/448) | [glvd/cluster#737](https://github.com/glvd/cluster/issues/737)
* REST API: Return allocation objects on pin/unpin | [glvd/cluster#843](https://github.com/glvd/cluster/issues/843)
* REST API: Support request logging | [glvd/cluster#574](https://github.com/glvd/cluster/issues/574) | [glvd/cluster#894](https://github.com/glvd/cluster/issues/894)
* Adder: improve error handling. Keep adding while at least one allocation works | [glvd/cluster#852](https://github.com/glvd/cluster/issues/852) | [glvd/cluster#871](https://github.com/glvd/cluster/issues/871)
* Adder: support user-given allocations for the `Add` operation | [glvd/cluster#761](https://github.com/glvd/cluster/issues/761) | [glvd/cluster#890](https://github.com/glvd/cluster/issues/890)
* ctl: support adding pin metadata | [glvd/cluster#670](https://github.com/glvd/cluster/issues/670) | [glvd/cluster#891](https://github.com/glvd/cluster/issues/891)


##### Bug fixes

* REST API: Fix `/allocations` when filter unset | [glvd/cluster#762](https://github.com/glvd/cluster/issues/762)
* REST API: Fix DELETE returning 500 when pin does not exist | [glvd/cluster#742](https://github.com/glvd/cluster/issues/742) | [glvd/cluster#854](https://github.com/glvd/cluster/issues/854)
* REST API: Return JSON body on 404s | [glvd/cluster#657](https://github.com/glvd/cluster/issues/657) | [glvd/cluster#879](https://github.com/glvd/cluster/issues/879)
* service: connectivity fixes | [glvd/cluster#787](https://github.com/glvd/cluster/issues/787) | [glvd/cluster#792](https://github.com/glvd/cluster/issues/792)
* service: fix using `/dnsaddr` peers | [glvd/cluster#818](https://github.com/glvd/cluster/issues/818)
* service: reading empty lines on peerstore panics | [glvd/cluster#886](https://github.com/glvd/cluster/issues/886)
* service/ctl: fix parsing string lists | [glvd/cluster#876](https://github.com/glvd/cluster/issues/876) | [glvd/cluster#841](https://github.com/glvd/cluster/issues/841)
* IPFSConnector: `pin/ls` does handle base32 and base58 cids properly | [glvd/cluster#808](https://github.com/glvd/cluster/issues/808) [glvd/cluster#809](https://github.com/glvd/cluster/issues/809)
* configs: some config keys not matching ENV vars names | [glvd/cluster#837](https://github.com/glvd/cluster/issues/837) | [glvd/cluster#778](https://github.com/glvd/cluster/issues/778)
* raft: delete removed raft peers from peerstore | [glvd/cluster#840](https://github.com/glvd/cluster/issues/840) | [glvd/cluster#846](https://github.com/glvd/cluster/issues/846)
* cluster: peers forgotten after being down | [glvd/cluster#648](https://github.com/glvd/cluster/issues/648) | [glvd/cluster#860](https://github.com/glvd/cluster/issues/860)
* cluster: State sync should not keep tracking when queue is full | [glvd/cluster#377](https://github.com/glvd/cluster/issues/377) | [glvd/cluster#901](https://github.com/glvd/cluster/issues/901)
* cluster: avoid random order on peer lists and listen multiaddresses | [glvd/cluster#327](https://github.com/glvd/cluster/issues/327) | [glvd/cluster#878](https://github.com/glvd/cluster/issues/878)
* cluster: fix recover and allocation re-assignment to existing pins | [glvd/cluster#912](https://github.com/glvd/cluster/issues/912) | [glvd/cluster#888](https://github.com/glvd/cluster/issues/888)

##### Other changes

* cluster: Dependency updates | [glvd/cluster#769](https://github.com/glvd/cluster/issues/769) | [glvd/cluster#789](https://github.com/glvd/cluster/issues/789) | [glvd/cluster#795](https://github.com/glvd/cluster/issues/795) | [glvd/cluster#822](https://github.com/glvd/cluster/issues/822) | [glvd/cluster#823](https://github.com/glvd/cluster/issues/823) | [glvd/cluster#828](https://github.com/glvd/cluster/issues/828) | [glvd/cluster#830](https://github.com/glvd/cluster/issues/830) | [glvd/cluster#853](https://github.com/glvd/cluster/issues/853) | [glvd/cluster#839](https://github.com/glvd/cluster/issues/839)
* cluster: Set `[]peer.ID` as type for user allocations | [glvd/cluster#767](https://github.com/glvd/cluster/issues/767)
* cluster: RPC: Split services among components | [glvd/cluster#773](https://github.com/glvd/cluster/issues/773)
* cluster: Multiple improvements to tests | [glvd/cluster#360](https://github.com/glvd/cluster/issues/360) | [glvd/cluster#502](https://github.com/glvd/cluster/issues/502) | [glvd/cluster#779](https://github.com/glvd/cluster/issues/779) | [glvd/cluster#833](https://github.com/glvd/cluster/issues/833) | [glvd/cluster#863](https://github.com/glvd/cluster/issues/863) | [glvd/cluster#883](https://github.com/glvd/cluster/issues/883) | [glvd/cluster#884](https://github.com/glvd/cluster/issues/884) | [glvd/cluster#797](https://github.com/glvd/cluster/issues/797) | [glvd/cluster#892](https://github.com/glvd/cluster/issues/892)
* cluster: Remove Gx | [glvd/cluster#765](https://github.com/glvd/cluster/issues/765) | [glvd/cluster#781](https://github.com/glvd/cluster/issues/781)
* cluster: Use `/p2p/` instead of `/ipfs/` in multiaddresses | [glvd/cluster#431](https://github.com/glvd/cluster/issues/431) | [glvd/cluster#877](https://github.com/glvd/cluster/issues/877)
* cluster: consolidate parsing of pin options | [glvd/cluster#913](https://github.com/glvd/cluster/issues/913)
* REST API: Replace regexps with `strings.HasPrefix` | [glvd/cluster#806](https://github.com/glvd/cluster/issues/806) | [glvd/cluster#807](https://github.com/glvd/cluster/issues/807)
* docker: use GOPROXY to build containers | [glvd/cluster#872](https://github.com/glvd/cluster/issues/872)
* docker: support `IPFS_CLUSTER_CONSENSUS` flag and other improvements | [glvd/cluster#882](https://github.com/glvd/cluster/issues/882)
* ctl: increase space for peernames | [glvd/cluster#887](https://github.com/glvd/cluster/issues/887)
* ctl: improve replication factor 0 explanation | [glvd/cluster#755](https://github.com/glvd/cluster/issues/755) | [glvd/cluster#909](https://github.com/glvd/cluster/issues/909)

#### Upgrading notices


##### Configuration changes

This release introduces a number of backwards-compatible configuration changes:

* The `service.json` file no longer includes `ID` and `PrivateKey`, which are
  now part of an `identity.json` file. However things should work as before if
  they do. Running `ipfs-cluster-service daemon` on a older configuration will
  automatically write an `identity.json` file with the old credentials so that
  things do not break when the compatibility hack is removed.

* The `service.json` can use a new single top-level `source` field which can
  be set to an HTTP url pointing to a full `service.json`. When present,
  this will be read and used when starting the daemon. `ipfs-cluster-service
  init http://url` produces this type of "remote configuration" file.

* `cluster` section:
  * A new, hidden `follower_mode` option has been introduced in the main
    `cluster` configuration section. When set, the cluster peer will provide
    clear errors when pinning or unpinning. This is a UI feature. The capacity
    of a cluster peer to pin/unpin depends on whether it is trusted by other
    peers, not on settin this hidden option.
  * A new `pin_recover_interval` option to controls how often pins in error
    states are retried.
  * A new `mdns_interval` controls the time between mDNS broadcasts to
    discover other peers in the network. Setting it to 0 disables mDNS
    altogether (default is 10 seconds).
  * A new `connection_manager` object can be used to limit the number of
    connections kept by the libp2p host:

```js
"connection_manager": {
    "high_water": 400,
    "low_water": 100,
    "grace_period": "2m0s"
},
```


* `consensus` section:
  * Only one configuration object is allowed inside the `consensus` section,
    and it must be either the `crdt` or the `raft` one. The presence of one or
    another is used to autoselect the consensus component to be used when
    running the daemon or performing `ipfs-cluster-service state`
    operations. `ipfs-cluster-service init` receives an optional `--consensus`
    flag to select which one to produce. By default it is the `crdt`.

* `ipfs_connector/ipfshttp` section:
  * The `pin_timeout` in the `ipfshttp` section is now starting from the last
    block received. Thus it allows more flexibility for things which are
    pinning very slowly, but still pinning.
  * The `pin_method` option has been removed, as go-ipfs does not do a
    pin-global-lock anymore. Therefore `pin add` will be called directly, can
    be called multiple times in parallel and should be faster than the
    deprecated `refs -r` way.
  * The `ipfshttp` section has a new (hidden) `unpin_disable` option
    (boolean). The component will refuse to unpin anything from IPFS when
    enabled. It can be used as a failsafe option to make sure cluster peers
    never unpin content.

* `datastore` section:
  * The configuration has a new `datastore/badger` section, which is relevant
    when using the `crdt` consensus component. It allows full control of the
    [Badger configuration](https://godoc.org/github.com/dgraph-io/badger#Options),
    which is particuarly important when running on systems with low memory:
  

```
  "datastore": {
    "badger": {
      "badger_options": {
        "dir": "",
        "value_dir": "",
        "sync_writes": true,
        "table_loading_mode": 2,
        "value_log_loading_mode": 2,
        "num_versions_to_keep": 1,
        "max_table_size": 67108864,
        "level_size_multiplier": 10,
        "max_levels": 7,
        "value_threshold": 32,
        "num_memtables": 5,
        "num_level_zero_tables": 5,
        "num_level_zero_tables_stall": 10,
        "level_one_size": 268435456,
        "value_log_file_size": 1073741823,
        "value_log_max_entries": 1000000,
        "num_compactors": 2,
        "compact_l_0_on_close": true,
        "read_only": false,
        "truncate": false
      }
    }
    }
```

* `pin_tracker/maptracker` section:
  * The `max_pin_queue_size` parameter has been hidden for default
    configurations and the default has been set to 1000000. 

* `api/restapi` section:
  * A new `http_log_file` options allows to redirect the REST API logging to a
    file. Otherwise, it is logged as part of the regular log. Lines follow the
    Apache Common Log Format (CLF).

##### REST API

The `POST /pins/{cid}` and `DELETE /pins/{cid}` now returns a pin object with
`200 Success` rather than an empty `204 Accepted` response.

Using an unexistent route will now correctly return a JSON object along with
the 404 HTTP code, rather than text.

##### Go APIs

There have been some changes to Go APIs. Applications integrating Cluster
directly will be affected by the new signatures of Pin/Unpin:

* The `Pin` and `Unpin` methods now return an object of `api.Pin` type, along with an error.
* The `Pin` method takes a CID and `PinOptions` rather than an `api.Pin` object wrapping
those.
* A new `PinUpdate` method has been introduced.

Additionally:

* The Consensus Component interface has changed to accommodate peer-trust operations.
* The IPFSConnector Component interface `Pin` method has changed to take an `api.Pin` type.


##### Other

* The IPFS Proxy now hijacks the `/api/v0/pin/update` and makes a Cluster PinUpdate.
* `ipfs-cluster-service init` now takes a `--consensus` flag to select between
  `crdt` (default) and `raft`. Depending on the values, the generated
  configuration will have the relevant sections for each.
* The Dockerfiles have been updated to:
  * Support the `IPFS_CLUSTER_CONSENSUS` flag to determine which consensus to
  use for the automatic `init`.
  * No longer use `IPFS_API` environment variable to do a `sed` replacement on
    the config, as `CLUSTER_IPFSHTTP_NODEMULTIADDRESS` is the canonical one to
    use.
  * No longer use `sed` replacement to set the APIs listen IPs to `0.0.0.0`
    automatically, as this can be achieved with environment variables
    (`CLUSTER_RESTAPI_HTTPLISTENMULTIADDRESS` and
    `CLUSTER_IPFSPROXY_LISTENMULTIADDRESS`) and can be dangerous for containers
    running in `net=host` mode.
  * The `docker-compose.yml` has been updated and simplified to launch a CRDT
    3-peer TEST cluster
* Cluster now uses `/p2p/` instead of `/ipfs/` for libp2p multiaddresses by
  default, but both protocol IDs are equivalent and interchangeable.
* Pinning an already existing pin will re-submit it to the consensus layer in
  all cases, meaning that pins in error states will start pinning again
  (before, sometimes this was only possible using recover). Recover stays as a
  broadcast/sync operation to trigger pinning on errored items. As a reminder,
  pin is consensus/async operation.
    
---


### v0.10.1 - 2019-04-10

#### Summary

This release is a maintenance release with a number of bug fixes and a couple of small features.


#### List of changes

##### Features

* Switch to go.mod | [glvd/cluster#706](https://github.com/glvd/cluster/issues/706) | [glvd/cluster#707](https://github.com/glvd/cluster/issues/707) | [glvd/cluster#708](https://github.com/glvd/cluster/issues/708)
* Remove basic monitor | [glvd/cluster#689](https://github.com/glvd/cluster/issues/689) | [glvd/cluster#726](https://github.com/glvd/cluster/issues/726)
* Support `nocopy` when adding URLs | [glvd/cluster#735](https://github.com/glvd/cluster/issues/735)

##### Bug fixes

* Mitigate long header attack | [glvd/cluster#636](https://github.com/glvd/cluster/issues/636) | [glvd/cluster#712](https://github.com/glvd/cluster/issues/712)
* Fix download link in README | [glvd/cluster#723](https://github.com/glvd/cluster/issues/723)
* Fix `peers ls` error when peers down | [glvd/cluster#715](https://github.com/glvd/cluster/issues/715) | [glvd/cluster#719](https://github.com/glvd/cluster/issues/719)
* Nil pointer panic on `ipfs-cluster-ctl add` | [glvd/cluster#727](https://github.com/glvd/cluster/issues/727) | [glvd/cluster#728](https://github.com/glvd/cluster/issues/728)
* Fix `enc=json` output on `ipfs-cluster-ctl add` | [glvd/cluster#729](https://github.com/glvd/cluster/issues/729)
* Add SSL CAs to Docker container | [glvd/cluster#730](https://github.com/glvd/cluster/issues/730) | [glvd/cluster#731](https://github.com/glvd/cluster/issues/731)
* Remove duplicate import | [glvd/cluster#734](https://github.com/glvd/cluster/issues/734)
* Fix version json object | [glvd/cluster#743](https://github.com/glvd/cluster/issues/743) | [glvd/cluster#752](https://github.com/glvd/cluster/issues/752)

#### Upgrading notices



##### Configuration changes

There are no configuration changes on this release.

##### REST API

The `/version` endpoint now returns a version object with *lowercase* `version` key.

##### Go APIs

There are no changes to the Go APIs.

##### Other

Since we have switched to Go modules for dependency management, `gx` is no
longer used and the maintenance of Gx dependencies has been dropped. The
`Makefile` has been updated accordinly, but now a simple `go install
./cmd/...` works.

---

### v0.10.0 - 2019-03-07

#### Summary

As we get ready to introduce a new CRDT-based "consensus" component to replace
Raft, IPFS Cluster 0.10.0 prepares the ground with substancial under-the-hood
changes. many performance improvements and a few very useful features.

First of all, this release **requires** users to run `state upgrade` (or start
their daemons with `ipfs-cluster-service daemon --upgrade`). This is the last
upgrade in this fashion as we turn to go-datastore-based storage. The next
release of IPFS Cluster will not understand or be able to upgrade anything
below 0.10.0.

Secondly, we have made some changes to internal types that should greatly
improve performance a lot, particularly calls involving large collections of
items (`pin ls` or `status`). There are also changes on how the state is
serialized, avoiding unnecessary in-memory copies. We have also upgraded the
dependency stack, incorporating many fixes from libp2p.

Thirdly, our new great features:

* `ipfs-cluster-ctl pin add/rm` now supports IPFS paths (`/ipfs/Qmxx.../...`,
  `/ipns/Qmxx.../...`, `/ipld/Qm.../...`) which are resolved automatically
  before pinning.
* All our configuration values can now be set via environment variables, and
these will be reflected when initializing a new configuration file.
* Pins can now specify a list of "priority allocations". This allows to pin
items to specific Cluster peers, overriding the default allocation policy.
* Finally, the REST API supports adding custom metadata entries as `key=value`
  (we will soon add support in `ipfs-cluster-ctl`). Metadata can be added as
  query arguments to the Pin or PinPath endpoints: `POST
  /pins/<cid-or-path>?meta-key1=value1&meta-key2=value2...`

Note that on this release we have also removed a lot of backwards-compatiblity
code for things older than version 0.8.0, which kept things working but
printed respective warnings. If you're upgrading from an old release, consider
comparing your configuration with the new default one.


#### List of changes

##### Features

  * Add full support for environment variables in configurations and initialization | [glvd/cluster#656](https://github.com/glvd/cluster/issues/656) | [glvd/cluster#663](https://github.com/glvd/cluster/issues/663) | [glvd/cluster#667](https://github.com/glvd/cluster/issues/667)
  * Switch to codecov | [glvd/cluster#683](https://github.com/glvd/cluster/issues/683)
  * Add auto-resolving IPFS paths | [glvd/cluster#450](https://github.com/glvd/cluster/issues/450) | [glvd/cluster#634](https://github.com/glvd/cluster/issues/634)
  * Support user-defined allocations | [glvd/cluster#646](https://github.com/glvd/cluster/issues/646) | [glvd/cluster#647](https://github.com/glvd/cluster/issues/647)
  * Support user-defined metadata in pin objects | [glvd/cluster#681](https://github.com/glvd/cluster/issues/681)
  * Make normal types serializable and remove `*Serial` types | [glvd/cluster#654](https://github.com/glvd/cluster/issues/654) | [glvd/cluster#688](https://github.com/glvd/cluster/issues/688) | [glvd/cluster#700](https://github.com/glvd/cluster/issues/700)
  * Support IPFS paths in the IPFS proxy | [glvd/cluster#480](https://github.com/glvd/cluster/issues/480) | [glvd/cluster#690](https://github.com/glvd/cluster/issues/690)
  * Use go-datastore as backend for the cluster state | [glvd/cluster#655](https://github.com/glvd/cluster/issues/655)
  * Upgrade dependencies | [glvd/cluster#675](https://github.com/glvd/cluster/issues/675) | [glvd/cluster#679](https://github.com/glvd/cluster/issues/679) | [glvd/cluster#686](https://github.com/glvd/cluster/issues/686) | [glvd/cluster#687](https://github.com/glvd/cluster/issues/687)
  * Adopt MIT+Apache 2 License (no more sign-off required) | [glvd/cluster#692](https://github.com/glvd/cluster/issues/692)
  * Add codecov configurtion file | [glvd/cluster#693](https://github.com/glvd/cluster/issues/693)
  * Additional tests for basic auth | [glvd/cluster#645](https://github.com/glvd/cluster/issues/645) | [glvd/cluster#694](https://github.com/glvd/cluster/issues/694)

##### Bug fixes

  * Fix docker compose tests | [glvd/cluster#696](https://github.com/glvd/cluster/issues/696)
  * Hide `ipfsproxy.extract_headers_ttl` and `ipfsproxy.extract_headers_path` options by default | [glvd/cluster#699](https://github.com/glvd/cluster/issues/699)

#### Upgrading notices

This release needs an state upgrade before starting the Cluster daemon. Run `ipfs-cluster-service state upgrade` or run it as `ipfs-cluster-service daemon --upgrade`. We recommend backing up the `~/.ipfs-cluster` folder or exporting the pinset with `ipfs-cluster-service state export`.

##### Configuration changes

Configurations now respects environment variables for all sections. They are
in the form:

`CLUSTER_COMPONENTNAME_KEYNAMEWITHOUTSPACES=value`

Environment variables will override `service.json` configuration options when
defined and the Cluster peer is started. `ipfs-cluster-service init` will
reflect the value of any existing environment variables in the new
`service.json` file.

##### REST API

The main breaking change to the REST API corresponds to the JSON
representation of CIDs in response objects:

* Before: `"cid": "Qm...."`
* Now: `"cid": { "/": "Qm...."}`

The new CID encoding is the default as defined by the `cid`
library. Unfortunately, there is no good solution to keep the previous
representation without copying all the objects (an innefficient technique we
just removed). The new CID encoding is otherwise aligned with the rest of the
stack.

The API also gets two new "Path" endpoints:

* `POST /pins/<ipfs|ipns|ipld>/<path>/...` and
* `DELETE /pins/<ipfs|ipns|ipld>/<path>/...`

Thus, it is equivalent to pin a CID with `POST /pins/<cid>` (as before) or
with `POST /pins/ipfs/<cid>`.

The calls will however fail when a non-compliant IPFS path is provided: `POST
/pins/<cid>/my/path` will fail because all paths must start with the `/ipfs`,
`/ipns` or `/ipld` components.

##### Go APIs

This release introduces lots of changes to the Go APIs, including the Go REST
API client, as we have started returning pointers to objects rather than the
objects directly. The `Pin` will now take `api.PinOptions` instead of
different arguments corresponding to the options. It is aligned with the new
`PinPath` and `UnpinPath`.

##### Other

As pointed above, 0.10.0's state migration is a required step to be able to
use future version of IPFS Cluster.

---

### v0.9.0 - 2019-02-18

#### Summary

IPFS Cluster version 0.9.0 comes with one big new feature, [OpenCensus](https://opencensus.io) support! This allows for the collection of distributed traces and metrics from the IPFS Cluster application as well as supporting libraries. Currently, we support the use of [Jaeger](https://jaegertracing.io) as the tracing backend and [Prometheus](https://prometheus.io) as the metrics backend. Support for other [OpenCensus backends](https://opencensus.io/exporters/) will be added as requested by the community.

#### List of changes

##### Features

  * Integrate [OpenCensus](https://opencensus.io) tracing and metrics into IPFS Cluster codebase | [glvd/cluster#486](https://github.com/glvd/cluster/issues/486) | [glvd/cluster#658](https://github.com/glvd/cluster/issues/658) | [glvd/cluster#659](https://github.com/glvd/cluster/issues/659) | [glvd/cluster#676](https://github.com/glvd/cluster/issues/676) | [glvd/cluster#671](https://github.com/glvd/cluster/issues/671) | [glvd/cluster#674](https://github.com/glvd/cluster/issues/674)

##### Bug Fixes

No bugs were fixed from the previous release.

##### Deprecated

  * The snap distribution of IPFS Cluster has been removed | [glvd/cluster#593](https://github.com/glvd/cluster/issues/593) | [glvd/cluster#649](https://github.com/glvd/cluster/issues/649).

#### Upgrading notices

##### Configuration changes

No changes to the existing configuration.

There are two new configuration sections with this release:

###### `tracing` section

The `tracing` section configures the use of Jaeger as a tracing backend.

```js
    "tracing": {
      "enable_tracing": false,
      "jaeger_agent_endpoint": "/ip4/0.0.0.0/udp/6831",
      "sampling_prob": 0.3,
      "service_name": "cluster-daemon"
    }
```

###### `metrics` section

The `metrics` section configures the use of Prometheus as a metrics collector.

```js
    "metrics": {
      "enable_stats": false,
      "prometheus_endpoint": "/ip4/0.0.0.0/tcp/8888",
      "reporting_interval": "2s"
    }
```

##### REST API

No changes to the REST API.

##### Go APIs

The Go APIs had the minor change of having a `context.Context` parameter added as the first argument 
to those that didn't already have it. This was to enable the proporgation of tracing and metric
values.

The following is a list of interfaces and their methods that were affected by this change:
 - Component
    - Shutdown
 - Consensus
    - Ready
    - LogPin
    - LogUnpin
    - AddPeer
    - RmPeer
    - State
    - Leader
    - WaitForSync
    - Clean
    - Peers
 - IpfsConnector
    - ID
    - ConnectSwarm
    - SwarmPeers
    - RepoStat
    - BlockPut
    - BlockGet
 - Peered
    - AddPeer
    - RmPeer
 - PinTracker
    - Track
    - Untrack
    - StatusAll
    - Status
    - SyncAll
    - Sync
    - RecoverAll
    - Recover
 - Informer
    - GetMetric
 - PinAllocator
    - Allocate
 - PeerMonitor
    - LogMetric
    - PublishMetric
    - LatestMetrics
 - state.State
    - Add
    - Rm
    - List
    - Has
    - Get
    - Migrate
 - rest.Client
    - ID
    - Peers
    - PeerAdd
    - PeerRm
    - Add
    - AddMultiFile
    - Pin
    - Unpin
    - Allocations
    - Allocation
    - Status
    - StatusAll
    - Sync
    - SyncAll
    - Recover
    - RecoverAll
    - Version
    - IPFS
    - GetConnectGraph
    - Metrics

These interface changes were also made in the respective implementations.
All export methods of the Cluster type also had these changes made.


##### Other

No other things.

---

### v0.8.0 - 2019-01-16

#### Summary

IPFS Cluster version 0.8.0 comes with a few useful features and some bugfixes.
A significant amount of work has been put to correctly handle CORS in both the
REST API and the IPFS Proxy endpoint, fixing some long-standing issues (we
hope once are for all).

There has also been heavy work under the hood to separate the IPFS HTTP
Connector (the HTTP client to the IPFS daemon) from the IPFS proxy, which is
essentially an additional Cluster API. Check the configuration changes section
below for more information about how this affects the configuration file.

Finally we have some useful small features:

* The `ipfs-cluster-ctl status --filter` option allows to just list those
items which are still `pinning` or `queued` or `error` etc. You can combine
multiple filters. This translates to a new `filter` query parameter in the
`/pins` API endpoint.
* The `stream-channels=false` query parameter for the `/add` endpoint will let
the API buffer the output when adding and return a valid JSON array once done,
making this API endpoint behave like a regular, non-streaming one.
`ipfs-cluster-ctl add --no-stream` acts similarly, but buffering on the client
side. Note that this will cause in-memory buffering of potentially very large
responses when the number of added files is very large, but should be
perfectly fine for regular usage.
* The `ipfs-cluster-ctl add --quieter` flag now applies to the JSON output
too, allowing the user to just get the last added entry JSON object when
adding a file, which is always the root hash.

#### List of changes

##### Features

  * IPFS Proxy extraction to its own `API` component: `ipfsproxy` | [glvd/cluster#453](https://github.com/glvd/cluster/issues/453) | [glvd/cluster#576](https://github.com/glvd/cluster/issues/576) | [glvd/cluster#616](https://github.com/glvd/cluster/issues/616) | [glvd/cluster#617](https://github.com/glvd/cluster/issues/617)
  * Add full CORS handling to `restapi` | [glvd/cluster#639](https://github.com/glvd/cluster/issues/639) | [glvd/cluster#640](https://github.com/glvd/cluster/issues/640)
  * `restapi` configuration section entries can be overriden from environment variables | [glvd/cluster#609](https://github.com/glvd/cluster/issues/609)
  * Update to `go-ipfs-files` 2.0 | [glvd/cluster#613](https://github.com/glvd/cluster/issues/613)
  * Tests for the `/monitor/metrics` endpoint | [glvd/cluster#587](https://github.com/glvd/cluster/issues/587) | [glvd/cluster#622](https://github.com/glvd/cluster/issues/622)
  * Support `stream-channels=fase` query parameter in `/add` | [glvd/cluster#632](https://github.com/glvd/cluster/issues/632) | [glvd/cluster#633](https://github.com/glvd/cluster/issues/633)
  * Support server side `/pins` filtering  | [glvd/cluster#445](https://github.com/glvd/cluster/issues/445) | [glvd/cluster#478](https://github.com/glvd/cluster/issues/478) | [glvd/cluster#627](https://github.com/glvd/cluster/issues/627)
  * `ipfs-cluster-ctl add --no-stream` option | [glvd/cluster#632](https://github.com/glvd/cluster/issues/632) | [glvd/cluster#637](https://github.com/glvd/cluster/issues/637)
  * Upgrade dependencies and libp2p to version 6.0.29 | [glvd/cluster#624](https://github.com/glvd/cluster/issues/624)

##### Bug fixes

 * Respect IPFS daemon response headers on non-proxied calls | [glvd/cluster#382](https://github.com/glvd/cluster/issues/382) | [glvd/cluster#623](https://github.com/glvd/cluster/issues/623) | [glvd/cluster#638](https://github.com/glvd/cluster/issues/638)
 * Fix `ipfs-cluster-ctl` usage with HTTPs and `/dns*` hostnames | [glvd/cluster#626](https://github.com/glvd/cluster/issues/626)
 * Minor fixes in sharness | [glvd/cluster#641](https://github.com/glvd/cluster/issues/641) | [glvd/cluster#643](https://github.com/glvd/cluster/issues/643)
 * Fix error handling when parsing the configuration | [glvd/cluster#642](https://github.com/glvd/cluster/issues/642)



#### Upgrading notices

This release comes with some configuration changes that are important to notice,
even though the peers will start with the same configurations as before.

##### Configuration changes

##### `ipfsproxy` section

This version introduces a separate `ipfsproxy` API component. This is
reflected in the `service.json` configuration, which now includes a new
`ipfsproxy` subsection under the `api` section. By default it looks like:

```js
    "ipfsproxy": {
      "node_multiaddress": "/ip4/127.0.0.1/tcp/5001",
      "listen_multiaddress": "/ip4/127.0.0.1/tcp/9095",
      "read_timeout": "0s",
      "read_header_timeout": "5s",
      "write_timeout": "0s",
      "idle_timeout": "1m0s"
   }
```

We have however added the necessary safeguards to keep backwards compatibility
for this release. If the `ipfsproxy` section is empty, it will be picked up from
the `ipfshttp` section as before. An ugly warning will be printed in this case.

Based on the above, the `ipfshttp` configuration section loses the
proxy-related options. Note that `node_multiaddress` stays in both component
configurations and should likely be the same in most cases, but you can now
potentially proxy requests to a different daemon than the one used by the
cluster peer.

Additional hidden configuration options to manage custom header extraction
from the IPFS daemon (for power users) have been added to the `ipfsproxy`
section but are not shown by default when initializing empty
configurations. See the documentation for more details.

###### `restapi` section

The introduction of proper CORS handling in the `restapi` component introduces
a number of new keys:

```js
      "cors_allowed_origins": [
        "*"
      ],
      "cors_allowed_methods": [
        "GET"
      ],
      "cors_allowed_headers": [],
      "cors_exposed_headers": [
        "Content-Type",
        "X-Stream-Output",
        "X-Chunked-Output",
        "X-Content-Length"
      ],
      "cors_allow_credentials": true,
      "cors_max_age": "0s"
```

Note that CORS will be essentially unconfigured when these keys are not
defined.

The `headers` key, which was used before to add some CORS related headers
manually, takes a new empty default. **We recommend emptying `headers` from
any CORS-related value.**


##### REST API

The REST API is fully backwards compatible:

* The `GET /pins` endpoint takes a new `?filter=<filter>` option. See
  `ipfs-cluster-ctl status --help` for acceptable values.
* The `POST /add` endpoint accepts a new `?stream-channels=<true|false>`
  option. By default it is set to `true`.

##### Go APIs

The signature for the `StatusAll` method in the REST `client` module has
changed to include a `filter` parameter.

There may have been other minimal changes to internal exported Go APIs, but
should not affect users.

##### Other

Proxy requests which are handled by the Cluster peer (`/pin/ls`, `/pin/add`,
`/pin/rm`, `/repo/stat` and `/add`) will now attempt to fully mimic ipfs
responses to the header level. This is done by triggering CORS pre-flight for
every hijacked request along with an occasional regular request to `/version`
to extract other headers (and possibly custom ones).

The practical result is that the proxy now behaves correctly when dropped
instead of IPFS into CORS-aware contexts (like the browser).

---

### v0.7.0 - 2018-11-01

#### Summary

IPFS Cluster version 0.7.0 is a maintenance release that includes a few bugfixes and some small features.

Note that the REST API response format for the `/add` endpoint has changed. Thus all clients need to be upgraded to deal with the new format. The `rest/api/client` has been accordingly updated.

#### List of changes

##### Features

  * Clean (rotate) the state when running `init` | [glvd/cluster#532](https://github.com/glvd/cluster/issues/532) | [glvd/cluster#553](https://github.com/glvd/cluster/issues/553)
  * Configurable REST API headers and CORS defaults | [glvd/cluster#578](https://github.com/glvd/cluster/issues/578)
  * Upgrade libp2p and other deps | [glvd/cluster#580](https://github.com/glvd/cluster/issues/580) | [glvd/cluster#590](https://github.com/glvd/cluster/issues/590) | [glvd/cluster#592](https://github.com/glvd/cluster/issues/592) | [glvd/cluster#598](https://github.com/glvd/cluster/issues/598) | [glvd/cluster#599](https://github.com/glvd/cluster/issues/599)
  * Use `gossipsub` to broadcast metrics | [glvd/cluster#573](https://github.com/glvd/cluster/issues/573)
  * Download gx and gx-go from IPFS preferentially | [glvd/cluster#577](https://github.com/glvd/cluster/issues/577) | [glvd/cluster#581](https://github.com/glvd/cluster/issues/581)
  * Expose peer metrics in the API + ctl commands | [glvd/cluster#449](https://github.com/glvd/cluster/issues/449) | [glvd/cluster#572](https://github.com/glvd/cluster/issues/572) | [glvd/cluster#589](https://github.com/glvd/cluster/issues/589) | [glvd/cluster#587](https://github.com/glvd/cluster/issues/587)
  * Add a `docker-compose.yml` template, which creates a two peer cluster | [glvd/cluster#585](https://github.com/glvd/cluster/issues/585) | [glvd/cluster#588](https://github.com/glvd/cluster/issues/588)
  * Support overwriting configuration values in the `cluster` section with environmental values | [glvd/cluster#575](https://github.com/glvd/cluster/issues/575) | [glvd/cluster#596](https://github.com/glvd/cluster/issues/596)
  * Set snaps to `classic` confinement mode and revert it since approval never arrived | [glvd/cluster#579](https://github.com/glvd/cluster/issues/579) | [glvd/cluster#594](https://github.com/glvd/cluster/issues/594)
* Use Go's reverse proxy library in the proxy endpoint | [glvd/cluster#570](https://github.com/glvd/cluster/issues/570) | [glvd/cluster#605](https://github.com/glvd/cluster/issues/605)


##### Bug fixes

  * `/add` endpoints improvements and IPFS Companion compatiblity | [glvd/cluster#582](https://github.com/glvd/cluster/issues/582) | [glvd/cluster#569](https://github.com/glvd/cluster/issues/569)
  * Fix adding with spaces in the name parameter | [glvd/cluster#583](https://github.com/glvd/cluster/issues/583)
  * Escape filter query parameter | [glvd/cluster#586](https://github.com/glvd/cluster/issues/586)
  * Fix some race conditions | [glvd/cluster#597](https://github.com/glvd/cluster/issues/597)
  * Improve pin deserialization efficiency | [glvd/cluster#601](https://github.com/glvd/cluster/issues/601)
  * Do not error remote pins | [glvd/cluster#600](https://github.com/glvd/cluster/issues/600) | [glvd/cluster#603](https://github.com/glvd/cluster/issues/603)
  * Clean up testing folders in `rest` and `rest/client` after tests | [glvd/cluster#607](https://github.com/glvd/cluster/issues/607)

#### Upgrading notices

##### Configuration changes

The configurations from previous versions are compatible, but a new `headers` key has been added to the `restapi` section. By default it gets CORS headers which will allow read-only interaction from any origin.

Additionally, all fields from the main `cluster` configuration section can now be overwrriten with environment variables. i.e. `CLUSTER_SECRET`, or  `CLUSTER_DISABLEREPINNING`.

##### REST API

The `/add` endpoint stream now returns different objects, in line with the rest of the API types.

Before:

```
type AddedOutput struct {
	Error
	Name  string
	Hash  string `json:",omitempty"`
	Bytes int64  `json:",omitempty"`
	Size  string `json:",omitempty"`
}
```

Now:

```
type AddedOutput struct {
	Name  string `json:"name"`
	Cid   string `json:"cid,omitempty"`
	Bytes uint64 `json:"bytes,omitempty"`
	Size  uint64 `json:"size,omitempty"`
}
```

The `/add` endpoint no longer reports errors as part of an AddedOutput object, but instead it uses trailer headers (same as `go-ipfs`). They are handled in the `client`.

##### Go APIs

The `AddedOutput` object has changed, thus the `api/rest/client` from older versions will not work with this one.

##### Other

No other things.

---

### v0.6.0 - 2018-10-03

#### Summary

IPFS version 0.6.0 is a new minor release of IPFS Cluster.

We have increased the minor release number to signal changes to the Go APIs after upgrading to the new `cid` package, but, other than that, this release does not include any major changes.

It brings a number of small fixes and features of which we can highlight two useful ones:

* the first is the support for multiple cluster daemon versions in the same cluster, as long as they share the same major/minor release. That means, all releases in the `0.6` series (`0.6.0`, `0.6.1` and so on...) will be able to speak among each others, allowing partial cluster upgrades.
* the second is the inclusion of a `PeerName` key in the status (`PinInfo`) objects. `ipfs-cluster-status` will now show peer names instead of peer IDs, making it easy to identify the status for each peer.

Many thanks to all the contributors to this release: @lanzafame, @meiqimichelle, @kishansagathiya, @cannium, @jglukasik and @mike-ngu.

#### List of changes

##### Features

  * Move commands to the `cmd/` folder | [glvd/cluster#485](https://github.com/glvd/cluster/issues/485) | [glvd/cluster#521](https://github.com/glvd/cluster/issues/521) | [glvd/cluster#556](https://github.com/glvd/cluster/issues/556)
  * Dependency upgrades: `go-dot`, `go-libp2p`, `cid` | [glvd/cluster#533](https://github.com/glvd/cluster/issues/533) | [glvd/cluster#537](https://github.com/glvd/cluster/issues/537) | [glvd/cluster#535](https://github.com/glvd/cluster/issues/535) | [glvd/cluster#544](https://github.com/glvd/cluster/issues/544) | [glvd/cluster#561](https://github.com/glvd/cluster/issues/561)
  * Build with go-1.11 | [glvd/cluster#558](https://github.com/glvd/cluster/issues/558)
  * Peer names in `PinInfo` | [glvd/cluster#446](https://github.com/glvd/cluster/issues/446) | [glvd/cluster#531](https://github.com/glvd/cluster/issues/531)
  * Wrap API client in an interface | [glvd/cluster#447](https://github.com/glvd/cluster/issues/447) | [glvd/cluster#523](https://github.com/glvd/cluster/issues/523) | [glvd/cluster#564](https://github.com/glvd/cluster/issues/564)
  * `Makefile`: add `prcheck` target and fix `make all` | [glvd/cluster#536](https://github.com/glvd/cluster/issues/536) | [glvd/cluster#542](https://github.com/glvd/cluster/issues/542) | [glvd/cluster#539](https://github.com/glvd/cluster/issues/539)
  * Docker: speed up [re]builds | [glvd/cluster#529](https://github.com/glvd/cluster/issues/529)
  * Re-enable keep-alives on servers | [glvd/cluster#548](https://github.com/glvd/cluster/issues/548) | [glvd/cluster#560](https://github.com/glvd/cluster/issues/560)

##### Bugfixes

  * Fix adding to cluster with unhealthy peers | [glvd/cluster#543](https://github.com/glvd/cluster/issues/543) | [glvd/cluster#549](https://github.com/glvd/cluster/issues/549)
  * Fix Snap builds and pushes: multiple architectures re-enabled | [glvd/cluster#520](https://github.com/glvd/cluster/issues/520) | [glvd/cluster#554](https://github.com/glvd/cluster/issues/554) | [glvd/cluster#557](https://github.com/glvd/cluster/issues/557) | [glvd/cluster#562](https://github.com/glvd/cluster/issues/562) | [glvd/cluster#565](https://github.com/glvd/cluster/issues/565)
  * Docs: Typos in Readme and some improvements | [glvd/cluster#547](https://github.com/glvd/cluster/issues/547) | [glvd/cluster#567](https://github.com/glvd/cluster/issues/567)
  * Fix tests in `stateless` PinTracker | [glvd/cluster#552](https://github.com/glvd/cluster/issues/552) | [glvd/cluster#563](https://github.com/glvd/cluster/issues/563)

#### Upgrading notices

##### Configuration changes

There are no changes to the configuration file on this release.

##### REST API

There are no changes to the REST API.

##### Go APIs

We have upgraded to the new version of the `cid` package. This means all `*cid.Cid` arguments are now `cid.Cid`.

##### Other

We are now using `go-1.11` to build and test cluster. We recommend using this version as well when building from source.

---


### v0.5.0 - 2018-08-23

#### Summary

IPFS Cluster version 0.5.0 is a minor release which includes a major feature: **adding content to IPFS directly through Cluster**.

This functionality is provided by `ipfs-cluster-ctl add` and by the API endpoint `/add`. The upload format (multipart) is similar to the IPFS `/add` endpoint, as well as the options (chunker, layout...). Cluster `add` generates the same DAG as `ipfs add` would, but it sends the added blocks directly to their allocations, pinning them on completion. The pin happens very quickly, as content is already locally available in the allocated peers.

The release also includes most of the needed code for the [Sharding feature](https://cluster.ipfs.io/developer/rfcs/dag-sharding-rfc/), but it is not yet usable/enabled, pending features from go-ipfs.

The 0.5.0 release additionally includes a new experimental PinTracker implementation: the `stateless` pin tracker. The stateless pin tracker relies on the IPFS pinset and the cluster state to keep track of pins, rather than keeping an in-memory copy of the cluster pinset, thus reducing the memory usage when having huge pinsets. It can be enabled with `ipfs-cluster-service daemon --pintracker stateless`.

The last major feature is the use of a DHT as routing layer for cluster peers. This means that peers should be able to discover each others as long as they are connected to one cluster peer. This simplifies the setup requirements for starting a cluster and helps avoiding situations which make the cluster unhealthy.

This release requires a state upgrade migration. It can be performed with `ipfs-cluster-service state upgrade` or simply launching the daemon with `ipfs-cluster-service daemon --upgrade`.

#### List of changes

##### Features

  * Libp2p upgrades (up to v6) | [glvd/cluster#456](https://github.com/glvd/cluster/issues/456) | [glvd/cluster#482](https://github.com/glvd/cluster/issues/482)
  * Support `/dns` multiaddresses for `node_multiaddress` | [glvd/cluster#462](https://github.com/glvd/cluster/issues/462) | [glvd/cluster#463](https://github.com/glvd/cluster/issues/463)
  * Increase `state_sync_interval` to 10 minutes | [glvd/cluster#468](https://github.com/glvd/cluster/issues/468) | [glvd/cluster#469](https://github.com/glvd/cluster/issues/469)
  * Auto-interpret libp2p addresses in `rest/client`'s `APIAddr` configuration option | [glvd/cluster#498](https://github.com/glvd/cluster/issues/498)
  * Resolve `APIAddr` (for `/dnsaddr` usage) in `rest/client` | [glvd/cluster#498](https://github.com/glvd/cluster/issues/498)
  * Support for adding content to Cluster and sharding (sharding is disabled) | [glvd/cluster#484](https://github.com/glvd/cluster/issues/484) | [glvd/cluster#503](https://github.com/glvd/cluster/issues/503) | [glvd/cluster#495](https://github.com/glvd/cluster/issues/495) | [glvd/cluster#504](https://github.com/glvd/cluster/issues/504) | [glvd/cluster#509](https://github.com/glvd/cluster/issues/509) | [glvd/cluster#511](https://github.com/glvd/cluster/issues/511) | [glvd/cluster#518](https://github.com/glvd/cluster/issues/518)
  * `stateless` PinTracker [glvd/cluster#308](https://github.com/glvd/cluster/issues/308) | [glvd/cluster#460](https://github.com/glvd/cluster/issues/460)
  * Add `size-only=true` to `repo/stat` calls | [glvd/cluster#507](https://github.com/glvd/cluster/issues/507)
  * Enable DHT-based peer discovery and routing for cluster peers | [glvd/cluster#489](https://github.com/glvd/cluster/issues/489) | [glvd/cluster#508](https://github.com/glvd/cluster/issues/508)
  * Gx-go upgrade | [glvd/cluster#517](https://github.com/glvd/cluster/issues/517)

##### Bugfixes

  * Fix type for constants | [glvd/cluster#455](https://github.com/glvd/cluster/issues/455)
  * Gofmt fix | [glvd/cluster#464](https://github.com/glvd/cluster/issues/464)
  * Fix tests for forked repositories | [glvd/cluster#465](https://github.com/glvd/cluster/issues/465) | [glvd/cluster#472](https://github.com/glvd/cluster/issues/472)
  * Fix resolve panic on `rest/client` | [glvd/cluster#498](https://github.com/glvd/cluster/issues/498)
  * Fix remote pins stuck in error state | [glvd/cluster#500](https://github.com/glvd/cluster/issues/500) | [glvd/cluster#460](https://github.com/glvd/cluster/issues/460)
  * Fix running some tests with `-race` | [glvd/cluster#340](https://github.com/glvd/cluster/issues/340) | [glvd/cluster#458](https://github.com/glvd/cluster/issues/458)
  * Fix ipfs proxy `/add` endpoint | [glvd/cluster#495](https://github.com/glvd/cluster/issues/495) | [glvd/cluster#81](https://github.com/glvd/cluster/issues/81) | [glvd/cluster#505](https://github.com/glvd/cluster/issues/505)
  * Fix ipfs proxy not hijacking `repo/stat` | [glvd/cluster#466](https://github.com/glvd/cluster/issues/466) | [glvd/cluster#514](https://github.com/glvd/cluster/issues/514)
  * Fix some godoc comments | [glvd/cluster#519](https://github.com/glvd/cluster/issues/519)

#### Upgrading notices

##### Configuration files

**IMPORTANT**: `0s` is the new default for the `read_timeout` and `write_timeout` values in the `restapi` configuration section, as well as `proxy_read_timeout` and `proxy_write_timeout` options in the `ipfshttp` section. Adding files to cluster (via the REST api or the proxy) is likely to timeout otherwise.

The `peerstore` file (in the configuration folder), no longer requires listing the multiaddresses for all cluster peers when initializing the cluster with a fixed peerset. It only requires the multiaddresses for one other cluster peer. The rest will be inferred using the DHT. The peerstore file is updated only on clean shutdown, and will store all known multiaddresses, even if not pertaining to cluster peers.

The new `stateless` PinTracker implementation uses a new configuration subsection in the `pin_tracker` key. This is only generated with `ipfs-cluster-service init`. When not present, a default configuration will be used (and a warning printed).

The `state_sync_interval` default has been increased to 10 minutes, as frequent syncing is not needed with the improvements in the PinTracker. Users are welcome to update this setting.


##### REST API

The `/add` endpoint has been added. The `replication_factor_min` and `replication_factor_max` options (in `POST allocations/<cid>`) have been deprecated and subsititued for `replication-min` and `replication-max`, although backwards comaptibility is kept.

Keep Alive has been disabled for the HTTP servers, as a bug in Go's HTTP client implementation may result adding corrupted content (and getting corrupted DAGs). However, while the libp2p API endpoint also suffers this, it will only close libp2p streams. Thus the performance impact on the libp2p-http endpoint should be minimal.

##### Go APIs

The `Config.PeerAddr` key in the `rest/client` module is deprecated. `APIAddr` should be used for both HTTP and LibP2P API endpoints. The type of address is automatically detected.

The IPFSConnector `Pin` call now receives an integer instead of a `Recursive` flag. It indicates the maximum depth to which something should be pinned. The only supported value is `-1` (meaning recursive). `BlockGet` and `BlockPut` calls have been added to the IPFSConnector component.

##### Other

As noted above, upgrade to `state` format version 5 is needed before starting the cluster service.

---

### v0.4.0 - 2018-05-30

#### Summary

The IPFS Cluster version 0.4.0 includes breaking changes and a considerable number of new features causing them. The documentation (particularly that affecting the configuration and startup of peers) has been updated accordingly in https://cluster.ipfs.io . Be sure to also read it if you are upgrading.

There are four main developments in this release:

* Refactorings around the `consensus` component, removing dependencies to the main component and allowing separate initialization: this has prompted to re-approach how we handle the peerset, the peer addresses and the peer's startup when using bootstrap. We have gained finer control of Raft, which has allowed us to provide a clearer configuration and a better start up procedure, specially when bootstrapping. The configuration file no longer mutates while cluster is running.
* Improvements to the `pintracker`: our pin tracker is now able to cancel ongoing pins when receiving an unpin request for the same CID, and vice-versa. It will also optimize multiple pin requests (by only queuing and triggering them once) and can now report
whether an item is pinning (a request to ipfs is ongoing) vs. pin-queued (waiting for a worker to perform the request to ipfs).
* Broadcasting of monitoring metrics using PubSub: we have added a new `monitor` implementation that uses PubSub (rather than RPC broadcasting). With the upcoming improvements to PubSub this means that we can do efficient broadcasting of metrics while at the same time not requiring peers to have RPC permissions, which is preparing the ground for collaborative clusters.
* We have launched the IPFS Cluster website: https://cluster.ipfs.io . We moved most of the documentation over there, expanded it and updated it.

#### List of changes

##### Features

  * Consensus refactorings | [glvd/cluster#398](https://github.com/glvd/cluster/issues/398) | [glvd/cluster#371](https://github.com/glvd/cluster/issues/371)
  * Pintracker revamp | [glvd/cluster#308](https://github.com/glvd/cluster/issues/308) | [glvd/cluster#383](https://github.com/glvd/cluster/issues/383) | [glvd/cluster#408](https://github.com/glvd/cluster/issues/408) | [glvd/cluster#415](https://github.com/glvd/cluster/issues/415) | [glvd/cluster#421](https://github.com/glvd/cluster/issues/421) | [glvd/cluster#427](https://github.com/glvd/cluster/issues/427) | [glvd/cluster#432](https://github.com/glvd/cluster/issues/432)
  * Pubsub monitoring | [glvd/cluster#400](https://github.com/glvd/cluster/issues/400)
  * Force killing cluster with double CTRL-C | [glvd/cluster#258](https://github.com/glvd/cluster/issues/258) | [glvd/cluster#358](https://github.com/glvd/cluster/issues/358)
  * 3x faster testsuite | [glvd/cluster#339](https://github.com/glvd/cluster/issues/339) | [glvd/cluster#350](https://github.com/glvd/cluster/issues/350)
  * Introduce `disable_repinning` option | [glvd/cluster#369](https://github.com/glvd/cluster/issues/369) | [glvd/cluster#387](https://github.com/glvd/cluster/issues/387)
  * Documentation moved to website and fixes | [glvd/cluster#390](https://github.com/glvd/cluster/issues/390) | [glvd/cluster#391](https://github.com/glvd/cluster/issues/391) | [glvd/cluster#393](https://github.com/glvd/cluster/issues/393) | [glvd/cluster#347](https://github.com/glvd/cluster/issues/347)
  * Run Docker container with `daemon --upgrade` by default | [glvd/cluster#394](https://github.com/glvd/cluster/issues/394)
  * Remove the `ipfs-cluster-ctl peers add` command (bootstrap should be used to add peers) | [glvd/cluster#397](https://github.com/glvd/cluster/issues/397)
  * Add tests using HTTPs endpoints | [glvd/cluster#191](https://github.com/glvd/cluster/issues/191) | [glvd/cluster#403](https://github.com/glvd/cluster/issues/403)
  * Set `refs` as default `pinning_method` and `10` as default `concurrent_pins` | [glvd/cluster#420](https://github.com/glvd/cluster/issues/420)
  * Use latest `gx` and `gx-go`. Be more verbose when installing | [glvd/cluster#418](https://github.com/glvd/cluster/issues/418)
  * Makefile: Properly retrigger builds on source change | [glvd/cluster#426](https://github.com/glvd/cluster/issues/426)
  * Improvements to StateSync() | [glvd/cluster#429](https://github.com/glvd/cluster/issues/429)
  * Rename `ipfs-cluster-data` folder to `raft` | [glvd/cluster#430](https://github.com/glvd/cluster/issues/430)
  * Officially support go 1.10 | [glvd/cluster#439](https://github.com/glvd/cluster/issues/439)
  * Update to libp2p 5.0.17 | [glvd/cluster#440](https://github.com/glvd/cluster/issues/440)

##### Bugsfixes:

  * Don't keep peers /ip*/ addresses if we know DNS addresses for them | [glvd/cluster#381](https://github.com/glvd/cluster/issues/381)
  * Running cluster with wrong configuration path gives misleading error | [glvd/cluster#343](https://github.com/glvd/cluster/issues/343) | [glvd/cluster#370](https://github.com/glvd/cluster/issues/370) | [glvd/cluster#373](https://github.com/glvd/cluster/issues/373)
  * Do not fail when running with `daemon --upgrade` and no state is present | [glvd/cluster#395](https://github.com/glvd/cluster/issues/395)
  * IPFS Proxy: handle arguments passed as part of the url | [glvd/cluster#380](https://github.com/glvd/cluster/issues/380) | [glvd/cluster#392](https://github.com/glvd/cluster/issues/392)
  * WaitForUpdates() may return before state is fully synced | [glvd/cluster#378](https://github.com/glvd/cluster/issues/378)
  * Configuration mutates no more and shadowing is no longer necessary | [glvd/cluster#235](https://github.com/glvd/cluster/issues/235)
  * Govet fixes | [glvd/cluster#417](https://github.com/glvd/cluster/issues/417)
  * Fix release changelog when having RC tags
  * Fix lock file not being removed on cluster force-kill | [glvd/cluster#423](https://github.com/glvd/cluster/issues/423) | [glvd/cluster#437](https://github.com/glvd/cluster/issues/437)
  * Fix indirect pins not being correctly parsed | [glvd/cluster#428](https://github.com/glvd/cluster/issues/428) | [glvd/cluster#436](https://github.com/glvd/cluster/issues/436)
  * Enable NAT support in libp2p host | [glvd/cluster#346](https://github.com/glvd/cluster/issues/346) | [glvd/cluster#441](https://github.com/glvd/cluster/issues/441)
  * Fix pubsub monitor not working on ARM | [glvd/cluster#433](https://github.com/glvd/cluster/issues/433) | [glvd/cluster#443](https://github.com/glvd/cluster/issues/443)

#### Upgrading notices

##### Configuration file

This release introduces **breaking changes to the configuration file**. An error will be displayed if `ipfs-cluster-service` is started with an old configuration file. We recommend re-initing the configuration file altogether.

* The `peers` and `bootstrap` keys have been removed from the main section of the configuration
* You might need to provide Peer multiaddresses in a text file named `peerstore`, in your `~/.ipfs-cluster` folder (one per line). This allows your peers how to contact other peers.
* A `disable_repinning` option has been added to the main configuration section. Defaults to `false`.
* A `init_peerset` has been added to the `raft` configuration section. It should be used to define the starting set of peers when a cluster starts for the first time and is not bootstrapping to an existing running peer (otherwise it is ignored). The value is an array of peer IDs.
* A `backups_rotate` option has been added to the `raft` section and specifies how many copies of the Raft state to keep as backups when the state is cleaned up.
* An `ipfs_request_timeout` option has been introduced to the `ipfshttp` configuration section, and controls the timeout of general requests to the ipfs daemon. Defaults to 5 minutes.
* A `pin_timeout` option has been introduced to the `ipfshttp` section, it controls the timeout for Pin requests to ipfs. Defaults to 24 hours.
* An `unpin_timeout` option has been introduced to the `ipfshttp` section. it controls the timeout for Unpin requests to ipfs. Defaults to 3h.
* Both `pinning_timeout` and `unpinning_timeout` options have been removed from the `maptracker` section.
* A `monitor/pubsubmon` section configures the new PubSub monitoring component. The section is identical to the existing `monbasic`, its only option being `check_interval` (defaults to 15 seconds).

The `ipfs-cluster-data` folder has been renamed to `raft`. Upon `ipfs-cluster-service daemon` start, the renaming will happen automatically if it exists. Otherwise it will be created with the new name.

##### REST API

There are no changes to REST APIs in this release.

##### Go APIs

Several component APIs have changed: `Consensus`, `PeerMonitor` and `IPFSConnector` have added new methods or changed methods signatures.

##### Other

Calling `ipfs-cluster-service` without subcommands no longer runs the peer. It is necessary to call `ipfs-cluster-service daemon`. Several daemon-specific flags have been made subcommand flags: `--bootstrap` and `--alloc`.

The `--bootstrap` flag can now take a list of comma-separated multiaddresses. Using `--bootstrap` will automatically run `state clean`.

The `ipfs-cluster-ctl` no longer has a `peers add` subcommand. Peers should not be added this way, but rather bootstrapped to an existing running peer.

---

### v0.3.5 - 2018-03-29

This release comes full with new features. The biggest ones are the support for parallel pinning (using `refs -r` rather than `pin add` to pin things in IPFS), and the exposing of the http endpoints through libp2p. This allows users to securely interact with the HTTP API without having to setup SSL certificates.

* Features
  * `--no-status` for `ipfs-cluster-ctl pin add/rm` allows to speed up adding and removing by not fetching the status one second afterwards. Useful for ingesting pinsets to cluster | [glvd/cluster#286](https://github.com/glvd/cluster/issues/286) | [glvd/cluster#329](https://github.com/glvd/cluster/issues/329)
  * `--wait` flag for `ipfs-cluster-ctl pin add/rm` allows to wait until a CID is fully pinned or unpinned [glvd/cluster#338](https://github.com/glvd/cluster/issues/338) | [glvd/cluster#348](https://github.com/glvd/cluster/issues/348) | [glvd/cluster#363](https://github.com/glvd/cluster/issues/363)
  * Support `refs` pinning method. Parallel pinning | [glvd/cluster#326](https://github.com/glvd/cluster/issues/326) | [glvd/cluster#331](https://github.com/glvd/cluster/issues/331)
  * Double default timeouts for `ipfs-cluster-ctl` | [glvd/cluster#323](https://github.com/glvd/cluster/issues/323) | [glvd/cluster#334](https://github.com/glvd/cluster/issues/334)
  * Better error messages during startup | [glvd/cluster#167](https://github.com/glvd/cluster/issues/167) | [glvd/cluster#344](https://github.com/glvd/cluster/issues/344) | [glvd/cluster#353](https://github.com/glvd/cluster/issues/353)
  * REST API client now provides an `IPFS()` method which returns a `go-ipfs-api` shell instance pointing to the proxy endpoint | [glvd/cluster#269](https://github.com/glvd/cluster/issues/269) | [glvd/cluster#356](https://github.com/glvd/cluster/issues/356)
  * REST http-api-over-libp2p. Server, client, `ipfs-cluster-ctl` support added | [glvd/cluster#305](https://github.com/glvd/cluster/issues/305) | [glvd/cluster#349](https://github.com/glvd/cluster/issues/349)
  * Added support for priority pins and non-recursive pins (sharding-related) | [glvd/cluster#341](https://github.com/glvd/cluster/issues/341) | [glvd/cluster#342](https://github.com/glvd/cluster/issues/342)
  * Documentation fixes | [glvd/cluster#328](https://github.com/glvd/cluster/issues/328) | [glvd/cluster#357](https://github.com/glvd/cluster/issues/357)

* Bugfixes
  * Print lock path in logs | [glvd/cluster#332](https://github.com/glvd/cluster/issues/332) | [glvd/cluster#333](https://github.com/glvd/cluster/issues/333)

There are no breaking API changes and all configurations should be backwards compatible. The `api/rest/client` provides a new `IPFS()` method.

We recommend updating the `service.json` configurations to include all the new configuration options:

* The `pin_method` option has been added to the `ipfshttp` section. It supports `refs` and `pin` (default) values. Use `refs` for parallel pinning, but only if you don't run automatic GC on your ipfs nodes.
* The `concurrent_pins` option has been added to the `maptracker` section. Only useful with `refs` option in `pin_method`.
* The `listen_multiaddress` option in the `restapi` section should be renamed to `http_listen_multiaddress`.

This release will require a **state upgrade**. Run `ipfs-cluster-service state upgrade` in all your peers, or start cluster with `ipfs-cluster-service daemon --upgrade`.

---

### v0.3.4 - 2018-02-20

This release fixes the pre-built binaries.

* Bugfixes
  * Pre-built binaries panic on start | [glvd/cluster#320](https://github.com/glvd/cluster/issues/320)

---

### v0.3.3 - 2018-02-12

This release includes additional `ipfs-cluster-service state` subcommands and the connectivity graph feature.

* Features
  * `ipfs-cluster-service daemon --upgrade` allows to automatically run migrations before starting | [glvd/cluster#300](https://github.com/glvd/cluster/issues/300) | [glvd/cluster#307](https://github.com/glvd/cluster/issues/307)
  * `ipfs-cluster-service state version` reports the shared state format version | [glvd/cluster#298](https://github.com/glvd/cluster/issues/298) | [glvd/cluster#307](https://github.com/glvd/cluster/issues/307)
  * `ipfs-cluster-service health graph` generates a .dot graph file of cluster connectivity | [glvd/cluster#17](https://github.com/glvd/cluster/issues/17) | [glvd/cluster#291](https://github.com/glvd/cluster/issues/291) | [glvd/cluster#311](https://github.com/glvd/cluster/issues/311)

* Bugfixes
  * Do not upgrade state if already up to date | [glvd/cluster#296](https://github.com/glvd/cluster/issues/296) | [glvd/cluster#307](https://github.com/glvd/cluster/issues/307)
  * Fix `ipfs-cluster-service daemon` failing with `unknown allocation strategy` error | [glvd/cluster#314](https://github.com/glvd/cluster/issues/314) | [glvd/cluster#315](https://github.com/glvd/cluster/issues/315)

APIs have not changed in this release. The `/health/graph` endpoint has been added.

---

### v0.3.2 - 2018-01-25

This release includes a number of bufixes regarding the upgrade and import of state, along with two important features:
  * Commands to export and import the internal cluster state: these allow to perform easy and human-readable dumps of the shared cluster state while offline, and eventually restore it in a different peer or cluster.
  * The introduction of `replication_factor_min` and `replication_factor_max` parameters for every Pin (along with the deprecation of `replication_factor`). The defaults are specified in the configuration. For more information on the usage and behavour of these new options, check the IPFS cluster guide.

* Features
  * New `ipfs-cluster-service state export/import/cleanup` commands | [glvd/cluster#240](https://github.com/glvd/cluster/issues/240) | [glvd/cluster#290](https://github.com/glvd/cluster/issues/290)
  * New min/max replication factor control | [glvd/cluster#277](https://github.com/glvd/cluster/issues/277) | [glvd/cluster#292](https://github.com/glvd/cluster/issues/292)
  * Improved migration code | [glvd/cluster#283](https://github.com/glvd/cluster/issues/283)
  * `ipfs-cluster-service version` output simplified (see below) | [glvd/cluster#274](https://github.com/glvd/cluster/issues/274)
  * Testing improvements:
    * Added tests for Dockerfiles | [glvd/cluster#200](https://github.com/glvd/cluster/issues/200) | [glvd/cluster#282](https://github.com/glvd/cluster/issues/282)
    * Enabled Jenkins testing and made it work | [glvd/cluster#256](https://github.com/glvd/cluster/issues/256) | [glvd/cluster#294](https://github.com/glvd/cluster/issues/294)
  * Documentation improvements:
    * Guide contains more details on state upgrade procedures | [glvd/cluster#270](https://github.com/glvd/cluster/issues/270)
    * ipfs-cluster-ctl exit status are documented on the README | [glvd/cluster#178](https://github.com/glvd/cluster/issues/178)

* Bugfixes
  * Force cleanup after sharness tests | [glvd/cluster#181](https://github.com/glvd/cluster/issues/181) | [glvd/cluster#288](https://github.com/glvd/cluster/issues/288)
  * Fix state version validation on start | [glvd/cluster#293](https://github.com/glvd/cluster/issues/293)
  * Wait until last index is applied before attempting snapshot on shutdown | [glvd/cluster#275](https://github.com/glvd/cluster/issues/275)
  * Snaps from master not pushed due to bad credentials
  * Fix overpinning or underpinning of CIDs after re-join | [glvd/cluster#222](https://github.com/glvd/cluster/issues/222)
  * Fix unmarshaling state on top of an existing one | [glvd/cluster#297](https://github.com/glvd/cluster/issues/297)
  * Fix catching up on imported state | [glvd/cluster#297](https://github.com/glvd/cluster/issues/297)

These release is compatible with previous versions of ipfs-cluster on the API level, with the exception of the `ipfs-cluster-service version` command, which returns `x.x.x-shortcommit` rather than `ipfs-cluster-service version 0.3.1`. The former output is still available as `ipfs-cluster-service --version`.

The `replication_factor` option is deprecated, but still supported and will serve as a shortcut to set both `replication_factor_min` and `replication_factor_max` to the same value. This affects the configuration file, the REST API and the `ipfs-cluster-ctl pin add` command.

---

### v0.3.1 - 2017-12-11

This release includes changes around the consensus state management, so that upgrades can be performed when the internal format changes. It also comes with several features and changes to support a live deployment and integration with IPFS pin-bot, including a REST API client for Go.

* Features
 * `ipfs-cluster-service state upgrade` | [glvd/cluster#194](https://github.com/glvd/cluster/issues/194)
 * `ipfs-cluster-test` Docker image runs with `ipfs:master` | [glvd/cluster#155](https://github.com/glvd/cluster/issues/155) | [glvd/cluster#259](https://github.com/glvd/cluster/issues/259)
 * `ipfs-cluster` Docker image only runs `ipfs-cluster-service` (and not the ipfs daemon anymore) | [glvd/cluster#197](https://github.com/glvd/cluster/issues/197) | [glvd/cluster#155](https://github.com/glvd/cluster/issues/155) | [glvd/cluster#259](https://github.com/glvd/cluster/issues/259)
 * Support for DNS multiaddresses for cluster peers | [glvd/cluster#155](https://github.com/glvd/cluster/issues/155) | [glvd/cluster#259](https://github.com/glvd/cluster/issues/259)
 * Add configuration section and options for `pin_tracker` | [glvd/cluster#155](https://github.com/glvd/cluster/issues/155) | [glvd/cluster#259](https://github.com/glvd/cluster/issues/259)
 * Add `local` flag to Status, Sync, Recover endpoints which allows to run this operations only in the peer receiving the request | [glvd/cluster#155](https://github.com/glvd/cluster/issues/155) | [glvd/cluster#259](https://github.com/glvd/cluster/issues/259)
 * Add Pin names | [glvd/cluster#249](https://github.com/glvd/cluster/issues/249)
 * Add Peer names | [glvd/cluster#250](https://github.com/glvd/cluster/issues/250)
 * New REST API Client module `github.com/glvd/cluster/api/rest/client` allows to integrate against cluster | [glvd/cluster#260](https://github.com/glvd/cluster/issues/260) | [glvd/cluster#263](https://github.com/glvd/cluster/issues/263) | [glvd/cluster#266](https://github.com/glvd/cluster/issues/266)
 * A few rounds addressing code quality issues | [glvd/cluster#264](https://github.com/glvd/cluster/issues/264)

This release should stay backwards compatible with the previous one. Nevertheless, some REST API endpoints take the `local` flag, and matching new Go public functions have been added (`RecoverAllLocal`, `SyncAllLocal`...).

---

### v0.3.0 - 2017-11-15

This release introduces Raft 1.0.0 and incorporates deep changes to the management of the cluster peerset.

* Features
  * Upgrade Raft to 1.0.0 | [glvd/cluster#194](https://github.com/glvd/cluster/issues/194) | [glvd/cluster#196](https://github.com/glvd/cluster/issues/196)
  * Support Snaps | [glvd/cluster#234](https://github.com/glvd/cluster/issues/234) | [glvd/cluster#228](https://github.com/glvd/cluster/issues/228) | [glvd/cluster#232](https://github.com/glvd/cluster/issues/232)
  * Rotating backups for ipfs-cluster-data | [glvd/cluster#233](https://github.com/glvd/cluster/issues/233)
  * Bring documentation up to date with the code [glvd/cluster#223](https://github.com/glvd/cluster/issues/223)

Bugfixes:
  * Fix docker startup | [glvd/cluster#216](https://github.com/glvd/cluster/issues/216) | [glvd/cluster#217](https://github.com/glvd/cluster/issues/217)
  * Fix configuration save | [glvd/cluster#213](https://github.com/glvd/cluster/issues/213) | [glvd/cluster#214](https://github.com/glvd/cluster/issues/214)
  * Forward progress updates with IPFS-Proxy | [glvd/cluster#224](https://github.com/glvd/cluster/issues/224) | [glvd/cluster#231](https://github.com/glvd/cluster/issues/231)
  * Delay ipfs connect swarms on boot and safeguard against panic condition | [glvd/cluster#238](https://github.com/glvd/cluster/issues/238)
  * Multiple minor fixes | [glvd/cluster#236](https://github.com/glvd/cluster/issues/236)
    * Avoid shutting down consensus in the middle of a commit
    * Return an ID containing current peers in PeerAdd
    * Do not shut down libp2p host in the middle of peer removal
    * Send cluster addresses to the new peer before adding it
    * Wait for configuration save on init
    * Fix error message when not enough allocations exist for a pin

This releases introduces some changes affecting the configuration file and some breaking changes affecting `go` and the REST APIs:

* The `consensus.raft` section of the configuration has new options but should be backwards compatible.
* The `Consensus` component interface has changed, `LogAddPeer` and `LogRmPeer` have been replaced by `AddPeer` and `RmPeer`. It additionally provides `Clean` and `Peers` methods. The `consensus/raft` implementation has been updated accordingly.
* The `api.ID` (used in REST API among others) object key `ClusterPeers` key is now a list of peer IDs, and not a list of multiaddresses as before. The object includes a new key `ClusterPeersAddresses` which includes the multiaddresses.
* Note that `--bootstrap` and `--leave` flags when calling `ipfs-cluster-service` will be stored permanently in the configuration (see [glvd/cluster#235](https://github.com/glvd/cluster/issues/235)).

---

### v0.2.1 - 2017-10-26

This is a maintenance release with some important bugfixes.

* Fixes:
  * Dockerfile runs `ipfs-cluster-service` instead of `ctl` | [glvd/cluster#194](https://github.com/glvd/cluster/issues/194) | [glvd/cluster#196](https://github.com/glvd/cluster/issues/196)
  * Peers and bootstrap entries in the configuration are ignored | [glvd/cluster#203](https://github.com/glvd/cluster/issues/203) | [glvd/cluster#204](https://github.com/glvd/cluster/issues/204)
  * Informers do not work on 32-bit architectures | [glvd/cluster#202](https://github.com/glvd/cluster/issues/202) | [glvd/cluster#205](https://github.com/glvd/cluster/issues/205)
  * Replication factor entry in the configuration is ignored | [glvd/cluster#208](https://github.com/glvd/cluster/issues/208) | [glvd/cluster#209](https://github.com/glvd/cluster/issues/209)

The fix for 32-bit architectures has required a change in the `IPFSConnector` interface (`FreeSpace()` and `Reposize()` return `uint64` now). The current implementation by the `ipfshttp` module has changed accordingly.


---

### v0.2.0 - 2017-10-23

* Features:
  * Basic authentication support added to API component | [glvd/cluster#121](https://github.com/glvd/cluster/issues/121) | [glvd/cluster#147](https://github.com/glvd/cluster/issues/147) | [glvd/cluster#179](https://github.com/glvd/cluster/issues/179)
  * Copy peers to bootstrap when leaving a cluster | [glvd/cluster#170](https://github.com/glvd/cluster/issues/170) | [glvd/cluster#112](https://github.com/glvd/cluster/issues/112)
  * New configuration format | [glvd/cluster#162](https://github.com/glvd/cluster/issues/162) | [glvd/cluster#177](https://github.com/glvd/cluster/issues/177)
  * Freespace disk metric implementation. It's now the default. | [glvd/cluster#142](https://github.com/glvd/cluster/issues/142) | [glvd/cluster#99](https://github.com/glvd/cluster/issues/99)

* Fixes:
  * IPFS Connector should use only POST | [glvd/cluster#176](https://github.com/glvd/cluster/issues/176) | [glvd/cluster#161](https://github.com/glvd/cluster/issues/161)
  * `ipfs-cluster-ctl` exit status with error responses | [glvd/cluster#174](https://github.com/glvd/cluster/issues/174)
  * Sharness tests and update testing container | [glvd/cluster#171](https://github.com/glvd/cluster/issues/171)
  * Update Dockerfiles | [glvd/cluster#154](https://github.com/glvd/cluster/issues/154) | [glvd/cluster#185](https://github.com/glvd/cluster/issues/185)
  * `ipfs-cluster-service`: Do not run service with unknown subcommands | [glvd/cluster#186](https://github.com/glvd/cluster/issues/186)

This release introduces some breaking changes affecting configuration files and `go` integrations:

* Config: The old configuration format is no longer valid and cluster will fail to start from it. Configuration file needs to be re-initialized with `ipfs-cluster-service init`.
* Go: The `restapi` component has been renamed to `rest` and some of its public methods have been renamed.
* Go: Initializers (`New<Component>(...)`) for most components have changed to accept a `Config` object. Some initializers have been removed.

---

Note, when adding changelog entries, write links to issues as `@<issuenumber>` and then replace them with links with the following command:

```
sed -i -r 's/@([0-9]+)/[ipfs\/ipfs-cluster#\1](https:\/\/github.com\/ipfs\/ipfs-cluster\/issues\/\1)/g' CHANGELOG.md
```
