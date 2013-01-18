package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"os/exec"
)

func main() {
	block := make(chan int)

	log_file, err := os.OpenFile("/tmp/resume_watcher.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer log_file.Close()

	log.SetOutput(log_file)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			err = watcher.Watch("/home/ryan/Dropbox/resume/markdown/resume.md")
			//err = watcher.Watch("/tmp/test.txt")
			if err != nil {
				log.Fatal(err)
			}

			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
				log.Println("Running build script...")
				err := exec.Command("/home/ryan/Dropbox/resume/markdown/build.sh").Run()
				if err != nil {
					log.Println("Error running the build script!")
				}

				log.Println("Build script ran successfully")
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	<-block

	log.Println("I should not be here")
}
