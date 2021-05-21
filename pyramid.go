package main
import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"syscall"
)


/* ###############################
###### GLOBAL VARIABLES!!!! ######
##################################*/

var ldflags string = `"-H=windowsgui -s -w"`
var PROMPT string = "\033[0;32m" + " PRISM > " + "\033[0m"
/*
var bot_url string = ""
var agent_name string = ""
var agent_url string = ""
var dropper_name string = ""
*/

func clearScreen() {
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func red() {
  fmt.Printf("\033[1;31m");
}

func green() {
  fmt.Printf("\033[0;32m");
}

func reset() {
  fmt.Printf("\033[0m");
}

func banner() {
	banner := `			      .
                             /_\
                     :      /_|_\
                    :::    /|__|_\
                   ::.::  /|_|__|_\      :
                  ::.:.::/__|_|__|_\    :.:
                 :..:.:./_|__|__|__|\  :.:.:
                :.:..:./|__|___|__|__\:.:..::
 ..............::..:../__|___|__|___|_\..:..::................
    ..........:..:..:/_|__|___|___|___|\:..:..::::::::::::::::::::
::::::::::::::.:..:./___|___|___|___|___\....................
        .........../..!...!...!...!...!..\...............
					  `
	banner2 := `-Pyramid v2.0-`
	fmt.Printf(banner)
	red()
	fmt.Println(banner2)
	reset()
}

func endbanner() {
	banner := `	    ,   ,
         ,-'{-'/
      ,-~ , \ {-~~-,	    Run bot.exe on target
    ,~  ,   ,',-~~-,',      Pilot bots via master.exe		
  ,'   ,   { {      } }     Godspeed..        -Oberon       }/
 ;     ,--/'\ \    / /                                     }/      /,/
;  ,-./      \ \  { {  (                                  /,;    ,/ ,/
; /   '       } } ', '-'-.___                            / ',  ,/  ',/
 \|         ,','    '~.___,---}                         / ,',,/  ,',;
  '        { {                                     __  /  ,'/   ,',;
        /   \ \                                 _,', '{  ',{   ',';'
       {     } }       /~\         .-:::-.     (--,   ;\ ',}  ',';
       \\._./ /      /' , \      ,:::::::::,     '~;   \},/  ',';     ,-=-
        '-..-'      /. '  .\_   ;:::::::::::;  __,{     '/  ',';     {
                   / , ~ . ^ '~'\:::::::::::<<~>-,,',    '-,  '',_    }
                /~~ . '  . ~  , .'~~\:::::::;    _-~  ;__,        ',-'
       /'\    /~,  . ~ , '  '  ,  .' \::::;'   <<<~'''   ''-,,__   ;
      /' .'\ /' .  ^  ,  ~  ,  . ' . ~\~                       \\, ',__
     / ' , ,'\.  ' ~  ,  ^ ,  '  ~ . . ''~~~',                   '-'--, \
    / , ~ . ~ \ , ' .  ^  '  , . ^   .   , ' .'-,___,---,__            ''
  /' ' . ~ . ' '\ '  ~  ,  .  ,  '  ,  . ~  ^  ,  .  ~  , .'~---,___
/' . '  ,  . ~ , \  '  ~  ,  .  ^  ,  ~  .  '  ,  ~  .  ^  ,  ~  .  '-,`
	clearScreen()
	//red()
	fmt.Println("\n" + banner + "\n")
	//reset()
}

func main() {
	currentWorkingDirectory, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
	}
	clearScreen()
	banner()

	fmt.Println("[~] COMPILING EXE FILES ")
	
	//builds master.exe
	testvar := "go build -o " + currentWorkingDirectory + `\master.exe ` + currentWorkingDirectory + `\master.go`
	cmd := exec.Command("powershell.exe", "/c", testvar)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	//builds bot.exe
	testvar = "go build -ldflags " + ldflags + ` -o ` + currentWorkingDirectory + `\bot.exe ` + currentWorkingDirectory + `\bot.go`
	cmd = exec.Command("powershell.exe", "/c", testvar)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Run()

	endbanner()

}