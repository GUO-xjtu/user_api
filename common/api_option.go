package common

type APIOptions struct {
	CheckUser       bool
	CheckAccountBan bool
	PureJSONResp    bool
}

func NewAPIOptions(options ...APIOption) *APIOptions {
	//set default values
	Opts := APIOptions{
		CheckUser:       true,
		CheckAccountBan: false,
		PureJSONResp:    false,
	}
	for _, opt := range options {
		opt.f(&Opts)
	}
	return &Opts
}

type APIOption struct {
	f func(*APIOptions)
}

func WithCheckUser(checkUser bool) APIOption {
	return APIOption{func(options *APIOptions) {
		options.CheckUser = checkUser
	}}
}

func WithCheckAccountBan(checkAccountBan bool) APIOption {
	return APIOption{func(options *APIOptions) {
		options.CheckAccountBan = checkAccountBan
	}}
}


func WithPureJSONResp(pureJSON bool) APIOption {
	return APIOption{func(options *APIOptions) {
		options.PureJSONResp = pureJSON
	}}
}

