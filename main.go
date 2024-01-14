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
	"strconv"
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
	fmt.Println("\nMoving Files to Target Location (This make take a minute)")

	moveFiles(fileArr, Location, NewLocation)

	fmt.Println("\nWould you like to download a moved files Log? (Y/N)")
	var response string
	fmt.Scanln(&response)
	if strings.ToUpper(response) == "Y" {
		fileData := getFileCopy(fileArr)
		writeToFile(fileData) //Calls func to create file with output in downloads folder
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

func moveFiles(arr []string, oldLocation string, newLocation string) {
	for i := range arr {
		sourceFile := filepath.Join(oldLocation, arr[i])
		destFile := filepath.Join(newLocation, arr[i])

		source, err := os.Open(sourceFile)
		if err != nil {
			log.Printf("Error opening file %s: %v", sourceFile, err)
			continue // Skip to the next iteration
		}
		defer source.Close()

		destination, err := os.Create(destFile)
		if err != nil {
			log.Printf("Error creating file %s: %v", destFile, err)
			continue // Skip to the next iteration
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			log.Printf("Error copying data for file %s: %v", sourceFile, err)
			continue // Skip to the next iteration
		}
	}
}

func getFileCopy(arr []string) string { // Puts array into format for Windows file explorer
	var indentCnt = 0
	var rowCnt = 1
	var fileLst string
	fileLst = fileLst + strconv.Itoa(rowCnt) + ") "
	for i := 0; i < len(arr); i++ {
		indentCnt++
		fileLst = fileLst + "\"" + arr[i] + "\" "
		if indentCnt > 16 {
			fileLst = fileLst + "\n\n"
			rowCnt++
			fileLst = fileLst + strconv.Itoa(rowCnt) + ") "
			indentCnt = 0

		}
	}
	return fileLst
}

func writeToFile(data string) { //Writes .txt placed in the downloads folder
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, "Downloads", "log.txt") //Change for where you want the output file to go
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	f.Write([]byte(data))

}
