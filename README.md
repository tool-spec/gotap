# gotap

[![DOI](https://zenodo.org/badge/1104655184.svg)](https://doi.org/10.5281/zenodo.18658795)

CLI shim for [tool-spec](https://tool-spec.github.io/tool-specs/) compliant tools. It runs inside Docker containers and handles metadata verification, conversion, input parsing, and parameter binding generation.

Part of the [tool-spec](https://github.com/tool-spec) ecosystem. See the [specification](https://tool-spec.github.io/tool-specs/) for context.

## Templates

These templates will use gotap (or already plan to):

- [tool_template_python](https://github.com/tool-spec/tool_template_python)
- [tool_template_r](https://github.com/tool-spec/tool_template_r)
- [tool_template_node](https://github.com/tool-spec/tool_template_node)
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

## Bindings

`gotap generate` produces typed parameter access code for Python, R, JavaScript, and MATLAB. Each binding runs `gotap parse` and returns structured parameters and data paths.

### Generate usage

```bash
gotap generate --target=python --output=src/parameters.py --spec-file=src/tool.yml
gotap generate --target=r --output=src/parameters.R --spec-file=src/tool.yml
gotap generate --target=javascript --output=src/parameters.ts --spec-file=src/tool.yml
gotap generate --target=matlab --output=get_parameters.m --spec-file=src/tool.yml
```

If the spec defines a single tool, you can omit the tool name. Otherwise pass it as the first argument or set `RUN_TOOL`.

The `generate` command overwrites the output file if it already exists.

### Development workflow

1. Build gotap:
   ```bash
   go build -o gotap .
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
go build -o gotap .
```
