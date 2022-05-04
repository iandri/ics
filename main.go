package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/spf13/pflag"

	"github.com/jordic/goics"
)

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

	var rss Rss

	if err := xml.Unmarshal(body, &rss); err != nil {
		log.Fatal(err)
	}

	t, err := time.Parse(time.RFC1123Z, rss.Channel.Item.Created)
	if err != nil {
		log.Fatal(err)
	}

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

func visitFlagSetMap(f *pflag.Flag) {
	if f.Value.String() == "" {
		log.Fatalf("%s is required", f.Name)
	}
}
