package schemas

type S3Prefix string

var (
	LandingS3Prefix      S3Prefix = "landing"
	ProfessionalS3Prefix S3Prefix = "professional"
	LocalS3Prefix        S3Prefix = "local"
	UserS3Prefix         S3Prefix = "user"
	ServiceS3Prefix      S3Prefix = "service"
	CommunityS3Prefix    S3Prefix = "community"
	TemplateS3Prefix     S3Prefix = "template"
)
