package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var m sync.Mutex
var Finished bool

func waitforprocess(p *os.Process) {
	p.Wait()
	m.Lock()
	Finished = true
	m.Unlock()
}

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Must supply at least the program to run!")
		return
	}
	str := ""
	for _, s := range os.Args[1:] {
		str += s + " "
	}
	log.Println("gorun " + str)
	execfile := os.Args[1]

again:

	log.Println("gorun: Launching new process")
	p, err := os.StartProcess(execfile, os.Args[2:], &os.ProcAttr{Env: nil, Dir: "", Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	if err != nil {
		panic(err)
	}
	log.Println("gorun: launched process")
	Finished = false
	initialStat, err := os.Stat(execfile)
	if err != nil {
		panic(err)
	}
	Relaunch := false
	go waitforprocess(p)
	for {
		time.Sleep(time.Second)
		m.Lock()
		if Finished {
			m.Unlock()
			Relaunch = false
			break
		}
		m.Unlock()
		stat, err := os.Stat(execfile)
		if err != nil {
			// assume this is file-not-found
			log.Println("Watched file not found... waiting ")
			time.Sleep(1 * time.Second)
			stat, err = os.Stat(execfile)
			if err != nil {
				log.Println("Watched file not found... waiting ")
				time.Sleep(2 * time.Second)
				stat, err = os.Stat(execfile)
				if err != nil {
					log.Println("Watched file not found... waiting ")
					time.Sleep(5 * time.Second)
					stat, err = os.Stat(execfile)
					if err != nil {
						log.Println("Watched file not found... killing!")
						p.Kill()
						p.Wait()
						p.Release()
						return
					}
				}
			}

		}
		//fmt.Println(stat.ModTime())
		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {

			Relaunch = true
			break
		}
	}
	if Relaunch {
		log.Println("***Restarting process***")
		log.Println("Killing existing process")
		p.Kill()
		p.Wait()
		p.Release()
		goto again
		// TODO
	}
	p.Release()
}
