package main

import (
	"time"

	"github.com/jordic/goics"
)

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
	Resolution  string
}

func (ev collection) EmitICal() goics.Componenter {

	c := goics.NewComponent()
	c.SetType("VCALENDAR")
	c.AddProperty("CALSCALE", "GREGORIAN")
	c.AddProperty("PRODID", "-//ZContent.net//Zap Calendar 1.0//EN")
	c.AddProperty("VERSION", ev.Version)

	s := goics.NewComponent()
	s.SetType("VEVENT")
	var dtend time.Time
	if ev.Resolution == "Unresolved" {
		now := time.Now()
		dtend = now.AddDate(1, 0, 0)
	} else {
		dtend = ev.Created.AddDate(0, 0, ev.Nits)
	}
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
