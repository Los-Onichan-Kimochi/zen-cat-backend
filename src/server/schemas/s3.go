package schemas

type S3Prefix string

var (
	LandingPrefix      S3Prefix = "landing"
	ProfessionalPrefix S3Prefix = "professional"
	LocalPrefix        S3Prefix = "local"
	UserPrefix         S3Prefix = "user"
	ServicePrefix      S3Prefix = "service"
	CommunityPrefix    S3Prefix = "community"
	TemplatePrefix     S3Prefix = "template"
)
