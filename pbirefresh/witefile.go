package pbirefresh

import (
    "fmt"
    "os"
)

// write the mesage string to provided filepath
// creates file is it does not exist, appends the string to given file, if it does.
func WriteToFile(filename string, message string) {
    fmt.Println("Writing to file")
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println("File does not exists or cannot be created")
        os.Exit(1)
    }
    defer file.Close()

    if _, err := file.WriteString(message);
    err != nil {
      fmt.Println(err)
    }
}
