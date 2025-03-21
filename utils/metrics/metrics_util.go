package metrics

import (
	"errors"
	"math/big"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/x/evm/types"
	metrics "github.com/hashicorp/go-metrics"
)

// Measures the time taken to execute a sudo msg
// Metric Names:
//
//	eni_sudo_duration_miliseconds
//	eni_sudo_duration_miliseconds_count
//	eni_sudo_duration_miliseconds_sum
func MeasureSudoExecutionDuration(start time.Time, msgType string) {
	metrics.MeasureSinceWithLabels(
		[]string{"eni", "sudo", "duration", "milliseconds"},
		start.UTC(),
		[]metrics.Label{telemetry.NewLabel("type", msgType)},
	)
}

// Measures failed sudo execution count
// Metric Name:
//
//	eni_sudo_error_count
func IncrementSudoFailCount(msgType string) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "sudo", "error", "count"},
		1,
		[]metrics.Label{telemetry.NewLabel("type", msgType)},
	)
}

// Gauge metric with enid version and git commit as labels
// Metric Name:
//
//	enid_version_and_commit
func GaugeenidVersionAndCommit(version string, commit string) {
	telemetry.SetGaugeWithLabels(
		[]string{"enid_version_and_commit"},
		1,
		[]metrics.Label{telemetry.NewLabel("enid_version", version), telemetry.NewLabel("commit", commit)},
	)
}

// eni_tx_process_type_count
func IncrTxProcessTypeCounter(processType string) {
	metrics.IncrCounterWithLabels(
		[]string{"eni", "tx", "process", "type"},
		1,
		[]metrics.Label{telemetry.NewLabel("type", processType)},
	)
}

// Measures the time taken to process a block by the process type
// Metric Names:
//
//	eni_process_block_miliseconds
//	eni_process_block_miliseconds_count
//	eni_process_block_miliseconds_sum
func BlockProcessLatency(start time.Time, processType string) {
	metrics.MeasureSinceWithLabels(
		[]string{"eni", "process", "block", "milliseconds"},
		start.UTC(),
		[]metrics.Label{telemetry.NewLabel("type", processType)},
	)
}

// Measures the time taken to execute a sudo msg
// Metric Names:
//
//	eni_tx_process_type_count
func IncrDagBuildErrorCounter(reason string) {
	metrics.IncrCounterWithLabels(
		[]string{"eni", "dag", "build", "error"},
		1,
		[]metrics.Label{telemetry.NewLabel("reason", reason)},
	)
}

// Counts the number of concurrent transactions that failed
// Metric Names:
//
//	eni_tx_concurrent_delivertx_error
func IncrFailedConcurrentDeliverTxCounter() {
	metrics.IncrCounterWithLabels(
		[]string{"eni", "tx", "concurrent", "delievertx", "error"},
		1,
		[]metrics.Label{},
	)
}

// Counts the number of operations that failed due to operation timeout
// Metric Names:
//
//	eni_log_not_done_after_counter
func IncrLogIfNotDoneAfter(label string) {
	metrics.IncrCounterWithLabels(
		[]string{"eni", "log", "not", "done", "after"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("label", label),
		},
	)
}

// Measures the time taken to execute a sudo msg
// Metric Names:
//
//	eni_deliver_tx_duration_miliseconds
//	eni_deliver_tx_duration_miliseconds_count
//	eni_deliver_tx_duration_miliseconds_sum
func MeasureDeliverTxDuration(start time.Time) {
	metrics.MeasureSince(
		[]string{"eni", "deliver", "tx", "milliseconds"},
		start.UTC(),
	)
}

// Measures the time taken to execute a batch tx
// Metric Names:
//
//	eni_deliver_batch_tx_duration_miliseconds
//	eni_deliver_batch_tx_duration_miliseconds_count
//	eni_deliver_batch_tx_duration_miliseconds_sum
func MeasureDeliverBatchTxDuration(start time.Time) {
	metrics.MeasureSince(
		[]string{"eni", "deliver", "batch", "tx", "milliseconds"},
		start.UTC(),
	)
}

// eni_oracle_vote_penalty_count
func SetOracleVotePenaltyCount(count uint64, valAddr string, penaltyType string) {
	metrics.SetGaugeWithLabels(
		[]string{"eni", "oracle", "vote", "penalty", "count"},
		float32(count),
		[]metrics.Label{
			telemetry.NewLabel("type", penaltyType),
			telemetry.NewLabel("validator", valAddr),
		},
	)
}

// eni_epoch_new
func SetEpochNew(epochNum uint64) {
	metrics.SetGauge(
		[]string{"eni", "epoch", "new"},
		float32(epochNum),
	)
}

// Measures throughput
// Metric Name:
//
//	eni_throughput_<metric_name>
func SetThroughputMetric(metricName string, value float32) {
	telemetry.SetGauge(
		value,
		"eni", "throughput", metricName,
	)
}

// Measures number of new websocket connects
// Metric Name:
//
//	eni_websocket_connect
func IncWebsocketConnects() {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "websocket", "connect"},
		1,
		nil,
	)
}

// Measures number of times a denom's price is updated
// Metric Name:
//
//	eni_oracle_price_update_count
func IncrPriceUpdateDenom(denom string) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "oracle", "price", "update"},
		1,
		[]metrics.Label{telemetry.NewLabel("denom", denom)},
	)
}

// Measures throughput per message type
// Metric Name:
//
//	eni_throughput_<metric_name>
func SetThroughputMetricByType(metricName string, value float32, msgType string) {
	telemetry.SetGaugeWithLabels(
		[]string{"eni", "loadtest", "tps", metricName},
		value,
		[]metrics.Label{telemetry.NewLabel("msg_type", msgType)},
	)
}

// Measures the number of times the total block gas wanted in the proposal exceeds the max
// Metric Name:
//
//	eni_failed_total_gas_wanted_check
func IncrFailedTotalGasWantedCheck(proposer string) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "failed", "total", "gas", "wanted", "check"},
		1,
		[]metrics.Label{telemetry.NewLabel("proposer", proposer)},
	)
}

// Measures the number of times the total block gas wanted in the proposal exceeds the max
// Metric Name:
//
//	eni_failed_total_gas_wanted_check
func IncrValidatorSlashed(proposer string) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "failed", "total", "gas", "wanted", "check"},
		1,
		[]metrics.Label{telemetry.NewLabel("proposer", proposer)},
	)
}

// Measures number of times a denom's price is updated
// Metric Name:
//
//	eni_oracle_price_update_count
func SetCoinsMinted(amount uint64, denom string) {
	telemetry.SetGaugeWithLabels(
		[]string{"eni", "mint", "coins"},
		float32(amount),
		[]metrics.Label{telemetry.NewLabel("denom", denom)},
	)
}

// Measures the number of times the total block gas wanted in the proposal exceeds the max
// Metric Name:
//
//	eni_tx_gas_counter
func IncrGasCounter(gasType string, value int64) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "tx", "gas", "counter"},
		float32(value),
		[]metrics.Label{telemetry.NewLabel("type", gasType)},
	)
}

// Measures the number of times optimistic processing runs
// Metric Name:
//
//	eni_optimistic_processing_counter
func IncrementOptimisticProcessingCounter(enabled bool) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "optimistic", "processing", "counter"},
		float32(1),
		[]metrics.Label{telemetry.NewLabel("enabled", strconv.FormatBool(enabled))},
	)
}

// Measures RPC endpoint request throughput
// Metric Name:
//
//	eni_rpc_request_counter
func IncrementRpcRequestCounter(endpoint string, connectionType string, success bool) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "rpc", "request", "counter"},
		float32(1),
		[]metrics.Label{
			telemetry.NewLabel("endpoint", endpoint),
			telemetry.NewLabel("connection", connectionType),
			telemetry.NewLabel("success", strconv.FormatBool(success)),
		},
	)
}

func IncrementErrorMetrics(scenario string, err error) {
	if err == nil {
		return
	}
	var assocErr types.AssociationMissingErr
	if errors.As(err, &assocErr) {
		IncrementAssociationError(scenario, assocErr)
		return
	}
	// add other error types to handle as metrics
}

func IncrementAssociationError(scenario string, err types.AssociationMissingErr) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "association", "error"},
		1,
		[]metrics.Label{
			telemetry.NewLabel("scenario", scenario),
			telemetry.NewLabel("type", err.AddressType()),
		},
	)
}

// Measures the RPC request latency in milliseconds
// Metric Name:
//
//	eni_rpc_request_latency_ms
func MeasureRpcRequestLatency(endpoint string, connectionType string, startTime time.Time) {
	metrics.MeasureSinceWithLabels(
		[]string{"eni", "rpc", "request", "latency_ms"},
		startTime.UTC(),
		[]metrics.Label{
			telemetry.NewLabel("endpoint", endpoint),
			telemetry.NewLabel("connection", connectionType),
		},
	)
}

// IncrProducerEventCount increments the counter for events produced.
// This metric counts the number of events produced by the system.
// Metric Name:
//
//	eni_loadtest_produce_count
func IncrProducerEventCount(msgType string) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "loadtest", "produce", "count"},
		1,
		[]metrics.Label{telemetry.NewLabel("msg_type", msgType)},
	)
}

// IncrConsumerEventCount increments the counter for events consumed.
// This metric counts the number of events consumed by the system.
// Metric Name:
//
//	eni_loadtest_consume_count
func IncrConsumerEventCount(msgType string) {
	telemetry.IncrCounterWithLabels(
		[]string{"eni", "loadtest", "consume", "count"},
		1,
		[]metrics.Label{telemetry.NewLabel("msg_type", msgType)},
	)
}

func AddHistogramMetric(key []string, value float32) {
	metrics.AddSample(key, value)
}

// Gauge for gas price paid for transactions
// Metric Name:
//
// eni_evm_effective_gas_price
func HistogramEvmEffectiveGasPrice(gasPrice *big.Int) {
	AddHistogramMetric(
		[]string{"eni", "evm", "effective", "gas", "price"},
		float32(gasPrice.Uint64()),
	)
}

// Gauge for block base fee
// Metric Name:
//
// eni_evm_block_base_fee
func GaugeEvmBlockBaseFee(baseFee *big.Int, blockHeight int64) {
	metrics.SetGauge(
		[]string{"eni", "evm", "block", "base", "fee"},
		float32(baseFee.Uint64()),
	)
}
