# XMLSyncGo

A Go tool for synchronizing job data between a remote XML file (via SFTP) and a local database (JSON files). It efficiently handles additions, updates, and removals of jobs, tracking changes and providing statistics. Ideal for automating data synchronization tasks.

## Features

* **SFTP Retrieval:** Connects to a remote SFTP server and downloads the XML file.  Only downloads if the remote file has been modified since the last sync.
* **XML Parsing:** Parses the XML file, extracting job data and flattening the structure for easier processing. Tracks the frequency of different XML tags.
* **Data Integrity Check:** Ensures that all jobs have a required identifier (external reference key) for proper matching.
* **Database Synchronization:**
    * **New Jobs:** Adds new jobs found in the XML file but not in the local database.
    * **Updated Jobs:** Updates existing jobs in the local database with data from the XML file. Tracks specific changes (added/removed/edited keys).
    * **Removed Jobs:** Marks jobs as deleted in the local database if they are present locally but not in the remote XML file.
* **Change Tracking:** Logs all changes to jobs, including timestamps and the type of change (added, removed, edited).
* **Statistics:** Tracks and displays statistics on the number of jobs added, removed, and updated.
* **Local Data Storage:** Stores job data and file metadata (remote file modification time) in local db file.
* **Configurable:** Uses a JSON configuration file for SFTP credentials, file paths, and other settings.
* **Scheduled Runs (Optional):** Can be run periodically to automate synchronization.
* **Cross-Platform:** Compiles to a single binary for easy cross-platform deployment.


## Getting Started

### Prerequisites
* Go 1.x
* Access to an SFTP server with the XML file to track.

### Installation
```sh
go get https://github.com/joaooliveirapro/xmlsyncgo
```

Create the `config.json` file.
```json
// config.json example
{
  "root": "data",
  "job_node": "job",
  "external_reference_key": "job_id",
  "remote_filepath": "/path/to/remote/jobs.xml",
  "filename": "jobs.xml",
  "sftp_credentials": {
    "hostname": "your_sftp_host",
    "port": "22",
    "username": "your_sftp_user",
    "password": "your_sftp_password"
  }
}
```

### License
The MIT License (MIT)

