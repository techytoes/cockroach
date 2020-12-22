// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package execinfrapb

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/util/humanizeutil"
	"github.com/cockroachdb/cockroach/pkg/util/optional"
	"github.com/cockroachdb/cockroach/pkg/util/tracing"
	"github.com/cockroachdb/cockroach/pkg/util/tracing/tracingpb"
	"github.com/dustin/go-humanize"
	"github.com/gogo/protobuf/types"
)

// ProcessorComponentID returns a ComponentID for the given processor in a flow.
func ProcessorComponentID(flowID FlowID, processorID int32) ComponentID {
	return ComponentID{
		FlowID: flowID,
		Type:   ComponentID_PROCESSOR,
		ID:     processorID,
	}
}

// StreamComponentID returns a ComponentID for the given stream in a flow.
func StreamComponentID(flowID FlowID, streamID StreamID) ComponentID {
	return ComponentID{
		FlowID: flowID,
		Type:   ComponentID_STREAM,
		ID:     int32(streamID),
	}
}

// FlowIDTagKey is the key used for flow id tags in tracing spans.
const (
	FlowIDTagKey = tracing.TagPrefix + "flowid"

	// StreamIDTagKey is the key used for stream id tags in tracing spans.
	StreamIDTagKey = tracing.TagPrefix + "streamid"

	// ProcessorIDTagKey is the key used for processor id tags in tracing spans.
	ProcessorIDTagKey = tracing.TagPrefix + "processorid"

	// StatTagPrefix is prefixed to all stats output in Span tags.
	StatTagPrefix = tracing.TagPrefix + "stat."
)

// StatsTags is part of the tracing.SpanStats interface.
func (s *ComponentStats) StatsTags() map[string]string {
	result := make(map[string]string, 4)
	if s.Component != (ComponentID{}) {
		result[FlowIDTagKey] = s.Component.FlowID.String()
		switch s.Component.Type {
		case ComponentID_PROCESSOR:
			result[ProcessorIDTagKey] = strconv.Itoa(int(s.Component.ID))
		case ComponentID_STREAM:
			result[StreamIDTagKey] = strconv.Itoa(int(s.Component.ID))
		case ComponentID_FLOW:
			// Nothing extra to set.
		}
	}

	s.formatStats(func(key string, value interface{}) {
		// The key becomes a tracing span tag. Replace spaces with dots and use
		// only lowercase characters.
		key = strings.ToLower(strings.ReplaceAll(key, " ", "."))
		result[StatTagPrefix+key] = fmt.Sprint(value)
	})
	return result
}

// StatsForQueryPlan returns the statistics as a list of strings that can be
// displayed in query plans and diagrams.
func (s *ComponentStats) StatsForQueryPlan() []string {
	result := make([]string, 0, 4)
	s.formatStats(func(key string, value interface{}) {
		result = append(result, fmt.Sprintf("%s: %v", key, value))
	})
	return result
}

// formatStats calls fn for each statistic that is set.
func (s *ComponentStats) formatStats(fn func(suffix string, value interface{})) {
	// Network Rx stats.
	if s.NetRx.Latency.HasValue() {
		fn("network latency", humanizeutil.Duration(s.NetRx.Latency.Value()))
	}
	if s.NetRx.WaitTime.HasValue() {
		fn("network wait time", humanizeutil.Duration(s.NetRx.WaitTime.Value()))
	}
	if s.NetRx.DeserializationTime.HasValue() {
		fn("deserialization time", humanizeutil.Duration(s.NetRx.DeserializationTime.Value()))
	}
	if s.NetRx.TuplesReceived.HasValue() {
		fn("network tuples received", humanizeutil.Count(s.NetRx.TuplesReceived.Value()))
	}
	if s.NetRx.BytesReceived.HasValue() {
		fn("network bytes received", humanize.IBytes(s.NetRx.BytesReceived.Value()))
	}
	if s.NetRx.MessagesReceived.HasValue() {
		fn("network messages received", humanizeutil.Count(s.NetRx.MessagesReceived.Value()))
	}

	// Network Tx stats.
	if s.NetTx.TuplesSent.HasValue() {
		fn("network tuples sent", humanizeutil.Count(s.NetTx.TuplesSent.Value()))
	}
	if s.NetTx.BytesSent.HasValue() {
		fn("network bytes sent", humanize.IBytes(s.NetTx.BytesSent.Value()))
	}
	if s.NetTx.MessagesSent.HasValue() {
		fn("network messages sent", humanizeutil.Count(s.NetTx.MessagesSent.Value()))
	}

	// Input stats.
	switch len(s.Inputs) {
	case 1:
		if s.Inputs[0].NumTuples.HasValue() {
			fn("input tuples", humanizeutil.Count(s.Inputs[0].NumTuples.Value()))
		}
		if s.Inputs[0].WaitTime.HasValue() {
			fn("input stall time", humanizeutil.Duration(s.Inputs[0].WaitTime.Value()))
		}

	case 2:
		if s.Inputs[0].NumTuples.HasValue() {
			fn("left tuples", humanizeutil.Count(s.Inputs[0].NumTuples.Value()))
		}
		if s.Inputs[0].WaitTime.HasValue() {
			fn("left stall time", humanizeutil.Duration(s.Inputs[0].WaitTime.Value()))
		}
		if s.Inputs[1].NumTuples.HasValue() {
			fn("right tuples", humanizeutil.Count(s.Inputs[1].NumTuples.Value()))
		}
		if s.Inputs[1].WaitTime.HasValue() {
			fn("right stall time", humanizeutil.Duration(s.Inputs[1].WaitTime.Value()))
		}
	}

	// KV stats.
	if s.KV.KVTime.HasValue() {
		fn("KV time", humanizeutil.Duration(s.KV.KVTime.Value()))
	}
	if s.KV.ContentionTime.HasValue() {
		fn("KV contention time", humanizeutil.Duration(s.KV.ContentionTime.Value()))
	}
	if s.KV.TuplesRead.HasValue() {
		fn("KV tuples read", humanizeutil.Count(s.KV.TuplesRead.Value()))
	}
	if s.KV.BytesRead.HasValue() {
		fn("KV bytes read", humanize.IBytes(s.KV.BytesRead.Value()))
	}

	// Exec stats.
	if s.Exec.ExecTime.HasValue() {
		fn("execution time", humanizeutil.Duration(s.Exec.ExecTime.Value()))
	}
	if s.Exec.MaxAllocatedMem.HasValue() {
		fn("max memory allocated", humanize.IBytes(s.Exec.MaxAllocatedMem.Value()))
	}
	if s.Exec.MaxAllocatedDisk.HasValue() {
		fn("max scratch disk allocated", humanize.IBytes(s.Exec.MaxAllocatedDisk.Value()))
	}

	// Output stats.
	if s.Output.NumBatches.HasValue() {
		fn("batches output", humanizeutil.Count(s.Output.NumBatches.Value()))
	}
	if s.Output.NumTuples.HasValue() {
		fn("tuples output", humanizeutil.Count(s.Output.NumTuples.Value()))
	}
}

// Union creates a new ComponentStats that contains all statistics in either the
// receiver (s) or the argument (other).
// If a statistic is set in both, the one in the receiver (s) is preferred.
func (s *ComponentStats) Union(other *ComponentStats) *ComponentStats {
	result := *s

	// Network Rx stats.
	if !result.NetRx.Latency.HasValue() {
		result.NetRx.Latency = other.NetRx.Latency
	}
	if !result.NetRx.WaitTime.HasValue() {
		result.NetRx.WaitTime = other.NetRx.WaitTime
	}
	if !result.NetRx.DeserializationTime.HasValue() {
		result.NetRx.DeserializationTime = other.NetRx.DeserializationTime
	}
	if !result.NetRx.TuplesReceived.HasValue() {
		result.NetRx.TuplesReceived = other.NetRx.TuplesReceived
	}
	if !result.NetRx.BytesReceived.HasValue() {
		result.NetRx.BytesReceived = other.NetRx.BytesReceived
	}
	if !result.NetRx.MessagesReceived.HasValue() {
		result.NetRx.MessagesReceived = other.NetRx.MessagesReceived
	}

	// Network Tx stats.
	if !result.NetTx.TuplesSent.HasValue() {
		result.NetTx.TuplesSent = other.NetTx.TuplesSent
	}
	if !result.NetTx.BytesSent.HasValue() {
		result.NetTx.BytesSent = other.NetTx.BytesSent
	}

	// Input stats. Make sure we don't reuse slices.
	result.Inputs = append([]InputStats(nil), s.Inputs...)
	result.Inputs = append(result.Inputs, other.Inputs...)

	// KV stats.
	if !result.KV.KVTime.HasValue() {
		result.KV.KVTime = other.KV.KVTime
	}
	if !result.KV.ContentionTime.HasValue() {
		result.KV.ContentionTime = other.KV.ContentionTime
	}
	if !result.KV.TuplesRead.HasValue() {
		result.KV.TuplesRead = other.KV.TuplesRead
	}
	if !result.KV.BytesRead.HasValue() {
		result.KV.BytesRead = other.KV.BytesRead
	}

	// Exec stats.
	if !result.Exec.ExecTime.HasValue() {
		result.Exec.ExecTime = other.Exec.ExecTime
	}
	if !result.Exec.MaxAllocatedMem.HasValue() {
		result.Exec.MaxAllocatedMem = other.Exec.MaxAllocatedMem
	}
	if !result.Exec.MaxAllocatedDisk.HasValue() {
		result.Exec.MaxAllocatedDisk = other.Exec.MaxAllocatedDisk
	}

	// Output stats.
	if !result.Output.NumBatches.HasValue() {
		result.Output.NumBatches = other.Output.NumBatches
	}
	if !result.Output.NumTuples.HasValue() {
		result.Output.NumTuples = other.Output.NumTuples
	}

	// Flow stats.
	if !result.FlowStats.MaxMemUsage.HasValue() {
		result.FlowStats.MaxMemUsage = other.FlowStats.MaxMemUsage
	}

	return &result
}

// MakeDeterministic is used only for testing; it modifies any non-deterministic
// statistics like elapsed time or exact number of bytes to fixed or
// manufactured values.
//
// Note that it does not modify which fields that are set. In other words, a
// field will have a non-zero protobuf value iff it had a non-zero protobuf
// value before. This allows tests to verify the set of stats that were
// collected.
func (s *ComponentStats) MakeDeterministic() {
	// resetUint resets an optional.Uint to 0, if it was set.
	resetUint := func(v *optional.Uint) {
		if v.HasValue() {
			v.Set(0)
		}
	}
	// timeVal resets a duration to 1ns, if it was set.
	timeVal := func(v *optional.Duration) {
		if v.HasValue() {
			v.Set(0)
		}
	}

	// NetRx.
	timeVal(&s.NetRx.Latency)
	timeVal(&s.NetRx.WaitTime)
	timeVal(&s.NetRx.DeserializationTime)
	if s.NetRx.BytesReceived.HasValue() {
		// BytesReceived can be non-deterministic because some message fields have
		// varying sizes across different runs (e.g. metadata). Override to a useful
		// value for tests.
		s.NetRx.BytesReceived.Set(8 * s.NetRx.TuplesReceived.Value())
	}
	if s.NetRx.MessagesReceived.HasValue() {
		// Override to a useful value for tests.
		s.NetRx.MessagesReceived.Set(s.NetRx.TuplesReceived.Value() / 2)
	}

	// NetTx.
	if s.NetTx.BytesSent.HasValue() {
		// BytesSent can be non-deterministic because some message fields have
		// varying sizes across different runs (e.g. metadata). Override to a useful
		// value for tests.
		s.NetTx.BytesSent.Set(8 * s.NetTx.TuplesSent.Value())
	}
	if s.NetTx.MessagesSent.HasValue() {
		// Override to a useful value for tests.
		s.NetTx.MessagesSent.Set(s.NetTx.TuplesSent.Value() / 2)
	}

	// KV.
	timeVal(&s.KV.KVTime)
	timeVal(&s.KV.ContentionTime)
	if s.KV.BytesRead.HasValue() {
		// BytesRead is overridden to a useful value for tests.
		s.KV.BytesRead.Set(8 * s.KV.TuplesRead.Value())
	}

	// Exec.
	timeVal(&s.Exec.ExecTime)
	resetUint(&s.Exec.MaxAllocatedMem)
	resetUint(&s.Exec.MaxAllocatedDisk)

	// Output.
	resetUint(&s.Output.NumBatches)

	// Inputs.
	for i := range s.Inputs {
		timeVal(&s.Inputs[i].WaitTime)
	}
}

// ExtractStatsFromSpans extracts all ComponentStats from a set of tracing
// spans.
func ExtractStatsFromSpans(
	spans []tracingpb.RecordedSpan, makeDeterministic bool,
) map[ComponentID]*ComponentStats {
	statsMap := make(map[ComponentID]*ComponentStats)
	for i := range spans {
		if spans[i].Stats == nil {
			continue
		}

		var stats ComponentStats
		if err := types.UnmarshalAny(spans[i].Stats, &stats); err != nil {
			continue
		}
		if stats.Component == (ComponentID{}) {
			continue
		}
		if makeDeterministic {
			stats.MakeDeterministic()
		}
		existing := statsMap[stats.Component]
		if existing == nil {
			statsMap[stats.Component] = &stats
		} else {
			// In the vectorized flow we can have multiple statistics entries for one
			// component. Merge the stats together.
			// TODO(radu): figure out a way to emit the statistics correctly in the
			// first place.
			statsMap[stats.Component] = existing.Union(&stats)
		}
	}
	return statsMap
}
