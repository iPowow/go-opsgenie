package opsgenie

import (
	"fmt"
)

func (client *Client) SendAlert(msg string, description string) (*Alert, error) {
	alert, _, err := client.alert.create(&AlertRequest{Message: msg, Description: description, ApiKey: client.config.apiKey})
	if err != nil {
		fmt.Println("ERROR: failed to create the alert:", err)
		return nil, err
	}
	return alert, err
}

func (client *Client) AcknowledgeAlert(alertId string) (*Alert, error) {
	alert, _, err := client.alert.acknowledge(&AlertRequest{AlertId: alertId, ApiKey: client.config.apiKey})
	if err != nil {
		fmt.Println("ERROR: failed to acknowledge the alert:", err)
		return nil, err
	}
	return alert, err
}

func (client *Client) GetAlert(alertId string) (*Alert, error) {
	alert, _, err := client.alert.get(&AlertRequest{Id: alertId, ApiKey: client.config.apiKey})
	if err != nil {
		fmt.Println("ERROR: failed to get the alert:", err)
		return nil, err
	}
	return alert, err
}

// <status> may take one of open, acked, unacked, seen, notseen, closed
func (client *Client) ListAlerts(status string) (*AlertList, error) {
	alert, _, err := client.alert.list(&AlertRequest{Status: status, ApiKey: client.config.apiKey})
	if err != nil {
		fmt.Println("ERROR: failed to list the alerts:", err)
		return nil, err
	}
	return alert, err
}
