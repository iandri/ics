package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jordic/goics"
)

var (
	filePath = flag.String("f", "", "input xml file")
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text        string `xml:",chardata"`
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		BuildInfo   struct {
			Text        string `xml:",chardata"`
			Version     string `xml:"version"`
			BuildNumber string `xml:"build-number"`
			BuildDate   string `xml:"build-date"`
		} `xml:"build-info"`
		Item struct {
			Text    string `xml:",chardata"`
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			Project struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
				Key  string `xml:"key,attr"`
			} `xml:"project"`
			Description string `xml:"description"`
			Environment string `xml:"environment"`
			Key         struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"key"`
			Summary string `xml:"summary"`
			Type    struct {
				Text    string `xml:",chardata"`
				ID      string `xml:"id,attr"`
				IconUrl string `xml:"iconUrl,attr"`
			} `xml:"type"`
			Priority struct {
				Text    string `xml:",chardata"`
				ID      string `xml:"id,attr"`
				IconUrl string `xml:"iconUrl,attr"`
			} `xml:"priority"`
			Status struct {
				Text        string `xml:",chardata"`
				ID          string `xml:"id,attr"`
				IconUrl     string `xml:"iconUrl,attr"`
				Description string `xml:"description,attr"`
			} `xml:"status"`
			StatusCategory struct {
				Text      string `xml:",chardata"`
				ID        string `xml:"id,attr"`
				Key       string `xml:"key,attr"`
				ColorName string `xml:"colorName,attr"`
			} `xml:"statusCategory"`
			Resolution struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"resolution"`
			Assignee struct {
				Text     string `xml:",chardata"`
				Username string `xml:"username,attr"`
			} `xml:"assignee"`
			Reporter struct {
				Text     string `xml:",chardata"`
				Username string `xml:"username,attr"`
			} `xml:"reporter"`
			Labels       string `xml:"labels"`
			Created      string `xml:"created"`
			Updated      string `xml:"updated"`
			Due          string `xml:"due"`
			Votes        string `xml:"votes"`
			Watches      string `xml:"watches"`
			Attachments  string `xml:"attachments"`
			Subtasks     string `xml:"subtasks"`
			Customfields struct {
				Text        string `xml:",chardata"`
				Customfield []struct {
					Text              string `xml:",chardata"`
					ID                string `xml:"id,attr"`
					Key               string `xml:"key,attr"`
					Customfieldname   string `xml:"customfieldname"`
					Customfieldvalues struct {
						Text             string `xml:",chardata"`
						Customfieldvalue struct {
							Text string `xml:",chardata"`
							Key  string `xml:"key,attr"`
						} `xml:"customfieldvalue"`
					} `xml:"customfieldvalues"`
				} `xml:"customfield"`
			} `xml:"customfields"`
		} `xml:"item"`
	} `xml:"channel"`
}

type collection struct {
	Reporter    string
	Created     time.Time
	Nits        int
	Extra       string
	Location    string
	Description string
	Summary     string
	Link        string
	Categories  string
	Dtstamp     string
	Version     string
}

func (ev collection) EmitICal() goics.Componenter {

	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCALE", "GREGORIAN")
	c.AddProperty("PRODID", "-//ZContent.net//Zap Calendar 1.0//EN")
	c.AddProperty("VERSION", ev.Version)

	s := goics.NewComponent()
	s.SetType("VEVENT")
	dtend := ev.Created.AddDate(0, 0, ev.Nits)
	k, v := goics.FormatDateField("DTEND", dtend)
	s.AddProperty(k, v)
	k, v = goics.FormatDateField("DTSTART", ev.Created)
	s.AddProperty(k, v)
	s.AddProperty("UID", ev.Reporter)
	s.AddProperty("DESCRIPTION", ev.Description)
	s.AddProperty("SUMMARY", ev.Summary)
	s.AddProperty("LOCATION", ev.Location)
	s.AddProperty("URL", ev.Link)
	s.AddProperty("CATEGORIES", ev.Categories)
	s.AddProperty("DTSTAMP", ev.Dtstamp)

	c.AddComponent(s)

	return c
}

func main() {
	flag.Parse()

	if *filePath == "" {
		log.Fatal("input file is missing")
	}

	f, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatal(err)
	}

	var rss Rss

	if err := xml.Unmarshal(f, &rss); err != nil {
		log.Fatal(err)
	}

	t, _ := time.Parse(time.RFC1123Z, rss.Channel.Item.Created)

	col := collection{
		Created:     t,
		Description: rss.Channel.Item.Description,
		Summary:     rss.Channel.Item.Summary,
		Reporter:    rss.Channel.Item.Reporter.Text,
		Location:    rss.Channel.Item.Key.Text,
		Link:        rss.Channel.Item.Link,
		Dtstamp:     "20150421T141403",
		Version:     "2.0",
		Categories:  "work",
	}

	w := &bytes.Buffer{}
	goics.NewICalEncode(w).Encode(col)

	// fmt.Println(w.String())

	file := fmt.Sprintf("%s.ics", rss.Channel.Item.Key.Text)

	if err := ioutil.WriteFile(file, w.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println(file, "created")
}
