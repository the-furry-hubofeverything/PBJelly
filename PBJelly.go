// PBJelly
// A uh.... "simple" script to update Porkbun DNS records using their API.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nrdcg/porkbun"
	"github.com/pelletier/go-toml"
)

var (
	c           config
	config_path = flag.String("c", "config.toml", "Path to config")
	domain      = os.Args[len(os.Args)-1]
	id          = flag.String("id", "", "ID of Porkbun DNS entry, use -l to view ID")
	interval    = flag.Duration("i", 1*time.Hour, "Time between updates")
	list        = flag.Bool("l", false, "List DNS records under domain")
	name        = flag.String("n", "", "Name of DNS record")
	oneshot     = flag.Bool("o", false, "Execute update only once (For cron or systemd)")
	recordType  = flag.String("t", "A", "Type of DNS record")
)

type config struct {
	Secret_key string
	API_key    string
	Optional   struct {
		Interval string
		ID       string
		Type     string
		Name     string
	}
}

func ReadConfig(filepath string) bool {
	// Reads config.toml, create it if it doesn't exist, and returns if permissions are 0600
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		d, err := toml.Marshal(config{})
		if err != nil {
			log.Fatal(err)
		}

		os.WriteFile(filepath, d, 0600)
		log.Print("Config file created!")
		os.Exit(0)
	} else if err != nil {
		log.Fatal(err)
	}

	if b, err := os.ReadFile(filepath); err != nil || toml.Unmarshal(b, &c) != nil {
		log.Fatal(err)
	}

	// Check permissions
	fi, err := os.Lstat(filepath)
	if err != nil {
		log.Fatal(err)
	}

	if c == (config{}) || c.API_key == "" || c.Secret_key == "" {
		log.Fatal("Incomplete Config file!")
	}

	return fi.Mode().Perm() == fs.FileMode(0600)
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <domain>\n ", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if ReadConfig(*config_path) {
		if flag.NArg() != 1 {
			flag.Usage()
			log.Fatal("Domain undefined")
		}

		// Reads config and fills in the optional fields if preexisting fields don't exist
		if c.Optional.Interval != "" {
			if t, err := time.ParseDuration(c.Optional.Interval); err == nil {
				interval = &t
				fmt.Println(*interval)
			} else {
				log.Fatal(err)
			}
		}
		if c.Optional.ID != "" && *id == "" {
			id = &c.Optional.ID
		}
		if c.Optional.Type != "" && *recordType == "" {
			recordType = &c.Optional.Type
		}
		if c.Optional.Name != "" && *name == "" {
			name = &c.Optional.Name
		}
	} else {
		log.Fatal("Config file permissions MUST be 600")
	}
}

func main() {
	client := porkbun.New(c.Secret_key, c.API_key)
	ctx := context.Background()

	for {
		// Change IP through the Porkbun API
		if pubIP, err := client.Ping(ctx); err == nil {
			records, err := client.RetrieveRecords(ctx, domain)
			if err != nil {
				log.Fatal(err)
			}

			for _, r := range records {
				if *list {
					log.Println("Current IP: " + pubIP)
					log.Printf("%+v\n", r)
					os.Exit(0)
				} else if (*id == "" || r.ID == *id) && (*name == "" || r.Name == *name) &&
					*recordType == r.Type && r.Content != pubIP {

					r.Content = pubIP
					if i, err := strconv.Atoi(r.ID); err != nil || client.EditRecord(ctx, domain, i, r) != nil {
						log.Panic(err)
					}
					log.Println("Updated " + *name + " to " + pubIP)

				}
			}

		} else {
			log.Fatal(err)
		}

		if *oneshot {
			break
		}
		time.Sleep(*interval)
	}
}
