package parser

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/joaooliveirapro/xmlsyncgo/initializers"
	"github.com/joaooliveirapro/xmlsyncgo/models"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

type ParserManager struct {
	file          *models.File
	jobs          []map[string]interface{}
	keysFrequency map[string]int
	stats         map[string]int
}

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
		return fmt.Errorf("❌ Failed to connect to SFTP: %v", err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("❌ Failed to create SFTP client: %v", err)
	}
	defer client.Close()

	// 2. Check remote modification time
	remoteModTime, err := client.Stat(pm.file.RemoteFilepath)
	if err != nil {
		return fmt.Errorf("❌ Error getting remote file mod time: %v", err)
	}
	if remoteModTime.ModTime().After(pm.file.RemoteModTime) { // Use After for time comparison
		// 3. Read the updated remote file
		fullFilePath := fmt.Sprintf("%s%s", pm.file.RemoteFilepath, pm.file.RemoteFilename)
		remoteFile, err := client.Open(fullFilePath)
		if err != nil {
			return fmt.Errorf("❌ Error opening remote file: %v", err)
		}
		defer remoteFile.Close()

		// 4. Create a local temporary file
		fullTempFilePath := fmt.Sprintf("./temp/%s", pm.file.RemoteFilename)
		localFile, err := os.Create(fullTempFilePath) // Create local file
		if err != nil {
			return fmt.Errorf("❌ Error creating local file '%s': %v", pm.file.RemoteFilename, err)
		}
		defer localFile.Close()

		// 5. Copy the file
		_, err = io.Copy(localFile, remoteFile)
		if err != nil {
			return fmt.Errorf("❌ Error downloading file: %v", err)
		}

		// Set the remote_mod_time to the newest mod time
		pm.file.RemoteModTime = remoteModTime.ModTime()
		result := initializers.DB.Save(&pm.file)
		if result.Error != nil {
			return fmt.Errorf("❌ Error updating file.remote_mod_time: %v", err)
		}

		// Log some audit trail
		models.NewAuditLog("✔  Successfully connected to SFTP. Remote file is newer than local data. Downloaded successfully.", pm.file.AuditIteration)
	} else {
		models.NewAuditLog("✔  Successfully connected to SFTP. Local data on track with remote data.", pm.file.AuditIteration)
	}
	return nil
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

// Parse XML file (downloaded to /temp/) using XML parser ETree
// Creates a list of jobs ([]map[string]interface{}) and a
// KeysFrequency map[string]int. Both properties of PM
func (pm *ParserManager) parseXML() error {
	// Ensure defaults data structures are created
	pm.jobs = make([]map[string]interface{}, 0)
	pm.keysFrequency = make(map[string]int)

	// 1. Read the downloaded file from the temp folder
	fullTempFilePath := fmt.Sprintf("./temp/%s", pm.file.RemoteFilename)
	f, err := os.Open(fullTempFilePath)
	if err != nil {
		return fmt.Errorf("❌ Error opening XML file: %v", err)
	}
	defer f.Close()

	// 2. Build the eTree
	tree := etree.NewDocument()
	if _, err := tree.ReadFrom(f); err != nil {
		return fmt.Errorf("❌ XML Parse Error: %v. Stopping", err)
	}
	root := tree.Root()

	// 3. For each element in the parsed file
	for _, jobElement := range root.SelectElements(pm.file.JobNodeKey) {
		// 4. Create a map[string]interface
		job := make(map[string]interface{})

		// 5. Iterate recursivelly over its keys
		iterateDescendants(jobElement, []string{}, func(element *etree.Element, path []string) {
			// 5.1 Build the flatten key path
			var pathStrBuilder strings.Builder
			for _, p := range path {
				pathStrBuilder.WriteString(p)
				pathStrBuilder.WriteString(".")
			}
			pathStr := pathStrBuilder.String()
			pathStr = pathStr[:pathStrBuilder.Len()-1]

			if element.Text() != "" {
				// 5.2 Once flattened, add the key and value to the job map object
				job[pathStr] = element.Text()
				// 5.3 Log the key frequency
				pm.keysFrequency[pathStr]++
			}
		})
		// 6. Append job processed to the list
		if len(job) > 0 {
			pm.jobs = append(pm.jobs, job)
		}
	}
	// Log some audit trail
	models.NewAuditLog("✔  Successfully parsed new XML file.", pm.file.AuditIteration)
	models.NewAuditLog(fmt.Sprintf("%+v", pm.keysFrequency), pm.file.AuditIteration)
	return nil
}

// Check if all the jobs parsed (before saving to DB) contain the main
// ExternalReferenceKey. If at least one fails, return error.
func (pm ParserManager) fileIntegrityCheck() error {
	fullExternlaReferenceKey := fmt.Sprintf("%s.%s", pm.file.JobNodeKey, pm.file.ExternalReferenceKey)
	if len(pm.jobs) != pm.keysFrequency[fullExternlaReferenceKey] {
		return fmt.Errorf("❌ There are jobs without a \"%s\" tag", fullExternlaReferenceKey)
	}
	return nil
}

// Computes a SHA256 hash of a map[string]any data object.
func (pm ParserManager) computeHash(data map[string]interface{}) string {
	jsonString, _ := json.Marshal(data)
	hasher := sha256.New()
	hasher.Write(jsonString)
	return hex.EncodeToString(hasher.Sum(nil))
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
		return err
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
					Ts:              time.Now().Unix(),
					RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
					Type:            "ADDED_KEY",
					Key:             newKey,
					Value:           value,
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
				Ts:              time.Now().Unix(),
				RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
				Type:            "REMOVED_KEY",
				Key:             removedKey,
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
							Ts:              time.Now().Unix(),
							RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
							Type:            "EDITED_KEY",
							Key:             remoteKey,
							NewValue:        newValue,
							OldValue:        oldValue,
						})
					}
				}
			}
		}
	}

	// Log some stats
	pm.stats["added_keys"] += len(addedKeys)
	pm.stats["removed_keys"] += len(removedKeys)
	pm.stats["edited_keys"] += editedKeysCounter
	pm.stats["jobs_processed"]++

	// Batch save all edits to db
	result := initializers.DB.Create(&edits)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pm ParserManager) processJobs() error {
	// Initialise stats map
	pm.stats = map[string]int{}

	// Hold remote jobs parsed to models.Job{} for batch insert into DB
	processedJobs := []models.Job{}
	// Holds all the remote jobs external reference ids to check for removed jobs from db
	remoteJobsIDs := map[string]bool{}

	// Loop through all remote jobs
	for i, remoteJob := range pm.jobs {
		// Skip if job doesn't contain JobExternalReferenceKey
		// This has already been checked but it's here again for sanity sake
		fullExternlaReferenceKey := fmt.Sprintf("%s.%s", pm.file.JobNodeKey, pm.file.ExternalReferenceKey)
		if _, ok := remoteJob[fullExternlaReferenceKey]; !ok {
			return fmt.Errorf("❌ Skipping job at position %d for missing \"%s\"", i, fullExternlaReferenceKey)
		}
		// Remote job external reference value
		remoteJobExternalReference := remoteJob[fullExternlaReferenceKey].(string)
		remoteJobsIDs[remoteJobExternalReference] = true // This job will either be added to DB or already exists

		// Check if DB contains this job
		var localJob models.Job
		result := initializers.DB.Where("external_reference = ?", remoteJobExternalReference).First(&localJob)

		var found bool = true
		// Not found is considered a type of error
		if result.Error != nil {
			// Not found
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				found = false
			} else { // Some other error
				// do something
			}
		}
		// Compute remote job hash
		remoteJobHash := pm.computeHash(remoteJob)
		// Convert remoteJob to JSON string
		remoteJobJson, err := json.Marshal(remoteJob)
		if err != nil {
			return err
		}
		if found {
			// If job has been edited (hash is different)
			if remoteJobHash != localJob.Hash {
				// Check what edits have been made and log them in DB
				err := pm.checkJobEdits(remoteJob, localJob)
				if err != nil {
					return err
				}
				// Update job record in DB to the new job model
				localJob.Content = string(remoteJobJson)
				// Update job record hash to new hash
				localJob.Hash = remoteJobHash
				// Save updated job record
				result := initializers.DB.Save(&localJob)
				if result.Error != nil {
					return result.Error
				}
				// Increment stats
				pm.stats["jobs_edited"]++
			} else {
				// Job hasn't been edited
			}
		} else {
			// Create job object (and edits)
			newJob := models.Job{
				ExternalReference: remoteJobExternalReference,
				Hash:              remoteJobHash,
				Content:           string(remoteJobJson),
				Edits: []models.Edit{{
					Ts:              time.Now().Unix(),
					RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
					Type:            "ADDED_JOB",
				}},
			}
			// Add new job to list for batch insert into db
			processedJobs = append(processedJobs, newJob)
			// Increment stats
			pm.stats["jobs_added"]++
		}
	}
	// Batch insert all processed jobs into DB
	result := initializers.DB.Create(&processedJobs)
	if result.Error != nil {
		return result.Error
	}
	// Check for removed jobs
	// Loop through all db jobs that are NOT marked deleted already
	var localJobs []models.Job
	result = initializers.DB.Select("external_reference", "id").Where("deleted_at IS NULL").Find(&localJobs)
	if result.Error != nil {
		return result.Error
	}
	for i := range localJobs {
		localJob := &localJobs[i] // Get a pointer to the actual job element
		// If the localjob id is not in the remoteJobs ids list, then
		// job has been deleted from source data.
		if _, ok := remoteJobsIDs[localJob.ExternalReference]; !ok {
			//	Delete job from db (this won't delete data, only sets the deleted_at)
			result := initializers.DB.Where("id = ?", localJob.ID).Delete(&models.Job{})
			if result.Error != nil {
				return result.Error
			}
			// Add edit to the job and save
			localJob.Edits = append(localJob.Edits, models.Edit{
				Type:            "REMOVED",
				Ts:              time.Now().Unix(),
				RemoteFileModTs: pm.file.RemoteModTime.Format("02/01/2006, 15:04:05"),
			})
			initializers.DB.Save(&localJob)
			// Increment stats
			pm.stats["jobs_removed"]++
		}
	}
	// Convert stat map to json string
	statJson, err := json.Marshal(pm.stats)
	if err != nil {
		return err
	}
	// Save stats to db
	stat := models.Stat{JsonStr: string(statJson)}
	result = initializers.DB.Save(&stat)
	if result.Error != nil {
		return result.Error
	}
	// Log some audit trail
	models.NewAuditLog(fmt.Sprintf("✔  Successfully processed %d jobs.", len(pm.jobs)), pm.file.AuditIteration)
	return nil
}

func (pm ParserManager) Run(clients *[]models.Client) {
	for _, client := range *clients {
		var err error
		for _, file := range client.Files {
			pm := ParserManager{file: &file}
			// SFTP Step
			if err = pm.sftpStep(); err != nil {
				models.NewAuditLog(err.Error(), pm.file.AuditIteration)
				continue
			}
			// Remote file parse step
			if err = pm.parseXML(); err != nil {
				models.NewAuditLog(err.Error(), pm.file.AuditIteration)
				continue
			}
			// After parsing the file, and before saving
			// anything to the db, check file data integrity
			if err = pm.fileIntegrityCheck(); err != nil {
				models.NewAuditLog(err.Error(), pm.file.AuditIteration)
				continue
			}
			// Process new jobs
			if err = pm.processJobs(); err != nil {
				models.NewAuditLog(err.Error(), pm.file.AuditIteration)
				continue
			}
			// Increment AuditIteration and save
			file.AuditIteration++
			result := initializers.DB.Save(&pm.file)
			if result.Error != nil {
				models.NewAuditLog(result.Error.Error(), pm.file.AuditIteration)
				continue
			}
		}
	}
}
