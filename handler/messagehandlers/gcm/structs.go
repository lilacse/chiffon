package gcm

type levelDetails struct {
	level      string
	cc         float64
	tapCount   int64
	holdCount  int64
	slideCount int64
	touchCount int64
	breakCount int64

	tapGreatLostRate float64
	tapGoodLostRate  float64
	tapMissLostRate  float64

	holdGreatLostRate float64
	holdGoodLostRate  float64
	holdMissLostRate  float64

	slideGreatLostRate float64
	slideGoodLostRate  float64
	slideMissLostRate  float64

	touchGreatLostRate float64
	touchGoodLostRate  float64
	touchMissLostRate  float64

	breakHighGreatLostRate float64
	breakMidGreatLostRate  float64
	breakLowGreatLostRate  float64
	breakGoodLostRate      float64
	breakMissLostRate      float64

	breakHighPerfectBonusLostRate float64
	breakLowPerfectBonusLostRate  float64
	breakGreatBonusLostRate       float64
	breakGoodBonusLostRate        float64
	breakMissBonusLostRate        float64
}
