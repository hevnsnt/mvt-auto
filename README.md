My attempt to automate the Mobile Verification Toolkit (https://github.com/mvt-project/mvt)

Mobile Verification Toolkit (MVT) is a collection of utilities to simplify and automate the process of gathering forensic traces helpful to identify a potential compromise of Android and iOS devices.

It has been developed and released by the [Amnesty International Security Lab](https://securitylab.amnesty.org) in July 2021 in the context of the [Pegasus Project](https://forbiddenstories.org/about-the-pegasus-project/) along with [a technical forensic methodology](https://www.amnesty.org/en/latest/research/2021/07/forensic-methodology-report-how-to-catch-nso-groups-pegasus/). It continues to be maintained by Amnesty International and other contributors.

> **Note**
> MVT is a forensic research tool intended for technologists and investigators. It requires understanding digital forensics and using command-line tools. This is not intended for end-user self-assessment. If you are concerned with the security of your device please seek reputable expert assistance.
>

### Indicators of Compromise

MVT supports using public [indicators of compromise (IOCs)](https://github.com/mvt-project/mvt-indicators) to scan mobile devices for potential traces of targeting or infection by known spyware campaigns. This includes IOCs published by [Amnesty International](https://github.com/AmnestyTech/investigations/) and other  research groups.


## Installation
Expectation to be run on a modern Apple Silicon, untested on any other platform
> brew install python3 pipx libusb sqlite3
(ibusb is not required if you intend to only use mvt-ios and not mvt-android.)
When working with Android devices you should additionally install Android SDK Platform Tools:

> brew install --cask android-platform-tools

Install pipx following the instructions above for your OS/distribution. Make sure to run pipx ensurepath and open a new terminal window.
> bash pipx install mvt

You now should have the mvt-ios and mvt-android utilities installed. If you run into problems with these commands not being found, ensure you have run pipx ensurepath and opened a new terminal window.
> go build backup_ios.go
