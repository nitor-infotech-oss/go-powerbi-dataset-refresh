# pbi-refresh-go

[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

# Objective

PowerBI provides two types of connections with the datasources.
  - Live Connection/Direct Query
  - Import Data - In this option, data is imported in PowerBI dataset and it needs to be refreshed on regular basis, when the data is updated in the database.

PowerBi provides timebased scheduling for refreshing the datasets, but this option may not be suitable in case of the syncronization with ETL Processes, refreshing the datawarehouse.

# Solution Approach

Microsoft PowerBI offers restful APIs Endpoints to programmatically work with PowerBI components like datasets, reports, etc.
This approach is configuration driven to refresh the list of PowerBI datasets and make the process event driven, instead of using time based scheduling.

# Configurations

  - Settings.json - This will work as configuration file for golang package and contains below information.
    Credentials

        {
            "jsonfilepath" : "dataset_config.json",
            "logfilepath"  : "refresh_history.log",
            "authorityurl" : "https://login.microsoftonline.com/",
            "resourceurl"  : "https://analysis.windows.net/powerbi/api",
            "clientid"     : "c1c04507-9857-44cc-xxxx-xxxxxxx",
            "username"     : "pbiusername",
            "password"     : "**********",
            "groupids"     : ["ba34f3fb-e648-4276-8fd5-66dfab523d1d"], <WorkspaceId>
            "webhookurl"   : "https://hooks.slack.com/services/*******" <sends notification on slack if configured, else no notification is sent>
        }

  - dataset_config.json - This will work as configuration file for golang package and contains below information.
        This file is reference point to pick which datasets needs to be refreshed. Here you can give details of one or multiple datasets as below.

            [{
                "ID" : 1,
                "SubjectArea" : "Power BI Samples",
                "DataSetName" : "Power BI Dataset Refresh Demo",
                "RefreshFrequencyType" : "Daily",
                "RefreshFrequency" : null,
                "IsActive" : 1
            }]

  - refresh_history.log - This stores log for the refreshes done so far.


### Prerequisites

GoLang packages dependencies:
* [go-autorest](https://github.com/Azure/go-autorest/autorest/adal)
* [jwt-go](https://github.com/dgrijalva/jwt-go)

### Imprting Package
```
import (
  "github.com/nitor-infotech-oss/go-powerbi-dataset-refresh/pbirefresh"
)
```

Package can be imported by simply using the path of repo where it resides.

### Example Usage

```
package main

import (
  "github.com/nitor-infotech-oss/go-powerbi-dataset-refresh/pbirefresh"
)

func main() {
  pbirefresh.RefreshDataset("settings.json")
}
```

### Todos

 - Write MORE Tests
