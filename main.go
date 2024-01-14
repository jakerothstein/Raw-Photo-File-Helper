// Created by Jake Rothstein

package main

import (
	"fmt"
	"github.com/sqweek/dialog"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	fmt.Println(" _______  __   __       _______    .___  ___.   ______   ____    ____  _______ .______      \n|   ____||  | |  |     |   ____|   |   \\/   |  /  __  \\  \\   \\  /   / |   ____||   _  \\     \n|  |__   |  | |  |     |  |__      |  \\  /  | |  |  |  |  \\   \\/   /  |  |__   |  |_)  |    \n|   __|  |  | |  |     |   __|     |  |\\/|  | |  |  |  |   \\      /   |   __|  |      /     \n|  |     |  | |  `----.|  |____    |  |  |  | |  `--'  |    \\    /    |  |____ |  |\\  \\----.\n|__|     |__| |_______||_______|   |__|  |__|  \\______/      \\__/     |_______|| _| `._____|\nBy Jake Rothstein")

	fmt.Println("\nEnter File Extension Type to Be Scanned (.JPG, .jpg, .png, etc)")
	var extenTyp string
	fmt.Scanln(&extenTyp)

	fmt.Println("\nEnter Target File Extension Type (.ORF, .CR3, .NEF, etc)")
	var desiredExtenTyp string
	fmt.Scanln(&desiredExtenTyp)

	NewLocation, err := dialog.Directory().Title("Enter File Location of Items to Be Scanned").Browse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fileArr := scanArray(NewLocation, extenTyp, desiredExtenTyp) //Gets all the files with the correct file ending and changes to desired file ending
	fmt.Println("\nFiles to be Indexed and Moved:")
	fmt.Println(fileArr)

	//RawLocation, err := dialog.Directory().Title("Enter Location of Files to be Copied").Browse()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}

	fmt.Println("\nItems Successfully Located & Indexed")
	Location, err := dialog.Directory().Title("Enter Target File Location").Browse()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("\nMoving Files to Target Location (This May Take a Minute)")

	FaledFiles := moveFiles(fileArr, Location, NewLocation)
	fmt.Println(FaledFiles)

	fmt.Println("\nWould you like to download a moved files Log? (Y/N)")
	var response string
	fmt.Scanln(&response)
	if strings.ToUpper(response) == "Y" {
		fileData, failedFiles := getFileCopy(fileArr, FaledFiles)
		writeToFile(fileData, failedFiles) //Calls func to create file with output in downloads folder
		fmt.Println("\nFile saved to Downloads folder - Closing program in 10 seconds")
	} else {
		fmt.Println("\nClosing Program in 10 seconds")
	}

	time.Sleep(10 * time.Second)
}

func scanArray(dir string, typ string, desiredTyp string) []string { //Scans Array for Files
	var photoList []string

	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		photoList = append(photoList, file.Name())
	}
	return searchArray(photoList, typ, desiredTyp)
}

func searchArray(arr []string, typ string, desiredTyp string) []string { //Filers scanned array for specified file types
	var fileList []string
	for i := 0; i < len(arr); i++ {
		if arr[i][len(arr[i])-(len(typ)):] == typ {
			fileList = append(fileList, arr[i][:len(arr[i])-(len(typ))]+desiredTyp)
		}
	}
	return fileList

}

func moveFiles(arr []string, oldLocation string, newLocation string) (failedFiles []string) {
	for i := range arr {
		sourceFile := filepath.Join(oldLocation, arr[i])
		destFile := filepath.Join(newLocation, arr[i])

		source, err := os.Open(sourceFile)
		if err != nil {
			log.Printf("Error opening file %s: %v", sourceFile, err)
			failedFiles = append(failedFiles, arr[i])
			continue // Skip to the next iteration
		}
		defer source.Close()

		destination, err := os.Create(destFile)
		if err != nil {
			log.Printf("Error creating file %s: %v", destFile, err)
			failedFiles = append(failedFiles, arr[i])
			continue // Skip to the next iteration
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			log.Printf("Error copying data for file %s: %v", sourceFile, err)
			failedFiles = append(failedFiles, arr[i])
			continue // Skip to the next iteration
		}
	}
	return failedFiles
}
func StringInSlice(str string, list []string) bool {
	for _, s := range list {
		if str == s {
			return true
		}
	}
	return false
}
func getFileCopy(arr, FailedFiles []string) (fileList, failedFilesOut string) { // Puts array into format for Windows file explorer
	var fileListStr string
	var failedFilesStr string
	fileListStr = fileListStr + "Files Successfully Moved:\n"
	failedFilesStr = failedFilesStr + "Files Failed to Move:\n"

	for i := 0; i < len(arr); i++ {
		if StringInSlice(arr[i], FailedFiles) {
			failedFilesStr += "\"" + arr[i] + "\","
		} else {
			fileListStr += "\"" + arr[i] + "\","
		}
	}

	// Remove the trailing commas
	fileListStr = strings.TrimRight(fileListStr, ",")
	failedFilesStr = strings.TrimRight(failedFilesStr, ",")

	return fileListStr, failedFilesStr
}

func writeToFile(data, failedFiles string) {
	homeDir, _ := os.UserHomeDir()
	downloadsPath := filepath.Join(homeDir, "Downloads")
	baseFileName := "log.txt"
	filePath := filepath.Join(downloadsPath, baseFileName)

	// Check if the file already exists
	_, err := os.Stat(filePath)
	counter := 1

	// If file exists, create a new file with a different name
	for !os.IsNotExist(err) {
		newFileName := fmt.Sprintf("log_%d.txt", counter)
		filePath = filepath.Join(downloadsPath, newFileName)
		counter++
		_, err = os.Stat(filePath)
	}

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	f.Write([]byte(data))
	f.Write([]byte("\n\n"))
	f.Write([]byte(failedFiles))
}
