package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/antonlindstrom/feedfinder"
	"github.com/olekukonko/tablewriter"
)

func main() {
	var verbose = flag.Bool("verbose", true, "Toggle verbose output")
	flag.Parse()

	client := feedfinder.New("comment|alt=rss")

	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}

	doc, err := client.DocumentFromURL(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *verbose {
		printTable(doc)
		return
	}

	for _, link := range doc.Links {
		if link.ResponseCode == 200 {
			fmt.Println(link.URL)
		}
	}
}

func printTable(doc *feedfinder.Document) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"URL", "Content-Type", "Response", "ETag"})

	for _, link := range doc.Links {
		table.Append([]string{link.URL, link.ContentType, fmt.Sprintf("%d", link.ResponseCode), link.ETag})
	}

	table.Render()
}
