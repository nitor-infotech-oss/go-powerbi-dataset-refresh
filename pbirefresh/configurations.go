package pbirefresh

import (
  "encoding/json"
  "fmt"
  "os"
  "strings"
)

// struct object to load dataset configurations
type DatasetConfig struct {
  Config []struct {
    ID                    string `json:"ID"`
    SubjectArea           string `json:"SubjectArea"`
    DataSetName           string `json:"DataSetName"`
    RefreshFrequencyType  string `json:"RefreshFrequencyType"`
    RefreshFrequency      string `json:"RefreshFrequency"`
    IsActive              string `json:"IsActive"`
  } `json:"config"`
}

// method that reads configuration file and converts it to DatasetConfig strut
func LoadConfiguration(file string) DatasetConfig {
  var config DatasetConfig
  configFile, err := os.Open(file)
  defer configFile.Close()
  if err != nil {
      fmt.Println(err.Error())
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&config)
  return config
}

// struct object to sotre dataset details
type DatasetDetails struct {
  ID string
  Name string
  GroupID string
}


// struct object to sotre frequency details
type FrequeryDetails struct {
  RefreshFrequencyType string
  RefreshFrequency string
}

// struct object to sotre powerbi response of the refresh url request
type DatasetsIdArray struct {
  ListIds []struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Status string `json:"status"`
    EndTime string `json:"endTime"`
  } `json:"value"`
}

// method to read the API response and convert it to strut object
func DatasetResponse(response string) DatasetsIdArray {
  var idList DatasetsIdArray
  jsonParser := json.NewDecoder(strings.NewReader(response))
  jsonParser.Decode(&idList)
  return idList
}

// method to get frequecy detail of a given dataset from the strut config object
func GetFrequencyFromConfig(config DatasetConfig, datasetName string) FrequeryDetails {
  var frequencyDetails FrequeryDetails
  for _, conf := range config.Config {
    if(conf.DataSetName == datasetName){
      fmt.Println(conf.DataSetName)
      return FrequeryDetails{conf.RefreshFrequencyType, conf.RefreshFrequency}
    }
  }
  return frequencyDetails
}
