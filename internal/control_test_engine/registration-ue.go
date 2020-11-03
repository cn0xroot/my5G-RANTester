package control_test_engine

import (
	"github.com/ishidawataru/sctp"
	"my5G-RANTester/internal/control_test_engine/nas_control"
	"my5G-RANTester/internal/control_test_engine/nas_control/mm_5gs"
	"my5G-RANTester/internal/control_test_engine/ngap_control/nas_transport"
	"my5G-RANTester/internal/control_test_engine/ngap_control/pdu_session_management"
	"my5G-RANTester/internal/control_test_engine/ngap_control/ue_context_management"
	"my5G-RANTester/internal/logging"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/openapi/models"
	"time"
)

func RegistrationUE(connN2 *sctp.SCTPConn, imsi string, ranUeId int64, ranIpAddr string, key string, opc string, amf string) (string, error) {

	// instance new ue.
	ue := &nas_control.RanUeContext{}

	// make initial UE message.
	err := nas_transport.InitialUEMessage(connN2, ue, imsi, ranUeId, key, opc, amf)
	if logging.Check_error(err, "send Initial Ue Message") {
		return ue.Supi, err
	}

	// receive NAS Authentication Request Msg
	ngapMsg, err := nas_transport.DownlinkNasTransport(connN2)
	if logging.Check_error(err, "receive DownlinkNasTransport/authentication request") {
		return ue.Supi, err
	}

	// send NAS Authentication Response
	pdu, err := mm_5gs.AuthenticationResponse(ue, ngapMsg)
	if logging.Check_error(err, "Authentication response worked fine") {
		return ue.Supi, err
	}

	err = nas_transport.UplinkNasTransport(connN2, ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	if logging.Check_error(err, "send UplinkNasTransport/Authentication Response") {
		return ue.Supi, err
	}

	// receive NAS Security Mode Command Msg
	_, err = nas_transport.DownlinkNasTransport(connN2)
	if logging.Check_error(err, "receive DownlinkNasTransport/Security Mode Command") {
		return ue.Supi, err
	}

	// send NAS Security Mode Complete Msg
	pdu, err = mm_5gs.SecurityModeComplete(ue)
	if logging.Check_error(err, "Security Mode Complete worked fine!") {
		return ue.Supi, err
	}
	err = nas_transport.UplinkNasTransport(connN2, ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	if logging.Check_error(err, "send UplinkNasTransport/Security Mode Complete Msg!") {
		return ue.Supi, err
	}

	// receive ngap Initial Context Setup Request Msg.
	_, err = nas_transport.DownlinkNasTransport(connN2)
	if logging.Check_error(err, "receive NGAP/Initial Context Setup Request") {
		return ue.Supi, err
	}

	// send ngap Initial Context Setup Response Msg
	err = ue_context_management.InitialContextSetupResponse(connN2, ue.AmfUeNgapId, ue.RanUeNgapId, ue.Supi)
	if logging.Check_error(err, "send NGAP/Initial context setup response message") {
		return ue.Supi, err
	}

	// send NAS Registration Complete Msg
	pdu, err = mm_5gs.RegistrationComplete(ue)
	if logging.Check_error(err, "NAS registration complete worked fine") {
		return ue.Supi, err
	}
	err = nas_transport.UplinkNasTransport(connN2, ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	if logging.Check_error(err, "send UplinkNasTransport/registration complete") {
		return ue.Supi, err
	}

	time.Sleep(100 * time.Millisecond)

	// called Single Network Slice Selection Assistance Information (S-NSSAI).
	sNssai := models.Snssai{
		Sst: 1, //The SST part of the S-NSSAI is mandatory and indicates the type of characteristics of the Network Slice.
		Sd:  "010203",
	}

	// send PduSessionEstablishmentRequest Msg
	pdu, err = mm_5gs.UlNasTransport(ue, uint8(ranUeId), nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)
	if logging.Check_error(err, "NAS UlNasTransport worked fine!") {
		return ue.Supi, err
	}

	err = nas_transport.UplinkNasTransport(connN2, ue.AmfUeNgapId, ue.RanUeNgapId, pdu)
	if logging.Check_error(err, "send UplinkNasTransport/Ul Nas Transport/PduSession Establishment request") {
		return ue.Supi, err
	}

	// receive 12. NGAP-PDU Session Resource Setup Request(DL nas transport((NAS msg-PDU session setup Accept)))
	_, err = nas_transport.DownlinkNasTransport(connN2)
	if logging.Check_error(err, "receive PDU Session Resource Setup Request/Dl Nas Transport/PDU establishment accept") {
		return ue.Supi, err
	}

	// send 14. NGAP-PDU Session Resource Setup Response.
	err = pdu_session_management.PDUSessionResourceSetupResponse(connN2, ue.AmfUeNgapId, ue.RanUeNgapId, ue.Supi, ranIpAddr)
	if logging.Check_error(err, "send PDU Session Resource Setup Response") {
		return ue.Supi, err
	}

	// time.Sleep(1 * time.Second)
	time.Sleep(100 * time.Millisecond)

	// function worked fine.
	return ue.Supi, nil
}