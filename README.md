# gotap

[![DOI](https://zenodo.org/badge/1104655184.svg)](https://doi.org/10.5281/zenodo.18658795)

CLI shim for [tool-spec](https://tool-spec.github.io/tool-specs/) compliant tools. It runs inside Docker containers and handles metadata verification, conversion, input parsing, and parameter binding generation.

Part of the [tool-spec](https://github.com/tool-spec) ecosystem. See the [specification](https://tool-spec.github.io/tool-specs/) for context.

## Templates

These templates will use gotap (or already plan to):

- [tool_template_python](https://github.com/tool-spec/tool_template_python)
- [tool_template_r](https://github.com/tool-spec/tool_template_r)
- [tool_template_node](https://github.com/tool-spec/tool_template_node)
- [tool_template_typescript](https://github.com/tool-spec/tool_template_typescript)
- [tool_template_octave](https://github.com/tool-spec/tool_template_octave)
- [tool_template_matlab](https://github.com/tool-spec/tool_template_matlab)
- [tool_template_jupyter](https://github.com/tool-spec/tool_template_jupyter)

## Commands

| Command     | Purpose |
|-------------|---------|
| `parse`     | Read `input.json` and `tool.yml`, validate, output validated parameters and data as JSON to stdout |
| `generate`  | Generate language-specific parameter bindings (Python, R, JavaScript, MATLAB) that call `gotap parse` |
| `metadata`  | Output tool metadata in schema.org or nfdi4earth format |
| `verify`    | Validate tool-spec metadata and report errors/warnings |
| `prepare`   | Create or update `inputs.json` via interactive CLI |
| `run`       | Validate inputs and execute the tool entrypoint |

Flags are inherited from the root where relevant: `--spec-file`, `--input-file`, `--citation-file`, `--license-file`, `--output-folder`.
Version check: `gotap --version` or `gotap -v`.

## Unified Logging

`gotap run` now creates a structured logging context for the tool process and injects it via environment variables:

- `GOTAP_RUN_ID`
- `GOTAP_LOG_FILE`
- `GOTAP_LOG_FORMAT`
- `GOTAP_LOG_LEVEL`

The canonical format is JSON Lines written to `_logs.jsonl` in the output folder. `stdout` and `stderr` are still captured separately as before.

### Developer workflow

Generated bindings now expose runtime helpers alongside parameter parsing:

- `get_parameters()`
- `get_data()` where supported
- `get_run_context()` / `getRunContext()`
- `get_logger()` / `getLogger()`

Example in Python:

```python
from parameters import get_parameters, get_logger

params = get_parameters()
log = get_logger()

log.info("start", "Loading inputs")
log.warn("missing-value", "3 cells were NA", field="elevation")
log.info("done", "Finished successfully")
```

Each log event is appended as one JSON object per line:

```json
{"ts":"2026-03-15T10:12:01Z","level":"info","event":"start","message":"Loading inputs","run_id":"abc123"}
```

## Bindings

`gotap generate` produces parameter access code for Python, R, JavaScript (plain or TypeScript), and MATLAB. Each binding runs `gotap parse` and returns structured parameters and data paths.

### Generate usage

```bash
gotap generate --target=python --output=src/parameters.py --spec-file=src/tool.yml
gotap generate --target=r --output=src/parameters.R --spec-file=src/tool.yml
gotap generate --target=node-js --output=src/parameters.js --spec-file=src/tool.yml
gotap generate --target=node-ts --output=src/parameters.ts --spec-file=src/tool.yml
gotap generate --target=matlab --output=get_parameters.m --spec-file=src/tool.yml
```

- `node-js`: plain JavaScript (ESM), for Node.js projects
- `node-ts`: TypeScript with interfaces, for Node.js TypeScript projects

If the spec defines a single tool, you can omit the tool name. Otherwise pass it as the first argument or set `RUN_TOOL`.

The `generate` command overwrites the output file if it already exists.

### Development workflow

1. Build gotap:
   ```bash
   go build -ldflags "-X github.com/hydrocode-de/gotap/cmd.Version=$(git describe --tags --always --dirty)" -o gotap .
   ```

2. Regenerate bindings when `tool.yml` changes:
   ```bash
   ./gotap generate --target=python --output=src/parameters.py --spec-file=src/tool.yml
   ```

3. Either:
   - Add generated files to `.gitignore` and regenerate as part of the build, or
   - Commit the generated files so they are always up to date and usable offline.

## Build

```bash
go build -ldflags "-X github.com/hydrocode-de/gotap/cmd.Version=$(git describe --tags --always --dirty)" -o gotap .
```
