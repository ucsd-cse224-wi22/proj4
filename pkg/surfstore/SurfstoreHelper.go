package surfstore

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/* Hash Related */
func GetBlockHashBytes(blockData []byte) []byte {
	h := sha256.New()
	h.Write(blockData)
	return h.Sum(nil)
}

func GetBlockHashString(blockData []byte) string {
	blockHash := GetBlockHashBytes(blockData)
	return hex.EncodeToString(blockHash)
}

/* File Path Related */
func ConcatPath(baseDir, fileDir string) string {
	return baseDir + "/" + fileDir
}

// NewFileMetaDataFromConfig returns a FileMetaData struct
// associated with one line in the local metadata file.
func NewFileMetaDataFromConfig(configString string) *FileMetaData {
	configItems := strings.Split(configString, CONFIG_DELIMITER)

	filename := configItems[FILENAME_INDEX]
	version, _ := strconv.Atoi(configItems[VERSION_INDEX])
	blockHashList := strings.Split(configItems[HASH_LIST_INDEX], HASH_DELIMITER)

	return &FileMetaData{
		Filename:      filename,
		Version:       int32(version),
		BlockHashList: blockHashList[:len(blockHashList)-1],
	}
}

// LoadMetaFromMetaFiles loads the local metadata file into a file meta map.
// The key is the file's name and the value is the file's metadata.
// You can use this function to load the index.txt file in this project.
func LoadMetaFromMetaFile(baseDir string) (fileMetaMap map[string]*FileMetaData, e error) {
	metaFilePath, _ := filepath.Abs(ConcatPath(baseDir, DEFAULT_META_FILENAME))

	fileMetaMap = make(map[string]*FileMetaData)

	metaFileStats, e := os.Stat(metaFilePath)
	if e != nil || metaFileStats.IsDir() {
		return fileMetaMap, nil
	}
	metaFD, e := os.Open(metaFilePath)
	if e != nil {
		log.Fatal("Error When Opening Meta")
	}
	defer metaFD.Close()

	leftOverContent := ""
	metaReader := bufio.NewReader(metaFD)
	for {
		lineContent, isPrefix, e := metaReader.ReadLine()
		if e != nil && e != io.EOF {
			log.Fatal("Error During Reading Meta")
		}

		leftOverContent += string(lineContent)
		if isPrefix {
			continue
		}

		if len(leftOverContent) == 0 {
			break
		}

		currFileMeta := NewFileMetaDataFromConfig(leftOverContent)

		leftOverContent = ""
		fileMetaMap[currFileMeta.Filename] = currFileMeta
	}

	return fileMetaMap, nil
}

// PrintMetaMap prints the contents of the metadata map.
// You might find this function useful for debugging.
func PrintMetaMap(metaMap map[string]*FileMetaData) {

	fmt.Println("--------BEGIN PRINT MAP--------")

	for _, filemeta := range metaMap {
		fmt.Println("\t", filemeta.Filename, filemeta.Version)
		for _, blockHash := range filemeta.BlockHashList {
			fmt.Println("\t", blockHash)
		}
	}

	fmt.Println("---------END PRINT MAP--------")

}
