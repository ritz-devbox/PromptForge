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
- `promptforge compile --explain` - Compile and write `prompt.ir.explain.json`
- `promptforge lint` - Lint `plan.md` and report issues
- `promptforge templates` - List available plan templates
- `promptforge migrate` - Upgrade `prompt.ir.json` to the current IR version
- `promptforge audit` - Validate `prompt.ir.json` integrity and schema sync

## Artifacts

- `plan.md` - Human-readable, editable, not authoritative
- `prompt.ir.json` - Machine-enforceable, authoritative
- `prompt.ir.schema.json` - JSON Schema for validating `prompt.ir.json`
- `prompt.ir.explain.json` - Mapping of plan sections to IR outputs

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

## Step-by-Step Guide

### Install Go (if needed)
1. Download Go from https://go.dev/dl/
2. Run the installer and keep defaults.
3. Verify in a new terminal:
   ```bash
   go version
   ```

### CLI Flow
1. **Build the CLI (once per change):**
   ```bash
   go build -o promptforge ./cmd/promptforge
   ```
2. **Initialize a plan:**
   ```bash
   ./promptforge init "Build a support triage assistant"
   ```
3. **Fill in plan.md:**
   ```text
   ## Constraints
   - Must classify tickets into one of: bug, billing, account, or feature
   - Must ask a clarifying question when category is unclear
   - Must return JSON only

   ## Out of Scope
   - Handling refunds
   - Changing user passwords
   - Accessing customer data
   ```
4. **Lint the plan:**
   ```bash
   ./promptforge lint
   ```
5. **Compile IR + schema:**
   ```bash
   ./promptforge compile
   ```
6. **Generate explain mapping:**
   ```bash
   ./promptforge compile --explain
   ```
7. **Audit outputs:**
   ```bash
   ./promptforge audit
   ```
8. **Migrate an older IR (if needed):**
   ```bash
   ./promptforge migrate
   ```

### No-Terminal Flow (Non-Technical)
1. **Open the UI:**
   - Double-click `PromptForge_Launcher.cmd`.
2. **Plan Builder mode:**
   - Enter a Goal, Constraints (one per line), and Out of Scope.
   - Click **Download plan.md** and save it.
3. **Inspector mode:**
   - Load your saved `plan.md`.
   - Optionally load `prompt.ir.json` and `prompt.ir.explain.json` to see mappings.

## Status

V0 - Initial implementation with deterministic compilation logic.

## Non-Technical Use

- Double-click `PromptForge_Launcher.cmd` to open the visual UI and load `plan.md`, `prompt.ir.json`, and `prompt.ir.explain.json` without using the terminal.

## Upgrade Tracker

- Schema validation + schema export (`prompt.ir.schema.json`)
- plan.md linting (`promptforge lint`)
- Explain compile mapping (`promptforge compile --explain`)
- IR versioning + migration (`version` field, `promptforge migrate`)
- Templates system (`promptforge templates`, `promptforge init --template`)
- VS Code extension enhancements (lint diagnostics, compile on save)
- Golden/snapshot test harness for deterministic compilation
- Audit command for IR integrity checks
- UI refresh for plan/IR exploration (`ui/index.html`)
- One-click UI launcher (`PromptForge_Launcher.cmd`)

