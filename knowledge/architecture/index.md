# Architecture

* [5-Layer System Architecture](layers.md) - Structural separation of concerns across the OKF specification, agent convention, skills, deterministic tooling, and knowledge corpus.
* [Go Single-Binary CLI & MCP Architecture Decision](tooling-decision.md) - Architectural decision to implement the deterministic OKF tooling layer as a standalone Go binary with dual CLI and MCP support.
* [Bundle Isolation and Mutation Security Boundaries](security-boundaries.md) - Defensive security architecture enforcing canonical bundle boundaries, symlink containment, path traversal prevention, and frontmatter injection defense.
* [CLI Optional Path Boundary](cli-argument-boundary.md) - CLI commands consume an optional bundle or target path only from the first remaining argument, preserving all subsequent flag values.
