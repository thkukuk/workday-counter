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
	"strings"
	"fmt"

	"github.com/rickar/cal"
)

var (
	version="unreleased"
	indexTemplate *template.Template
	title = "Workday-Counter"
	message = "Workdays"
	country = "Germany"
	state = "Bayern"
	workdaysTitle = ""
	workdaysLabel = "Workdays"
	startDateLabel = "Since"
	startDate time.Time
	endDateLabel = "Until"
	endDate time.Time
	country1 = ""
	state1 = ""
	workdays1Title = ""
	workdays1Label = "Workdays"
	startDate1Label = "Since"
	startDate1 time.Time
	endDate1Label = "Until"
	endDate1 time.Time
	country2 = ""
	state2 = ""
	workdays2Title = ""
	workdays2Label = "Workdays"
	startDate2Label = "Since"
	startDate2 time.Time
	endDate2Label = "Until"
	endDate2 time.Time
	dir *string
	zero_time time.Time
	displayWorkday = 0
)

func init() {
	dir = flag.String ("dir", "", "directory to read files from")
}

func CalcBusinessDays(country string, state string,
	from time.Time, to time.Time) (int64, *cal.Calendar) {

	c := cal.NewCalendar()
	// change the holiday calculation behavior
	c.Observed = cal.ObservedExact

	if strings.EqualFold(country, "Germany") {
		// add holidays for the business
		cal.AddGermanHolidays(c)
		if strings.EqualFold(state, "Bayern") {
			// Nuremberg does not have Maria Himmelfahrt...
			c.AddHoliday(
				cal.DEHeiligeDreiKoenige,
				cal.DEFronleichnam,
				// cal.DEMariaHimmelfahrt,
				cal.DEAllerheiligen,
				cal.DEReformationstag2017,
			)
		}
	} else if strings.EqualFold(country, "China") {
		// Chinese holiday definition is missing
	} else {
		log.Printf("Unknown Country: %s, ignoring holidays\n",
			country)
	}

	return c.CountWorkdays(from, to), c
}

func date2str(cal *cal.Calendar, date time.Time) string {
	if cal == nil {
		return ""
	}
	if cal.IsWorkday(date) {
		return date.Format("January 02, 2006")
	} else {
		return date.Format("January 02, 2006") + " (free)"
	}
}

func main() {
	flag.Parse()
	indexTemplate = template.Must(template.ParseFiles(*dir + "/index.template"))
	if len(os.Getenv("TITLE")) > 0 {
		title = os.Getenv("TITLE")
	}
	if len(os.Getenv("MESSAGE")) > 0 {
		message = os.Getenv("MESSAGE")
	}
	if len(os.Getenv("DISPLAY_WORKDAY")) > 0 {
		displayWorkday, _ = strconv.Atoi(os.Getenv("DISPLAY_WORKDAY"))
	}
	if len(os.Getenv("COUNTRY")) > 0 {
		country = os.Getenv("COUNTRY")
	}
	if len(os.Getenv("STATE")) > 0 {
		state = os.Getenv("STATE")
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
	if len(os.Getenv("COUNTRY1")) > 0 {
		country1 = os.Getenv("COUNTRY1")
	}
	if len(os.Getenv("STATE1")) > 0 {
		state1 = os.Getenv("STATE1")
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

	if len(os.Getenv("COUNTRY2")) > 0 {
		country2 = os.Getenv("COUNTRY2")
	}
	if len(os.Getenv("STATE2")) > 0 {
		state2 = os.Getenv("STATE2")
	}
	if len(os.Getenv("WORKDAYS2_TITLE")) > 0 {
		workdays2Title = os.Getenv("WORKDAYS2_TITLE")
	}
	if len(os.Getenv("WORKDAYS2_LABEL")) > 0 {
		workdays2Label = os.Getenv("WORKDAYS2_LABEL")
	}
	if len(os.Getenv("STARTDATE2")) > 0 {
		startDate2, _ = time.Parse("2006-01-02", os.Getenv("STARTDATE2"))
	}
	if len(os.Getenv("ENDDATE2")) > 0 {
		endDate2, _ = time.Parse("2006-01-02", os.Getenv("ENDDATE2"))
	}
	if len(os.Getenv("STARTDATE2_LABEL")) > 0 {
		startDate2Label = os.Getenv("STARTDATE2_LABEL")
	}
	if len(os.Getenv("ENDDATE2_LABEL")) > 0 {
		endDate2Label = os.Getenv("ENDDATE2_LABEL")
	}

	log.Printf("Workday-Counter %s started\n", version)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/fonts/",
		func (w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *dir + r.URL.Path)
		})
        http.HandleFunc("/logos/",
                func (w http.ResponseWriter, r *http.Request) {
                        http.ServeFile(w, r, *dir + r.URL.Path)
                })
	server := &http.Server{Addr: ":8080"}
	server.SetKeepAlivesEnabled(false)
	log.Fatal(server.ListenAndServe())
}

type TemplateArgs struct {
	Title           string
	Message         string
	MessageNr       string
	WorkdaysTitle   string
	WorkdaysLabel   string
        Workdays        string
	StartDateLabel  string
	StartDate       string
	EndDateLabel    string
	EndDate         string
	Weeks           string
	Workdays1Title  string
	Workdays1Label  string
        Workdays1       string
	StartDate1Label string
	StartDate1      string
	EndDate1Label   string
	EndDate1        string
	Weeks1          string
	Workdays2Title  string
	Workdays2Label  string
        Workdays2       string
	StartDate2Label string
	StartDate2      string
	EndDate2Label   string
	EndDate2        string
	Weeks2          string
}

func getDate(date time.Time) time.Time {
	if date == zero_time {
		return time.Now()
	}
	return date
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var workdays1 int64
	var workdays1_str string
	var cal1 *cal.Calendar
	var workdays2 int64
	var workdays2_str string
	var cal2 *cal.Calendar

	workdays, cal := CalcBusinessDays(country, state,
		getDate(startDate), getDate(endDate))

	if startDate1 != zero_time || endDate1 != zero_time {
		workdays1, cal1 = CalcBusinessDays(country1, state1,
			getDate(startDate1), getDate(endDate1))
		workdays1_str = strconv.FormatInt(workdays1, 10)
	}

	if startDate2 != zero_time  || endDate2 != zero_time {
		workdays2, cal2 = CalcBusinessDays(country2, state2,
			getDate(startDate2), getDate(endDate2))
		workdays2_str = strconv.FormatInt(workdays2, 10)
	}

	var messageNr int64
	switch displayWorkday {
	case 0:
		messageNr = workdays;
	case 1:
		messageNr = workdays1;
	case 2:
		messageNr = workdays2;
	}

	err := indexTemplate.Execute(w, TemplateArgs{
		Title:             title,
		Message:           message,
		MessageNr:         strconv.FormatInt(messageNr, 10),
		WorkdaysTitle:     workdaysTitle,
		WorkdaysLabel:     workdaysLabel,
		Workdays:          strconv.FormatInt(workdays, 10),
		StartDateLabel:    startDateLabel,
		StartDate:         date2str(cal, getDate(startDate)),
		EndDateLabel:      endDateLabel,
		EndDate:           date2str(cal, getDate(endDate)),
		Weeks:             fmt.Sprintf("%.1f", getDate(endDate).Sub(getDate(startDate)).Hours() / (7*24)),

		Workdays1Title:    workdays1Title,
		Workdays1Label:    workdays1Label,
		Workdays1:         workdays1_str,
		StartDate1Label:   startDate1Label,
		StartDate1:        date2str(cal1, getDate(startDate1)),
		EndDate1Label:     endDate1Label,
		EndDate1:          date2str(cal1, getDate(endDate1)),
		Weeks1:            fmt.Sprintf("%.1f", getDate(endDate1).Sub(getDate(startDate1)).Hours() / (7*24)),

		Workdays2Title:    workdays2Title,
		Workdays2Label:    workdays2Label,
		Workdays2:         workdays2_str,
		StartDate2Label:   startDate2Label,
		StartDate2:        date2str(cal2, getDate(startDate2)),
		EndDate2Label:     endDate2Label,
		EndDate2:          date2str(cal2, getDate(endDate2)),
		Weeks2:            fmt.Sprintf("%.1f", getDate(endDate2).Sub(getDate(startDate2)).Hours() / (7*24)),
	})
	if err != nil {
		log.Println("Error executing template:", err)
        }
}
