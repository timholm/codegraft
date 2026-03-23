# codegraft

**Language:** Go
**Source:** https://arxiv.org/abs/2603.15159
**Estimated lines:** 5200

## Problem

Enterprise engineering teams have internal SDKs and private libraries that AI coding assistants cannot use effectively — even RAG with perfect documentation fails because LLMs lack learned patterns for invoking unfamiliar APIs.

## Solution

Codegraft is a CLI tool that ingests a private library's source code and API surface, automatically synthesizes high-quality training examples using graph-based data evolution (inspired by PriCoder), and produces LoRA adapters for open-source models via Ollama integration. The key innovation is a fully automated pipeline from library source to fine-tuned model — no ML expertise required, no data leaves the organization.

## Expected Files

["main.go","cmd/root.go","cmd/ingest.go","cmd/synthesize.go","cmd/tune.go","cmd/serve.go","cmd/status.go","internal/parser/parser.go","internal/parser/symbols.go","internal/parser/languages.go","internal/graph/graph.go","internal/graph/evolution.go","internal/graph/pruning.go","internal/graph/scoring.go","internal/synth/generator.go","internal/synth/templates.go","internal/synth/diversity.go","internal/synth/validator.go","internal/adapter/ollama.go","internal/adapter/gguf.go","internal/adapter/lora.go","internal/server/server.go","internal/server/completion.go","internal/server/middleware.go","internal/config/config.go","internal/config/defaults.go","internal/store/store.go","internal/store/sqlite.go","internal/llm/client.go","internal/llm/prompt.go"]
