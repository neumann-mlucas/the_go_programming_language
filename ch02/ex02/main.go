package main

import (
	"fmt"
	"os"
	"strconv"
)

func argToLen(n float64) string {
	f := Feet(n)
	m := Meter(n)
	return fmt.Sprintf("%s = %s, %s = %s\n",
		f, FToM(f), m, MToF(m))
}

func argToTemp(n float64) string {
	c := Celsius(n)
	f := Fahrenheit(n)
	return fmt.Sprintf("%s = %s, %s = %s\n",
		c, CToF(c), f, FToC(f))
}

func argToWgth(n float64) string {
	k := Kilogram(n)
	p := Pound(n)
	return fmt.Sprintf("%s = %s, %s = %s\n",
		k, KGToP(k), p, PToKG(p))
}

func argToAllUnits(n float64) string {
	return argToLen(n) + argToTemp(n) + argToWgth(n)

}
func main() {
	for _, arg := range os.Args[1:] {
		n, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Printf(argToAllUnits(n))
	}
}
