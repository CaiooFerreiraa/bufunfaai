package dto

type DeviceMetadata struct {
	DeviceName  string `json:"device_name" validate:"omitempty,max=120"`
	Platform    string `json:"platform" validate:"omitempty,max=50"`
	AppVersion  string `json:"app_version" validate:"omitempty,max=50"`
	Fingerprint string `json:"fingerprint" validate:"omitempty,max=255"`
}

type DeviceOutput struct {
	ID         string `json:"id"`
	DeviceName string `json:"device_name"`
	Platform   string `json:"platform"`
	AppVersion string `json:"app_version"`
	LastSeenAt string `json:"last_seen_at"`
	CreatedAt  string `json:"created_at"`
}
