package main

import (
	"fmt"

	"github.com/joaooliveirapro/xmlsyncgo/src/models"
)

func main() {
	sm, err := models.NewSyncManager(
		"./data/config.json",
		"./data/file_table.json",
		"./data/job_table.json",
	)
	if err != nil {
		fmt.Println(err)
	}

	ok := sm.SftpStep()
	if !ok {
		fmt.Println("Opps...")
		return
	}
	// Read remote file (using locally stored to avoid mul connections to sftp while dev)
	err = sm.ParseXML()
	if err != nil {
		fmt.Println(err)
	}

	// Check integrity of new remote file
	ok = sm.FileIntegrityCheck()
	if !ok {
		fmt.Println("Oops...")
	}

	// Process new jobs
	sm.ProcessJobs()

	// Check for removed jobs
	sm.CheckRemovedJobs()

	sm.SaveData()
	// Print stats
	// sm.PrintStats()

	// Clean up temp files
	// sm.CleanUpTempFiles()
}
