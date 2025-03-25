# backup_ios

**backup_ios** automates the forensic backup and analysis process for iOS devices using the [Mobile Verification Toolkit (MVT)](https://github.com/mvt-project/mvt). This tool streamlines the workflow for forensic investigators by automating encrypted backups, decryption, and subsequent IOC-based analysis, all in one go.

## Overview

The Mobile Verification Toolkit (MVT) is a suite of utilities designed to simplify and automate the process of gathering forensic traces from mobile devices—traces that can help identify potential compromises on Android and iOS devices. Developed by the [Amnesty International Security Lab](https://securitylab.amnesty.org) in July 2021 (in the context of the [Pegasus Project](https://forbiddenstories.org/about-the-pegasus-project/)) along with a robust [forensic methodology](https://www.amnesty.org/en/latest/research/2021/07/forensic-methodology-report-how-to-catch-nso-groups-pegasus/), MVT continues to evolve with contributions from Amnesty International and the wider digital forensic community.

> **Note:**  
> MVT is a forensic research tool intended for technologists and investigators. It requires a solid understanding of digital forensics and command-line operations. This is **not** intended for end-user self-assessment. If you are concerned about the security of your device, please seek reputable expert assistance.

## Indicators of Compromise (IOCs)

MVT can leverage public [indicators of compromise (IOCs)](https://github.com/mvt-project/mvt-indicators) to scan mobile devices for traces of known spyware campaigns. These IOCs include datasets published by [Amnesty International](https://github.com/AmnestyTech/investigations/) and other research groups, enabling targeted forensic analysis.

## Features

- **Automated iOS Backup:**  
  Leverages `idevicebackup2` to create encrypted backups of iOS devices.
- **Encryption Management:**  
  Automatically enables and disables encryption on the device as part of the backup workflow.
- **Backup Decryption & Analysis:**  
  Uses `mvt-ios` to decrypt backups and run forensic checks.
- **IOC Integration:**  
  Updates and applies indicators of compromise to scan the decrypted backup for signs of compromise.
- **Realtime Output:**  
  Displays the backup process output in real time to provide transparency and progress feedback.

## Requirements

- **Platform:**  
  Designed for modern Apple Silicon macOS; untested on other platforms.
- **Software Dependencies:**  
  - [Python3](https://www.python.org/)
  - [pipx](https://github.com/pipxproject/pipx)
  - [libusb](https://libusb.info/)
  - [SQLite3](https://www.sqlite.org/index.html)
- **Additional for Android (Optional):**  
  - Android SDK Platform Tools

### Installing Dependencies via Homebrew

Install the required packages with:

```bash
brew install python3 pipx libusb sqlite3
```

For Android device support (if needed):
```bash
brew install --cask android-platform-tools
```

Next, ensure you have pipx installed and properly set up:
```bash
pipx ensurepath
```

Then install the Mobile Verification Toolkit:
```bash
pipx install mvt
```

This installs the mvt-ios and mvt-android utilities.

### Installation
Clone or download this repository, then build the tool using Go:
```bash
go build backup_ios.go
```
This will create the backup_ios executable.


### Usage
1. Connect and Trust Your Device: Make sure your iOS device is connected and trusted by your computer.
2. Run the Tool: Execute the command:

```bash
./backup_ios
```

3. Follow the Prompts:
* Select the external drive where the backup will be stored.
* If encryption is already enabled, provide the existing backup password when prompted.
* The tool will then:
  * Create a timestamped backup directory.
  * Enable encryption on your device.
  * Perform a realtime backup using idevicebackup2.
  * Disable encryption (unlock the phone) once the backup is complete.
  * Decrypt the backup using mvt-ios.
  * Update the IOC list and run a forensic scan on the decrypted backup.

4. Review Results:
The tool logs its progress and results, and outputs diagnostic messages throughout the process.
