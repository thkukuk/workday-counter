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
	workdaysLabel = "Workdays"
	startDateLabel = "Since"
	startDate = time.Now()
	endDateLabel = "Until"
	endDate = time.Now()
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
	server := &http.Server{Addr: ":8080"}
	server.SetKeepAlivesEnabled(false)
	log.Fatal(server.ListenAndServe())
}

type TemplateArgs struct {
	Title          string
	Message        string
	WorkdaysLabel  string
        Workdays       string
	StartDateLabel string
	StartDate      string
	EndDateLabel   string
	EndDate        string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	workdays := CalcBusinessDays(startDate,endDate)

	indexTemplate.Execute(w, TemplateArgs{
		Title:             title,
		Message:           message,
		WorkdaysLabel:     workdaysLabel,
		Workdays:          strconv.Itoa(workdays),
		StartDateLabel:    startDateLabel,
		StartDate:         startDate.Format("January 02, 2006"),
		EndDateLabel:      endDateLabel,
		EndDate:           endDate.Format("January 02, 2006"),
	})
}
