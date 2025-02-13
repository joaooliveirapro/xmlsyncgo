package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type FileTable struct {
	RemoteModTime int64 `json:"remote_mod_time"`
}

type SyncManager struct {
	config                  ConfigModel
	jobExternalReferenceKey string
	jobs                    []map[string]interface{} // Jobs from new remote file
	keysFrequency           map[string]int
	localTempFilepath       string
	remoteModTime           time.Time
	stats                   map[string]int
	fileTablePath           string
	fileTable               FileTable
	jobTablePath            string
	jobTable                []Job // Jobs from db
}

func NewSyncManager(configFilepath, fileTablePath, jobTablePath string) (*SyncManager, error) {
	if filepath.Ext(configFilepath) != ".json" {
		return nil, fmt.Errorf("invalid file extension. Must be '.json'")
	}

	// Read config file
	configFile, err := os.Open(configFilepath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	// Unmarshall config file
	var config ConfigModel
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return nil, err
	}
	// Ready file_table.json
	fileTableFile, err := os.Open(fileTablePath)
	if err != nil {
		return nil, err
	}
	defer fileTableFile.Close()
	// Unlarshall file_table.json
	var fileTable FileTable
	err = json.NewDecoder(fileTableFile).Decode(&fileTable)
	if err != nil {
		return nil, err
	}

	// Read job_table.json
	jobTableFile, err := os.Open(jobTablePath)
	if err != nil {
		return nil, err
	}
	defer jobTableFile.Close()
	// Unmarshall job_table.json
	var jobTable []Job
	err = json.NewDecoder(jobTableFile).Decode(&jobTable)
	if err != nil {
		return nil, err
	}

	// Return a new SyncManager with unmarshalled data
	return &SyncManager{
		config:                  config,
		jobExternalReferenceKey: fmt.Sprintf("%s.%s", config.JobNode, config.ExternalReferenceKey),
		keysFrequency:           make(map[string]int),
		stats:                   make(map[string]int),
		fileTablePath:           fileTablePath,
		fileTable:               fileTable,
		jobTablePath:            jobTablePath,
		jobTable:                jobTable,
	}, nil
}

// Recursive function to iterate over all descendants and build path
func iterateDescendants(element *etree.Element, currentPath []string, callback func(*etree.Element, []string)) {
	newPath := make([]string, len(currentPath))
	copy(newPath, currentPath) // Important: Create a copy!
	newPath = append(newPath, element.Tag)

	callback(element, newPath) // Call the callback

	for _, child := range element.ChildElements() {
		iterateDescendants(child, newPath, callback) // Recursive call
	}
}

func (sm *SyncManager) ParseXML() error {
	f, err := os.Open(sm.localTempFilepath)
	if err != nil {
		return fmt.Errorf("❌ Error opening XML file: %v", err)
	}
	defer f.Close()

	tree := etree.NewDocument()
	if _, err := tree.ReadFrom(f); err != nil {
		return fmt.Errorf("❌ XML Parse Error: %v. Stopping", err)
	}
	root := tree.Root()

	sm.jobs = make([]map[string]interface{}, 0)
	sm.keysFrequency = make(map[string]int)

	for _, jobElement := range root.SelectElements(sm.config.JobNode) {
		job := make(map[string]interface{})

		iterateDescendants(jobElement, []string{}, func(element *etree.Element, path []string) {
			var pathStrBuilder strings.Builder
			for _, p := range path {
				pathStrBuilder.WriteString(p)
				pathStrBuilder.WriteString(".")
			}
			pathStr := pathStrBuilder.String()
			pathStr = pathStr[:pathStrBuilder.Len()-1]

			if element.Text() != "" {
				job[pathStr] = element.Text()
				sm.keysFrequency[pathStr]++
			}
		})

		if len(job) > 0 {
			sm.jobs = append(sm.jobs, job)
		}
	}

	fmt.Printf("✔  Successfully parsed XML file at \"%s\".\n", sm.localTempFilepath)
	return nil
}

func (sm *SyncManager) SftpStep() bool {
	// 1. Connect to SFTP
	auth := ssh.Password(sm.config.SftpCredentials.Password) // Or use other auth methods
	config := &ssh.ClientConfig{
		User:            sm.config.SftpCredentials.Username,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, *never* do this!
	}
	addr := fmt.Sprintf("%s:%s", sm.config.SftpCredentials.Hostname, sm.config.SftpCredentials.Port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("❌ Failed to connect to SFTP: %v\n", err)
		return false
	}
	defer conn.Close() // Important: Close the connection

	client, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Printf("❌ Failed to create SFTP client: %v\n", err)
		return false
	}
	defer client.Close() // Important: Close the SFTP client

	fmt.Println("✔  Successfully connected to SFTP.")

	// 2. Check remote modification time
	remoteModTime, err := client.Stat(sm.config.RemoteFilepath)
	if err != nil {
		fmt.Printf("❌ Error getting remote file mod time: %v\n", err)
		return false
	}

	sm.remoteModTime = remoteModTime.ModTime()
	localModTime := time.Unix(sm.fileTable.RemoteModTime, 0) // Convert from int64 to time.Time

	if localModTime.IsZero() || sm.remoteModTime.After(localModTime) { // Use After for time comparison
		fmt.Println("🤩 Remote file is newer than local data.")

		// 3. Download the file
		remoteFile, err := client.Open(sm.config.RemoteFilepath)
		if err != nil {
			fmt.Printf("❌ Error opening remote file: %v\n", err)
			return false
		}
		defer remoteFile.Close()

		localFile, err := os.Create(sm.config.LocalFileTempFolder) // Create local file
		if err != nil {
			fmt.Printf("❌ Error creating local file '%s': %v\n", sm.config.RemoteFilename, err)
			return false
		}
		defer localFile.Close()

		_, err = io.Copy(localFile, remoteFile) // Download
		if err != nil {
			fmt.Printf("❌ Error downloading file: %v\n", err)
			return false
		}

		sm.localTempFilepath = sm.config.RemoteFilename      // Set the local temp file path
		sm.fileTable.RemoteModTime = sm.remoteModTime.Unix() // Store as Unix timestamp

		return true
	} else {
		fmt.Println("😁 Local data is up to date with remote file. Stopping.")
		return false
	}
}

func (sm *SyncManager) FileIntegrityCheck() bool {
	if len(sm.jobs) != sm.keysFrequency[sm.jobExternalReferenceKey] {
		fmt.Printf("❌ There are jobs without a \"%s\" tag. Aborting.\n", sm.jobExternalReferenceKey)
		return false
	}
	fmt.Println("✔  File integrity checked.")
	return true
}

func (sm *SyncManager) CleanUpTempFiles() {
	err := os.Remove(sm.localTempFilepath)
	if err != nil {
		fmt.Printf("❌ Error removing temp file: %v\n", err)
	} else {
		fmt.Printf("✔  Successfully removed temp file at \"%s\"\n", sm.localTempFilepath)
	}
}

func (sm *SyncManager) CheckRemovedJobs() {
	remoteJobIDs := make(map[string]bool)

	// Efficiently populate the remoteJobIDs map
	for _, remoteJob := range sm.jobs {
		if id, ok := remoteJob[sm.jobExternalReferenceKey].(string); ok {
			remoteJobIDs[id] = true
		}
	}

	for i := range sm.jobTable { // Iterate by index for direct modification
		localJob := &sm.jobTable[i] // Get a pointer to the job for direct modification

		if id, ok := localJob.Content[sm.jobExternalReferenceKey].(string); ok {
			if _, found := remoteJobIDs[id]; !found { // Check for *absence* in remoteJobIDs
				if !localJob.Deleted {
					localJob.Deleted = true
					localJob.Edits = append(localJob.Edits, Edit{
						Type:            "REMOVED",
						Ts:              time.Now().Unix(),
						RemoteFileModTs: sm.remoteModTime.Format("02/01/2006, 15:04:05"),
					})
					sm.stats["jobs_removed"]++
				}
			}
		}
	}

	fmt.Println("✔  Successfully checked for jobs removed.")
	if sm.stats["jobs_removed"] > 0 {
		fmt.Printf("sm.stats[\"jobs_removed\"]: %v\n", sm.stats["jobs_removed"])
	}
}

func (sm *SyncManager) ComputeHash(jsonBlob map[string]interface{}) string {
	jsonString, _ := json.Marshal(jsonBlob)
	hasher := sha256.New()
	hasher.Write(jsonString)
	return hex.EncodeToString(hasher.Sum(nil))
}

func (sm *SyncManager) lookupDB(remoteJob map[string]interface{}) (Job, int, bool) {
	for i, localJob := range sm.jobTable {
		if localJob.Content[sm.jobExternalReferenceKey] == remoteJob[sm.jobExternalReferenceKey] {
			return localJob, i, true
		}
	}
	return Job{}, -1, false
}

func (sm *SyncManager) PrintStats() {
	fmt.Printf("len(sm.jobs): %v\n", len(sm.jobs))
	// fmt.Printf("sm.keysFrequency: %v\n", sm.keysFrequency)
}

func (sm *SyncManager) checkJobEdits(remoteJob map[string]interface{}, localJob Job) []Edit {
	edits := make([]Edit, 0)

	addedKeys := make(map[string]bool)
	for k := range remoteJob {
		addedKeys[k] = true
	}
	for k := range localJob.Content {
		if _, ok := remoteJob[k]; ok {
			delete(addedKeys, k) // Remove common keys
		}
	}

	if len(addedKeys) > 0 {
		for k := range addedKeys {
			if val, ok := remoteJob[k].(string); ok { // Type assertion
				edits = append(edits, Edit{
					Ts:              time.Now().Unix(),
					RemoteFileModTs: sm.remoteModTime.Format("02/01/2006, 15:04:05"),
					Type:            "ADDED_KEY",
					Key:             k,
					Value:           val,
				})
			}
		}
	}

	removedKeys := make(map[string]bool)
	for k := range localJob.Content {
		removedKeys[k] = true
	}
	for k := range remoteJob {
		if _, ok := localJob.Content[k]; ok {
			delete(removedKeys, k) // Remove common keys
		}
	}

	if len(removedKeys) > 0 {
		for k := range removedKeys {
			edits = append(edits, Edit{
				Ts:              time.Now().Unix(),
				RemoteFileModTs: sm.remoteModTime.Format("02/01/2006, 15:04:05"),
				Type:            "REMOVED_KEY",
				Key:             k,
			})
		}
	}

	editedKeys := 0
	for k, v := range remoteJob {
		if localV, ok := localJob.Content[k]; ok { // Check if key exists in local job
			if v != localV {
				editedKeys++
				if newV, okNew := v.(string); okNew {
					if oldV, okOld := localV.(string); okOld {
						edits = append(edits, Edit{
							Ts:              time.Now().Unix(),
							RemoteFileModTs: sm.remoteModTime.Format("02/01/2006, 15:04:05"),
							Type:            "EDITED_KEY",
							Key:             k,
							NewValue:        newV,
							OldValue:        oldV,
						})
					}
				}
			}
		}
	}

	sm.stats["added_keys"] += len(addedKeys)
	sm.stats["removed_keys"] += len(removedKeys)
	sm.stats["edited_keys"] += editedKeys
	sm.stats["jobs_processed"]++

	return edits
}

func (sm *SyncManager) SaveData() error {
	jobTableJSON, err := json.MarshalIndent(sm.jobTable, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(sm.jobTablePath, jobTableJSON, 0644); err != nil {
		return err
	}

	fileTableJSON, err := json.MarshalIndent(sm.fileTable, "", "    ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(sm.fileTablePath, fileTableJSON, 0644); err != nil {
		return err
	}

	fmt.Println("✔  Successfully saved data.")
	return nil
}

func (sm *SyncManager) ProcessJobs() {
	for i, remoteJob := range sm.jobs {
		// Skip if job doesn't contain JobExternalReferenceKey
		if _, ok := remoteJob[sm.jobExternalReferenceKey]; !ok {
			fmt.Printf("\t🤔 Skipping job at position %d for missing \"%s\".\n", i, sm.jobExternalReferenceKey)
			continue
		}

		localJob, index, found := sm.lookupDB(remoteJob)

		if found {
			remoteJobHash := sm.ComputeHash(remoteJob)
			if remoteJobHash != localJob.Hash {
				fmt.Printf("\t[EDITED_JOB] \"%s\"\n", remoteJob[sm.jobExternalReferenceKey])

				jobEdits := sm.checkJobEdits(remoteJob, localJob)

				localJob.Content = remoteJob
				localJob.Hash = remoteJobHash
				localJob.Edits = append(localJob.Edits, jobEdits...) // Append edits
				sm.jobTable[index] = localJob
				sm.stats["jobs_edited"]++

			}
		} else {
			fmt.Printf("\t[ADDED_JOB] \"%s\"\n", remoteJob[sm.jobExternalReferenceKey])
			newJob := Job{
				Hash:    sm.ComputeHash(remoteJob),
				Content: remoteJob,
				Edits: []Edit{{
					Ts:              time.Now().Unix(),
					RemoteFileModTs: sm.remoteModTime.Format("02/01/2006, 15:04:05"),
					Type:            "ADDED_JOB",
				}},
			}
			sm.jobTable = append(sm.jobTable, newJob)
			sm.stats["jobs_added"]++
		}
	}
	fmt.Printf("✔  Successfully processed %d jobs.\n", len(sm.jobs))
}
