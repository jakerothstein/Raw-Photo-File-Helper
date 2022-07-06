// Created by Jake Rothstein

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {

	fmt.Println(" _______  __   __       _______    .___  ___.   ______   ____    ____  _______ .______      \n|   ____||  | |  |     |   ____|   |   \\/   |  /  __  \\  \\   \\  /   / |   ____||   _  \\     \n|  |__   |  | |  |     |  |__      |  \\  /  | |  |  |  |  \\   \\/   /  |  |__   |  |_)  |    \n|   __|  |  | |  |     |   __|     |  |\\/|  | |  |  |  |   \\      /   |   __|  |      /     \n|  |     |  | |  `----.|  |____    |  |  |  | |  `--'  |    \\    /    |  |____ |  |\\  \\----.\n|__|     |__| |_______||_______|   |__|  |__|  \\______/      \\__/     |_______|| _| `._____|\nBy Jake Rothstein")

	fmt.Println("\nEnter File Extension Type to Be Scanned (.JPG, .jpg, .png, etc)")
	var extenTyp string
	fmt.Scanln(&extenTyp)

	fmt.Println("\nEnter File Extension Type to Be Changed to (.ORF, .CR3, .NEF, etc)")
	var desiredExtenTyp string
	fmt.Scanln(&desiredExtenTyp)

	fmt.Println("\nEnter File Location of Items to Be Scanned (C:\\\\Folder\\\\Folder\\\\Folder)")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	Location := scanner.Text()
	//oldLocation = "F:\\Photography\\Temp Dumps" //For testing

	fileArr := scanArray(Location, extenTyp, desiredExtenTyp) //Gets all the files with the correct file ending and changes to desired file ending

	fmt.Println("\nItems Successfully Located & Indexed\n")

	fileData := getFileCopy(fileArr)
	writeToFile(fileData) //Calls func to create file with output in downloads folder

	fmt.Println("\nFile saved to Downloads folder - Closing program in 10 seconds")

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
	path := filepath.Join(homeDir, "Downloads", "output.txt") //Change for where you want the output file to go
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	f.Write([]byte(data))

}
