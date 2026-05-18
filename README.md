# hchi

A small Go utility for showing system information, memory usage, and setting up automated health checks.

## Compile

From the repository root:

```bash
go build -o hchi .
```

This creates a `hchi` executable in the current directory.

## Install

If you want to install it to your Go bin directory:

```bash
go install .
```

Then make sure your Go bin directory is on your `PATH`.

> On Windows PowerShell, the current directory is not searched by default. That means `hchi` will only work if the installed executable is on your PATH. If you are in the project folder and have not installed it to PATH yet, use:
>
> ```powershell
> .\hchi help
> ```

## Usage

Run the binary directly, with `go run`, or after install from a PATH directory.

```bash
./hchi help
```

If the executable is installed and on your PATH, you can also run:

```bash
hchi help
```

### Verify install

After installing, run this to confirm:

```bash
hchi help
```

You should see output like:

```text
hchi - A simple system information and environment variable fetcher
Usage:
  hchi check       Run all health checks
  hchi cron        Setup hourly automated check
  hchi env <VAR>   Show value of environment variable VAR
  hchi help        Show this help message
  hchi login       Add hchi to ~/.zshrc for shell startup
  hchi mem         Show detailed memory usage
  hchi sys         Show system information
```

### Available commands

- `hchi check` — Run all health checks
- `hchi cron` — Setup hourly automated check (cron on macOS/Linux or Task Scheduler on Windows)
- `hchi env <VAR>` — Show value of an environment variable
- `hchi help` — Show usage information
- `hchi login` — Add hchi to `~/.zshrc` for shell startup
- `hchi mem` — Show detailed memory usage
- `hchi sys` — Show basic system information

## Notes

- On Linux/macOS, `hchi cron` appends output to `~/.hchi.log`.
- The `login` command appends a startup hook to `~/.zshrc` so the tool runs on new shell launch.
