package settings

type UpdateSettingsDTO struct {
	TradeName      string `json:"trade_name,omitempty"`
	LegalName      string `json:"legal_name,omitempty"`
	CNPJ           string `json:"cnpj,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Email          string `json:"email,omitempty"`
	Street         string `json:"street,omitempty"`
	Number         string `json:"number,omitempty"`
	Complement     string `json:"complement,omitempty"`
	Neighborhood   string `json:"neighborhood,omitempty"`
	City           string `json:"city,omitempty"`
	State          string `json:"state,omitempty"`
	ZipCode        string `json:"zip_code,omitempty"`
	PrimaryColor   string `json:"primary_color,omitempty"`
	SecondaryColor string `json:"secondary_color,omitempty"`
	LogoURL        string `json:"logo_url,omitempty"`
}

type SettingsDTO struct {
	ID                string `json:"id"`
	TradeName         string `json:"trade_name"`
	LegalName         string `json:"legal_name"`
	CNPJ              string `json:"cnpj"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	Street            string `json:"street"`
	Number            string `json:"number"`
	Complement        string `json:"complement,omitempty"`
	Neighborhood      string `json:"neighborhood"`
	City              string `json:"city"`
	State             string `json:"state"`
	ZipCode           string `json:"zip_code"`
	PrimaryColor      string `json:"primary_color"`
	SecondaryColor    string `json:"secondary_color"`
	LogoURL           string `json:"logo_url,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
} 