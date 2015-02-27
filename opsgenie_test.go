package opsgenie

import (
	"fmt"
	"testing"
)

const (
	apiKey = "YOUR_OPS_GENIE_KEY"
)

func TestSendAlertWithWrongApiKey(t *testing.T) {
	fmt.Println("TestSendAlertWithWrongApiKey")

	opsGenie := New("12345-54321-12345-43212")

	alert, err := opsGenie.SendAlert("Problem on PROD")
	if err == nil {
		t.Errorf("SendAlert should not have been successful")
	}
	if alert != nil {
		t.Errorf("SendAlert should not have returned an Alert message")
	}
}

func TestSendRealAlert1(t *testing.T) {
	fmt.Println("TestSendRealAlert1")

	opsGenie := New(apiKey)

	alert, err := opsGenie.SendAlert("Problem on PROD")
	if err != nil {
		t.Errorf("SendAlert should have been successful")
	}

	alertId := alert.AlertId

	if alertId == "" {
		t.Errorf("SendAlert should have return an AlertId")
	}
}

func TestAcknowledgeInexixtingAlert(t *testing.T) {
	opsGenie := New(apiKey)

	a, err := opsGenie.AcknowledgeAlert("UNKNOWN")
	if err == nil {
		t.Errorf("AcknowledgeAlert should not have been successful")
	}
	if a != nil {
		t.Errorf("AcknowledgeAlert should not have returned an Alert message")
	}
}

func TestGetInexixtingAlert(t *testing.T) {
	fmt.Println("TestGetInexixtingAlert")
	opsGenie := New(apiKey)

	a, err := opsGenie.GetAlert("UNKNOWN")
	if err == nil {
		t.Errorf("GetAlert should not have been successful")
	}
	if a != nil {
		t.Errorf("GetAlert should not have returned an Alert message")
	}
}

func TestSendAndAcknowledgeAlert(t *testing.T) {
	fmt.Println("TestSendAndAcknowledgeAlert")

	opsGenie := New(apiKey)

	alert, err := opsGenie.SendAlert("Problem on PROD")
	if err != nil {
		t.Errorf("SendAlert should have been successful")
	}
	alertId := alert.AlertId

	alertDetail, err := opsGenie.GetAlert(alertId)
	if err != nil {
		t.Errorf("GetAlert should have been successful")
	}

	fmt.Println("Alert details:", alertDetail)

	alertList, err := opsGenie.ListAlerts("open")
	if err != nil {
		t.Errorf("ListAlerts should have been successful")
	}

	if foundAlert(alertList.Alerts, alertId) {
		t.Errorf("AcknowledgeAlert should have found my alert")
	}

	_, err = opsGenie.AcknowledgeAlert(alertId)
	if err != nil {
		t.Errorf("AcknowledgeAlert should have been successful")
	}

	_, err = opsGenie.AcknowledgeAlert(alertId)
	if err == nil {
		t.Errorf("AcknowledgeAlert should not be able to acknowledge an already ackowledged alert")
	}
}

func TestListOpenAlerts(t *testing.T) {
	fmt.Println("TestListOpenAlerts")

	opsGenie := New(apiKey)

	alertList, err := opsGenie.ListAlerts("open")
	if err != nil {
		t.Errorf("ListAlerts should have been successful")
	}
	if alertList == nil {
		t.Errorf("ListAlerts should have returned something")
	}
}

func TestAcknowledgeOpenAlerts(t *testing.T) {
	fmt.Println("TestAcknowledgeOpenAlerts")

	opsGenie := New(apiKey)

	_, err := opsGenie.SendAlert("New Problem on PROD")

	alertList, err := opsGenie.ListAlerts("unacked")

	for _, alert := range alertList.Alerts {
		//fmt.Printf("Got alert: %+v", alert)
		alertId := alert.Id
		if alertId == "" {
			t.Errorf("AlertId should have been available. Alert:", alert)
			return
		}
		fmt.Println("Acknowledging alert #", alertId)
		_, err = opsGenie.AcknowledgeAlert(alertId)
		if err != nil {
			t.Errorf("AcknowledgeAlert should have been successful")
			return
		}
	}
}

func foundAlert(alerts []Alert, alertId string) bool {
	for _, elem := range alerts {
		if alertId == elem.AlertId {
			return true
		}
	}
	return false
}

func main() {

}
