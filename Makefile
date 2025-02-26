watch:
	air --build.cmd "go build -o xmlsyncgo.exe main.go" --build.bin "xmlsyncgo.exe" --build.exclude_dir "src\adminfe,src\migrate,src\web,temp,tmp"
