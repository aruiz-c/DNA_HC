package main

import (
	"fmt"
	"time"

	"collector/aci"
)

func getClient(host, usr, pwd string) (aci.Client, error) {
	client, err := aci.NewClient(
		host, usr, pwd,
		aci.RequestTimeout(600),
	)
	if err != nil {
		return aci.Client{}, fmt.Errorf("failed to create ACI client: %v", err)
	}

	// Authenticate
	log.Info().Str("host", host).Msg("APIC host")
	log.Info().Str("user", usr).Msg("APIC username")
	log.Info().Msg("Authenticating to the APIC...")
	if err := client.Login(); err != nil {
		return aci.Client{}, fmt.Errorf("cannot authenticate to the APIC at %s: %v", host, err)
	}
	return client, nil
}

// Fetch data via API.
func fetchResource(client aci.Client, req *Request, arc archiveWriter) error {
	startTime := time.Now()
	log.Debug().Time("start_time", startTime).Msgf("begin: %s", req.prefix)

	log.Info().Str("resource", req.prefix).Msg("fetching resource...")
	log.Debug().Str("url", req.path).Msg("requesting resource")

	res, err := client.Get(req.path, req.mods...)
	if err != nil {
		return fmt.Errorf("failed to make request for %s: %v", req.path, err)
	}
	err = arc.add(req.prefix+".json", []byte(res.Get(req.filter).Raw))
	if err != nil {
		return err
	}
	log.Debug().
		TimeDiff("elapsed_time", time.Now(), startTime).
		Msgf("done: %s", req.prefix)
	return nil
}

func pause(msg string) {
	fmt.Println("Press enter to exit.")
	var throwaway string
	fmt.Scanln(&throwaway)
}
