package config

var (
	TEMP  = "Temperature in 0.005 degree increments"
	VOLT  = "Battery voltage between approximately 1.6 - 3.6 V"
	HUM   = "Humidity between 0 - 100 % in 0.0025 % increments "
	PRES  = "Atmospheric pressure between 50000 - 115536 Pa in 1 Pa increments"
	ACCEL = "Acceleration on the device's Z axis up to 2 or 16 G"
	TX    = "Bluetooth transmit power between -40 - 22 dBm in 2 dBm increments"
	MOV   = "One byte rolling counter which is triggered by the device's movement"
	SEQ   = "Sequence counter between 0 - 65534 which is increased by one after each measurement"
)
