package watcher

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func WatchModules(path string, registeredModules []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					if strings.HasSuffix(event.Name, ".module.go") {
						checkModule(event.Name, registeredModules)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watcher error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	select {}
}

func checkModule(file string, registered []string) {
	base := filepath.Base(file)
	modName := strings.TrimSuffix(base, ".module.go")

	for _, r := range registered {
		if r == modName {
			return
		}
	}

	fmt.Printf("❌ Module '%s' exists but is NOT registered in AppModule!\n", modName)
}
