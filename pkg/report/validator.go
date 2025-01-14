package report

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/wallarm/gotestwaf/internal/payload/encoder"
	"github.com/wallarm/gotestwaf/internal/payload/placeholder"
)

var (
	gtwVersionRegex = regexp.MustCompile(`(v\d+\.\d+\.\d+(\-\d+\-g[a-f0-9]{7})?|unknown)`)
	fpRegex         = regexp.MustCompile(`[a-f0-9]{32}`)
	markRegex       = regexp.MustCompile(`(N/A|[A-F][\+\-]?)`)
	suffixRegex     = regexp.MustCompile(`(na|[a-f])`)
)

func validateGtwVersion(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	result := gtwVersionRegex.MatchString(version)

	return result
}

func validateFp(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	result := fpRegex.MatchString(version)

	return result
}

func validateMark(fl validator.FieldLevel) bool {
	mark := fl.Field().String()
	result := markRegex.MatchString(mark)

	return result
}

func validateCssSuffix(fl validator.FieldLevel) bool {
	suffix := fl.Field().String()
	result := suffixRegex.MatchString(suffix)

	return result
}

func validateEncoders(fl validator.FieldLevel) bool {
	encoders := fl.Field().MapKeys()

	if _, ok := encoders[0].Interface().(string); !ok {
		return false
	}

	for _, e := range encoders {
		if _, ok := encoder.Encoders[e.String()]; !ok {
			return false
		}
	}

	return true
}

func validatePlaceholders(fl validator.FieldLevel) bool {
	placeholders := fl.Field().MapKeys()

	if _, ok := placeholders[0].Interface().(string); !ok {
		return false
	}

	for _, p := range placeholders {
		if _, ok := placeholder.Placeholders[p.String()]; !ok {
			return false
		}
	}

	return true
}

// ValidateReportData validates report data
func ValidateReportData(reportData *HtmlReport) error {
	validate := validator.New()
	validate.RegisterValidation("gtw_version", validateGtwVersion)
	validate.RegisterValidation("fp", validateFp)
	validate.RegisterValidation("mark", validateMark)
	validate.RegisterValidation("css_suffix", validateCssSuffix)
	validate.RegisterValidation("encoders", validateEncoders)
	validate.RegisterValidation("placeholders", validatePlaceholders)

	err := validate.Struct(reportData)
	if err != nil {
		return errors.Wrap(err, "found invalid values in the report data")
	}

	return nil
}
