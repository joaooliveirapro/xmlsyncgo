package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/beevik/etree"
	"github.com/joaooliveirapro/xmlsyncgo/src/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/src/models"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

type ParserManager struct {
	tx            *gorm.DB
	file          *models.File
	jobs          []map[string]interface{}
	keysFrequency map[string]int
	stats         *models.Stat
}

const (
	ADDED_KEYS   = "ADDED_KEYS"
	EDITED_KEYS  = "EDITED_KEYS"
	REMOVED_KEYS = "REMOVED_KEYS"
	JOBS_ADDED   = "JOBS_ADDED"
	JOBS_REMOVED = "JOBS_REMOVED"
	JOBS_EDITED  = "JOBS_EDITED"
)

// Log into the SFTP, check if the data is newer and download it
// to the temp folder for further processing.
func (pm ParserManager) sftpStep() error {
	// 1. Connect to SFTP
	auth := ssh.Password(pm.file.Password)
	config := &ssh.ClientConfig{
		User:            pm.file.Username,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%s", pm.file.Hostname, pm.file.Port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("[processJobs] ftpStep] Failed to connect to SFTP: %w", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("[sftpStep] Failed to create SFTP client: %w", err)
	}
	defer client.Close()

	// 2. Check remote modification time
	remoteModTime, err := client.Stat(pm.file.RemoteFilepath)
	if err != nil {
		return fmt.Errorf("[sftpStep] Error getting remote file mod time: %w", err)
	}
	if remoteModTime.ModTime().After(pm.file.RemoteModTime) { // Use After for time comparison
		// 3. Read the updated remote file
		fullFilePath := fmt.Sprintf("%s%s", pm.file.RemoteFilepath, pm.file.RemoteFilename)
		remoteFile, err := client.Open(fullFilePath)
		if err != nil {
			return fmt.Errorf("[sftpStep] Error opening remote file: %w", err)
		}
		defer remoteFile.Close()
		// 4. Create a local temporary file
		fullTempFilePath := fmt.Sprintf("./temp/%s", pm.file.RemoteFilename)
		localFile, err := os.Create(fullTempFilePath) // Create local file
		if err != nil {
			return fmt.Errorf("[sftpStep] Error creating local file '%s': %w", pm.file.RemoteFilename, err)
		}
		defer localFile.Close()
		// 5. Copy the file
		_, err = io.Copy(localFile, remoteFile)
		if err != nil {
			return fmt.Errorf("[sftpStep] Error downloading file: %w", err)
		}
		// Set the remote_mod_time to the newest mod time
		pm.file.RemoteModTime = remoteModTime.ModTime()
		result := pm.tx.Save(&pm.file)
		if result.Error != nil {
			return fmt.Errorf("[sftpStep] Error updating file.remote_mod_time: %w", err)
		}
	}
	return nil
}

func iterateDescendants(element *etree.Element, currentPath []string, callback func(*etree.Element, string)) {
	newPath := make([]string, len(currentPath), len(currentPath)+10)
	copy(newPath, currentPath)
	newPath = append(newPath, element.Tag)
	pathStr := strings.Join(newPath, ".")
	callback(element, pathStr)
	for _, child := range element.ChildElements() {
		iterateDescendants(child, newPath, callback)
	}
}

// Optimized parseXML function with Mutex
func (pm *ParserManager) parseXML() error {
	// Read temp file (downloaded from source)
	f, err := ReadTempFile(pm.file.RemoteFilename)
	if err != nil {
		return err // Descriptive internally
	}
	defer f.Close()
	// Build XML ET
	tree := etree.NewDocument()
	if _, err := tree.ReadFrom(f); err != nil {
		return fmt.Errorf("[parseXML] XML Parse Error: %w. Stopping", err)
	}
	root := tree.Root()
	var wg sync.WaitGroup // Wg group
	var mu sync.Mutex     // Mutex for keysFrequency
	jobChan := make(chan map[string]interface{})
	// loop all elements in XML file and build a map[string]string of
	// a flat representation of the XML file keys and values
	for _, jobElement := range root.SelectElements(pm.file.JobNodeKey) {
		wg.Add(1)
		go func(jobElement *etree.Element) {
			defer wg.Done()
			job := make(map[string]interface{})
			iterateDescendants(jobElement, []string{}, func(element *etree.Element, pathStr string) {
				if element.Text() != "" {
					job[pathStr] = element.Text()
					mu.Lock() // Lock before writing to keysFrequency
					pm.keysFrequency[pathStr]++
					mu.Unlock() // Unlock after writing
				}
			})
			if len(job) > 0 {
				jobChan <- job
			}
		}(jobElement)
	}
	// Wait for all go routines to finish and close the channel
	go func() {
		wg.Wait()
		close(jobChan)
	}()
	// Each job added to the channel, append to pm.jobs
	for job := range jobChan {
		pm.jobs = append(pm.jobs, job)
	}
	// Log keys frequency map to db
	pm.tx.Create(&models.AuditLog{
		Text:           fmt.Sprintf("%+v", pm.keysFrequency),
		FileID:         pm.file.ID,
		AuditIteration: pm.file.AuditIteration,
	})
	// Data integrity check - it's important all jobs have the ExternalReferenceKey
	fullExternlaReferenceKey := fmt.Sprintf("%s.%s", pm.file.JobNodeKey, pm.file.ExternalReferenceKey)
	if len(pm.jobs) != pm.keysFrequency[fullExternlaReferenceKey] {
		return fmt.Errorf("[fileIntegrityCheck] There are jobs without a \"%s\" tag", fullExternlaReferenceKey)
	}
	return nil
}

// Check job edits compares remote job keys with the job data stored in db
// Only called if the hash of remote job is different from the hash of local data
func (pm ParserManager) checkJobEdits(remoteJob map[string]interface{}, localJob models.Job) error {
	// Hold all edits found
	edits := make([]models.Edit, 0)

	// Unmarshall localJob.Content to map[string]interface{}
	var localJobcontentJson map[string]interface{}
	err := json.Unmarshal([]byte(localJob.Content), &localJobcontentJson)
	if err != nil {
		return fmt.Errorf("[checkJobEdits] Error unmarshalling localJob.Content=%s | error=%s", localJob.Content, err.Error())
	}

	// Check if new keys have been added.
	// These keys are in remoteJob but not in localJob
	addedKeys := make(map[string]bool)
	for key := range remoteJob {
		// If key not in localJob, then it's a new key
		if _, ok := localJobcontentJson[key]; !ok {
			addedKeys[key] = true
		}
	}
	if len(addedKeys) > 0 {
		for newKey := range addedKeys {
			if value, ok := remoteJob[newKey].(string); ok {
				edits = append(edits, models.Edit{
					Type:            "ADDED_KEY",
					Ts:              time.Now().Unix(),
					RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
					Key:             newKey,
					Value:           value,
					JobID:           localJob.ID,
				})
			}
		}
	}
	// Check for keys present in localJob but removed from remoteJob
	// These keys have been removed
	removedKeys := make(map[string]bool)
	for key := range localJobcontentJson {
		// If key not in remoteJob, then it has been removed
		if _, ok := remoteJob[key]; !ok {
			removedKeys[key] = true
		}
	}
	if len(removedKeys) > 0 {
		for removedKey := range removedKeys {
			edits = append(edits, models.Edit{
				Type:            "REMOVED_KEY",
				Ts:              time.Now().Unix(),
				RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
				Key:             removedKey,
				JobID:           localJob.ID,
			})
		}
	}
	// Check for edited keys
	editedKeysCounter := 0
	for remoteKey, remoteValue := range remoteJob {
		// Check if key exists in localJob
		if localValue, ok := localJobcontentJson[remoteKey]; ok {
			// check if values are different
			if remoteValue != localValue {
				// Increment counter (for stats)
				editedKeysCounter++
				// Add new Edit to []edits
				if newValue, okNew := remoteValue.(string); okNew {
					if oldValue, okOld := localValue.(string); okOld {
						edits = append(edits, models.Edit{
							Type:            "EDITED_KEY",
							Ts:              time.Now().Unix(),
							RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
							Key:             remoteKey,
							NewValue:        newValue,
							OldValue:        oldValue,
							JobID:           localJob.ID,
						})
					}
				}
			}
		}
	}

	// Log some stats
	pm.stats.IncrementKey(ADDED_KEYS, len(addedKeys))
	pm.stats.IncrementKey(REMOVED_KEYS, len(removedKeys))
	pm.stats.IncrementKey(EDITED_KEYS, editedKeysCounter)

	// Batch save all edits to db
	result := pm.tx.Create(&edits)
	if result.Error != nil {
		return fmt.Errorf("[checkJobEdits] Error batch saving edits to DB | error=%w", result.Error)
	}
	return nil
}

func (pm ParserManager) processJobs() error {
	// Hold remote jobs parsed to models.Job{} for batch insert into DB
	processedJobs := make([]models.Job, 0, len(pm.jobs)) // Pre-allocate
	// Holds all the remote jobs external reference ids to check for removed jobs from db
	remoteJobsIDs := make(map[string]bool, len(pm.jobs)) // Pre-allocate
	// Processing errors
	processingErrors := make([]string, 0) // Pre-allocate
	// Loop through all remote jobs
	if len(pm.jobs) > 0 {
		// Full External Reference Key
		fullExternlaReferenceKey := fmt.Sprintf("%s.%s", pm.file.JobNodeKey, pm.file.ExternalReferenceKey)
		//
		externalReferences := make([]string, 0, len(pm.jobs))
		// Loop through all remote jobs
		for _, remoteJob := range pm.jobs {
			// Remote job external reference value fileIntegrityCheck has already made sure all jobs contain fullExternlaReferenceKey
			remoteJobExternalReference := remoteJob[fullExternlaReferenceKey].(string)
			externalReferences = append(externalReferences, remoteJobExternalReference)
			remoteJobsIDs[remoteJobExternalReference] = true // This job will either be added to DB or already exists
		}

		// Batch fetch existing jobs from the database
		var existingJobs []models.Job
		result := pm.tx.Where("external_reference IN ?", externalReferences).Find(&existingJobs)
		if result.Error != nil {
			return fmt.Errorf("[processJobs] error fetching existing jobs: %w", result.Error)
		}

		existingJobsMap := make(map[string]models.Job, len(existingJobs)) // For quick lookup
		for i := range existingJobs {
			existingJobsMap[existingJobs[i].ExternalReference] = existingJobs[i]
		}

		for _, remoteJob := range pm.jobs {
			var localJob models.Job
			remoteJobExternalReference := remoteJob[fullExternlaReferenceKey].(string)
			// Remote job hash
			remoteJobHash := ComputeSHA256Hash(remoteJob)
			// Convert job to json
			remoteJobJson, err := json.Marshal(remoteJob)
			if err != nil {
				processingErrors = append(processingErrors, fmt.Errorf("[processJobs] error converting remoteJob to JSON: %w", err).Error())
				continue
			}

			localJob, found := existingJobsMap[remoteJobExternalReference]
			if found {
				// If job has been edited (hash is different)
				if remoteJobHash != localJob.Hash {
					// Check what edits have been made and log them in DB
					err := pm.checkJobEdits(remoteJob, localJob)
					if err != nil {
						processingErrors = append(processingErrors, fmt.Errorf("[processJobs] error during checkJobEdits %w", err).Error())
						continue
					}
					// Update job record in DB to the new job model
					localJob.Content = string(remoteJobJson)
					// Update job record hash to new hash
					localJob.Hash = remoteJobHash
					// Save updated job record
					result := pm.tx.Save(&localJob)
					if result.Error != nil {
						return fmt.Errorf("[processJobs] error saving localJob %w", result.Error)
					}
					// Increment stats
					pm.stats.IncrementKey(JOBS_EDITED, 1)
				}
			} else {
				// Create job object (and edits)
				newJob := models.Job{
					ExternalReference: remoteJobExternalReference,
					Content:           string(remoteJobJson),
					Hash:              remoteJobHash,
					FileID:            pm.file.ID,
					Edits: []models.Edit{{
						Type:            "ADDED_JOB",
						Ts:              time.Now().Unix(),
						RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
						// JobID:           localJob.ID,
					}},
				}
				// Add new job to list for batch insert into db
				processedJobs = append(processedJobs, newJob)
				// Increment stats
				pm.stats.IncrementKey(JOBS_ADDED, 1)
			}

		}
		// Batch insert all processed jobs into DB
		if len(processedJobs) > 0 {
			result := pm.tx.Create(&processedJobs)
			if result.Error != nil {
				return fmt.Errorf("[processJobs] error creating processedJobs %w", result.Error)
			}
		}
		// Return all errors found in one joined string
		if len(processingErrors) > 0 {
			return fmt.Errorf("%s", strings.Join(processingErrors, "\n"))
		}
	}

	// Check for removed jobs
	// Loop through all db jobs that are NOT marked deleted already
	var localJobs []models.Job
	result := pm.tx.Where("deleted", false).Find(&localJobs)
	if result.Error != nil {
		return fmt.Errorf("[processJobs] error finding []localJobs with deleted = false %w", result.Error)
	}
	if len(localJobs) > 0 {
		for i := range localJobs {
			localJob := &localJobs[i] // Get a pointer to the actual job element
			// If the localjob id is not in the remoteJobs ids list, then
			// job has been deleted from source data.
			if _, ok := remoteJobsIDs[localJob.ExternalReference]; !ok {
				// Add edit to the job and save
				localJob.Edits = append(localJob.Edits, models.Edit{
					Type:            "REMOVED_JOB",
					Ts:              time.Now().Unix(),
					RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
					JobID:           localJob.ID,
				})
				//	Delete job from db (this won't delete data; set Deleted=true)
				pm.tx.Model(&localJob).Updates(models.Job{Deleted: true})
				// Increment stats
				pm.stats.IncrementKey(JOBS_REMOVED, 1)
			}
		}
	}

	// Save stats to db (pass the transaction object)
	err := pm.stats.SaveToDB(pm.tx)
	if err != nil {
		return err // Error descriptive internally
	}
	// Log some audit trail
	pm.tx.Create(&models.AuditLog{
		Text:           fmt.Sprintf("✔  Successfully processed %d jobs.", len(pm.jobs)),
		FileID:         pm.file.ID,
		AuditIteration: pm.file.AuditIteration,
	})
	// Increment AuditIteration and save
	pm.file.AuditIteration++
	result = pm.tx.Save(&pm.file)
	if result.Error != nil {
		return fmt.Errorf("[processJobs] error saving pm.file %w", result.Error)
	}
	return nil
}

func Main() {
	statKeys := []string{ADDED_KEYS, EDITED_KEYS, REMOVED_KEYS, JOBS_ADDED, JOBS_REMOVED, JOBS_EDITED}
	var clients []models.Client                     // Get all clients from DB
	initializers.DB.Preload("Files").Find(&clients) // Preload "Files" for each client
	if len(clients) > 0 {
		for _, client := range clients {
			var err error
			if len(client.Files) > 0 {
				for _, file := range client.Files {
					// Surround all in a single db tx return err will rollback()
					err = initializers.DB.Transaction(func(tx *gorm.DB) error {
						// Initialise parse manager for each file
						pm := ParserManager{
							tx:            tx,
							file:          &file,
							jobs:          make([]map[string]interface{}, 0, 1000),
							keysFrequency: make(map[string]int),
							stats:         models.NewStat(statKeys, file.ID),
						}
						// SFTP Step
						if err = pm.sftpStep(); err != nil {
							return err
						}
						// Remote file parse step
						if err = pm.parseXML(); err != nil {
							return err
						}
						// Process new jobs
						if err = pm.processJobs(); err != nil {
							return err
						}
						return nil // Commit tx
					})
					if err != nil {
						initializers.DB.Create(&models.AuditLog{
							Text:           err.Error(),
							FileID:         file.ID,
							AuditIteration: file.AuditIteration,
						})
					}
				}
			}
		}
	}
}
