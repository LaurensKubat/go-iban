package iban

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// countrySettings contains length for IBAN and format for bban
type CountrySettings struct {
	// Length of IBAN code for this country
	Length int

	// Format of bban part of IBAN for this country
	Format string

	// Membership of country
	Sepa bool
}

// IBAN struct
type IBAN struct {
	// Full code
	code string

	// Full code prettyfied for printing on paper
	printCode string

	// Country code
	countryCode string

	// Check digits
	checkDigits string

	// Country settings
	countrySettings *CountrySettings

	// Country specific bban part
	bban string
}

/*
	Taken from http://www.tbg5-finance.org/ code example
*/
var countries = map[string]CountrySettings{
	"AD": CountrySettings{Length: 24, Format: "F04F04A12", 		Sepa: false},
	"AE": CountrySettings{Length: 23, Format: "F03F16", 		Sepa: false},
	"AL": CountrySettings{Length: 28, Format: "F08A16", 		Sepa: false}, //8!n16!c
	"AT": CountrySettings{Length: 20, Format: "F05F11", 		Sepa: true},
	"AZ": CountrySettings{Length: 28, Format: "U04A20", 		Sepa: false},
	"BA": CountrySettings{Length: 20, Format: "F03F03F08F02", 	Sepa: false},
	"BE": CountrySettings{Length: 16, Format: "F03F07F02", 		Sepa: true},
	"BG": CountrySettings{Length: 22, Format: "U04F04F02A08", 	Sepa: true},
	"BH": CountrySettings{Length: 22, Format: "U04A14", 		Sepa: false},
	"BR": CountrySettings{Length: 29, Format: "F08F05F10U01A01", Sepa: false},
	"CH": CountrySettings{Length: 21, Format: "F05A12", 		Sepa: true},
	"CR": CountrySettings{Length: 21, Format: "F03F14", 		Sepa: false},
	"CY": CountrySettings{Length: 28, Format: "F03F05A16", 		Sepa: false},
	"CZ": CountrySettings{Length: 24, Format: "F04F06F10", 		Sepa: true},
	"DE": CountrySettings{Length: 22, Format: "F08F10", 		Sepa: true},
	"DK": CountrySettings{Length: 18, Format: "F04F09F01", 		Sepa: true},
	"DO": CountrySettings{Length: 28, Format: "U04F20", 		Sepa: false},
	"EE": CountrySettings{Length: 20, Format: "F02F02F11F01", 	Sepa: true},
	"ES": CountrySettings{Length: 24, Format: "F04F04F01F01F10", Sepa: true},
	"FI": CountrySettings{Length: 18, Format: "F06F07F01", 		Sepa: true},
	"FO": CountrySettings{Length: 18, Format: "F04F09F01", 		Sepa: true},
	"FR": CountrySettings{Length: 27, Format: "F05F05A11F02", 	Sepa: true},
	"GB": CountrySettings{Length: 22, Format: "U04F06F08", 		Sepa:true},
	"GE": CountrySettings{Length: 22, Format: "U02F16", 		Sepa:false},
	"GI": CountrySettings{Length: 23, Format: "U04A15", 		Sepa:true},
	"GL": CountrySettings{Length: 18, Format: "F04F09F01", 		Sepa:true},
	"GR": CountrySettings{Length: 27, Format: "F03F04A16", 		Sepa:true},
	"GT": CountrySettings{Length: 28, Format: "A04A20", 		Sepa:false},
	"HR": CountrySettings{Length: 21, Format: "F07F10", 		Sepa:false},
	"HU": CountrySettings{Length: 28, Format: "F03F04F01F15F01", Sepa:true},
	"IE": CountrySettings{Length: 22, Format: "U04F06F08", 		Sepa:true},
	"IL": CountrySettings{Length: 23, Format: "F03F03F13", 		Sepa:false},
	"IS": CountrySettings{Length: 26, Format: "F04F02F06F10", 	Sepa:true},
	"IT": CountrySettings{Length: 27, Format: "U01F05F05A12", 	Sepa:true},
	"JO": CountrySettings{Length: 30, Format: "U04F04A18", 		Sepa:false},
	"KW": CountrySettings{Length: 30, Format: "U04A22", 		Sepa:false},
	"KZ": CountrySettings{Length: 20, Format: "F03A13", 		Sepa:false},
	"LB": CountrySettings{Length: 28, Format: "F04A20", 		Sepa:false},
	"LC": CountrySettings{Length: 32, Format: "U04A24", 		Sepa:false},
	"LI": CountrySettings{Length: 21, Format: "F05A12", 		Sepa:true},
	"LT": CountrySettings{Length: 20, Format: "F05F11", 		Sepa:true},
	"LU": CountrySettings{Length: 20, Format: "F03A13", 		Sepa:true},
	"LV": CountrySettings{Length: 21, Format: "U04A13", 		Sepa:true},
	"MC": CountrySettings{Length: 27, Format: "F05F05A11F02", 	Sepa:true},
	"MD": CountrySettings{Length: 24, Format: "A20", 			Sepa:false},
	"ME": CountrySettings{Length: 22, Format: "F03F13F02", 		Sepa:false},
	"MK": CountrySettings{Length: 19, Format: "F03A10F02", 		Sepa:false},
	"MR": CountrySettings{Length: 27, Format: "F05F05F11F02", 	Sepa:false},
	"MT": CountrySettings{Length: 31, Format: "U04F05A18", 		Sepa:true},
	"MU": CountrySettings{Length: 30, Format: "U04F02F02F12F03U03", Sepa:false},
	"NL": CountrySettings{Length: 18, Format: "U04F10", 		Sepa:true},
	"NO": CountrySettings{Length: 15, Format: "F04F06F01", 		Sepa:true},
	"PK": CountrySettings{Length: 24, Format: "U04A16", 		Sepa:false},
	"PL": CountrySettings{Length: 28, Format: "F08F16", 		Sepa:true},
	"PS": CountrySettings{Length: 29, Format: "U04A21", 		Sepa:false},
	"PT": CountrySettings{Length: 25, Format: "F04F04F11F02", 	Sepa:true},
	"QA": CountrySettings{Length: 29, Format: "U04A21", 		Sepa:false},
	"RO": CountrySettings{Length: 24, Format: "U04A16", 		Sepa:true},
	"RS": CountrySettings{Length: 22, Format: "F03F13F02", 		Sepa:false},
	"SA": CountrySettings{Length: 24, Format: "F02A18", 		Sepa:false},
	"SC": CountrySettings{Length: 31, Format: "U04F02F02F16U03", Sepa:false},
	"SE": CountrySettings{Length: 24, Format: "F03F16F01", 		Sepa:true},
	"SI": CountrySettings{Length: 19, Format: "F05F08F02", 		Sepa:true},
	"SK": CountrySettings{Length: 24, Format: "F04F06F10", 		Sepa:true},
	"SM": CountrySettings{Length: 27, Format: "U01F05F05A12", 	Sepa:true},
	"ST": CountrySettings{Length: 25, Format: "F08F11F02", 		Sepa:false},
	"TL": CountrySettings{Length: 23, Format: "F03F14F02", 		Sepa:false},
	"TN": CountrySettings{Length: 24, Format: "F02F03F13F02", 	Sepa:false},
	"TR": CountrySettings{Length: 26, Format: "F05A01A16", 		Sepa:false},
	"UA": CountrySettings{Length: 29, Format: "F06A19", 		Sepa:false},
	"VG": CountrySettings{Length: 24, Format: "U04F16", 		Sepa:false},
	"XK": CountrySettings{Length: 20, Format: "F04F10F02", 		Sepa:false},
}

func (i *IBAN)Validate() (error) {
	err1 := i.validateBban()
	err2 := i.validateCheckDigits()
	err := ""
	if err1 != nil {
		err = err + err1.Error()
	}
	if err2 != nil{
		err = err + err2.Error()
	}
	return errors.New(err)
}

func (i *IBAN)PrintCode() string {
	return i.printCode
}


func (i *IBAN)validateCheckDigits() error {
	// Move the four initial characters to the end of the string
	iban := i.code[4:] + i.code[:4]
	// Replace each letter in the string with two digits, thereby expanding the string, where A = 10, B = 11, ..., Z = 35
	mods := ""
	for _, c := range iban {
		// Get character code point value
		i := int(c)

		// Check if c is characters A-Z (codepoint 65 - 90)
		if i > 64 && i < 91 {
			// A=10, B=11 etc...
			i -= 55
			// Add int as string to mod string
			mods += strconv.Itoa(i)
		} else {
			mods += string(c)
		}
	}

	// Create bignum from mod string and perform module
	bigVal, success := new(big.Int).SetString(mods, 10)
	if !success {
		return errors.New("IBAN check digits validation failed")
	}

	modVal := new(big.Int).SetInt64(97)
	resVal := new(big.Int).Mod(bigVal, modVal)

	// Check if module is equal to 1
	if resVal.Int64() != 1 {
		return errors.New("IBAN has incorrect check digits")
	}

	return nil
}

func (i *IBAN) validateBban() error {
	bban := i.bban
	format := i.countrySettings.Format

	// Format regex to get parts
	frx, err := regexp.Compile(`[ABCFLUW]\d{2}`)
	if err != nil {
		return fmt.Errorf("Failed to validate bban: %v", err.Error())
	}

	// Get format part strings
	fps := frx.FindAllString(format, -1)

	// Create regex from format parts
	bbr := ""

	for _, ps := range fps {
		switch ps[:1] {
		case "F":
			bbr += "[0-9]"
		case "L":
			bbr += "[a-z]"
		case "U":
			bbr += "[A-Z]"
		case "A":
			bbr += "[0-9A-Za-z]"
		case "B":
			bbr += "[0-9A-Z]"
		case "C":
			bbr += "[A-Za-z]"
		case "W":
			bbr += "[0-9a-z]"
		}

		// Get repeat factor for group
		repeat, atoiErr := strconv.Atoi(ps[1:])
		if atoiErr != nil {
			return fmt.Errorf("Failed to validate bban: %v", atoiErr.Error())
		}

		// Add to regex
		bbr += fmt.Sprintf("{%d}", repeat)
	}

	// Compile regex and validate bban
	bbrx, err := regexp.Compile(bbr)
	if err != nil {
		return fmt.Errorf("Failed to validate bban: %v", err.Error())
	}

	if !bbrx.MatchString(bban) {
		return errors.New("bban part of IBAN is not formatted according to country specification")
	}

	return nil
}

// NewIBAN create new IBAN with validation
func NewIBAN(s string) (*IBAN, error) {
	iban := IBAN{}

	// Prepare string: remove spaces and convert to upper case
	s = strings.ToUpper(strings.Replace(s, " ", "", -1))
	iban.code = s

	// Validate characters
	r, err := regexp.Compile(`^[0-9A-Z]*$`)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate IBAN: %v", err.Error())
	}

	if !r.MatchString(s) {
		return nil, errors.New("IBAN can contain only alphanumeric characters")
	}

	// Get country code and check digits
	r, err = regexp.Compile(`^\D\D\d\d`)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate IBAN: %v", err.Error())
	}

	hs := r.FindString(s)
	if hs == "" {
		return nil, errors.New("IBAN must start with country code (2 characters) and check digits (2 digits)")
	}

	iban.countryCode = hs[0:2]
	iban.checkDigits = hs[2:4]

	// Get country settings for country code
	cs, ok := countries[iban.countryCode]
	if !ok {
		return nil, fmt.Errorf("Unsupported country code %v", iban.countryCode)
	}

	iban.countrySettings = &cs

	// Validate code length
	if len(s) != cs.Length {
		return nil, fmt.Errorf("IBAN length %d does not match length %d specified for country code %v", len(s), cs.Length, iban.countryCode)
	}

	// Set and validate bban part, the part after the language code and check digits
	iban.bban = s[4:]

	err = iban.validateBban()
	if err != nil {
		return nil, err
	}

	// Validate check digits with mod97
	err = iban.validateCheckDigits()
	if err != nil {
		return nil, err
	}

	// Generate print code from code (splits code in sections of 4 characters)
	prc := ""
	for len(s) > 4 {
		prc += s[:4] + " "
		s = s[4:]
	}

	iban.printCode = prc + s

	return &iban, nil
}
