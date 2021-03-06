# metric-tank

is a multi-tenant timeries metrics data base. (aka TSDB)



# http interface

## data querying
* `http://localhost:6063/get` either POST or GET, with the following parameters:
  * `target` mandatory. can be specified multiple times to request several series. Supported formats:
    * simply the raw id of a metric. like `1.2345foobar`
    * `consolidateBy(<id>,'<function>')`. single quotes only. accepted functions are avg, average, last, min, max, sum.
       example: `consolidateBy(1.2345foobar,'average')`.
  * `maxDataPoints`: max points to be returned. runs runtime consolidation when needed. optional
  * `from` and `to` unix timestamps. optional
    * from is inclusive, to is exclusive. you can also use 'until' but to takes precedence.
    * so from=x, to=y returns data that can include x and y-1 but not y.
    * from defaults to now-24h, to to now+1.
    * from can also be a human friendly pattern like -10min or -7d

* the response will id the series by the target used to request them

note:
* it just serves up the data that it has, in timestamp ascending order. it does no effort to try to fill in gaps.
* no support for wildcards, patterns, "magic" time specs like "-10min" etc.
* it is assumed that authorisation (by org-id) has already been performed.  (the graphite-raintank plugin does this)

## other useful endpoints exemplified through curl commands:

* `curl http://localhost:6063/` app status (OK if either primary or secondary that has been warmed up). good for loadbalancers.
* `curl http://localhost:6063/cluster` cluster status
* `curl -X POST -d primary=false http://localhost:6063/cluster` set primary true/false


# aggregations

MT can save various bands of aggregated data, using multiple consolidation functions per series. this works seamlessly with consolidateBy, unlike graphite.

TODO: you can currently write fake metrics with same key as aggregated metrics, which would conflict, we should probably blacklist such patterns

# clustering

run one primary which writes to cassandra
when one primary is down you need to be careful about when to promote a secondary to primary:

* after you see the "starting data consumption" log message for a primary, data consomuption starts. this timestamp is important.
* look at your largest chunkSpan. secondary can only be promoted when a new interval starts for the largest chunkSpan. intervals start when clock unix timestamp divides without remainder by chunkSpan. How long you should wait is also shown (in seconds) via the `cluster.promotion_wait` metric.
* of course there are other factors: any running primary should be depromoted and have saved its data to cassandra, all metricPersist message should have made it through NSQ into the about-to-be-promoted instance.


# graphite-api

* `/render` has a very, very limited subset of the graphite render api. basically you can specify targets by their graphite key, set from, to and maxDataPoints, and use consolidateBy.
No other function or parameter is currently supported.  Also we don't check org-id so don't expose this publically
* `/metrics/index.json` is like graphite.  Don't expose this publically


## design limitations to address at some point:

* NSQ does not provide ordering guarantees, we need ordering for optimal compression, aggregations. currently we drop out of order points which may result in gaps.
see https://github.com/raintank/raintank-metric/issues/41 for more info. also [it may also affect alerting](https://github.com/raintank/raintank-metric/issues/17). we're looking into kafka.


* rollups is a bit clunky:
  - for simplicity just reuses AggMetric but this is not a good fit. it keeps too many string id's in memory, too much Sprintf overhead.
  - also per-target-type aggregations (like counter -> last), not all aggregations always make sense for all types.
  - no need to take all raw inputs into each aggregator, they can instead take summaries from previous aggregators
  we should redo them at some point. 

* we don't have a list of all keys inside the tsdb. consequences: you can't get lists/search/autocomplete and for non-existant keys we still query cassandra


## index design

metric definitions are currently stored in ES as well as internally (other options can come later).
ES is the failsafe option used by graphite-raintank.py and such.
The index is used internally for the graphite-api and is experimental.  It's powered by a radix tree and trigram index.

note that any given metric may appear multiple times, under different organisations

definition id's are unique across the entire system and can be computed, so don't require coordination across distributed nodes.

there can be multiple definitions for each metric, if the interval changes for example
currently those all just stored individually in the radix tree and trigram index, which is a bit redundant
in the future, we might just index the metric names and then have a separate structure to resolve a name to its multiple metricdefs, which could be cheaper.
