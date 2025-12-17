package main

import (
	"chowski3/common/apiclients/plexclient"
	"fmt"
)

func main() {
	fmt.Println("Main!")

	// Test out plex API!
	plex, err := plexclient.Connect()
	if err != nil {
		fmt.Println("Error connecting to plex:", err)
		return
	}

	playlists, err := plex.ListPlaylists()
	if err != nil {
		fmt.Println("Error listing playlists:", err)
		return
	}
	fmt.Println("Playlists:", playlists)

	fmt.Println("Please select a playlist:")
	for i := range len(*playlists) {
		fmt.Printf("  %d) %s (%d tracks)\n", i, (*playlists)[i].Name, (*playlists)[i].NumTracks)
	}
	var selection int
	_, err = fmt.Scan(&selection)
	if err != nil || selection < 0 || selection >= len(*playlists) {
		fmt.Println("Invalid selection: ", selection, err)
		return
	}
	playlist := (*playlists)[selection]
	fmt.Println("Selected", playlist)

	tracks, err := plex.ListPlaylistContents(playlist)
	if err != nil {
		fmt.Println("Error listing playlist contents:", err)
	}
	fmt.Printf("Found %d tracks\n", len(*tracks))
}

