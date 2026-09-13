# ✋ sure

> **Your terminal's second opinion**

`sure` is an AI-powered CLI utility that intercepts terminal commands before they execute, analyzes them for potential destructive behavior, mistakes, or dangerous typos, and asks for confirmation before letting you proceed.

**The core rule:** `sure` NEVER automatically modifies, corrects, replaces, or rewrites the user's command. It only asks whether you want to proceed.

---

## 🌟 What's New in v0.2.0

- 🧩 **Modular Block Architecture:** All core subsystems (`Provider`, `Parser`, `Pipeline`, `Storage`, `UI`, `Tracker`) are cleanly abstracted interfaces ("blocks"). Fork developers can effortlessly swap out any block.
- 🆓 **Free-Tier Default Model:** Defaults to `gemini-2.5-flash`, which is eligible for the free tier on Google AI Studio without billing surprises.
- 💻 **Local Model Support:** First-class support for **Ollama**, **LM Studio**, and local OpenAI-compatible endpoints (`llama3.2`, `mistral`, `qwen2.5-coder`) with zero external API calls or fees.
- 📊 **Mistakes Caught & Accuracy Tracking:** Track every intercepted command, accuracy percentage, and typing skill improvement over time via `sure stats`.
- 🌐 **Web Dashboard & Showcase Badge:** Launch `sure web` to view an interactive dark-mode dashboard and embed a live SVG badge in your GitHub README / dotfiles repo!

---

## 🚀 Quick Install

### Homebrew (macOS / Linux)

```bash
brew tap sidx1-scratch/homebrew-tap
brew install sure
```

### From Source

```bash
git clone https://github.com/sidx1-scratch/sure.git
cd sure
make install
```

---

## 🛠️ Setup Instructions

1. **Run the interactive setup wizard:**
   ```bash
   sure setup
   ```
   Provide your Gemini API key (free-tier eligible, securely stored in OS keyring) or press Enter if using local models.

2. **Add shell integration to your profile:**

   **Bash (`~/.bashrc`):**
   ```bash
   eval "$(sure shell-init --shell bash)"
   ```

   **Zsh (`~/.zshrc`):**
   ```bash
   eval "$(sure shell-init --shell zsh)"
   ```

3. **Restart your shell:**
   ```bash
   exec $SHELL
   ```

---

## 💻 Using Local Models (Ollama / LM Studio)

You can run `sure` 100% offline with zero API keys or costs:

```bash
# Switch provider to local / ollama
sure config set provider local

# Configure local endpoint (default: http://localhost:11434)
sure config set local.endpoint http://localhost:11434

# Select model (e.g. llama3.2, mistral, qwen2.5-coder)
sure config set local.model llama3.2
```

Verify with:
```bash
sure doctor
```

---

## 📊 Viewing Mistakes Caught & Typing Improvement

Check your command accuracy and typing improvement directly in your terminal:

```bash
sure stats
```

Example output:
```text
📊 sure — Terminal Command Safety & Typing Improvement
=======================================================
Total Commands Analyzed:    142
🛑 Mistakes Caught:         4 (intercepted & avoided)
⚠️  Warnings Overridden:     1 (user confirmed Y)
✨ Safe Commands:           137 (clean execution)
🎯 Typing Accuracy Score:   97.2%
📈 Improvement Trend:       improving
```

Launch the web dashboard:
```bash
sure web
```
Open **http://localhost:7873** to see mistake history, accuracy trends, and embed your showcase badge:
```html
<img src="http://localhost:7873/api/badge" alt="sure safety badge" />
```

---

## 🧩 Modular Block Architecture (For Fork Developers)

`sure` is built as independent, pluggable "blocks" so open-source contributors can easily swap components without modifying core logic:

| Block Interface | Package | Default Implementation | Swappable For |
|---|---|---|---|
| `ai.Provider` | `internal/ai` | `GeminiProvider`, `LocalProvider` | Anthropic, Mistral, custom HTTP API |
| `parser.Parser` | `internal/parser` | `DefaultParser` (tokenizer) | `tree-sitter-bash`, `mvdan.cc/sh` |
| `analyzer.Pipeline` | `internal/analyzer` | `DefaultPipeline` | Custom heuristic filters, AST rules |
| `cache.Storage` | `internal/cache` | `FileStorage` (~/.cache/sure) | SQLite, Redis, BadgerDB |
| `tui.UI` | `internal/tui` | `TerminalUI` (ANSI Box) | Bubbletea, Gum, Desktop notifications |
| `stats.Tracker` | `internal/stats` | `FileTracker` (JSONL) | SQLite, InfluxDB, Prometheus |

Each block provides `SetGlobal<Block>()` or `RegisterProvider()` for trivial extension.

---

## 💡 Warning Examples

When a warning is necessary, `sure` provides a concise explanation of **what could happen and why**:

```text
⚠ This command may be destructive.

Command:
  rm -rf ~/Downloads/*

Why:
  This recursively deletes everything inside your Downloads
  folder without asking for confirmation.

Do you want to proceed?

[Y] Yes   [N] No
```

If you select **Yes**, `sure` executes the **EXACT original command**.
If you select **No**, it cancels execution and records a caught mistake in your accuracy stats.

---

## ⚙️ Configuration Reference

Stored at `~/.config/sure/config.yaml`:

```yaml
provider: gemini         # "gemini", "local", "ollama", "openai"
model: gemini-2.5-flash  # Free-tier eligible on Google AI Studio
sensitivity: medium      # "low", "medium", "high"
enabled: true
analyze_pipelines: true
local:
  endpoint: http://localhost:11434
  model: llama3.2
web_port: 7873
excluded_commands:
  - cd
  - ls
  - pwd
  - echo
always_confirm:
  - rm -rf /
  - mkfs
```

---

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.
