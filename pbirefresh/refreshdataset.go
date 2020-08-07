package pbirefresh

import (
  "time"
  "fmt"
  "strconv"
  "io/ioutil"
  "net/http"
  "github.com/Azure/go-autorest/autorest/adal"
  "github.com/nitor-infotech-oss/go-powerbi-dataset-refresh/slacknotifications"
)

// method to refresh the datasets based on the settings and config file
func RefreshDataset(configFile string) {
  // load gonfigurations from file
  settings := LoadSettings(configFile)
  config := LoadConfiguration(settings.JSONFilePath)
  dt := time.Now()
  tenantID := "common"

  // invoke oauth
  oauthConfig, err := adal.NewOAuthConfig(settings.AuthorityUrl, tenantID)

  // authorise and fetch token
  spt, err := adal.NewServicePrincipalTokenFromUsernamePassword(
  	*oauthConfig,
  	settings.ClientID,
  	settings.Username,
  	settings.Password,
  	settings.ResourceUrl)

  err  = spt.Refresh()
  var token string
  if (err == nil) {
      token = spt.Token().AccessToken
  }

  // create bearer for headers
  bearer := "Bearer " + token
  allDatasets := []DatasetDetails{}

  // iterate through the configured groupsis in config and get all the datasets for those groups.
  for _, groupId := range settings.GroupIDs {
      refreshURL := "https://api.powerbi.com/v1.0/myorg/groups/" + groupId + "/datasets"
      req, err := http.NewRequest("GET", refreshURL, nil)
      req.Header.Add("Authorization", bearer)
      client := &http.Client{}
      resp, err := client.Do(req)
      if err != nil {
          fmt.Println("Error on response.\n[ERRO] -", err)
      }
      bodyBytes, _ := ioutil.ReadAll(resp.Body)
      bodyString := string(bodyBytes)
      newDatasetArray := DatasetResponse(bodyString)
      for _, datasetId := range newDatasetArray.ListIds {
        allDatasets = append(allDatasets, DatasetDetails{datasetId.ID, datasetId.Name, groupId})
      }
  }

  // summaryString to be stored in logs and send to notifications
  var summaryString string
  summaryString = "\nPower BI Refresh Summary:\n"

  // iterate through datasets and check if it is from the configured ones
  for _, datasetDetails := range allDatasets {
    frequencyDetails := GetFrequencyFromConfig(config, datasetDetails.Name)
    if(frequencyDetails != FrequeryDetails{}){
      // check if the frequecy setting allows it to be refreshed at the execution time
      weeklyCondition := (frequencyDetails.RefreshFrequencyType == "Weekly" && frequencyDetails.RefreshFrequency==time.Now().Weekday().String())
      dailyCondition  := (frequencyDetails.RefreshFrequencyType == "Daily")
      monthlyCondition := (frequencyDetails.RefreshFrequencyType == "Monthly" && frequencyDetails.RefreshFrequency==strconv.Itoa(time.Now().Day()))
      if (weeklyCondition || dailyCondition || monthlyCondition){
          // call powerbi dataset refresh API
          finalRefreshURL := "https://api.powerbi.com/v1.0/myorg/groups/"+datasetDetails.GroupID+"/datasets/"+datasetDetails.ID+"/refreshes"
          req, err := http.NewRequest("POST", finalRefreshURL, nil)
          req.Header.Add("Authorization", bearer)
          client := &http.Client{}
          resp, err := client.Do(req)
          if err != nil {
              fmt.Println("Error on response.\n[ERRO] -", err, resp)
          }

          // wait untill the dataset refresh status is updated
          // check the status of the refresh by calling another endpoint
          running := 1
        	for running > 0 {
            refreshURLHistory := "https://api.powerbi.com/v1.0/myorg/datasets/"+datasetDetails.ID+"/refreshes?$top=1"
            fmt.Println(refreshURLHistory)
            req, err := http.NewRequest("GET", refreshURLHistory, nil)
            req.Header.Add("Authorization", bearer)
            client := &http.Client{}
            resp, err := client.Do(req)
            if err != nil {
                fmt.Println("Error on response.\n[ERRO] -", err)
            }
            bodyBytes, _ := ioutil.ReadAll(resp.Body)
            bodyString := string(bodyBytes)
            responseArray := DatasetResponse(bodyString)
            fmt.Println(responseArray)
            if(responseArray.ListIds[0].Status=="Failed"){
              summaryString = summaryString + datasetDetails.Name + " - " + responseArray.ListIds[0].Status
              running = 0
              break
            } else if(responseArray.ListIds[0].Status=="Completed"){
              summaryString = summaryString + datasetDetails.Name + " - " + responseArray.ListIds[0].Status + " - " + responseArray.ListIds[0].EndTime
              running = 0
              break
            }
            time.Sleep(30 * time.Second)
          }
        }

    }
  }
  // check if slack url is configured, send notification with summary if it is.
  if (settings.WebookURL ! =""){
    webhookUrl := settings.WebookURL
    err = slacknotifications.SendSlackNotification(webhookUrl, summaryString)
    if err != nil {
      fmt.Println("Error while sending summary to slack")
    }
  }
  // write the summary to the configured log file
  WriteToFile(settings.LogFilePath, summaryString)
}
