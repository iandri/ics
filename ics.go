package main

import (
	"time"

	"github.com/jordic/goics"
)

type collection struct {
	Reporter    string
	DTStart     time.Time
	DTEnd       time.Time
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
	k, v := goics.FormatDateTimeField("DTEND", ev.DTEnd)
	s.AddProperty(k, v)
	k, v = goics.FormatDateTimeField("DTSTART", ev.DTStart)
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
