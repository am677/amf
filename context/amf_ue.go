package context

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"free5gc/lib/UeauCommon"
	"free5gc/lib/fsm"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/logger"
	"reflect"
	"regexp"
	"time"
)

type OnGoingProcedure string

const (
	OnGoingProcedureNothing    OnGoingProcedure = "Nothing"
	OnGoingProcedurePaging     OnGoingProcedure = "Paging"
	OnGoingProcedureN2Handover OnGoingProcedure = "N2Handover"
)

const (
	NgRanCgiPresentNRCGI    int32 = 0
	NgRanCgiPresentEUTRACGI int32 = 1
)

const (
	RecommendRanNodePresentRanNode int32 = 0
	RecommendRanNodePresentTAI     int32 = 1
)

type AmfUe struct {
	/* Gmm State */
	Sm map[models.AccessType]*fsm.FSM
	/* Registration procedure related context */
	RegistrationType5GS                uint8
	IdentityTypeUsedForRegistration    uint8
	RegistrationRequest                *nasMessage.RegistrationRequest
	ServingAmfChanged                  bool
	DeregistrationTargetAccessType     uint8 // only used when deregistration procedure is initialized by the network
	RegistrationAcceptForNon3GPPAccess []byte
	/* Used for AMF relocation */
	TargetAmfProfile *models.NfProfile
	TargetAmfUri     string
	/* Ue Identity*/
	PlmnId              models.PlmnId
	Suci                string
	Supi                string
	UnauthenticatedSupi bool
	Gpsi                string
	Pei                 string
	Tmsi                int32 // 5G-Tmsi
	Guti                string
	GroupID             string
	EBI                 int32
	IsCleartext         bool
	/* Ue Identity*/
	EventSubscriptionsInfo map[string]*AmfUeEventSubscription
	/* User Location*/
	RatType                  models.RatType
	Location                 models.UserLocation
	Tai                      models.Tai
	LocationChanged          bool
	LastVisitedRegisteredTai models.Tai
	TimeZone                 string
	/* context about udm */
	UdmId                             string
	NudmUECMUri                       string
	NudmSDMUri                        string
	ContextValid                      bool
	Reachability                      models.UeReachability
	SubscribedData                    models.SubscribedData
	SmfSelectionData                  *models.SmfSelectionSubscriptionData
	UeContextInSmfData                *models.UeContextInSmfData
	TraceData                         *models.TraceData
	UdmGroupId                        string
	SubscribedNssai                   []models.SubscribedSnssai
	AccessAndMobilitySubscriptionData *models.AccessAndMobilitySubscriptionData
	/* contex abut ausf */
	AusfGroupId                       string
	AusfId                            string
	AusfUri                           string
	RoutingIndicator                  string
	AuthenticationCtx                 *models.UeAuthenticationCtx
	AuthFailureCauseSynchFailureTimes int
	ABBA                              []uint8
	Kseaf                             string
	Kamf                              string
	/* context about PCF */
	PcfId                        string
	PcfUri                       string
	PolicyAssociationId          string
	AmPolicyUri                  string
	AmPolicyAssociation          *models.PolicyAssociation
	RequestTriggerLocationChange bool // true if AmPolicyAssociation.Trigger contains RequestTrigger_LOC_CH
	ConfigurationUpdateMessage   []byte
	/* UeContextForHandover*/
	HandoverNotifyUri string
	/* N1N2Message */
	N1N2MessageIDGenerator          int
	N1N2Message                     *N1N2Message
	N1N2MessageSubscribeIDGenerator int
	N1N2SubscriptionID              string
	N1N2MessageSubscribeInfo        map[string]*models.UeN1N2InfoSubscriptionCreateData
	/* Pdu Sesseion */
	StoredSmContext map[int32]*StoredSmContext // for DUPLICATE PDU Session ID
	SmContextList   map[int32]*SmContext
	/* Related Context*/
	RanUe map[models.AccessType]*RanUe
	/* other */
	OnGoing                       map[models.AccessType]*OnGoing
	UeRadioCapability             string // OCTET string
	Capability5GMM                nasType.Capability5GMM
	ConfigurationUpdateIndication nasType.ConfigurationUpdateIndication
	/* context related to Paging */
	UeRadioCapabilityForPaging                 *UERadioCapabilityForPaging
	InfoOnRecommendedCellsAndRanNodesForPaging *InfoOnRecommendedCellsAndRanNodesForPaging
	UESpecificDRX                              uint8
	/* Security Context */
	SecurityContextAvailable bool
	SecurityCapabilities     UeSecurityCapabilities
	NasUESecurityCapability  nasType.UESecurityCapability // for security command
	NgKsi                    models.NgKsi
	MacFailed                bool
	KnasInt                  []uint8 // 16 byte
	KnasEnc                  []uint8 // 16 byte
	Kgnb                     []uint8 // 32 byte
	Kn3iwf                   []uint8 // 32 byte
	NH                       []uint8 // 32 byte
	NCC                      uint8   // 0..7
	ULCountOverflow          uint16
	ULCountSQN               uint8
	DLCount                  uint32
	CipheringAlg             uint8
	IntegrityAlg             uint8
	/* Registration Area */
	RegistrationArea map[models.AccessType][]models.Tai
	LadnInfo         []LADN
	/* Network Slicing related context and Nssf */
	NssfId                            string
	NssfUri                           string
	NetworkSliceInfo                  *models.AuthorizedNetworkSliceInfo
	AllowedNssai                      map[models.AccessType][]models.AllowedSnssai
	ConfiguredNssai                   []models.ConfiguredSnssai
	NetworkSlicingSubscriptionChanged bool
	/* T3513(Paging) */
	T3513            *time.Timer // for paging
	PagingRetryTimes int
	LastPagingPkg    []byte
	/* T3565(Notification) */
	T3565                  *time.Timer // for NAS Notification
	NotificationRetryTimes int
	LastNotificationPkg    []byte
	/* T3560 (for authentication request/security mode command retransmission) */
	T3560           *time.Timer
	T3560RetryTimes int
	/* T3550 (for registration accept retransmission) */
	T3550           *time.Timer
	T3550RetryTimes int
	/* Ue Context Release Cause */
	ReleaseCause map[models.AccessType]*CauseAll
	/* T3502 (Assigned by AMF, and used by UE to initialize registration procedure) */
	T3502Value                      int // Second
	T3512Value                      int // default 54 min
	Non3gppDeregistrationTimerValue int // default 54 min
	/* T3522 (for deregistration request) */
	T3522           *time.Timer
	T3522RetryTimes int
}

type AmfUeEventSubscription struct {
	Timestamp         time.Time
	AnyUe             bool
	RemainReports     *int32
	EventSubscription *models.AmfEventSubscription
}
type N1N2Message struct {
	Request     models.N1N2MessageTransferRequest
	Status      models.N1N2MessageTransferCause
	ResourceUri string
}
type OnGoing struct {
	Procedure OnGoingProcedure
	Ppi       int32 //Paging priority
}

type SmContext struct {
	SmfId             string
	SmfUri            string
	PlmnId            models.PlmnId
	UserLocation      models.UserLocation
	PduSessionContext *models.PduSessionContext
}
type StoredSmContext struct {
	SmfId             string
	SmfUri            string
	PduSessionContext *models.PduSessionContext
	AnType            models.AccessType
	Payload           []byte
}

type UERadioCapabilityForPaging struct {
	NR    string // OCTET string
	EUTRA string // OCTET string
}

// TS 38.413 9.3.1.100
type InfoOnRecommendedCellsAndRanNodesForPaging struct {
	RecommendedCells    []RecommendedCell  // RecommendedCellsForPaging
	RecommendedRanNodes []RecommendRanNode // RecommendedRanNodesForPaging
}

// TS 38.413 9.3.1.71
type RecommendedCell struct {
	NgRanCGI         NGRANCGI
	TimeStayedInCell *int64
}

// TS 38.413 9.3.1.101
type RecommendRanNode struct {
	Present         int32
	GlobalRanNodeId *models.GlobalRanNodeId
	Tai             *models.Tai
}

type NGRANCGI struct {
	Present  int32
	NRCGI    *models.Ncgi
	EUTRACGI *models.Ecgi
}

// TS 38.413 9.3.1.86
type UeSecurityCapabilities struct {
	NREncryptionAlgorithms             [2]byte // 2 byte hex string, (NEA1 NEA2 NEA3 ....)
	NRIntegrityProtectionAlgorithms    [2]byte // 2 byte hex string, (NIA1 NIA2 NIA3 ....)
	EUTRAEncryptionAlgorithms          [2]byte // 2 byte hex string, (EEA1 EEA2 EEA3 ....)
	EUTRAIntegrityProtectionAlgorithms [2]byte // 2 byte hex string, (EIA1 EIA2 EIA3 ....)
}

func (ue *AmfUe) init() {
	ue.UnauthenticatedSupi = true
	ue.EventSubscriptionsInfo = make(map[string]*AmfUeEventSubscription)
	ue.Sm = make(map[models.AccessType]*fsm.FSM)
	ue.SmContextList = make(map[int32]*SmContext)
	ue.StoredSmContext = make(map[int32]*StoredSmContext)
	ue.RanUe = make(map[models.AccessType]*RanUe)
	ue.RegistrationArea = make(map[models.AccessType][]models.Tai)
	ue.AllowedNssai = make(map[models.AccessType][]models.AllowedSnssai)
	ue.N1N2MessageIDGenerator = 1
	ue.N1N2MessageSubscribeInfo = make(map[string]*models.UeN1N2InfoSubscriptionCreateData)
	ue.OnGoing = make(map[models.AccessType]*OnGoing)
	ue.OnGoing[models.AccessType_NON_3_GPP_ACCESS] = new(OnGoing)
	ue.OnGoing[models.AccessType_NON_3_GPP_ACCESS].Procedure = OnGoingProcedureNothing
	ue.OnGoing[models.AccessType__3_GPP_ACCESS] = new(OnGoing)
	ue.OnGoing[models.AccessType__3_GPP_ACCESS].Procedure = OnGoingProcedureNothing
	ue.ReleaseCause = make(map[models.AccessType]*CauseAll)
	ue.SecurityCapabilities.NREncryptionAlgorithms = [2]byte{0x00, 0x00}
	ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms = [2]byte{0x00, 0x00}
	ue.SecurityCapabilities.EUTRAEncryptionAlgorithms = [2]byte{0x00, 0x00}
	ue.SecurityCapabilities.EUTRAIntegrityProtectionAlgorithms = [2]byte{0x00, 0x00}
	ue.IsCleartext = true
}

func (ue *AmfUe) CmConnect(anType models.AccessType) bool {
	if _, ok := ue.RanUe[anType]; !ok {
		return false
	}
	return true
}

func (ue *AmfUe) CmIdle(anType models.AccessType) bool {
	return !ue.CmConnect(anType)
}

func (ue *AmfUe) Remove() {
	for _, ranUe := range ue.RanUe {
		if err := ranUe.Remove(); err != nil {
			logger.ContextLog.Errorf("Remove RanUe error: %v", err)
		}
	}
	if len(ue.Supi) > 0 {
		delete(AMF_Self().UePool, ue.Supi)
	}
	delete(AMF_Self().TmsiPool, ue.Tmsi)
	delete(AMF_Self().GutiPool, ue.Guti)
}

func (ue *AmfUe) DetachRanUe(anType models.AccessType) {
	delete(ue.RanUe, anType)
}

func (ue *AmfUe) AttachRanUe(ranUe *RanUe) {
	ue.RanUe[ranUe.Ran.AnType] = ranUe
	ranUe.AmfUe = ue
}

func (ue *AmfUe) GetAnType() models.AccessType {
	if ue.CmConnect(models.AccessType__3_GPP_ACCESS) {
		return models.AccessType__3_GPP_ACCESS
	} else if ue.CmConnect(models.AccessType_NON_3_GPP_ACCESS) {
		return models.AccessType_NON_3_GPP_ACCESS
	}
	return ""
}

func (ue *AmfUe) GetCmInfo() (cmInfos []models.CmInfo) {
	var cmInfo models.CmInfo
	cmInfo.AccessType = models.AccessType__3_GPP_ACCESS
	if ue.CmConnect(cmInfo.AccessType) {
		cmInfo.CmState = models.CmState_CONNECTED
	} else {
		cmInfo.CmState = models.CmState_IDLE
	}
	cmInfos = append(cmInfos, cmInfo)
	cmInfo.AccessType = models.AccessType_NON_3_GPP_ACCESS
	if ue.CmConnect(cmInfo.AccessType) {
		cmInfo.CmState = models.CmState_CONNECTED
	} else {
		cmInfo.CmState = models.CmState_IDLE
	}
	cmInfos = append(cmInfos, cmInfo)
	return
}

func (ue *AmfUe) InAllowedNssai(targetSNssai models.Snssai, anType models.AccessType) bool {
	for _, sNssai := range ue.AllowedNssai[anType] {
		if reflect.DeepEqual(sNssai, targetSNssai) {
			return true
		}
	}
	return false
}

func (ue *AmfUe) InSubscribedNssai(targetSNssai models.Snssai) bool {
	for _, sNssai := range ue.SubscribedNssai {
		if reflect.DeepEqual(sNssai.SubscribedSnssai, targetSNssai) {
			return true
		}
	}
	return false
}

func (ue *AmfUe) GetNsiInformationFromSnssai(anType models.AccessType, snssai models.Snssai) *models.NsiInformation {
	for _, allowedSnssai := range ue.AllowedNssai[anType] {
		if reflect.DeepEqual(*allowedSnssai.AllowedSnssai, snssai) {
			// TODO: select NsiInformation based on operator policy
			if len(allowedSnssai.NsiInformationList) != 0 {
				return &allowedSnssai.NsiInformationList[0]
			}
		}
	}
	return nil
}

func (ue *AmfUe) TaiListInRegistrationArea(taiList []models.Tai) bool {
	for _, tai := range taiList {
		if !InTaiList(tai, ue.RegistrationArea[ue.GetAnType()]) {
			return false
		}
	}
	return true
}

func (ue *AmfUe) HasWildCardSubscribedDNN() bool {
	for _, snssaiInfo := range ue.SmfSelectionData.SubscribedSnssaiInfos {
		for _, dnnInfo := range snssaiInfo.DnnInfos {
			if dnnInfo.Dnn == "*" {
				return true
			}
		}
	}
	return false
}

func (ue *AmfUe) SecurityContextIsValid() bool {
	return ue.SecurityContextAvailable && ue.NgKsi.Ksi != nasMessage.NasKeySetIdentifierNoKeyIsAvailable && !ue.MacFailed
}

// Kamf Derivation function defined in TS 33.501 Annex A.7
func (ue *AmfUe) DerivateKamf() {

	supiRegexp, _ := regexp.Compile("(?:imsi|supi)-([0-9]{5,15})")
	groups := supiRegexp.FindStringSubmatch(ue.Supi)
	if groups == nil {
		logger.NasLog.Errorln("supi is not correct")
		return
	}

	P0 := []byte(groups[1])
	L0 := UeauCommon.KDFLen(P0)
	P1 := ue.ABBA
	L1 := UeauCommon.KDFLen(P1)

	KseafDecode, _ := hex.DecodeString(ue.Kseaf)
	KamfBytes := UeauCommon.GetKDFValue(KseafDecode, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
	ue.Kamf = hex.EncodeToString(KamfBytes)
}

// Algorithm key Derivation function defined in TS 33.501 Annex A.9
func (ue *AmfUe) DerivateAlgKey() {

	// Security Key
	P0 := []byte{N_NAS_ENC_ALG}
	L0 := UeauCommon.KDFLen(P0)
	P1 := []byte{ue.CipheringAlg}
	L1 := UeauCommon.KDFLen(P1)

	KamfBytes, _ := hex.DecodeString(ue.Kamf)
	kenc := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	ue.KnasEnc = kenc[16:32]

	// Integrity Key
	P0 = []byte{N_NAS_INT_ALG}
	L0 = UeauCommon.KDFLen(P0)
	P1 = []byte{ue.IntegrityAlg}
	L1 = UeauCommon.KDFLen(P1)

	kint := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	ue.KnasInt = kint[16:32]
}

// Access Network key Derivation function defined in TS 33.501 Annex A.9
func (ue *AmfUe) DerivateAnKey(anType models.AccessType) {

	accessType := ACCESS_TYPE_3GPP // Defalut 3gpp
	P0 := ue.GetSecurityULCount()
	L0 := UeauCommon.KDFLen(P0)
	if anType == models.AccessType_NON_3_GPP_ACCESS {
		accessType = ACCESS_TYPE_NON_3GPP
	}
	P1 := []byte{accessType}
	L1 := UeauCommon.KDFLen(P1)

	KamfBytes, _ := hex.DecodeString(ue.Kamf)
	key := UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_KGNB_KN3IWF_DERIVATION, P0, L0, P1, L1)
	switch accessType {
	case ACCESS_TYPE_3GPP:
		ue.Kgnb = key
	case ACCESS_TYPE_NON_3GPP:
		ue.Kn3iwf = key
	}
}

// NH Derivation function defined in TS 33.501 Annex A.10
func (ue *AmfUe) DerivateNH(syncInput []byte) {

	P0 := syncInput
	L0 := UeauCommon.KDFLen(P0)

	KamfBytes, _ := hex.DecodeString(ue.Kamf)
	ue.NH = UeauCommon.GetKDFValue(KamfBytes, UeauCommon.FC_FOR_NH_DERIVATION, P0, L0)
}

func (ue *AmfUe) GetSecurityULCount() []byte {
	return GetSecurityCount(ue.ULCountOverflow, ue.ULCountSQN)
}
func (ue *AmfUe) GetSecurityDLCount() []byte {
	var r = make([]byte, 4)
	binary.BigEndian.PutUint32(r, ue.DLCount&0xffffff)
	return r
}

func (ue *AmfUe) UpdateSecurityContext(anType models.AccessType) {
	ue.DerivateAnKey(anType)
	switch anType {
	case models.AccessType__3_GPP_ACCESS:
		ue.DerivateNH(ue.Kgnb)
	case models.AccessType_NON_3_GPP_ACCESS:
		ue.DerivateNH(ue.Kn3iwf)
	}
	ue.NCC = 1
}

func (ue *AmfUe) UpdateNH() {
	ue.NCC++
	ue.DerivateNH(ue.NH)
}

func (ue *AmfUe) SelectSecurityAlg(intOrder, encOrder []uint8) {
	ue.CipheringAlg = ALG_CIPHERING_128_NEA0
	ue.IntegrityAlg = ALG_INTEGRITY_128_NIA0
	for _, intAlg := range intOrder {
		if intAlg == 0 && ue.NasUESecurityCapability.GetIA0_5G() == 1 {
			break
		}
		match := ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms[0] & intAlg
		if match > 0 {
			switch match {
			case 0x80:
				ue.IntegrityAlg = ALG_INTEGRITY_128_NIA1
			case 0x40:
				ue.IntegrityAlg = ALG_INTEGRITY_128_NIA2
			case 0x20:
				ue.IntegrityAlg = ALG_INTEGRITY_128_NIA3
			}
			break
		}
	}
	for _, encAlg := range encOrder {
		if encAlg == 0 && ue.NasUESecurityCapability.GetEA0_5G() == 1 {
			break
		}
		match := ue.SecurityCapabilities.NREncryptionAlgorithms[0] & encAlg
		if match > 0 {
			switch match {
			case 0x80:
				ue.CipheringAlg = ALG_CIPHERING_128_NEA1
			case 0x40:
				ue.CipheringAlg = ALG_CIPHERING_128_NEA2
			case 0x20:
				ue.CipheringAlg = ALG_CIPHERING_128_NEA3
			}
			break
		}
	}
}

func (ue *AmfUe) ClearRegistrationRequestData() {
	ue.RegistrationRequest = nil
	ue.RegistrationType5GS = 0
	ue.IdentityTypeUsedForRegistration = 0
	ue.ServingAmfChanged = false
	ue.RegistrationAcceptForNon3GPPAccess = nil
}

func (ue *AmfUe) RemoveAmPolicyAssociation() {
	ue.AmPolicyAssociation = nil
	ue.PolicyAssociationId = ""
}

func (ue *AmfUe) CopyDataFromUeContextModel(ueContext models.UeContext) {
	if ueContext.Supi != "" {
		ue.Supi = ueContext.Supi
		ue.UnauthenticatedSupi = ueContext.SupiUnauthInd
	}

	if ueContext.Pei != "" {
		ue.Pei = ueContext.Pei
	}

	if ueContext.UdmGroupId != "" {
		ue.UdmGroupId = ueContext.UdmGroupId
	}

	if ueContext.AusfGroupId != "" {
		ue.AusfGroupId = ueContext.AusfGroupId
	}

	if ueContext.RoutingIndicator != "" {
		ue.RoutingIndicator = ueContext.RoutingIndicator
	}

	if ueContext.SubUeAmbr != nil {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		if ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr == nil {
			ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr = new(models.AmbrRm)
		}

		subAmbr := ue.AccessAndMobilitySubscriptionData.SubscribedUeAmbr
		subAmbr.Uplink = ueContext.SubUeAmbr.Uplink
		subAmbr.Downlink = ueContext.SubUeAmbr.Downlink
	}

	if ueContext.SubRfsp != 0 {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.RfspIndex = ueContext.SubRfsp
	}

	if len(ueContext.RestrictedRatList) > 0 {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.RatRestrictions = ueContext.RestrictedRatList
	}

	if len(ueContext.ForbiddenAreaList) > 0 {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.ForbiddenAreas = ueContext.ForbiddenAreaList
	}

	if ueContext.ServiceAreaRestriction != nil {
		if ue.AccessAndMobilitySubscriptionData == nil {
			ue.AccessAndMobilitySubscriptionData = new(models.AccessAndMobilitySubscriptionData)
		}
		ue.AccessAndMobilitySubscriptionData.ServiceAreaRestriction = ueContext.ServiceAreaRestriction
	}

	if ueContext.SeafData != nil {
		seafData := ueContext.SeafData

		ue.NgKsi = *seafData.NgKsi
		if seafData.KeyAmf != nil {
			if seafData.KeyAmf.KeyType == models.KeyAmfType_KAMF {
				ue.Kamf = seafData.KeyAmf.KeyVal
			}
		}
		ue.NH, _ = hex.DecodeString(seafData.Nh)
		ue.NCC = uint8(seafData.Ncc)
	}

	if ueContext.PcfId != "" {
		ue.PcfId = ueContext.PcfId
	}

	if ueContext.PcfAmPolicyUri != "" {
		ue.AmPolicyUri = ueContext.PcfAmPolicyUri
	}

	if len(ueContext.AmPolicyReqTriggerList) > 0 {
		if ue.AmPolicyAssociation == nil {
			ue.AmPolicyAssociation = new(models.PolicyAssociation)
		}
		for _, trigger := range ueContext.AmPolicyReqTriggerList {
			switch trigger {
			case models.AmPolicyReqTrigger_LOCATION_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_LOC_CH)
			case models.AmPolicyReqTrigger_PRA_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_PRA_CH)
			case models.AmPolicyReqTrigger_SARI_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_SERV_AREA_CH)
			case models.AmPolicyReqTrigger_RFSP_INDEX_CHANGE:
				ue.AmPolicyAssociation.Triggers = append(ue.AmPolicyAssociation.Triggers, models.RequestTrigger_RFSP_CH)
			}
		}
	}

	if len(ueContext.SessionContextList) > 0 {
		for _, pduSessionContext := range ueContext.SessionContextList {
			smContext := SmContext{
				PduSessionContext: &pduSessionContext,
			}
			ue.SmContextList[pduSessionContext.PduSessionId] = &smContext
		}
	}

	if len(ueContext.MmContextList) > 0 {
		for _, mmContext := range ueContext.MmContextList {
			if mmContext.AccessType == models.AccessType__3_GPP_ACCESS {
				if nasSecurityMode := mmContext.NasSecurityMode; nasSecurityMode != nil {
					switch nasSecurityMode.IntegrityAlgorithm {
					case models.IntegrityAlgorithm_NIA0:
						ue.IntegrityAlg = ALG_INTEGRITY_128_NIA0
					case models.IntegrityAlgorithm_NIA1:
						ue.IntegrityAlg = ALG_INTEGRITY_128_NIA1
					case models.IntegrityAlgorithm_NIA2:
						ue.IntegrityAlg = ALG_INTEGRITY_128_NIA2
					case models.IntegrityAlgorithm_NIA3:
						ue.IntegrityAlg = ALG_INTEGRITY_128_NIA3
					}

					switch nasSecurityMode.CipheringAlgorithm {
					case models.CipheringAlgorithm_NEA0:
						ue.CipheringAlg = ALG_CIPHERING_128_NEA0
					case models.CipheringAlgorithm_NEA1:
						ue.CipheringAlg = ALG_CIPHERING_128_NEA1
					case models.CipheringAlgorithm_NEA2:
						ue.CipheringAlg = ALG_CIPHERING_128_NEA2
					case models.CipheringAlgorithm_NEA3:
						ue.CipheringAlg = ALG_CIPHERING_128_NEA3
					}

					if mmContext.NasDownlinkCount != 0 {
						ue.DLCount = uint32(mmContext.NasDownlinkCount)
					}

					if mmContext.NasUplinkCount != 0 {
						ue.ULCountOverflow = uint16((mmContext.NasUplinkCount & 0x0ff0) >> 8)
						ue.ULCountSQN = uint8((mmContext.NasUplinkCount & 0x000f))
					}

					// TS 29.518 Table 6.1.6.3.2.1
					if mmContext.UeSecurityCapability != "" {
						// ue.SecurityCapabilities
						buf, _ := base64.StdEncoding.DecodeString(mmContext.UeSecurityCapability)
						ue.NasUESecurityCapability.Buffer = buf
						ue.NasUESecurityCapability.SetLen(uint8(len(buf)))
					}
				}
			}

			if mmContext.AllowedNssai != nil {
				for _, snssai := range mmContext.AllowedNssai {
					allowedSnssai := models.AllowedSnssai{
						AllowedSnssai: &snssai,
					}
					ue.AllowedNssai[mmContext.AccessType] = append(ue.AllowedNssai[mmContext.AccessType], allowedSnssai)
				}
			}
		}
	}
	if ueContext.TraceData != nil {
		ue.TraceData = ueContext.TraceData
	}
}
