package lang

import (
	"backend/internal/infrastructure/lang/eng"
	"backend/internal/infrastructure/lang/ru"
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type Language interface {
	GetAttribute(string) string
	GetMessages() map[string]string
	GetResponseMessage(key string) string
}

type Registry struct {
	languages   map[string]Language
	defaultLang string
}

func NewRegistry() *Registry {
	return &Registry{
		languages: map[string]Language{
			"en": &eng.Package{},
			"ru": &ru.Package{},
		},
	}
}

// GetLanguage возвращает языковый пакет по коду
func (r *Registry) GetLanguage(langCode string) Language {
	if lang, exists := r.languages[langCode]; exists {
		return lang
	}
	return r.languages[r.defaultLang]
}

func (r *Registry) GetResponseMessage(ctx context.Context, key string) string {
	locale, ok := ctx.Value("locale").(string)
	if !ok {
		locale = "en"
	}
	lang := r.GetLanguage(locale)
	message := lang.GetResponseMessage(key)
	return message
}

func (r *Registry) Validate(ctx context.Context, err error) (map[string]string, error) {
	if err != nil {
		errs := err.(validator.ValidationErrors)
		humanReadableErrors, err := r.translateValidation(ctx, errs)
		if err != nil {
			slog.Error("Error localizing validation messages: " + err.Error())

			return nil, err
		}

		return humanReadableErrors, err
	}

	return nil, nil
}

// translateValidation переводит validation сообщение
func (r *Registry) translateValidation(
	ctx context.Context,
	errs validator.ValidationErrors,
) (map[string]string, error) {
	locale, ok := ctx.Value("locale").(string)
	if !ok {
		locale = "en"
	}
	if locale == "" {
		return nil, errors.New("locale is not set")
	}

	lang := r.GetLanguage(locale)
	validationMessages := lang.GetMessages()
	validatedMessages := make(map[string]string)

	for _, err := range errs {
		var res string
		res = strings.ReplaceAll(
			validationMessages[err.Tag()],
			"{field}",
			lang.GetAttribute(strcase.ToSnake(err.Field())),
		)
		if err.Param() != "" {
			res = strings.ReplaceAll(
				res,
				"{param}",
				lang.GetAttribute(strcase.ToSnake(err.Param())),
			)
		}
		validatedMessages[strcase.ToSnake(err.Field())] = res
	}
	return validatedMessages, nil
}
