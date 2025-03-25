// This script works, but backup output is not shown.
package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GeneratePassword returns a secure random password of a given length.
func GeneratePassword(length int) (string, error) {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?"
	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		password[i] = chars[n.Int64()]
	}
	return string(password), nil
}

// Drive represents an external drive with its Path and Name.
type Drive struct {
	Path string
	Name string
}

// ListExternalDrives returns available external drives from /Volumes.
func ListExternalDrives() ([]Drive, error) {
	entries, err := os.ReadDir("/Volumes")
	if err != nil {
		return nil, err
	}
	var drives []Drive
	for _, e := range entries {
		if e.IsDir() && e.Name() != "Macintosh HD" {
			drives = append(drives, Drive{
				Path: filepath.Join("/Volumes", e.Name()),
				Name: e.Name(),
			})
		}
	}
	return drives, nil
}

// selectDrive prompts the user to choose from a list of drives.
func selectDrive(drives []Drive) Drive {
	fmt.Println("Available External Drives:")
	for i, drive := range drives {
		fmt.Printf("[%d] %s\n", i+1, drive.Name)
	}
	fmt.Print("Select drive by number: ")
	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(drives) {
		log.Fatalf("Invalid drive selection.")
	}
	return drives[choice-1]
}

// logAction appends a timestamped message to the specified log file.
func logAction(logFile, message string) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer f.Close()
	logEntry := fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), message)
	if _, err := f.WriteString(logEntry); err != nil {
		log.Println("Error writing log entry:", err)
	}
}

// runCommand logs and then executes a system command.
func runCommand(logFile, name string, args ...string) error {
	logAction(logFile, fmt.Sprintf("Executing: %s %s", name, strings.Join(args, " ")))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCommandWithOutput runs a command and returns its combined output.
func runCommandWithOutput(logFile, name string, args ...string) (string, error) {
	logAction(logFile, fmt.Sprintf("Executing: %s %s", name, strings.Join(args, " ")))
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	outStr := string(output)
	fmt.Print(outStr)
	return outStr, err
}

// checkDependencies verifies that required external programs are installed.
func checkDependencies() error {
	dependencies := []string{"idevicebackup2", "mvt-ios"}
	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			return fmt.Errorf("dependency '%s' not found. Install it before running this program.\nFor idevicebackup2: brew install libimobiledevice\nFor mvt-ios: pip install mvt", dep)
		}
	}
	fmt.Println("All dependencies are satisfied.")
	return nil
}

// isEncryptionEnabled checks whether backup encryption is currently enabled.
func isEncryptionEnabled() bool {
	// Use a wrong password to trigger a check.
	cmd := exec.Command("idevicebackup2", "backup", "--password", "wrongpassword", "/tmp/idevice_check")
	err := cmd.Run()
	os.RemoveAll("/tmp/idevice_check")
	return err == nil
}

// setupBackupDirectory creates a timestamped backup directory on the selected drive.
func setupBackupDirectory(basePath string) (string, error) {
	backupDir := filepath.Join(basePath, "ios_backup", time.Now().Format("20060102_150405"))
	err := os.MkdirAll(backupDir, 0755)
	return backupDir, err
}

// disableEncryption attempts to disable encryption using the given password.
func disableEncryption(logFile, password string) error {
	fmt.Println("Unlocking Phone...")
	if err := runCommand(logFile, "idevicebackup2", "encryption", "off", password); err != nil {
		logAction(logFile, fmt.Sprintf("Failed to disable encryption: %v", err))
		return err
	}
	fmt.Println("Encryption disabled successfully.")
	logAction(logFile, "Encryption disabled successfully.")
	return nil
}

// handleEncryption manages encryption: if enabled, disable it first, then enable with a new password.
func handleEncryption(logFile string) (string, error) {
	password, err := GeneratePassword(20)
	if err != nil {
		return "", fmt.Errorf("password generation failed: %w", err)
	}

	if isEncryptionEnabled() {
		fmt.Print("Backup encryption is enabled. Enter existing backup password: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		existingPassword := scanner.Text()
		if err := runCommand(logFile, "idevicebackup2", "encryption", "off", existingPassword); err != nil {
			// Log the error but continue to enable encryption with the new password.
			fmt.Println("Warning: Failed to disable existing encryption:", err)
			logAction(logFile, fmt.Sprintf("Warning: Failed to disable existing encryption: %v", err))
		}
	}

	if err := runCommand(logFile, "idevicebackup2", "encryption", "on", password); err != nil {
		return password, fmt.Errorf("failed to enable encryption: %w", err)
	}
	fmt.Printf("Encryption successfully enabled with password: %s\n", password)
	logAction(logFile, fmt.Sprintf("Encryption enabled with password: %s", password))
	return password, nil
}

// backupDevice executes the backup using idevicebackup2.
// A delay is added before starting the backup, and we capture the full output.
func backupDevice(logFile, backupDir, password string) error {
	// Allow the device to settle after encryption is enabled.
	time.Sleep(2 * time.Second)
	start := time.Now()
	fmt.Printf("Starting Backup of device...")
	_, err := runCommandWithOutput(logFile, "idevicebackup2", "backup", backupDir, "--password", password)
	if err != nil {
		return err
	}
	duration := time.Since(start)
	fmt.Printf("Backup completed successfully in %s\n", duration)
	logAction(logFile, fmt.Sprintf("Backup completed in %s", duration))
	return nil
}

// decryptBackup locates the actual backup subdirectory and decrypts the backup using mvt-ios.
func decryptBackup(logFile, backupDir, password string) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("no backup subdirectories found")
	}

	var backupPath string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), "._") {
			backupPath = filepath.Join(backupDir, entry.Name())
			break
		}
	}
	if backupPath == "" {
		return fmt.Errorf("no valid backup directory found")
	}

	decryptedDir := filepath.Join(backupDir, "decrypted")
	if err := os.MkdirAll(decryptedDir, 0755); err != nil {
		return err
	}
	if err := runCommand(logFile, "mvt-ios", "decrypt-backup", "-d", decryptedDir, "-p", password, backupPath); err != nil {
		return err
	}
	logAction(logFile, "Backup decrypted successfully.")
	return nil
}

// mainRun performs the backup workflow.
// After a successful backup it disables encryption, then decrypts the backup,
// updates the IOC list, and finally searches the decrypted backup.
func mainRun() error {
	// Check dependencies.
	if err := checkDependencies(); err != nil {
		return err
	}

	// List and select external drive.
	drives, err := ListExternalDrives()
	if err != nil || len(drives) == 0 {
		return fmt.Errorf("no external drives found or error accessing drives")
	}
	selectedDrive := selectDrive(drives)

	// Set up backup directory.
	backupDir, err := setupBackupDirectory(selectedDrive.Path)
	if err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	fmt.Printf("Backup directory is \"%s\"\n", backupDir)

	// Set up logging.
	logFile := filepath.Join(backupDir, "backup_log.txt")
	logAction(logFile, fmt.Sprintf("Selected backup directory: %s", backupDir))

	// Enable encryption and get the new password.
	encryptionPassword, err := handleEncryption(logFile)
	if err != nil {
		return err
	}

	// Save the encryption password to a file.
	passwordFile := filepath.Join(backupDir, "backup_password.txt")
	if err := os.WriteFile(passwordFile, []byte(encryptionPassword), 0600); err != nil {
		return fmt.Errorf("failed to write password file: %w", err)
	}
	logAction(logFile, fmt.Sprintf("Generated and saved new encryption password: %s", encryptionPassword))

	// Perform backup.
	if err := backupDevice(logFile, backupDir, encryptionPassword); err != nil {
		disableEncryption(logFile, encryptionPassword)
		return fmt.Errorf("backup failed: %w", err)
	}

	// After backup is complete, disable encryption.
	if err := disableEncryption(logFile, encryptionPassword); err != nil {
		return fmt.Errorf("failed to disable encryption after backup: %w", err)
	}

	// Decrypt the backup.
	if err := decryptBackup(logFile, backupDir, encryptionPassword); err != nil {
		return fmt.Errorf("backup decryption failed: %w", err)
	}

	// Update IOC list.
	fmt.Println("Updating IOC list...")
	if err := runCommand(logFile, "mvt-ios", "download-iocs"); err != nil {
		return fmt.Errorf("failed to update IOC list: %w", err)
	}
	fmt.Println("IOC list updated successfully.")
	logAction(logFile, "IOC list updated successfully.")

	// Search the decrypted backup.
	decryptedDir := filepath.Join(backupDir, "decrypted")
	fmt.Println("Searching decrypted backup...")
	if err := runCommand(logFile, "mvt-ios", "check-backup", decryptedDir); err != nil {
		return fmt.Errorf("failed to search decrypted backup: %w", err)
	}
	fmt.Println("Backup search completed.")
	logAction(logFile, "Backup search completed.")

	fmt.Println("Process completed successfully.")
	return nil
}

func main() {
	if err := mainRun(); err != nil {
		fmt.Println("Error encountered:", err)
		os.Exit(1)
	}
}
