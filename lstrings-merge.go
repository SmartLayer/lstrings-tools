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
     id     [65536]int32
     offset [65536]int32
     length [65536]int32
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

     if len(os.Args) != 3 || os.Args[1] == "-h" {
     	fmt.Println("Convert .ilstring or .dlstring to tab-separated-vector format.")
	fmt.Println("Provide a single file as input. (This is a golang learning project.)")
	return
     }

     var dic1, dic2 dictionary
     var count1, count2 int32
     var dataSize2, dataSize1 int32
     reader1 := bytes.NewReader(readfile(os.Args[1], &dic1, &count1, &dataSize1))
     reader2 := bytes.NewReader(readfile(os.Args[2], &dic2, &count2, &dataSize2))
     if count1 != count2 {
     	log.Fatal("As an experiment project we don't process files with different record count yet.")
     }


     // new dictionary
     var length int32
     var offset int32
     offset = 0
     for i := count1 - count2; i< count1 & count2; i++ {
     	 if dic1.id[i] != dic2.id[i] {
	      	log.Fatal("As an experiment project we don't process files with different ID sets yet.")
	 }
          reader1.Seek(int64(dic1.offset[i]),0)
     	  reader2.Seek(int64(dic2.offset[i]),0)
	  binary.Read(reader1, binary.LittleEndian, &length)
	  dic1.length[i] = length
	  binary.Read(reader2, binary.LittleEndian, &length)
	  dic2.length[i] = length
	  length = dic1.length[i] + dic2.length[i]
	  offset += length + 4
     }

     // write header and dictionary
     binary.Write(os.Stdout, binary.LittleEndian, count1)
     binary.Write(os.Stdout, binary.LittleEndian, offset)
     log.Printf("Datasize: %x\n",  offset)
     offset = 0
     for i := count1 - count2; i< count1 & count2; i++ {
          reader1.Seek(int64(dic1.offset[i]),0)
     	  reader2.Seek(int64(dic2.offset[i]),0)
	  binary.Read(reader1, binary.LittleEndian, &length)
	  dic1.length[i] = length
	  binary.Read(reader2, binary.LittleEndian, &length)
	  dic2.length[i] = length
	  length = dic1.length[i] + dic2.length[i]
	  binary.Write(os.Stdout, binary.LittleEndian, dic1.id[i])
	  binary.Write(os.Stdout, binary.LittleEndian, offset)
	  offset += length + 4
     }
     
     log.Printf("Final size: %x, %d\n", offset, offset)

     // write data
     var length1, length2 int32
     buffer := make([]byte, 1024) // let's hope no string is longer
     offset = 0
     for i := count1 - count2; i< count1 & count2; i++{
          reader1.Seek(int64(dic1.offset[i]),0)
	  reader2.Seek(int64(dic2.offset[i]),0)
	  binary.Read(reader1, binary.LittleEndian, &length1)
     	  reader1.Read(buffer)
	  //msg = string(buffer[:length1-1]) + "\n"
	  buffer[length1-1] = '\n'
	  binary.Read(reader2, binary.LittleEndian, &length2)
	  reader2.Read(buffer[length1:])
	  //msg = msg + string(buffer[:length-1])
	  length = length1 + length2
	  binary.Write(os.Stdout, binary.LittleEndian, length)
	  os.Stdout.Write(buffer[:length])
	  offset += length + 4
     }
     log.Printf("Final size: %x, %d\n", offset, offset)
}

func readfile(filename string, dic *dictionary, count *int32, dataSize *int32) []byte {
     file, err := os.Open(filename)
     log.Println("Reading", filename)
     if err != nil {
     	log.Fatal(err)
     }

     *count, *dataSize = readHeader(file)
     log.Println("The file has", *count, "records,")
     log.Println("Consisting", *dataSize, "bytes")

     for i := *count-*count; i < *count ; i++ {
     	 binary.Read(file, binary.LittleEndian, &((*dic).id[i]))
	 binary.Read(file, binary.LittleEndian, &((*dic).offset[i]))
     }
     
     data := make([]byte, *dataSize)
     n, err := file.Read(data)
     if int32(n) != *dataSize {
     	log.Fatal(err)
     }

     return data
}
