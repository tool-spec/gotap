# gotap

CLI shim for [tool-spec](https://tool-spec.github.io/tool-specs/) compliant tools. Runs inside Docker containers and handles metadata verification, conversion, input parsing, and parameter binding generation.

Part of the [tool-spec](https://github.com/tool-spec) ecosystem.

## Commands

- **parse** – Validate `input.json` against `tool.yml`, output parameters and data as JSON
- **generate** – Generate typed parameter bindings (Python, R, JavaScript, MATLAB)
- **metadata** – Output schema.org or nfdi4earth metadata
- **verify** – Validate tool-spec metadata
- **prepare** – Interactive CLI to create or update `inputs.json`
- **run** – Validate and execute the tool entrypoint (`--generate-bindings` to regenerate before run)

## Templates

Used by [tool_template_python](https://github.com/tool-spec/tool_template_python), [tool_template_r](https://github.com/tool-spec/tool_template_r), [tool_template_node](https://github.com/tool-spec/tool_template_node), [tool_template_octave](https://github.com/tool-spec/tool_template_octave), [tool_template_matlab](https://github.com/tool-spec/tool_template_matlab), [tool_template_jupyter](https://github.com/tool-spec/tool_template_jupyter).

## Build

```bash
go build -o gotap .
```
