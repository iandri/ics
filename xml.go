package main

import "encoding/xml"

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
			Labels     string `xml:"labels"`
			Created    string `xml:"created"`
			Updated    string `xml:"updated"`
			Due        string `xml:"due"`
			Votes      string `xml:"votes"`
			Watches    string `xml:"watches"`
			Issuelinks struct {
				Text          string `xml:",chardata"`
				Issuelinktype struct {
					Text        string `xml:",chardata"`
					ID          string `xml:"id,attr"`
					Name        string `xml:"name"`
					Inwardlinks struct {
						Text        string `xml:",chardata"`
						Description string `xml:"description,attr"`
						Issuelink   struct {
							Text     string `xml:",chardata"`
							Issuekey struct {
								Text string `xml:",chardata"`
								ID   string `xml:"id,attr"`
							} `xml:"issuekey"`
						} `xml:"issuelink"`
					} `xml:"inwardlinks"`
				} `xml:"issuelinktype"`
			} `xml:"issuelinks"`
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
						Customfieldvalue []struct {
							Text string `xml:",chardata"`
							Key  string `xml:"key,attr"`
						} `xml:"customfieldvalue"`
					} `xml:"customfieldvalues"`
				} `xml:"customfield"`
			} `xml:"customfields"`
		} `xml:"item"`
	} `xml:"channel"`
}
