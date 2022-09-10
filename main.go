package main

import (
	"fmt"
	"gdyndns/pkg/configuration"
	gandiclient "gdyndns/pkg/gandi_client"
	ipapiclient "gdyndns/pkg/ip_api_client"

	"github.com/gookit/goutil/arrutil"
)

func main() {
	config := configuration.New()
	ipApiClient := ipapiclient.New(ipapiclient.IpApiClientConfig{})
	gandiClient := gandiclient.New(gandiclient.GandiClientConfig{
		ApiKey: config.GandiV5ApiKey,
	})

	currentIp, ipError := ipApiClient.GetExternalIp()
	if ipError != nil {
		panic(fmt.Errorf("failed to lookup current ip address: %w", ipError))
	}

	for _, record := range config.Records {
		log := fmt.Sprintf("(%s, %s, %s):", record.Zone, record.Name, record.Type)

		fmt.Printf("%s checking record\n", log)
		info, err := gandiClient.GetRecord(record)
		if err != nil {
			panic(fmt.Errorf("%s %w", log, err))
		}

		if arrutil.StringsHas(info.Values, currentIp) {
			fmt.Printf("%s no need to update record\n", log)
			continue
		}

		fmt.Printf("%s updating record with new ip address %s\n", log, currentIp)
		err = gandiClient.UpdateRecord(record, gandiclient.UpdateRecordPayload{
			Values: []string{currentIp},
			TTL: record.TTL,
		})
		if err != nil {
			fmt.Printf("%s failed to update record: %s\n", log, err)
		}
	}
}
