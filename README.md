![image](https://user-images.githubusercontent.com/63486672/119204908-992b5700-ba4b-11eb-80d4-5450cf4d48b4.png)

## 🪐 Golang WIP POC Windows botnet written for practice
##### *Say THAT five times fast..*
### Features
* Built for target and control via Windows OS ⌨️
* Unencrypted IRC communications 🐶
* Powershell command execution 🐚
* Mass (all bots) or targeted control 🚀
* Built to utilize freenode, can be tweaked to use TOR 🦊
##### *Only can target windows, the master MIGHT be able to run on linux*
## 📝 Installation

After installing [golang](https://golang.org/doc/install), use go to fetch the required dependencies

```bash
go get github.com/gianarb/go-irc
go get gopkg.in/sorcix/irc.v2
```
## 👽 Configuration

Create a channel on [freenode](https://webchat.freenode.net/), and enter it in both go files
###### The master name is not important, it is simply the name that your master.exe uses to connect
#### To edit the bot.go file
| ![image](https://user-images.githubusercontent.com/63486672/119205486-30dd7500-ba4d-11eb-8c47-ca1d89aca6f5.png) |
| :------: |
Edit line 23
#### To edit the master.go file
| ![image](https://user-images.githubusercontent.com/63486672/119205280-b44a9680-ba4c-11eb-9c1b-b05896176602.png) |
| :------: |
Edit lines 25 & 26
##### ⚠️ Ensure that both channels match between these files before compilation

## 🌔 Compilation
Afterwards, compile, or simply run prism.go to generate the bot and master stub for your system
```bash
go build prism.go -o prism.exe
/// Or, to skip compilation
go run prism.go
```
❗ Note, if compiling, the produced prism.exe file must be ran from a terminal window (powershell/cmd)

## :octocat: Usage
To view or control bots, simply run the master.exe from within a terminal

```bash
./master.exe
```
| ![image](https://user-images.githubusercontent.com/63486672/119274294-07932500-bbc4-11eb-9572-e345004ec7bf.png) |
| :------: |
🙊 🙈 🙉
## License
[MIT](https://choosealicense.com/licenses/mit/)
