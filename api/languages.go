package api

// Languages is the list of supported target languages for MVP.
var Languages = []Language{
	{Code: "es", Name: "Spanish"},
	{Code: "fr", Name: "French"},
	{Code: "de", Name: "German"},
	{Code: "ja", Name: "Japanese"},
	{Code: "it", Name: "Italian"},
	{Code: "pt", Name: "Portuguese"},
	{Code: "ko", Name: "Korean"},
	{Code: "zh", Name: "Chinese"},
}

// LanguageByCode returns a language by its code, or nil if not found.
func LanguageByCode(code string) *Language {
	for i := range Languages {
		if Languages[i].Code == code {
			return &Languages[i]
		}
	}
	return nil
}
