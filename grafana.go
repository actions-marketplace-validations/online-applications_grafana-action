package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"strings"
	"os"
	"fmt"
	"time"
	"bytes"
)

	
type Annotation struct {
	Dashboardid int      `json:"dashboardId"`
	Panelid     int      `json:"panelId"`
	Time        int64    `json:"time"`
	Timeend     int64    `json:"timeEnd"`
	Tags        []string `json:"tags"`
	Text        string   `json:"text"`
}

	
type AnnotationResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

type Folders []struct {
	ID    int    `json:"id"`
	UID   string `json:"uid"`
	Title string `json:"title"`
}

type Dashboards []struct {
	ID          int           `json:"id"`
	UID         string        `json:"uid"`
	Title       string        `json:"title"`
	URI         string        `json:"uri"`
	URL         string        `json:"grafanaBaseUrl"`
	Slug        string        `json:"slug"`
	Type        string        `json:"type"`
	Tags        []interface{} `json:"tags"`
	Isstarred   bool          `json:"isStarred"`
	Folderid    int           `json:"folderId"`
	Folderuid   string        `json:"folderUid"`
	Foldertitle string        `json:"folderTitle"`
	Folderurl   string        `json:"folderUrl"`
}

type Dashboard struct {
	Dashboard struct {
		ID           int           `json:"id"`
		Panels       []struct {
			ID        int `json:"id"`
			Title     string      `json:"title"`
			Type      string      `json:"type"`
			
		} `json:"panels"`
		Refresh       bool          `json:"refresh"`
		Schemaversion int           `json:"schemaVersion"`
		Style         string        `json:"style"`
		Tags          []interface{} `json:"tags"`
		Templating    struct {
			List []interface{} `json:"list"`
		} `json:"templating"`
		Time struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"time"`
		Timepicker struct {
			RefreshIntervals []string `json:"refresh_intervals"`
		} `json:"timepicker"`
		Timezone string `json:"timezone"`
		Title    string `json:"title"`
		UID      string `json:"uid"`
		Version  int    `json:"version"`
	} `json:"dashboard"`
}

var client http.Client = http.Client{}
var grafanaBaseUrl string = os.Getenv("GRAFANA_URL")
var projectName string = os.Getenv("PROJECT_NAME")
var folderName string = os.Getenv("TEAM")
var grafanaToken string = os.Getenv("GRAFANA_API_TOKEN")


func callGrafana(path string)([]byte, error){
	req, err := http.NewRequest("GET", grafanaBaseUrl + path, nil)
	req.Header.Set("Authorization", "Bearer " + grafanaToken)
	res, err := client.Do(req)
	
	if err != nil {
		log.Fatal( err )
		return nil, err
	}
	data, err := ioutil.ReadAll( res.Body )

	if err != nil {
		log.Fatal( err )
		return nil, err
	}
	// res.Body.Close()
	return data , nil
}

func main(){
	foldersData, err := callGrafana("/folders")

	if err != nil {
		log.Fatal( err )
	}

	folders := Folders{}
    json.Unmarshal(foldersData, &folders)
	log.Println(folders)
	for _, folder := range folders {
		if strings.ToLower(folder.Title) == strings.ToLower(folderName) {
			log.Println("Grafana Folder Found")
			path := fmt.Sprintf("/search?folderIds=%d&query=%s", folder.ID, strings.ToUpper(projectName))
			dashboardsData, err := callGrafana(path)
			if err != nil {
				log.Fatal( err )
			}
			dashboards := Dashboards{}
			json.Unmarshal(dashboardsData, &dashboards)
			for _, dashboard := range dashboards {
				path := fmt.Sprintf("/dashboards/uid/%s", dashboard.UID)
				dashboardData, err := callGrafana(path)
				if err != nil {
					log.Fatal( err )
				}
				dashboard := Dashboard{}
				json.Unmarshal(dashboardData, &dashboard)
				log.Println("Grafana Dashboard Found")
				for _, panel := range dashboard.Dashboard.Panels {
					log.Println("Grafana Panel Found")
					now := time.Now()
					nanos := now.UnixNano()
					millis := nanos / 1000000
					annotation := Annotation{
						Dashboardid: dashboard.Dashboard.ID,
						Panelid: panel.ID,
						Time: millis,
						Timeend: millis,
						Tags: []string{"GitHub_Action", "deploy", "cicd", projectName},
						Text: "Deploy success",
					}
					log.Println(annotation)
					payload, err := json.Marshal(annotation)
					if err != nil {
						log.Fatal( err )
					}

					req, err := http.NewRequest("POST", grafanaBaseUrl + "/annotations" , bytes.NewBuffer(payload))
					if err != nil {
						log.Fatal( err )
					}
					req.Header.Set("Content-Type", "application/json")
					req.Header.Set("Authorization", "Bearer "+ grafanaToken)
					res, err := client.Do(req)
				
					if err != nil {
						log.Fatal( err )
					}
					body, err := ioutil.ReadAll(res.Body)
				
					var result AnnotationResponse
					err = json.Unmarshal([]byte(body), &result)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("Grafana New Annotation Added")
				}
			}
		}
	}
}