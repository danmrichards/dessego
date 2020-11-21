package transport

// ResponseType indicates to Demon's Souls the type of response being returned.
type ResponseType int

var (
	ResponseLogin                  ResponseType = 0x02
	ResponseAddQWCData             ResponseType = 0x09
	ResponseCharacterTendency      ResponseType = 0x0e
	ResponseGetWanderingGhost      ResponseType = 0x11
	ResponseGeneric                ResponseType = 0x17
	ResponseAddBloodMsg            ResponseType = 0x1d
	ResponseReplayData             ResponseType = 0x1e
	ResponseListData               ResponseType = 0x1f
	ResponseUpdateMsgGrade         ResponseType = 0x2a
	ResponseTimeMsg                ResponseType = 0x22
	ResponseDeleteBloodMsg         ResponseType = 0x27
	ResponseCharacterMPGrade       ResponseType = 0x28
	ResponseCharacterBloodMsgGrade ResponseType = 0x29
)