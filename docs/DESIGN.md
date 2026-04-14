# Melrōse Architecture & Design Decisions

This document details the architectural and design decisions behind the `melrōse` project.

## 1. Custom Domain-Specific Language (DSL)
Instead of adopting a general-purpose language (like Lua or Python) or a data format (like JSON or YAML), Melrōse uses a custom declarative functional language.
* **Why:** Music composition is highly mathematical and compositional. A functional DSL allows expressing complex musical structures (sequences, chords, parallel tracks) concisely. 
* **Implementation:** It leverages `github.com/expr-lang/expr` to evaluate expressions securely and natively in Go. This allows defining custom functions (e.g., `sequence`, `chord`, `transpose`) that immediately map to Go structs, avoiding complex parsing boilerplate while giving high performance.

## 2. Abstraction over Core Musical Primitives
The core domain model revolves around a unified interface: `Sequenceable`.
* **Why:** Everything in Melrōse that makes a sound fundamentally resolves down to a sequence of notes (`[][]Note`, representing successive groups of simultaneously played notes). 
* **Implementation:** Functions like `note()`, `chord()`, `sequence()`, and operations like `transpose()`, `reverse()`, and `stretch()` all implement `Sequenceable` or return objects that do. This allows infinite functional composition (e.g., `transpose(2, stretch(2, sequence('c e g')))`).

## 3. Immutable and Functional Operations
Operations in the `op` package do not mutate the incoming sequence. They create lightweight wrappers or generate new sequences.
* **Why:** Immutability guarantees that when a sequence is played on multiple tracks or used in multiple loops, applying a modifier (like `dynamic`) to one instance does not unpredictably alter the original variable. It ensures safety during live-coding.

## 4. Timeline-Based Event Scheduling
Playback does not rely on simple procedural `time.Sleep` calls.
* **Why:** If the system just slept between notes, it would block execution and make it impossible to interrupt or dynamically alter playback (crucial for live-coding). 
* **Implementation:** The `core.Timeline` stores `TimelineEvent`s scheduled in the future. A background `Beatmaster` goroutine sweeps the timeline to dispatch events (like `note_on` and `note_off`) to the `AudioDevice`. This separates evaluation from playback, allowing a user to run new code without glitching the audio.

## 5. Live Variable Delegation
Loops and tracks evaluate their contents lazily during playback.
* **Why:** For live performance, if a user redefines a variable `p = sequence("c d e")` to `p = sequence("c d e f")`, a running loop `loop(p)` must adapt instantly on its next iteration.
* **Implementation:** The DSL parser wraps variables in `VariableStorage`. Loops point to these dynamic variable proxies rather than statically copied `Sequence` instances. This forms the foundation of the live-coding experience.

## 6. Abstracted Audio Transport Layer
The `AudioDevice` and `midi/transport` interfaces hide MIDI driver specifics.
* **Why:** Standardizing the transport layer makes the core engine platform-agnostic. 
* **Implementation:** Melrōse can be compiled using `gitlab.com/gomidi/rtmididrv` for native macOS/Linux/Windows playback, but also features a `wasm.go` build tag implementation. This allows the same composition engine to run inside a web browser or WebAssembly runtime without native C bindings.

## 7. Client-Server Ideology
Melrōse provides both an interactive CLI (`ui/cli` via `liner`) and a headless HTTP Server (`server/lang.go`).
* **Why:** Live-coding in a terminal is limited. By exposing a `/v1/statements` endpoint, Melrōse delegates heavy IDE features (syntax highlighting, hot reloading, snippet execution) to external clients, such as the official Visual Studio Code extension or an MCP (Model Context Protocol) server.
* **Implementation:** The Go binary acts as a background runtime engine, receiving AST expressions via HTTP, evaluating them in the shared `Context`, and piping output to the MIDI timeline.

## 8. Encapsulated Context
A `core.Context` object is threaded through almost all evaluators and playback components.
* **Why:** It contains the `VariableStorage`, `AudioDevice`, `LoopController` (Beatmaster), and environment limits. Instead of relying on global variables, this design makes it trivial to write deterministic unit tests, run multiple independent sessions, or tear down a live-coding setup cleanly.
