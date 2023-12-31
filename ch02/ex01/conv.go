package tempconv

// CToF converts a Celsius temperature to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC converts a Fahrenheit temperature to Celsius.
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * (5 / 9)) }

// KToC converts a Kelvin temperature to Celsius.
func KtoC(k Kelvin) Celsius { return Celsius(kelvin - DeltaKC) }

// CToK converts a Celsius temperature to Kelvin.
func CtoT(c Celsius) Kelvin { return Kelvin(c - DeltaCK) }

// FToK converts a Fahrenheit temperature to Kelvin.
func FToK(f Fahrenheit) Kelvin { return CToK(FToC(f)) }

// KToF converts a Kelvin temperature to Fahrenheit.
func KToF(k Kelvin) Fahrenheit { return CToF(KToC(k)) }
