# ✋ sure

> **Your terminal's second opinion**

`sure` is an AI-powered CLI tool that intercepts terminal commands before they execute, analyzes them for potential destructive behavior, mistakes, or dangerous typos, and asks for confirmation before letting you proceed. It's like having a senior engineer looking over your shoulder before you hit Enter on that `rm -rf` command.

## 🚀 Quick Install

### Homebrew (macOS/Linux)
```bash
brew tap sidx1/homebrew-tap
brew install sure
```

### Go Install
```bash
go install github.com/sidx1/sure@latest
```

### From Source
```bash
git clone https://github.com/sidx1/sure.git
cd sure
make install
```

## 🛠️ Setup Instructions

1. Run the interactive setup command:
   ```bash
   sure setup
   ```
   This will prompt you for your API key and configure default models.

2. Add shell integration. Add the appropriate line to your shell configuration file:

   **Bash (`~/.bashrc`):**
   ```bash
   eval "$(sure shell-init --shell bash)"
   ```

   **Zsh (`~/.zshrc`):**
   ```zsh
   eval "$(sure shell-init --shell zsh)"
   ```

3. Restart your shell or run `source ~/.bashrc` (or `~/.zshrc`).

## 💡 Usage Examples

When you type a potentially dangerous command, `sure` intercepts it:

```bash
$ rm -rf /etc/nginx/conf.d
✋ Warning: This command will recursively delete Nginx configuration files without prompting.
Are you sure you want to proceed? [y/N]: N
Command aborted.
```

Or for a potentially destructive git command:

```bash
$ git push origin master --force
✋ Warning: You are forcefully pushing to the master branch. This will overwrite remote history and may break the repository for other collaborators.
Are you sure you want to proceed? [y/N]: N
Command aborted.
```

## ⚙️ Configuration

You can manually edit the configuration file located at `~/.config/sure/config.yaml`:

```yaml
provider: anthropic
api_key: sk-ant-api03...
model: claude-3-5-haiku-20241022
temperature: 0.1
auto_allow: true      # Automatically allow commands that are deemed safe
exclude:              # Commands to skip analysis for (Regex)
  - "^ls"
  - "^git status"
```

## 📖 Commands Reference

- `sure setup`: Interactive setup for API keys and configuration.
- `sure status`: Check if shell integration is active and check configuration status.
- `sure config`: Manage configuration options (e.g., `sure config set provider openai`).
- `sure scan <script.sh>`: Analyze a shell script file for potential issues without running it.
- `sure doctor`: Verify installation, connectivity, and shell hooks.
- `sure shell-init --shell [bash|zsh]`: Output the shell initialization script.

## 🔌 Temporarily Disabling

If you want to temporarily disable `sure` without removing it from your config:

```bash
sure_disable
```

To re-enable:

```bash
sure_enable
```

Alternatively, you can set the environment variable `SURE_DISABLED=1` for a single command or session.

## 🏗️ Architecture & How It Works

`sure` uses shell pre-execution hooks to intercept commands before they are sent to the operating system. 

```
User types command -> Shell preexec hook -> `sure analyze` -> AI Model Evaluation
                                                                    |
    +---------------------------------------------------------------+
    |
    v
Safe? -> Yes -> Shell executes command
    |
    v
 No -> Prompt User -> User confirms -> Shell executes command
                   |
                   v
              User denies -> Command aborted
```

- **Bash**: Uses the `DEBUG` trap with `shopt -s extdebug`.
- **Zsh**: Uses the built-in `add-zsh-hook preexec`.

When a command is intercepted, `sure` checks local rules. If the command isn't trivially safe (like `cd` or `echo`), it sends the command to an LLM via your configured provider to assess its potential impact.

## 🧘 Design Philosophy

1. **Fast by default**: Shells need to be snappy. `sure` uses fast models (like Claude 3.5 Haiku or GPT-4o-mini) and avoids calling the API for harmless commands.
2. **Fail open**: If the API is down or there's a timeout, `sure` defaults to warning you but letting the command through, rather than locking you out of your terminal.
3. **Unobtrusive**: `sure` aims to be invisible until you're about to make a mistake.

## 🤝 Contributing

Contributions are welcome! Please open an issue or submit a Pull Request on the GitHub repository.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
