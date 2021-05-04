package polkadotTransaction

const (
	//Balannce_Transfer      = "0400" //测试网
	Balannce_Transfer      = "0500" //主网
	Balannce_Transfer_name = "transfer"
	Default_Period         = 50
	SigningBitV4           = byte(0x84)
	SigningBitV4Westend    = byte(0x84)

	Compact_U32      = "Compact<u32>"
	AccounntIDFollow = false
	GenesisHash      = "b0a8d493285c2df73290dfb7e61f870f17b41801197a149ca93654499ea3dafe"
	SpecVersion      = 1058
	AddrPrefix       = 0x00
	TestAddrPrefix   = 42

	DOT_Balannce_Transfer       = "0500"
	DOT_AddrPrefix        uint8 = 0
	DOT_AccounntIDFollow        = false

	KSM_Balannce_Transfer       = "0400"
	KSM_AddrPrefix        uint8 = 2
	KSM_AccounntIDFollow        = false

	PLM_Balannce_Transfer       = "0303"
	PLM_AddrPrefix        uint8 = 5
	PLM_AccounntIDFollow        = true

	FIS_Balannce_Transfer       = "0603"
	FIS_AddrPrefix        uint8 = 20
	FIS_AccounntIDFollow        = true

	EDG_Balannce_Transfer       = "0603"
	EDG_AddrPrefix        uint8 = 7
	EDG_AccounntIDFollow        = false
)

const (
	modeBits                  = 2
	singleMode           byte = 0
	twoByteMode          byte = 1
	fourByteMode         byte = 2
	bigIntMode           byte = 3
	singleModeMaxValue        = 63
	twoByteModeMaxValue       = 16383
	fourByteModeMaxValue      = 1073741823
)

var modeToNumOfBytes = map[byte]uint{
	singleMode:   1,
	twoByteMode:  2,
	fourByteMode: 4,
}