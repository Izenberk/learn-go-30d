# learn-go-30d

> 30-day structured deep-dive into production-grade Go.

---

## Philosophy

This isn't a tutorial repo. It's a **deliberate practice log** — no copy-paste, no skimming.

Every concept is:
1. **Read** — one concept, one source, understood from first principles
2. **Broken** — intentionally written wrong to find the edges
3. **Applied** — mapped back to a real project (see below)

The goal isn't to finish 30 days. The goal is to think like a Go engineer by the end of them.

---

## Anchor Project

Each concept learned here gets applied to the **Health-Check API** — a Go service with Prometheus metrics and Grafana dashboards, running on local Docker Compose.

It's intentionally small in scope. That's the point — the complexity comes from doing it *correctly*, not from feature count.

| Layer | Stack |
|-------|-------|
| API | Go (net/http) |
| Metrics | Prometheus |
| Dashboards | Grafana |
| Runtime | Docker Compose (local) |

---

## Structure

```
learn-go-30d/
├── ring1/          # Days 01–10 — Language core
│   ├── day01/      # Interfaces — implicit, small, composable
│   ├── day02/      # Errors — wrap, sentinel, Is/As
│   └── ...
├── ring2/          # Days 11–20 — Idiomatic patterns
│   └── ...
├── ring3/          # Days 21–30 — Production internals
│   └── ...
└── README.md
```

Each day folder contains:
- `main.go` — the concept implemented, broken, and understood
- `notes.md` — mental models, analogies, Health-Check API connections

---

## The 30-Day Map

### Ring 1 — Language Core (Days 1–10)

| Day | Topic | Status |
|-----|-------|--------|
| 01 | Interfaces — implicit, small, composable | ✅ |
| 02 | Errors — wrap, sentinel, Is/As | ✅ |
| 03 | Goroutines — GMP scheduler, leaks | ⬜ |
| 04 | Channels — buffered, nil, close rules | ⬜ |
| 05 | Select — multi-wait, nil trick, default | ⬜ |
| 06 | Context — cancel, timeout, values | ⬜ |
| 07 | sync — Mutex, WaitGroup, Once, atomic | ⬜ |
| 08 | Testing — table-driven, benchmarks | ⬜ |
| 09 | Modules — go.mod, workspace, tidy | ⬜ |
| 10 | Toolchain — vet, race detector, build flags | ⬜ |

### Ring 2 — Idiomatic Patterns (Days 11–20)

| Day | Topic | Status |
|-----|-------|--------|
| 11 | Functional options — WithX pattern | ⬜ |
| 12 | Embedding — composition over inheritance | ⬜ |
| 13 | Worker pool — bounded, graceful shutdown | ⬜ |
| 14 | Pipeline — fan-out, fan-in, stages | ⬜ |
| 15 | Generics — constraints, when to use | ⬜ |
| 16 | Dependency injection — no framework | ⬜ |
| 17 | Middleware — HTTP + gRPC interceptors | ⬜ |
| 18 | Structured logging — slog, fields, levels | ⬜ |
| 19 | Graceful shutdown — signal, drain | ⬜ |
| 20 | Mocking — interfaces, testify, stubs | ⬜ |

### Ring 3 — Production Internals (Days 21–30)

| Day | Topic | Status |
|-----|-------|--------|
| 21 | Escape analysis — stack vs heap, gcflags -m | ⬜ |
| 22 | pprof — CPU, heap, goroutine dump, flame graph | ⬜ |
| 23 | GC — tricolor mark-sweep, GOGC, GOMEMLIMIT | ⬜ |
| 24 | Memory model — happens-before, sync guarantees | ⬜ |
| 25 | Benchmarks — b.N, AllocsPerOp, ResetTimer | ⬜ |
| 26 | Tracing — runtime/trace, execution tracer | ⬜ |
| 27 | Mutex contention — pprof mutex profile | ⬜ |
| 28 | Allocation reduction — sync.Pool, reuse patterns | ⬜ |
| 29 | Build constraints — tags, cross-compilation | ⬜ |
| 30 | Apply all — profile Health-Check API end-to-end | ⬜ |

---

## Daily Method — Structured Tinkering

```
Read     (20 min) — one concept, one source
Break it (20 min) — write the wrong version, run go test -race, hit the boundary
Apply    (rest)   — find where it belongs in the Health-Check API
```

---

## Key Principles

- No copy-paste — understand the *why* behind every line
- Accept interfaces, return concrete types
- Errors are values — explicit, wrapped with context, inspectable
- Concurrency is not parallelism — goroutines are cheap, leaks are not

---

*The Health-Check API is the textbook — every concept maps back to real code.*