package main

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/Danoxv/storkvid/modules"
)

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
)

func main() {
	crf := flag.Int("crf", 23, "Constant Rate Factor (качество видео)")
	noAudio := flag.Bool("no-audio", false, "Выключить аудио")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Printf("%s 🗄️  Передайте видеофайл! %s\n\n", RED, RESET)
		return
	}
	inputPath := args[0]

	if !modules.IsVideo(inputPath) {
		fmt.Printf("%s 🗄️  Ожидаю видеофайл! %s\n\n", RED, RESET)
		return
	}

	fmt.Printf("%s 🚀 Начинаем обработку видео... %s\n\n", GREEN, RESET)

	fmt.Printf("%s 🎛  Параметры: %d\n%s", YELLOW, *crf, RESET)
	fmt.Printf("   %s • CRF: %d\n%s", YELLOW, *crf, RESET)
	fmt.Printf("   %s • Аудио: %t\n\n%s", YELLOW, !*noAudio, RESET)

	arrayCommand := []string{
		"-i", inputPath,
		"-c:v", "h264",
		"-crf", strconv.Itoa(*crf),
	}

	if *noAudio {
		arrayCommand = append(arrayCommand, "-an")
	}

	outputPath := filepath.Clean(inputPath) + ".mp4"
	arrayCommand = append(arrayCommand, outputPath)

	cmd := exec.Command("ffmpeg", arrayCommand...)

	modules.Execute(cmd, inputPath)
}
