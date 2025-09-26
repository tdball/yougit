package main

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func GetTopLevelChanges() []string {
	topLevelFiles := []string{}
	diff := exec.Command("git", "diff", "--cached", "--name-only")
	output, err := diff.Output()
	if err != nil {
		logger.Error(string(output))
		logger.Error("Unable to get git diff")
		return topLevelFiles
	}

	fileNames := strings.SplitSeq(string(output), "\n")
	for fileName := range fileNames {
		topLevelFileName := strings.Split(fileName, "/")[0]
		if len(topLevelFileName) > 0 {
			topLevelFiles = append(topLevelFiles, topLevelFileName)
		}
	}
	return topLevelFiles
}

func IsDir(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		logger.Error("Unable to open", "error", err, "path", path)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.Error("Unable to describe file", "error", err)
	}
	return fileInfo.IsDir()
}

func ExecuteJustTest(path string) error {
	if IsDir(path) {
		logger.Info("Changing directory", "directory", path)
		os.Chdir(path)
		defer os.Chdir("..")
	}

	test := exec.Command("just", "test")
	output, err := test.CombinedOutput()
	if err != nil {
		logger.Error(
			"Failed to execute 'just test'",
			"error", err,
		)
		os.Exit(1)
	}
	logger.Info("Just Results:", "output", output)
	return nil
}

func main() {
	files := GetTopLevelChanges()
	logger.Info("Files Changed", "count", len(files), "files", strings.Join(files, ", "))
	for _, file := range files {
		ExecuteJustTest(file)
	}
}
