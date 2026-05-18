# WiscKey in Go — Project Plan

> A summary of our conversation and full 4-week project roadmap.

---

## Background

After completing MIT's distributed systems course (building Raft in Go), the goal is a new pet project that involves:

- Persistence
- Distribution / network problems
- A clean paper implementation worth blogging about

**Chosen project:** Implement **WiscKey** (USENIX ATC '16) in Go — a storage engine that separates keys from values to reduce write amplification in LSM-trees.

---

## Why WiscKey?

Three other candidates were considered:

| Project | Verdict |
|---|---|
| **WiscKey** | ✅ Best fit — meaty, scoped, great benchmark story |
| **Bitcask** | Good but well-trodden, less impressive to blog about |
| **Anna / CRDTs** | Harder to scope to 1 month without feeling half-baked |

WiscKey wins because:
- Touches persistence, crash recovery, GC, and compaction — all meaty problems
- You build a real key-value store you can benchmark vs a naive LSM
- The blog story writes itself: write amplification before/after, shown in graphs
- The paper is clean and implementable — no hand-wavy parts

**Paper:** [WiscKey: Separating Keys from Values in SSD-Conscious Storage](https://www.usenix.org/system/files/conference/fast16/fast16-papers-lu.pdf) — USENIX FAST '16

---

## Language Choice

**Go is the right tool.** Here's the honest breakdown:

### C — most credible for a storage engine blog post
WiscKey's story is about I/O performance, write amplification, and fsync correctness. C gives you direct control over memory layout, page alignment, and syscalls with no runtime in between. Every serious storage engine (LevelDB, RocksDB, SQLite, LMDB) is in C/C++.

**Downside:** you'll fight `malloc`/`free` and pointer arithmetic. Valgrind and AddressSanitizer become essential.

### Go — solid choice, faster to finish
Go was basically designed for this kind of project. The standard library has everything — `os.File`, `sync.RWMutex`, `encoding/binary` — and goroutines make background compaction/GC natural to express.

**Honest tradeoff:** Go's runtime (GC, goroutine scheduler) adds a layer of indirection that makes low-level performance results slightly murkier to explain. Not a dealbreaker, but you'll occasionally hit latency spikes from the GC, not your design.

### C# — wrong tool
Not because C# is bad, but because the .NET runtime is far from the metal. The blog story gets muddier and you'd spend energy fighting the framework.

### Recommendation
**Start in Go.** You already know it, WiscKey is complex enough that getting a *correct* implementation matters more than a fast one at first, and a correct Go implementation with good benchmarks is a genuinely impressive blog post. If you want to go deeper, porting the hot path to C is a natural follow-up post.

---

## Shaking off the Rust (Go edition)

Don't do a tutorial. You already know the language — you just need your fingers to remember it.

**Build a toy in 1-2 days:**

A simple TCP server where multiple clients can connect and send messages. The server fans out each message to all other connected clients (basically a tiny chat server). Persist messages to an append-only log file on disk, reload on restart.

This forces you to touch:
- Goroutines and channels
- `net` package
- `os.File`
- `encoding/binary` or `encoding/json`
- Mutexes
- Graceful shutdown

All of which you'll use in the real project. **Throw it away when done.**

---

## 4-Week Project Plan

### Week 0 — Days 1–2: Rust shakeoff
Build the toy TCP server described above. Don't skip this.

### Week 0 — Days 3–4: Read the paper

Read the WiscKey paper in **two passes**:

**Pass 1 (1 hour):** Read abstract, intro, and section 5 (evaluation) first. Understand *what they claim* before *how*. Sketch the architecture on paper.

**Pass 2 (2 hours):** Read sections 3 and 4 carefully (design + implementation). Every time you hit something unclear, write a question in a doc. You won't answer them all now — that's fine.

Also skim:
- The [LevelDB README](https://github.com/google/leveldb) for LSM context
- This [LSM-tree explainer](https://www.cs.umb.edu/~poneil/lsmtree.pdf) if LSM internals feel hazy

---

### Week 1 — Core Storage Primitives

Get the boring foundational pieces right — everything else sits on them.

- **WAL** — write-ahead log, append-only, with fsync
- **MemTable** — in-memory sorted map (`btree` package or sorted slice)
- **SSTable** — write a sorted file to disk, read it back
- No compaction yet, no vLog yet

Flow: write goes to WAL + memtable → flush memtable to SSTable when full → reads check memtable then SSTables.

**End of week 1:** a naive LSM key-value store. Ugly, but correct.

---

### Week 2 — The WiscKey Part

This is the actual paper contribution.

**Introduce the vLog** — a separate append-only file where values live:
- SSTable entries now store `(key → {vlog_offset, value_size})` instead of the value
- Reads: SSTable lookup → get offset → read value from vLog
- Writes: append value to vLog → insert `(key, offset)` into LSM

**Then handle vLog GC.** The vLog grows forever without cleanup. WiscKey's GC:
1. Reads a chunk from the tail
2. Checks which keys are still valid via LSM lookup
3. Rewrites live values to the head
4. Advances the tail pointer

This is where crashes hurt — think carefully about what happens if GC is interrupted.

---

### Week 3 — Crash Recovery + Correctness

This is what separates a toy from something you're proud of.

- On startup: replay WAL to rebuild memtable
- Persist the vLog tail pointer durably (fsync before updating)
- Write a **chaos test:** open the DB, do writes, kill the process with `os.Exit(1)` at a random point, reopen, verify no corruption and no lost acknowledged writes
- Add `DB.Stats()` reporting write amplification — this becomes your blog graph

**End of week 3:** run it, kill it mid-write 100 times, always recover cleanly.

---

### Week 4 — Polish + Benchmarks + Blog

- Benchmark vs your week 1 pure-LSM version: write throughput, read latency, space amplification
- Write the blog post alongside cleaning the code

**Blog post structure:**
1. What's write amplification and why does it hurt
2. How LSM-trees work
3. WiscKey's insight
4. My implementation
5. Benchmark results

- Publish the repo with a good README

---

## Key Go Packages

```
os, io              — file I/O, fsync
encoding/binary     — serialize integers to bytes
sync                — RWMutex for concurrent reads
github.com/google/btree  — sorted in-memory structure
testing, testing/quick   — for chaos tests
```

---

## Starting Interface

Don't design everything upfront. Start with the simplest interface:

```go
type DB interface {
    Put(key, value []byte) error
    Get(key []byte) ([]byte, error)
    Close() error
}
```

Everything else is an implementation detail. Refactor as you learn what the paper actually requires.

---

## Appendix: What are CRDTs?

*(Came up as an alternative project direction)*

CRDTs (Conflict-free Replicated Data Types) solve the fundamental distributed systems tension: how do multiple nodes update shared data without coordinating?

**The core insight:** design your data structure so *any two states can always be merged* with no conflicts — ever. The math behind it is **lattices** — a partial order where there's always a "join" (least upper bound) of any two states.

**Example — G-Counter:**
```
Node A: [3, 0, 0]
Node B: [0, 2, 0]
Merge:  [3, 2, 0]  ← take max of each slot
Total = 5  ✓ always correct, no coordination needed
```

**Common CRDTs:**

| Type | What it models | Merge rule |
|---|---|---|
| G-Counter | increment-only counter | max per slot |
| PN-Counter | increment + decrement | two G-counters |
| G-Set | add-only set | union |
| 2P-Set | add + remove (once) | union of add/remove sets |
| LWW-Register | single value, last write wins | max timestamp |
| OR-Set | add + remove (repeatedly) | unique tag per add |
| RGA / LSEQ | ordered sequence (collaborative text) | complex |

Used in production by: Redis Enterprise, Riak, Figma, Notion, Linear.

**Why not for this project:** harder to scope to 1 month without feeling half-baked. WiscKey has a cleaner finish line.
