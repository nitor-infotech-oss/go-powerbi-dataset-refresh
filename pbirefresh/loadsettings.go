package pbirefresh

import (
  "encoding/json"
  "fmt"
  "os"
)

// struct object to load settings configured
type Settings struct {
  JSONFilePath  string `json:"jsonfilepath"`
  LogFilePath   string `json:"logfilepath"`
  AuthorityUrl  string `json:"authorityurl"`
  ResourceUrl   string `json:"resourceurl"`
  ClientID      string `json:"clientid"`
  Username      string `json:"username"`
  Password      string `json:"password"`
  GroupIDs      []string `json:"groupids"`
  WebookURL     string `json:"webhookurl"`
}

// method that reads settings file and converts it to Settings strut
func LoadSettings(file string) Settings {
  var settings Settings
  configFile, err := os.Open(file)
  defer configFile.Close()
  if err != nil {
      fmt.Println(err.Error())
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&settings)
  return settings
}
