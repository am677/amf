package message_test

import (
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/context"
	gmm_message "free5gc/src/amf/gmm/message"
	"free5gc/src/amf/handler"
	"free5gc/src/amf/util"
	"testing"
	"time"
)

func init() {
	go handler.Handle()
	TestAmf.SctpSever()

}

func TestSendIdentityRequest(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	message.SendIdentityRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMessage.MobileIdentity5GSTypeSuci)
	TestAmf.Conn.Close()
}

func TestSendAuthenticationRequest(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.NgKsi.Tsc = models.ScType_NATIVE
	ue.NgKsi.Ksi = 1
	ue.AuthenticationCtx = &models.UeAuthenticationCtx{
		AuthType: models.AuthType__5_G_AKA,
		Var5gAuthData: map[string]interface{}{
			"rand":      "0123456789abcdef0123456789abcdef",
			"hxresStar": "0123456789abcdef0123456789abcdef",
		},
	}
	message.SendAuthenticationRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS])

	TestAmf.Conn.Close()
}

func TestSendAuthenticationReject(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]

	message.SendAuthenticationReject(ue.RanUe[models.AccessType__3_GPP_ACCESS], "")
	TestAmf.Conn.Close()
}

func TestSendNotification(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	nasMsg, err := message.BuildNotification(ue, nasMessage.AccessTypeNon3GPP)
	if err != nil {
		t.Error(err.Error())
	}
	message.SendNotification(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMsg)
	// time.Sleep(1 * time.Second)
	util.ClearT3565(ue)
	TestAmf.Conn.Close()

}

func TestSendAuthenticationResult(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.ABBA = []uint8{0x00, 0x00}
	message.SendAuthenticationResult(ue.RanUe[models.AccessType__3_GPP_ACCESS], false, "AQIACTIMAUAA")
	TestAmf.Conn.Close()
}

func TestSendServiceAccept(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	var test [16]bool
	message.SendServiceAccept(ue.RanUe[models.AccessType__3_GPP_ACCESS], nil, &test, []uint8{10}, []uint8{92})
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestSendServiceReject(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	message.SendServiceReject(ue.RanUe[models.AccessType__3_GPP_ACCESS], nil, nasMessage.Cause5GMMUEIdentityCannotBeDerivedByTheNetwork)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestSendRegistrationReject(t *testing.T) {
	time.Sleep(100 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.T3502Value = 12 * 60
	message.SendRegistrationReject(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMessage.Cause5GMMIllegalUE, "")
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestSendSecurityModeCommand(t *testing.T) {
	time.Sleep(200 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.NasUESecurityCapability.SetLen(8)
	copy(ue.NasUESecurityCapability.Octet[:], []uint8{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00})
	copy(ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms[:], []uint8{0xe0, 0x00})
	copy(ue.SecurityCapabilities.NREncryptionAlgorithms[:], []uint8{0xe0, 0x00})
	message.SendSecurityModeCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], false, "")
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestSendDeregistrationRequestUETerminated(t *testing.T) {
	time.Sleep(200 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	// TODO: fill fake data to test if needed
	message.SendDeregistrationRequest(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMessage.AccessType3GPP, false, nasMessage.Cause5GMMImplicitlyDeregistered)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestSendDeregistrationAcceptUEOriginating(t *testing.T) {
	time.Sleep(200 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	message.SendDeregistrationAccept(ue.RanUe[models.AccessType__3_GPP_ACCESS])
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestSendRegistrationAccept(t *testing.T) {
	time.Sleep(200 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.Guti = "20893cafe0000000001"
	ue.T3502Value = context.DefaultT3502
	ue.T3512Value = 3600
	TestAmf.TestAmf.AllocateRegistrationArea(ue, models.AccessType__3_GPP_ACCESS)
	TestAmf.TestAmf.PlmnSupportList = append(TestAmf.TestAmf.PlmnSupportList, context.PlmnSupportItem{
		PlmnId: models.PlmnId{
			Mcc: "466",
			Mnc: "683",
		},
	})

	ue.LadnInfo = append(ue.LadnInfo, context.LADN{
		Ladn: "free5gc",
	})
	ue.LadnInfo[0].TaiLists = append(ue.LadnInfo[0].TaiLists, models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "000001",
	})
	message.SendRegistrationAccept(ue, models.AccessType__3_GPP_ACCESS, nil, nil, nil, nil, nil)
	time.Sleep(10 * time.Millisecond)
}

func TestSendStatus5GMM(t *testing.T) {
	time.Sleep(200 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	message.SendStatus5GMM(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMessage.Cause5GMMIllegalUE)
	time.Sleep(10 * time.Millisecond)
	TestAmf.Conn.Close()
}

func TestSendConfigurationUpdateCommand(t *testing.T) {
	time.Sleep(200 * time.Millisecond)
	TestAmf.AmfInit()
	TestAmf.SctpConnectToServer(models.AccessType__3_GPP_ACCESS)
	ue := TestAmf.TestAmf.UePool["imsi-2089300007487"]
	ue.Guti = "20893cafe0000000001"
	TestAmf.TestAmf.AllocateRegistrationArea(ue, models.AccessType__3_GPP_ACCESS)
	ue.LadnInfo = append(ue.LadnInfo, context.LADN{
		Ladn: "free5gc",
	})
	ue.LadnInfo[0].TaiLists = append(ue.LadnInfo[0].TaiLists, models.Tai{
		PlmnId: &models.PlmnId{
			Mcc: "208",
			Mnc: "93",
		},
		Tac: "000001",
	})
	message.SendConfigurationUpdateCommand(ue, models.AccessType__3_GPP_ACCESS, nil)
}
