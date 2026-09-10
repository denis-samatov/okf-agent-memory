---
type: Convention
title: Offset-Aware Timestamp Format
description: OKF timestamp fields use RFC3339 datetimes with explicit UTC offsets and stale lifecycle checks compare parsed instants.
tags: [timestamps, rfc3339, validation, lifecycle]
generated: { by: agent/codex, at: "2026-09-10T04:09:21Z" }
---

# Offset-Aware Timestamp Format

Timestamp-valued metadata uses the RFC3339 profile of ISO 8601 with an explicit UTC offset. Staleness compares instants, not strings. See [Knowledge Lifecycle and Review Workflow](lifecycle.md).
