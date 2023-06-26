// Package gtranslate provides the Gtranslate-backed translation service.
package gtranslate

//go:generate mockgen -source=gtranslate.go -destination=./mocks/gtranslate.go

import (
	"context"
	"fmt"
	"strings"

	"github.com/bregydoc/gtranslate"
)

// New returns a new Gtranslate translation service.
//
// Use WithClientOptions() to configure the *gtranslate.Client:
//
//	New("auth-key", WithClientOptions(gtranslate.BaseURL("https://example.com")))
//
// Use WithTranslateOptions() to append gtranslate.TranslateOptions to every request
// that is made through *Service.Translate():
//
//	New("auth-key", WithTranslateOptions(gtranslate.Formality(gtranslate.MoreFormal)))
func New(authKey string, opts ...Option) *Service {
	client := gtranslate.New(authKey)
	svc := NewWithClient(client, opts...)
	for _, opt := range svc.clientOpts {
		opt(client)
	}
	return svc
}

// NewWithClient does the same as New(), but accepts an existing *gtranslate.Client.
//
// WithClientOptions() has no effect in this case.
func NewWithClient(client Client, opts ...Option) *Service {
	svc := &Service{client: client}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

// Option is a service option.
type Option func(*Service)

// WithClientOptions configures the created *gtranslate.Client.
func WithClientOptions(opts ...gtranslate.ClientOption) Option {
	return func(svc *Service) {
		svc.clientOpts = append(svc.clientOpts, opts...)
	}
}

// WithTranslateOptions adds translation options to every request.
func WithTranslateOptions(opts ...gtranslate.TranslateOption) Option {
	return func(svc *Service) {
		svc.translateOpts = append(svc.translateOpts, opts...)
	}
}

// Client is an interface for *gtranslate.Client.
type Client interface {
	Translate(
		ctx context.Context,
		text string,
		targetLang gtranslate.Language,
		opts ...gtranslate.TranslateOption,
	) (string, gtranslate.Language, error)
}

// Client returns the underlying *gtranslate.Client.
func (svc *Service) Client() Client {
	return svc.client
}

// Service is the Gtranslate translation service.
//
// It delegates translation requests to the underlying *gtranslate.Client (https://github.com/bounoable/gtranslate).
//
// The gtranslate.SourceLang() and gtranslate.PreserveFormatting() options will be used automatically.
type Service struct {
	client        Client
	clientOpts    []gtranslate.ClientOption
	translateOpts []gtranslate.TranslateOption
}

// Translate translates the given text from sourceLang to targetLang.
func (svc *Service) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	opts := append([]gtranslate.TranslateOption{
		gtranslate.SourceLang(gtranslate.Language(strings.ToUpper(sourceLang))),
		gtranslate.PreserveFormatting(true),
		gtranslate.SplitSentences(gtranslate.SplitNoNewlines),
	}, svc.translateOpts...)

	translated, _, err := svc.client.Translate(ctx, text, gtranslate.Language(strings.ToUpper(targetLang)), opts...)
	if err != nil {
		return translated, fmt.Errorf("gtranslate: %w", err)
	}

	return translated, nil
}
