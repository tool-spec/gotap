package input

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hydrocode-de/gotap/internal/config"
	toolspec "github.com/hydrocode-de/tool-spec-go"
	"github.com/shirou/gopsutil/v3/process"
)

type ResolvedCommand struct {
	Command    string
	Executable string
	Extension  string
}

type ExecutionResult struct {
	Stdout        []byte        `json:"-"`
	Stderr        []byte        `json:"-"`
	ExitCode      int           `json:"exit_code"`
	UserTime      time.Duration `json:"user_time"`
	SystemTime    time.Duration `json:"system_time"`
	MemoryMax     uint64        `json:"memory_max_bytes"`
	MemoryAverage uint64        `json:"memory_average_bytes"`
	CPUMax        uint64        `json:"cpu_max_permille"`
	CPUAverage    uint64        `json:"cpu_average_permille"`
	ReadBytesSum  uint64        `json:"read_bytes_sum"`
	WriteBytesSum uint64        `json:"write_bytes_sum"`
	RunID         string        `json:"run_id,omitempty"`
	LogFile       string        `json:"log_file,omitempty"`
	LogFormat     string        `json:"log_format,omitempty"`
	LogLevel      string        `json:"log_level,omitempty"`
}

func isExecutable(path string) bool {
	_, err := exec.LookPath(path)
	return err == nil
}

func isFile(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func parseCommand(command string) (ResolvedCommand, error) {
	var chunks []string
	for _, c := range strings.Split(command, " ") {
		if c != "" {
			chunks = append(chunks, c)
		}
	}
	if len(chunks) == 0 {
		return ResolvedCommand{}, fmt.Errorf("the command is empty")
	}

	executable := ""
	fileToken := ""
	if isExecutable(chunks[0]) {
		executable = chunks[0]
		for _, chunk := range chunks[1:] {
			if isFile(chunk) {
				fileToken = chunk
				break
			}
		}
	} else if isFile(chunks[0]) {
		executable = ""
		fileToken = chunks[0]
	}

	if executable == "" {
		return ResolvedCommand{}, fmt.Errorf("the command is not a valid executable or file")
	}

	return ResolvedCommand{
		Command:    command,
		Executable: executable,
		Extension:  filepath.Ext(fileToken),
	}, nil
}

func parseMatch(match string) (ResolvedCommand, error) {
	fileExtension := strings.ToLower(filepath.Ext(match))
	executable := ""
	switch fileExtension {
	case ".sh":
		executable = "sh"
	case ".py":
		executable = "python3"
	case ".R":
		executable = "Rscript"
	case ".jl":
		executable = "julia"
	case ".pl":
		executable = "perl"
	case ".m", ".matlab":
		if isExecutable("matlab") {
			executable = "matlab"
		} else {
			executable = "octave"
		}
	case ".js":
		executable = "node"
	default:
		executable = match
	}

	command := ""
	if executable == match {
		command = executable
	} else {
		command = fmt.Sprintf("%s %s", executable, match)
	}

	return ResolvedCommand{
		Command:    command,
		Executable: executable,
		Extension:  fileExtension,
	}, nil
}

func InferBindingTarget(executable string) (target string, outputFile string, ok bool) {
	exe := strings.ToLower(strings.TrimSpace(executable))
	switch exe {
	case "python", "python3":
		return "python", "parameters.py", true
	case "rscript":
		return "r", "parameters.R", true
	case "node":
		return "node-js", "parameters.js", true
	case "tsx", "ts-node":
		return "node-ts", "parameters.ts", true
	case "matlab", "octave":
		return "matlab", "get_parameters.m", true
	default:
		return "", "", false
	}
}

func ResolveCommand(spec toolspec.ToolSpec) (ResolvedCommand, error) {
	if spec.Command != "" {
		return parseCommand(spec.Command)
	}
	if cmd := os.Getenv("TAP_COMMAND"); cmd != "" {
		return parseCommand(cmd)
	}

	// now we search for ./ run, run.sh, run.* in that order
	directories := make([]string, 0, 2)
	wd, err := os.Getwd()
	if err == nil {
		directories = append(directories, wd)
	}
	specFilePath := config.GetViper().GetString("spec_file")
	directories = append(directories, filepath.Dir(specFilePath))

	for _, directory := range directories {
		matches, err := filepath.Glob(filepath.Join(directory, "run*"))
		if err != nil {
			continue
		}
		var lastMatch ResolvedCommand
		foundAny := false
		foundBash := false
		for _, match := range matches {
			resolved, err := parseMatch(match)
			if err != nil {
				continue
			}
			if resolved.Extension == "" {
				return resolved, nil
			}
			if resolved.Extension == ".sh" {
				foundBash = true
				lastMatch = resolved
			} else if !foundAny && !foundBash {
				foundAny = true
				lastMatch = resolved
			}
		}
		if foundAny || foundBash {
			return lastMatch, nil
		}
	}

	// if we reach this, we never could parse a match
	return ResolvedCommand{}, fmt.Errorf("the command could not be found. Consider adding it to your tool.yml")
}

func ExecuteCommand(command ResolvedCommand, extraEnv []string) (ExecutionResult, error) {
	cmd := exec.Command("sh", "-c", command.Command)
	cmd.Env = append(os.Environ(), extraEnv...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("failed to execute command: %w", err)
	}

	pid := int32(cmd.Process.Pid)
	proc, err := process.NewProcess(pid)
	if err != nil {
		return ExecutionResult{}, fmt.Errorf("failed to create process: %w", err)
	}

	var memSamples []uint64
	var cpuSamples []uint64

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// this actually runs in a goroutine and waits for the command to finish
	done := make(chan bool)
	go func() {
		cmd.Wait()
		close(done)
	}()

	sampling := true
	for sampling {
		select {
		case <-done:
			ticker.Stop()
			sampling = false
		case <-ticker.C:
			mem, err := proc.MemoryInfo()
			if err == nil {
				memSamples = append(memSamples, mem.RSS)
			}
			cpu, err := proc.CPUPercent()
			if err == nil {
				cpuSamples = append(cpuSamples, uint64(cpu*1000))
			}
		}
	}

	var readBytesSum, writeBytesSum uint64
	ioCounters, err := proc.IOCounters()
	if err == nil {
		readBytesSum = ioCounters.ReadBytes
		writeBytesSum = ioCounters.WriteBytes
	}

	exitCode := cmd.ProcessState.ExitCode()
	return ExecutionResult{
		Stdout:        stdout.Bytes(),
		Stderr:        stderr.Bytes(),
		ExitCode:      exitCode,
		UserTime:      cmd.ProcessState.UserTime(),
		SystemTime:    cmd.ProcessState.SystemTime(),
		MemoryMax:     calcualateMax(memSamples),
		MemoryAverage: calcualteAverage(memSamples),
		CPUMax:        calcualateMax(cpuSamples),
		CPUAverage:    calcualteAverage(cpuSamples),
		ReadBytesSum:  readBytesSum,
		WriteBytesSum: writeBytesSum,
	}, nil
}

func calcualateMax(samples []uint64) uint64 {
	if len(samples) == 0 {
		return 0 // Return 0 for missing CPU data
	}
	max := samples[0]
	for _, sample := range samples {
		if sample > max {
			max = sample
		}
	}
	return max
}

func calcualteSum(samples []uint64) uint64 {
	sum := uint64(0)
	for _, sample := range samples {
		sum += sample
	}
	return sum // Sum is safe even with empty array (returns 0)
}

func calcualteAverage(samples []uint64) uint64 {
	if len(samples) == 0 {
		return 0 // Return 0 for missing CPU data
	}
	sum := uint64(0)
	for _, sample := range samples {
		sum += sample
	}
	return sum / uint64(len(samples))
}
