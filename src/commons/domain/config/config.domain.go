package configDomain

type Config struct {
	BACKUP_KEEP        int
	BACKUP_SRC         string
	BACKUP_DST         string
	BACKUP_PREFiX_NAME string
	BACKUP_COMMAND     string

	SMTP_HOST string
	SMTP_PORT int
	SMTP_USER string
	SMTP_PASS string
	SMTP_FROM string

	AWS_SECRET_ACCESS_KEY string
	AWS_ACCESS_KEY_ID     string
	AWS_BUCKET            string
	AWS_REGION            string
	AWS_ENDPOINT          string
	AWS_URL_PUBLIC        string
}
