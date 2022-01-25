package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nono/cozy-desktop-experiments/client"
	"github.com/nono/cozy-desktop-experiments/localfs"
	"github.com/nono/cozy-desktop-experiments/platform"
	"github.com/nono/cozy-desktop-experiments/state"
)

// Start will listen for http requests, and serve them.
func Start(port string) error {
	http.Handle("/", http.FileServer(http.Dir("./web/assets")))
	http.HandleFunc("/run", Run)
	fmt.Println("Starting server at port " + port)
	return http.ListenAndServe(":"+port, nil)
}

func Run(w http.ResponseWriter, r *http.Request) {
	localDir := "."
	localFS := localfs.NewDirFS(localDir)
	// localFS := localfs.NewMemFS()

	// remoteClient := client.New("http://cozy.localhost:8080/")
	remoteClient := client.NewFake("http://cozy.localhost:8080/")
	remoteClient.AddInitialDocs()

	fmt.Println("Start")
	platform := platform.New(localFS, remoteClient)
	if err := state.Sync(platform); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
