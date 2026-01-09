# PromptForge - AI Prompt Contract Generator

Transform your ideas into machine-enforceable prompt contracts for AI agents. Simple form-based editor that generates JSON-ready contracts.

## Features

- 🎯 **Simple Form-Based Editor** - No markdown syntax needed
- ⚡ **Auto-Compilation** - Generates JSON automatically from your plan
- 📄 **JSON Output** - Ready-to-use contracts for AI agents/LLMs
- 🔄 **Smart Parsing** - Extracts rules from your constraints
- 🚫 **Failure Modes** - Generates error handling from out-of-scope items

## Quick Start

### 1. Create a New Plan

1. Press `Ctrl + Shift + P` (or `Cmd + Shift + P` on Mac)
2. Type: `PromptForge: Create New Plan`
3. Enter your rough idea (e.g., "Create a chatbot that helps debug API errors")
4. The editor opens automatically

### 2. Fill in Your Plan

- **Goal**: What should the AI do?
- **Constraints**: What rules must it follow? (one per line)
- **Out of Scope**: What should it NOT do? (one per line)

### 3. Generate JSON

1. Click "💾 Save & Generate JSON"
2. JSON appears below the form
3. Click "📋 Copy JSON" to use with your AI agent

## Example

**Goal:**
```
Create a chatbot that helps users debug API errors
```

**Constraints:**
```
- Must provide code examples in any programming language requested
- Must validate API endpoints before suggesting fixes
- Must not execute code, only show examples
```

**Out of Scope:**
```
- Does not handle authentication issues
- Does not debug network connectivity problems
```

**Generated JSON** includes:
- System role based on your Goal
- Rules from your Constraints
- Failure modes from your Out of Scope

## Usage

### Create New Plan
- Command Palette: `PromptForge: Create New Plan`
- Or right-click in Explorer → "Create New Prompt Plan"

### Edit Existing Plan
- Right-click on `plan.md` → "Open PromptForge Editor"
- Or Command Palette: `PromptForge: Open Editor`

## Requirements

- VS Code 1.74.0 or higher
- PromptForge CLI (optional, for terminal usage)

## How It Works

1. **You enter your idea** in simple terms
2. **Extension creates plan.md** with your content
3. **Compiler parses** your plan and extracts:
   - Goal → System Role
   - Constraints → Rules
   - Out of Scope → Failure Modes
4. **JSON is generated** ready for AI agents

## Files Created

- `promptforge/plan.md` - Your human-readable plan (editable)
- `prompt.ir.json` - Machine-enforceable contract (generated)

## Use with AI Agents

Copy the generated JSON and use it with:
- OpenAI API
- Anthropic Claude
- LangChain
- Any LLM that accepts structured prompts

## License

MIT License - Free to use, modify, and distribute

## Support

For issues, feature requests, or contributions, visit:
https://github.com/promptforge/promptforge

