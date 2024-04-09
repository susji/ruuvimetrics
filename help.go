package main

var (
	HELP_TEMP  = "Temperature in 0.005 degree increments"
	HELP_VOLT  = "Battery voltage between approximately 1.6 - 3.6 V"
	HELP_HUM   = "Humidity between 0 - 100 % in 0.0025 % increments "
	HELP_PRES  = "Atmospheric pressure between 50000 - 115536 Pa in 1 Pa increments"
	HELP_ACCEL = "Acceleration on the device's Z axis up to 2 or 16 G"
	HELP_TX    = "Bluetooth transmit power between -40 - 22 dBm in 2 dBm increments"
	HELP_MOV   = "One byte rolling counter which is triggered by the device's movement"
	HELP_SEQ   = "Sequence counter between 0 - 65534 which is increased by one after each measurement"
)
