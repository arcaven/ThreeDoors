# OTEL Snapshot Field Mapping (Phase 3 Marvel)

This document maps `quota-usage.jsonl` snapshot fields to [R-016 OTEL metrics](../../_bmad-output/planning-artifacts/tdd-performance-otel-research.md) for Phase 3 portability.

## Snapshot → OTEL Metric Mapping

| Snapshot Field | OTEL Metric | Type | Labels/Attributes |
|---|---|---|---|
| `window_usage_pct` | `dark_factory.phase.token_usage` | Gauge | `window="5h"` |
| `agent_breakdown[].billed_tokens` | `dark_factory.phase.token_usage` | Gauge | `agent=<name>` |
| `agent_breakdown[].input_tokens` | `gen_ai.client.token.usage` | Counter | `agent=<name>, token_type="input"` |
| `agent_breakdown[].output_tokens` | `gen_ai.client.token.usage` | Counter | `agent=<name>, token_type="output"` |
| `total_billed_tokens` | `dark_factory.phase.token_usage` | Gauge | `scope="total"` |
| `threshold_tier` | `dark_factory.quota.tier` | — | Attribute on all metrics |
| `peak_hours` | `dark_factory.quota.peak_hours` | — | Attribute on all metrics |
| `window_start_time` | Span attribute | — | `window_start` |
| `estimated_reset_time` | Span attribute | — | `window_reset` |

## OTEL Semantic Conventions Reference

Per the [OTEL GenAI Semantic Conventions](https://opentelemetry.io/docs/specs/semconv/gen-ai/) (experimental):

- `gen_ai.client.token.usage` — standard token counter (input/output split)
- Custom `dark_factory.*` metrics use the namespace established in R-016

## Integration Notes

- Snapshots are written as JSONL for Phase 1 (file-based monitoring)
- Phase 3 collector reads JSONL and emits OTEL metrics via OTLP gRPC/HTTP
- The `type: "quota_snapshot"` discriminator allows mixed-type JSONL files
- No schema changes needed for OTEL export — the JSONL fields map 1:1
