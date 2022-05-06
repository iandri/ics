package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"

	"github.com/jordic/goics"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	ctx := context.Background()

	serverFlag := pflag.StringP("server", "s", "", "server url, or export env JIRA_SERVER")
	username := pflag.StringP("username", "u", "", "username, or export env JIRA_USERNAME")
	password := pflag.StringP("password", "p", "", "password, or export env JIRA_PASSWORD")
	ticket := pflag.StringP("ticket", "t", "", "ticker number, (ex. OP-7378)")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	if v, ok := os.LookupEnv("JIRA_SERVER"); ok {
		pflag.Set("server", v)
	}

	if v, ok := os.LookupEnv("JIRA_USERNAME"); ok {
		pflag.Set("username", v)
	}

	if v, ok := os.LookupEnv("JIRA_PASSWORD"); ok {
		pflag.Set("password", v)
	}

	pflag.VisitAll(visitFlagSetMap)

	body, err := GetTicket(ctx, *serverFlag, *username, *password, *ticket)
	if err != nil {
		log.Fatal(err)
	}

	var rss *Rss

	if err := xml.Unmarshal(body, &rss); err != nil {
		log.Fatal(err)
	}

	var eventStart time.Time
	var eventEnd time.Time
	resolution := rss.Channel.Item.Resolution.Text
	customTimeFormat := "Mon, _2 Jan 2006 15:04:05 -0700"
	description := rss.Channel.Item.Description

	if t, found := rss.getCustomFieldValue("customfield_16771"); found {
		fmt.Println(t)
		if ts, err := time.Parse(customTimeFormat, t); err == nil {
			eventStart = ts
		} else {
			log.Fatal(err)
		}
	} else {
		if ts, err := time.Parse(time.RFC1123Z, rss.Channel.Item.Created); err != nil {
			eventStart = ts
		} else {
			log.Fatal(err)
		}
	}

	if t, found := rss.getCustomFieldValue("customfield_16775"); found {
		minutes, err := strconv.ParseFloat(t, 32)
		if err != nil {
			log.Fatal(err)
		}
		eventEnd = eventStart.Add(time.Duration(minutes) * time.Minute)
	} else {
		if resolution == "Unresolved" {
			eventEnd = time.Now().AddDate(1, 0, 0)
		} else {
			eventEnd = eventStart
		}
	}

	if v, found := rss.getCustomFieldValue("customfield_16782"); found {
		description = rss.Channel.Item.Description + " \n" + v
	}

	col := collection{
		DTStart:     eventStart,
		DTEnd:       eventEnd,
		Description: description,
		Summary:     rss.Channel.Item.Summary,
		Reporter:    rss.Channel.Item.Reporter.Text,
		Location:    rss.Channel.Item.Key.Text,
		Link:        rss.Channel.Item.Link,
		Dtstamp:     "20150421T141403",
		Version:     "2.0",
		Categories:  "work",
		Resolution:  rss.Channel.Item.Resolution.Text,
	}

	writer := &bytes.Buffer{}
	goics.NewICalEncode(writer).Encode(col)

	file := fmt.Sprintf("%s.ics", rss.Channel.Item.Key.Text)

	if err := ioutil.WriteFile(file, writer.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println(file, "created")
}

func (rss *Rss) getCustomFieldValue(id string) (string, bool) {
	customFields := rss.Channel.Item.Customfields
	for _, customField := range customFields.Customfield {
		if customField.ID == id {
			return customField.Customfieldvalues.Customfieldvalue[0].Text, true
		}
	}
	return "", false
}

func visitFlagSetMap(f *pflag.Flag) {
	if f.Value.String() == "" {
		log.Fatalf("%s is required", f.Name)
	}
}
