package db

import (
	"encoding/json"
	"fmt"

	eureka_client "github.com/xuanbo/eureka-client"
)

func PingEureka() {
	client := eureka_client.NewClient(&eureka_client.Config{
		DefaultZone:           "http://discoveryServer:8761/eureka/",
		App:                   "security-service",
		Port:                  8090,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODCUT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	})

	client.Start()

	apps := client.Applications
	b, _ := json.Marshal(apps)
	fmt.Println(b)
}