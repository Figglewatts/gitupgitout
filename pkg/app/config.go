package app

type Option func(*App)

func (app *App) configure(opts ...Option) {
	for _, opt := range opts {
		opt(app)
	}
}

func WithConcurrency(concurrency int) Option {
	return func(app *App) {
		app.concurrency = concurrency
	}
}
