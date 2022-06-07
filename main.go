// Created by Jake Rothstein

package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	fmt.Println(" _______  __   __       _______    .___  ___.   ______   ____    ____  _______ .______      \n|   ____||  | |  |     |   ____|   |   \\/   |  /  __  \\  \\   \\  /   / |   ____||   _  \\     \n|  |__   |  | |  |     |  |__      |  \\  /  | |  |  |  |  \\   \\/   /  |  |__   |  |_)  |    \n|   __|  |  | |  |     |   __|     |  |\\/|  | |  |  |  |   \\      /   |   __|  |      /     \n|  |     |  | |  `----.|  |____    |  |  |  | |  `--'  |    \\    /    |  |____ |  |\\  \\----.\n|__|     |__| |_______||_______|   |__|  |__|  \\______/      \\__/     |_______|| _| `._____|\nBy Jake Rothstein")

	fmt.Println("\nEnter File Extension Type to Be Scanned (.raw, .CR3, .png, etc)")
	var extenTyp string
	fmt.Scanln(&extenTyp)

	fmt.Println("\nEnter File Extension Type to Be Changed to (.raw, .CR3, .png, etc)")
	var desiredExtenTyp string
	fmt.Scanln(&desiredExtenTyp)

	fmt.Println("\nEnter File Location of Items to Be Scanned")
	var oldLocation string
	fmt.Scanln(&oldLocation)

	fileArr := scanArray(oldLocation, extenTyp, desiredExtenTyp) //Gets all the files with the correct file ending and changes to desired file ending

	fmt.Println("\nItems Successfully Located & Indexed")

	copyPaste := getFileCopy(fileArr)
	fmt.Println(copyPaste)

}

func scanArray(dir string, typ string, desiredTyp string) []string {
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

func searchArray(arr []string, typ string, desiredTyp string) []string {
	var fileList []string
	for i := 0; i < len(arr); i++ {
		if arr[i][len(arr[i])-(len(typ)):] == typ {
			fileList = append(fileList, arr[i][:len(arr[i])-(len(typ))]+desiredTyp)
		}
	}
	return fileList

}

func getFileCopy(arr []string) string {
	var fileLst string
	for i := 0; i < len(arr); i++ {
		fileLst = fileLst + "\"" + arr[i] + "\" "
	}
	return fileLst
}
