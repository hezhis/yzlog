package yzlog

// An Option func configures a Logger.
type Option func(log *Logger)

// WithDevelopment sets the development mode.
func WithDevelopment(flag bool) Option {
	return func(log *Logger) {
		log.development = flag
	}
}

func WithDisableCaller(flag bool) Option {
	return func(log *Logger) {
		log.disableCaller = flag
	}
}
