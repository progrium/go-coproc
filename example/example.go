package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"bitbucket.org/kardianos/osext"
	"github.com/ActiveState/tail"
	"github.com/enr/go-coproc"
)

func main() {
	var name string
	var output func(string)
	if len(os.Args) == 3 {
		name = os.Args[2]
		output = func(str string) {
			os.Stdout.Write(append([]byte(str), '\n'))
		}

	} else {
		name = "master"
		output = func(str string) {
			log.Println(str)
		}
		count, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		exe, err := osext.Executable()
		if err != nil {
			panic(err)
		}
		group := new(coproc.Group)
		for i := 1; i <= count; i++ {
			name := "coproc" + strconv.Itoa(i)
			p := &coproc.Process{
				Name:    name,
				Command: exe,
				Args:    []string{"0", name},
				Pidfile: coproc.Pidfile(name + ".pid"),
				Logfile: name + ".log",
				Respawn: 3,
			}
			t, _ := tail.TailFile(name+".log", tail.Config{Follow: true})
			p.Start()
			go func() {
				for line := range t.Lines {
					output(line.Text)
				}
			}()
			group.Add(p)
		}
	}

	for {
		output(name)
		time.Sleep(3 * time.Second)
	}
}
