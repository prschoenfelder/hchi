package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func checkMemory(summary bool) {
	// For platform independent memory check, we use runtime package
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Convert to megabytes
	allocMB := m.Alloc / 1024 / 1024

	if summary {
		if allocMB > 100 {
			fmt.Println("❌ [MEM]  High tool memory usage detected")
		} else {
			fmt.Println("✔ [MEM]  Tool memory footprint is stable")
		}
	} else {
		fmt.Printf("--- Memory ---\n")
		fmt.Printf("RAM used by tool:     %d MB\n", allocMB)
		fmt.Printf("System architecture:  %s\n", runtime.GOARCH)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "check": // run the health checks
		runAllChecks()
	case "cron": // create an hourly cron/task entry
		setupAutomatedCheck()
	case "env": // display environment variable value
		if len(args) < 2 {
			fmt.Println("Please provide an environment variable name.")
			os.Exit(1)
		}
		showEnv(args[1])
	case "help": // print usage details
		printHelp()
	case "login": // add hchi to zsh startup
		setupShellStartup()
	case "mem": // show detailed memory usage
		checkMemory(false)
	case "sys": // show basic system info
		showSystemInfo()
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("hchi - A simple system information and environment variable fetcher")
	fmt.Println("Usage:")
	fmt.Println("  hchi check       Run all health checks")
	fmt.Println("  hchi cron        Setup hourly automated check")
	fmt.Println("  hchi env <VAR>   Show value of environment variable VAR")
	fmt.Println("  hchi help        Show this help message")
	fmt.Println("  hchi login       Add hchi to ~/.zshrc for shell startup")
	fmt.Println("  hchi mem         Show detailed memory usage")
	fmt.Println("  hchi sys         Show system information")
}

func runAllChecks() {
	fmt.Println("=== SYSTEM HEALTH CHECK ===")
	checkMemory(true)
}

func setupAutomatedCheck() {
	fmt.Println("Setting up an automated health check...")

	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error finding executable path: %v\n", err)
		return
	}

	switch runtime.GOOS {
	case "windows":
		setupWindowsTask(exePath)
	case "linux", "darwin":
		setupCronJob(exePath)
	default:
		fmt.Printf("Automated checks are not supported on this platform: %s\n", runtime.GOOS)
	}
}

func setupCronJob(exePath string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error finding home directory: %v\n", err)
		return
	}

	logPath := fmt.Sprintf("%s/.hchi.log", home)
	cronEntry := fmt.Sprintf("0 * * * * \"%s\" check >> %s 2>&1\n", exePath, logPath)

	cmd := exec.Command("crontab", "-l")
	currentCrontab, err := cmd.Output()
	if err != nil {
		currentCrontab = []byte{}
	}

	if strings.Contains(string(currentCrontab), fmt.Sprintf("\"%s\" check", exePath)) {
		fmt.Println("✅ Cron job already exists. hchi will run the 'check' command every hour.")
		return
	}

	newCrontab := string(currentCrontab)
	if len(newCrontab) > 0 && !strings.HasSuffix(newCrontab, "\n") {
		newCrontab += "\n"
	}
	newCrontab += cronEntry

	setCmd := exec.Command("crontab", "-")
	stdin, err := setCmd.StdinPipe()
	if err != nil {
		fmt.Printf("❌ Error opening crontab pipeline: %v\n", err)
		return
	}

	err = setCmd.Start()
	if err != nil {
		fmt.Printf("❌ Error starting crontab command: %v\n", err)
		return
	}

	_, err = stdin.Write([]byte(newCrontab))
	if err != nil {
		fmt.Printf("❌ Error writing to crontab: %v\n", err)
		return
	}
	stdin.Close()

	err = setCmd.Wait()
	if err != nil {
		fmt.Printf("❌ Error saving crontab: %v\n", err)
		return
	}

	fmt.Printf("✅ Cron job created successfully. hchi will run the 'check' command every hour and append output to %s.\n", logPath)
}

func setupShellStartup() {
	home, _ := os.UserHomeDir()
	zshrcPath := fmt.Sprintf("%s/.zshrc", home)
	exePath, _ := os.Executable()
	entry := fmt.Sprintf("\n# Run hchi on terminal startup\n\"%s\" check\n", exePath)

	contents, err := os.ReadFile(zshrcPath)
	if err == nil && strings.Contains(string(contents), entry) {
		fmt.Println("✅ hchi is already configured in ~/.zshrc.")
		return
	}

	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		f.WriteString(entry)
		fmt.Println("✅ Added to ~/.zshrc. You will see system stats every time you open a terminal!")
	} else {
		fmt.Printf("❌ Could not update ~/.zshrc: %v\n", err)
	}
}

func setupWindowsTask(exePath string) {
	cmd := exec.Command("schtasks", "/create", "/sc", "hourly", "/tn", "HchiCheck", "/tr", "\""+exePath+"\" check")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Error creating scheduled task: %v\n", err)
		fmt.Println(string(output))
		fmt.Println("Please run this command with administrator privileges.")
		return
	}

	fmt.Println("✅ Scheduled task created successfully. hchi will run the 'check' command every hour.")
}

func showEnv(key string) {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Environment variable '%s' is empty or not set.\n", key)
		return
	}
	fmt.Printf("'%s': %s\n", key, value)
}

func showSystemInfo() {
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Printf("Architecture:     %s\n", runtime.GOARCH)
	fmt.Println("CPU Cores:       ", runtime.NumCPU())
}
