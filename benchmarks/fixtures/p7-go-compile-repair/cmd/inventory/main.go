package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"inventory/api"
	"inventory/report"
	"inventory/store"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	show := flag.String("show", "", "print one SKU and exit")
	flag.Parse()
	st := store.New()
	st.Put(store.Item{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25})
	if *show != "" {
		it, err := st.Get(*show)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		out, _ := report.Render([]store.Item{it}, "table")
		fmt.Print(out)
		return
	}
	log.Fatal(http.ListenAndServe(*addr, (&api.Server{Store: st}).Routes()))
}
