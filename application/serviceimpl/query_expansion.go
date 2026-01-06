package serviceimpl

import (
	"strings"
	"unicode"
)

// Thai provinces list (77 provinces)
var thaiProvinces = map[string]bool{
	// Central Region
	"กรุงเทพ":       true,
	"กรุงเทพมหานคร": true,
	"bangkok":       true,
	"นนทบุรี":       true,
	"nonthaburi":    true,
	"ปทุมธานี":      true,
	"pathumthani":   true,
	"พระนครศรีอยุธยา": true,
	"อยุธยา":        true,
	"ayutthaya":     true,
	"อ่างทอง":       true,
	"angthong":      true,
	"ลพบุรี":        true,
	"lopburi":       true,
	"สิงห์บุรี":     true,
	"singburi":      true,
	"ชัยนาท":        true,
	"chainat":       true,
	"สระบุรี":       true,
	"saraburi":      true,
	"นครนายก":       true,
	"nakhonnayok":   true,
	"นครปฐม":        true,
	"nakhonpathom":  true,
	"สมุทรปราการ":   true,
	"samutprakan":   true,
	"สมุทรสาคร":     true,
	"samutsakhon":   true,
	"สมุทรสงคราม":   true,
	"samutsongkhram": true,

	// Eastern Region
	"ชลบุรี":     true,
	"chonburi":   true,
	"ระยอง":      true,
	"rayong":     true,
	"จันทบุรี":   true,
	"chanthaburi": true,
	"ตราด":       true,
	"trat":       true,
	"ฉะเชิงเทรา": true,
	"chachoengsao": true,
	"ปราจีนบุรี": true,
	"prachinburi": true,
	"สระแก้ว":    true,
	"sakaeo":     true,

	// Western Region
	"ราชบุรี":       true,
	"ratchaburi":    true,
	"กาญจนบุรี":     true,
	"kanchanaburi":  true,
	"สุพรรณบุรี":    true,
	"suphanburi":    true,
	"เพชรบุรี":      true,
	"phetchaburi":   true,
	"ประจวบคีรีขันธ์": true,
	"prachuapkhirikhan": true,

	// Northern Region
	"เชียงใหม่":   true,
	"chiangmai":   true,
	"เชียงราย":    true,
	"chiangrai":   true,
	"ลำปาง":       true,
	"lampang":     true,
	"ลำพูน":       true,
	"lamphun":     true,
	"แม่ฮ่องสอน":  true,
	"maehongson":  true,
	"น่าน":        true,
	"nan":         true,
	"พะเยา":       true,
	"phayao":      true,
	"แพร่":        true,
	"phrae":       true,
	"อุตรดิตถ์":   true,
	"uttaradit":   true,
	"ตาก":         true,
	"tak":         true,
	"สุโขทัย":     true,
	"sukhothai":   true,
	"พิษณุโลก":    true,
	"phitsanulok":  true,
	"พิจิตร":      true,
	"phichit":     true,
	"กำแพงเพชร":   true,
	"kamphaengphet": true,
	"เพชรบูรณ์":   true,
	"phetchabun":  true,
	"นครสวรรค์":   true,
	"nakhonsawan": true,
	"อุทัยธานี":   true,
	"uthaithani":  true,

	// Northeastern Region (Isan)
	"นครราชสีมา": true,
	"โคราช":      true,
	"nakhonratchasima": true,
	"korat":      true,
	"บุรีรัมย์":  true,
	"buriram":    true,
	"สุรินทร์":   true,
	"surin":      true,
	"ศรีสะเกษ":   true,
	"sisaket":    true,
	"อุบลราชธานี": true,
	"ubonratchathani": true,
	"ยโสธร":      true,
	"yasothon":   true,
	"ชัยภูมิ":    true,
	"chaiyaphum": true,
	"อำนาจเจริญ": true,
	"amnatcharoen": true,
	"หนองบัวลำภู": true,
	"nongbualamphu": true,
	"ขอนแก่น":    true,
	"khonkaen":   true,
	"อุดรธานี":   true,
	"udonthani":  true,
	"เลย":        true,
	"loei":       true,
	"หนองคาย":    true,
	"nongkhai":   true,
	"มหาสารคาม":  true,
	"mahasarakham": true,
	"ร้อยเอ็ด":   true,
	"roiet":      true,
	"กาฬสินธุ์":  true,
	"kalasin":    true,
	"สกลนคร":     true,
	"sakonnakhon": true,
	"นครพนม":     true,
	"nakhonphanom": true,
	"มุกดาหาร":   true,
	"mukdahan":   true,
	"บึงกาฬ":     true,
	"buengkan":   true,

	// Southern Region
	"นครศรีธรรมราช": true,
	"nakhonsithammarat": true,
	"กระบี่":       true,
	"krabi":        true,
	"พังงา":        true,
	"phangnga":     true,
	"ภูเก็ต":       true,
	"phuket":       true,
	"สุราษฎร์ธานี": true,
	"suratthani":   true,
	"ระนอง":        true,
	"ranong":       true,
	"ชุมพร":        true,
	"chumphon":     true,
	"สงขลา":        true,
	"songkhla":     true,
	"สตูล":         true,
	"satun":        true,
	"ตรัง":         true,
	"trang":        true,
	"พัทลุง":       true,
	"phatthalung":  true,
	"ปัตตานี":      true,
	"pattani":      true,
	"ยะลา":         true,
	"yala":         true,
	"นราธิวาส":     true,
	"narathiwat":   true,
}

// Tourism-related keywords that indicate no expansion needed
var tourismKeywords = []string{
	"ท่องเที่ยว",
	"เที่ยว",
	"สถานที่",
	"ที่เที่ยว",
	"แหล่งท่องเที่ยว",
	"travel",
	"tourism",
	"tourist",
	"attraction",
	"place",
	"visit",
	"วัด",
	"temple",
	"น้ำตก",
	"waterfall",
	"ภูเขา",
	"mountain",
	"ทะเล",
	"beach",
	"sea",
	"หาด",
	"อุทยาน",
	"park",
	"ตลาด",
	"market",
	"ร้านอาหาร",
	"restaurant",
	"โรงแรม",
	"hotel",
	"ที่พัก",
	"resort",
	"คาเฟ่",
	"cafe",
}

// ExpandSearchQuery expands a search query if it's just a province name
// This helps get better search results when users only type province names
// lang: "en" for English expansion, "th" or empty for Thai expansion
func ExpandSearchQuery(query string, lang string) string {
	if query == "" {
		return query
	}

	// Determine expansion suffix based on language
	expansionSuffix := " สถานที่ท่องเที่ยว"
	if lang == "en" {
		expansionSuffix = " tourist attractions"
	}

	// Normalize query for comparison
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	// Check if query already contains tourism-related keywords
	for _, keyword := range tourismKeywords {
		if strings.Contains(normalizedQuery, strings.ToLower(keyword)) {
			// Already has tourism context, no expansion needed
			return query
		}
	}

	// Check if query is just a province name (possibly with minor variations)
	queryWords := strings.Fields(normalizedQuery)

	// For single-word or two-word queries that match a province
	if len(queryWords) <= 2 {
		// Check exact match
		if thaiProvinces[normalizedQuery] {
			return query + expansionSuffix
		}

		// Check each word
		for _, word := range queryWords {
			if thaiProvinces[word] {
				return query + expansionSuffix
			}
		}

		// Check without spaces (for Thai text)
		noSpaces := strings.ReplaceAll(normalizedQuery, " ", "")
		if thaiProvinces[noSpaces] {
			return query + expansionSuffix
		}
	}

	// Check if query starts with "จังหวัด" (province prefix)
	if strings.HasPrefix(normalizedQuery, "จังหวัด") {
		provinceName := strings.TrimPrefix(normalizedQuery, "จังหวัด")
		provinceName = strings.TrimSpace(provinceName)
		if thaiProvinces[provinceName] || len(provinceName) > 0 {
			return query + expansionSuffix
		}
	}

	return query
}

// IsThaiProvince checks if the given text is a Thai province name
func IsThaiProvince(text string) bool {
	normalized := strings.ToLower(strings.TrimSpace(text))
	return thaiProvinces[normalized]
}

// ContainsThai checks if string contains Thai characters
func ContainsThai(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Thai, r) {
			return true
		}
	}
	return false
}
