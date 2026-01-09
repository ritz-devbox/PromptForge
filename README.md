# PromptForge

A CLI-first developer tool that compiles human intent into strict, machine-enforceable prompt contracts for LLMs.

## Purpose

PromptForge separates planning (human-facing) from compilation (machine-facing), producing deterministic, schema-bound prompt artifacts.

## Architecture

- **Clean separation of concerns**: CLI, compiler, and domain models are isolated
- **Compiler pattern**: Intent (plan.md) → IR (prompt.ir.json)
- **Deterministic behavior**: No LLM calls, no network calls, no inferred behavior
- **Testable**: All components are designed for unit testing

## Commands

- `promptforge init` - Initialize a new project (creates `plan.md`)
- `promptforge compile` - Compile `plan.md` to `prompt.ir.json`

## Artifacts

- `plan.md` - Human-readable, editable, not authoritative
- `prompt.ir.json` - Machine-enforceable, authoritative

## Building

```bash
go build -o promptforge ./cmd/promptforge
```

## Testing

### Run Unit Tests
```bash
go test ./...
```

### Manual Testing
1. **Initialize a project:**
   ```bash
   ./promptforge init
   ```
   Creates `promptforge/plan.md` with a template.

2. **Compile the project:**
   ```bash
   ./promptforge compile
   ```
   Generates `prompt.ir.json` in the repository root.

See [TESTING.md](TESTING.md) for detailed testing instructions.

## Status

V0 - Initial implementation with deterministic compilation logic.

