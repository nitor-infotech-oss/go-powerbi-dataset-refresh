
package main

import (
  "github.com/nitor-infotech-oss/go-powerbi-dataset-refresh/pbirefresh"
)

func main() {
  pbirefresh.RefreshDataset("settings.json")
}
