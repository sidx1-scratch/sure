# ✋ sure

> **Your terminal's second opinion**

`sure` is an AI-powered CLI utility that intercepts terminal commands before they execute, analyzes them for potential destructive behavior, mistakes, or dangerous typos, and asks for confirmation before letting you proceed.

**The most important rule:** `sure` NEVER automatically modifies, corrects, replaces, or rewrites the user's command. It only asks whether you want to proceed.

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
   This displays a clean TUI box wizard asking for your Gemini API key (stored securely in the OS keyring).

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

For a suspected typo:

```text
⚠ I think you may have typed this command incorrectly.

Command:
  git chekcout main

Why:
  `chekcout` does not appear to be a valid Git subcommand.

Do you want to proceed?

[Y] Yes   [N] No
```

If you select **Yes**, `sure` executes the **EXACT original command**.
If you select **No**, it cancels execution and returns directly to the shell.

---

## ⚙️ Configuration

Configuration is stored at `~/.config/sure/config.yaml`:

```yaml
provider: gemini
model: gemini-2.0-flash
sensitivity: medium      # low, medium, high
enabled: true
analyze_pipelines: true
excluded_commands:
  - cd
  - ls
  - pwd
  - echo
always_confirm:
  - rm -rf /
  - mkfs
```

### CLI Config Commands

```bash
sure status                    # Inspect provider, key status, cache, & shell hook
sure config show               # View current YAML settings
sure config set sensitivity high # Change setting
sure config reset              # Reset to default configuration
sure scan                      # Discover and cache local CLI documentation
sure doctor                    # Diagnose environment, API key, and connectivity
```

### Temporarily Disabling

Temporarily toggle interception directly in your active shell session:

```bash
sure_disable   # Interception paused
sure_enable    # Interception resumed
```

---

## 🏗️ Architecture & How It Works

```text
User enters command
        ↓
Shell Hook (Bash DEBUG extdebug / Zsh preexec)
        ↓
sure analyze --command "..."
        ↓
Command Parser & Tokenizer
        ↓
CLI Documentation Discovery & Cache (~/.cache/sure/docs)
        ↓
AI Analysis (Gemini API with structured JSON output)
        ↓
Warning + Explanation (if issue detected)
        ↓
User choice [Y/N]
        ↓
Original unmodified command executed by Shell
```

---

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.
