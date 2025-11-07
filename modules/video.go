package modules

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
)

func IsVideo(filepath string) bool {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v",
		"-count_frames",
		"-show_entries", "stream=nb_read_frames",
		"-of", "default=nokey=1:noprint_wrappers=1",
		filepath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	out := strings.TrimSpace(string(output))
	if out == "" {
		return false
	}
	frame, _ := strconv.Atoi(out)
	return frame > 1
}

func getDurationVideo(filepath string) float64 {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filepath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0
	}

	out := strings.TrimSpace(string(output))
	if out == "" {
		return 0
	}
	duration, err := strconv.ParseFloat(out, 64)
	if err != nil {
		return 0
	}

	return duration
}

func Execute(command *exec.Cmd, inputFile string) int {
	stderrPipe, _ := command.StderrPipe()
	reader := bufio.NewReader(stderrPipe)

	if err := command.Start(); err != nil {
		fmt.Println("Обработчик не запустился:", err)
		return 1
	}

	duration := getDurationVideo(inputFile)

	for {
		line, err := reader.ReadString('\r')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Ошибка чтения stderr: %v\n", err)
			}
			break
		}
		timeToString := getTimeToString(line)
		if timeToString == "" {
			continue
		}
		timeToSeconds(timeToString)

		progressBar(timeToSeconds(timeToString), duration)
	}

	if err := command.Wait(); err != nil {
		return 1
	}
	progressBar(duration, duration)

	fmt.Printf("\n\n%s Успешно обработано!%s\n", GREEN, RESET)
	return 0
}

func getTimeToString(line string) string {
	pattern := regexp.MustCompile(`time=(\d{2}:\d{2}:\d{2}\.\d{2})`)
	matches := pattern.FindStringSubmatch(line)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func timeToSeconds(time string) float64 {
	time = strings.TrimSpace(time)
	parts := strings.Split(time, ":")
	if len(parts) == 3 {
		h, _ := strconv.ParseFloat(parts[0], 64)
		m, _ := strconv.ParseFloat(parts[1], 64)
		s, _ := strconv.ParseFloat(parts[2], 64)
		return h*3600 + m*60 + s
	}
	if len(parts) == 2 {
		m, _ := strconv.ParseFloat(parts[0], 64)
		s, _ := strconv.ParseFloat(parts[1], 64)
		return m*60 + s
	}
	s, _ := strconv.ParseFloat(parts[0], 64)
	return s
}

func progressBar(partTime, fullTime float64) {
	lengthBar := 75
	progress := partTime / fullTime
	if progress >= 1 {
		progress = 1
	}

	filled := int(progress * float64(lengthBar))
	bar := strings.Repeat("█", filled) + strings.Repeat("-", lengthBar-filled)
	fmt.Printf("\r [%s] %.1f%%", bar, progress*100)
}
