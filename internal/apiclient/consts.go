package apiclient

const DDDNS_HOST = "dominikmatic.com"
const DDDNS_PORT = "53535"
const DEFAULT_TOKEN_PATH = "./token"

const MODE_UPDATE = "update"
const MODE_DELETE = "delete"
const DEFAULT_MODE = MODE_UPDATE

var ALLOWED_MODES = []string{
	MODE_UPDATE,
	MODE_DELETE,
}

var PUBLIC_IP_PROVIDERS = []string{
	"https://ifconfig.me",
	"https://icanhazip.com",
}
