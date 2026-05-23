// Package redactor provides pattern-based redaction of sensitive data within
// log field values before they are written to output.
//
// A Redactor holds a set of compiled regular-expression rules. Each rule maps
// a pattern to a replacement string. Rules are applied in an unspecified order;
// if ordering matters, construct multiple Redactors and chain them manually.
//
// Example usage:
//
//	r, err := redactor.New(map[string]string{
//		`(?i)password=\S+`: "password=[REDACTED]",
//		`\b\d{16}\b`:       "[CARD]",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	sanitised := r.RedactMap(parsedFields)
//
// DefaultRules returns a convenience set of patterns covering common secrets
// such as passwords in query strings and credit-card numbers.
package redactor
