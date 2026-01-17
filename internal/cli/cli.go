package cli

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Execute(args []string) {
	if len(args) == 0 {
		help()
		return
	}

	switch args[0] {
	case "add":
		addConnect()
	case "connect":
		connect()
	default:
		help()
	}
}

func help() {
	fmt.Println("Usage sgirl <add --ip=<..> --username=<...> --password=<...> --name=<...> |connect <*nameOfConnect>")
}

func addConnect() {
	add := flag.NewFlagSet("add", flag.ExitOnError)
	ip := add.String("ip", "", "IP-address of target device")
	user := add.String("user", "", "username to login in target device")
	passwd := add.String("password", "", "password of $user to login")
	name := add.String("name", "", "number or name connect for you")
	//Проверка на наличие всех обязательных аргументов
	if len(os.Args) < 3 {
		fmt.Println("You should use: sgirl add --ip=<...> --username=<...> --password=<...> | --name=<...>")
		os.Exit(1)
	}
	if len(*name) == 0 {
		checkConfigDir, _ := os.ReadDir("../config")
		*name = fmt.Sprintf("num_%s", len(checkConfigDir)+1)
	}

	add.Parse(os.Args[2:])
	config := Config{
		NAME:     *name,
		IP:       *ip,
		USERNAME: *user,
		PASSWORD: *passwd,
		Created:  time.Now().Format(time.RFC3339),
	}
	configYaml, _ := yaml.Marshal(config)

	//Запись в конфиг
	os.WriteFile(fmt.Sprintf("../config/%s_config.yaml", *name), configYaml, 0644)

	//Вывод сведений в термнинал
	fmt.Printf("%-20s %-15s %-10s %s\n", "IP", "USER", "PprintUsageASSWORD", "NAME")
	fmt.Printf("IP: %s, User: %s, Password: %s, ConnectName: %s\n", *ip, *user, strings.Repeat("*", len(*passwd)), *name)
	fmt.Printf("File saved on ../config/%s_config.yaml", *name)
}

func connect() {
	connect := flag.NewFlagSet("connect", flag.ExitOnError)
	if len(os.Args) < 2 {
		fmt.Println("Please  usage: sgirl add <name of connect>")
		os.Exit(1)
	}
	connect.Parse(os.Args[2:])
	nameConn := connect.Arg(2)
	open, _ := os.ReadFile(fmt.Sprintf("../config/%s_config.yaml", nameConn))
	var cfg Config
	yaml.Unmarshal(open, &cfg)
	cmd := exec.Command("ssh", cfg.USERNAME, cfg.IP)
	stdin, _ := cmd.StdinPipe()
	cmd.Start()
	go func() {
		time.Sleep(1 * time.Second)
		stdin.Write([]byte(cfg.PASSWORD + "\r"))
		stdin.Close()
	}()
	cmd.Wait()
	fmt.Println("Connected to %s", nameConn)
}

type Config struct {
	NAME     string `yaml:"name"`
	IP       string `yaml:"ip"`
	USERNAME string `yaml:"user"`
	PASSWORD string `yaml:"password"`
	Created  string `yaml:"created"`
}
