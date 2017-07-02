package main

// This is also a learning project for using file / buffer in golang
// WW Z

import (
	"os"
	"io"
	"encoding/binary"
	"fmt"
	"log"
	"bytes"
)

type dictionary struct {
     id [65536]int32
     offset [65536]int32
}

func readHeader(file io.Reader) (int32, int32){
     var count int32
     var dataSize int32
     err := binary.Read(file, binary.LittleEndian, &count)
     if err != nil {
	fmt.Println("binary.Read failed:", err)
     }
     binary.Read(file, binary.LittleEndian, &dataSize)
     return count, dataSize
}


func main() {
     var dic1 dictionary

     if len(os.Args) != 2 || os.Args[1] == "-h" {
     	fmt.Println("Convert .ilstring or .dlstring to tab-separated-vector format.")
	fmt.Println("Provide a single file as input. (This is a golang learning project.)")
	return
     }
     
     filename := os.Args[1]
     file, err := os.Open(filename)
     log.Println("Reading", filename)
     if err != nil {
     	log.Fatal(err)
     }
     count, dataSize := readHeader(file)
     log.Println("The file has", count, "records,")
     log.Println("Consisting", dataSize, "bytes")

     var i int32
     for i = 0; i < count; i++ {
     	 binary.Read(file, binary.LittleEndian, &dic1.id[i])
	 binary.Read(file, binary.LittleEndian, &dic1.offset[i])
     }
     
     data := make([]byte, dataSize)
     n, err := file.Read(data)
     if int32(n) != dataSize {
     	log.Fatal(err)
     }

     reader := bytes.NewReader(data)

     buffer := make([]byte, 1024)
     var length int32
     for i = 0; i < count; i++ {
          reader.Seek(int64(dic1.offset[i]),0)
	  binary.Read(reader, binary.LittleEndian, &length)
	  length-- // last byte is \0
	  reader.Read(buffer)
	  fmt.Println(dic1.id[i], string(buffer[:length]))
     }
}

