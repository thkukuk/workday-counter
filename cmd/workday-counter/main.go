// Copyright 2020 Thorsten Kukuk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
	"strconv"
)

var (
	version="unreleased"
	indexTemplate *template.Template
	title = "Workday-Counter"
	message = "Workdays"
	workdaysTitle = ""
	workdaysLabel = "Workdays"
	startDateLabel = "Since"
	startDate = time.Now()
	endDateLabel = "Until"
	endDate = time.Now()
	workdays1Title = ""
	workdays1Label = "Workdays"
	startDate1Label = "Since"
	startDate1 time.Time
	endDate1Label = "Until"
	endDate1 = time.Now()
	dir *string
)

func init() {
	dir = flag.String ("dir", "", "directory to read files from")
}

func CalcBusinessDays(from time.Time, to time.Time) int {
	totalDays := float32(to.Sub(from) / (24 * time.Hour))
	weekDays := float32(from.Weekday()) - float32(to.Weekday())
	businessDays := int(1 + (totalDays*5-weekDays*2)/7)

	if to.Weekday() == time.Saturday {
		businessDays--
	}

	if from.Weekday() == time.Sunday {
		businessDays--
	}

	return businessDays
}

func main() {
	flag.Parse()
	if len(*dir) > 0 {
		*dir = *dir + "/"
	}
	indexTemplate = template.Must(template.ParseFiles(*dir + "index.template"))
	if len(os.Getenv("TITLE")) > 0 {
		title = os.Getenv("TITLE")
	}
	if len(os.Getenv("MESSAGE")) > 0 {
		message = os.Getenv("MESSAGE")
	}
	if len(os.Getenv("WORKDAYS_TITLE")) > 0 {
		workdaysTitle = os.Getenv("WORKDAYS_TITLE")
	}
	if len(os.Getenv("WORKDAYS_LABEL")) > 0 {
		workdaysLabel = os.Getenv("WORKDAYS_LABEL")
	}
	if len(os.Getenv("STARTDATE")) > 0 {
		startDate, _ = time.Parse("2006-01-02", os.Getenv("STARTDATE"))
	}
	if len(os.Getenv("ENDDATE")) > 0 {
		endDate, _ = time.Parse("2006-01-02", os.Getenv("ENDDATE"))
	}
	if len(os.Getenv("STARTDATE_LABEL")) > 0 {
		startDateLabel = os.Getenv("STARTDATE_LABEL")
	}
	if len(os.Getenv("ENDDATE_LABEL")) > 0 {
		endDateLabel = os.Getenv("ENDDATE_LABEL")
	}
	if len(os.Getenv("WORKDAYS1_TITLE")) > 0 {
		workdays1Title = os.Getenv("WORKDAYS1_TITLE")
	}
	if len(os.Getenv("WORKDAYS1_LABEL")) > 0 {
		workdays1Label = os.Getenv("WORKDAYS1_LABEL")
	}
	if len(os.Getenv("STARTDATE1")) > 0 {
		startDate1, _ = time.Parse("2006-01-02", os.Getenv("STARTDATE1"))
	}
	if len(os.Getenv("ENDDATE1")) > 0 {
		endDate1, _ = time.Parse("2006-01-02", os.Getenv("ENDDATE1"))
	}
	if len(os.Getenv("STARTDATE1_LABEL")) > 0 {
		startDate1Label = os.Getenv("STARTDATE1_LABEL")
	}
	if len(os.Getenv("ENDDATE1_LABEL")) > 0 {
		endDate1Label = os.Getenv("ENDDATE1_LABEL")
	}

	log.Printf("Workday-Counter %s started\n", version)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/openSUSE-Kubic-Logo.png",
		func (w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *dir + "openSUSE-Kubic-Logo.png")
		})
	http.HandleFunc("/SUSE-Logo.png",
		func (w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *dir + "SUSE-Logo.png")
		})
        http.HandleFunc("/SUSE-Logo.svg",
                func (w http.ResponseWriter, r *http.Request) {
                        http.ServeFile(w, r, *dir + "SUSE-Logo.svg")
                })
	server := &http.Server{Addr: ":8080"}
	server.SetKeepAlivesEnabled(false)
	log.Fatal(server.ListenAndServe())
}

type TemplateArgs struct {
	Title           string
	Message         string
	WorkdaysTitle   string
	WorkdaysLabel   string
        Workdays        string
	StartDateLabel  string
	StartDate       string
	EndDateLabel    string
	EndDate         string
	Workdays1Title  string
	Workdays1Label  string
        Workdays1       string
	StartDate1Label string
	StartDate1      string
	EndDate1Label   string
	EndDate1        string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	workdays := CalcBusinessDays(startDate,endDate)

	var workdays1 int
	var workdays1_str string
	var t time.Time
	if t != startDate1 {
		workdays1 = CalcBusinessDays(startDate1, endDate1)
		workdays1_str = strconv.Itoa(workdays1)
	}

	err := indexTemplate.Execute(w, TemplateArgs{
		Title:             title,
		Message:           message,
		WorkdaysTitle:     workdaysTitle,
		WorkdaysLabel:     workdaysLabel,
		Workdays:          strconv.Itoa(workdays),
		StartDateLabel:    startDateLabel,
		StartDate:         startDate.Format("January 02, 2006"),
		EndDateLabel:      endDateLabel,
		EndDate:           endDate.Format("January 02, 2006"),
		Workdays1Title:    workdays1Title,
		Workdays1Label:    workdays1Label,
		Workdays1:         workdays1_str,
		StartDate1Label:   startDate1Label,
		StartDate1:        startDate1.Format("January 02, 2006"),
		EndDate1Label:     endDate1Label,
		EndDate1:          endDate1.Format("January 02, 2006"),
	})
	if err != nil {
		log.Println("Error executing template:", err)
        }
}
