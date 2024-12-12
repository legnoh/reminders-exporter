# reminders-exporter

[![Static Badge](https://img.shields.io/badge/homebrew-legnoh%2Fetc%2Freminders--exporter-orange?logo=apple)](https://github.com/legnoh/homebrew-etc/blob/main/Formula/reminders-exporter.rb)

This tool provides daemon service for Apple Reminder app data exporter for Prometheus.

## Usage

Install, init, and start. That's it !

All configs are provided from `~/.reminders-exporter/config.yml` file.  
Create a configuration file with the following command and edit it.

- Config sample: [`sample/configs.yml`](./cmd/sample/configs.yml).

### macOS

This app is **ONLY** available on macOS and depends on [`reminders-cli`](https://github.com/keith/reminders-cli) to retrieve data.


```sh
# install
brew tap keith/formulae
brew install legnoh/etc/reminders-exporter

# init & edit
reminders-exporter init
vi ~/.reminders-exporter/config.yml

# start
brew services start reminders-exporter
```

> [!NOTE]
> The first time you retrieve data, a prompt will appear asking for permission to access your reminder data.
> <img width="372" alt="privary setting" src="https://github.com/user-attachments/assets/d4e775eb-dfbc-45e3-9b34-3ba8964385b2" />
