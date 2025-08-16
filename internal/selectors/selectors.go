package selectors

import (
	"regexp"
	"strings"

	"github.com/sriharip316/tablo/internal/flatten"
)

type Expr struct {
	Raw   string
	parts []segment
}

type segment struct {
	pattern *regexp.Regexp // supports * and ? translated
	literal string         // fast-path for exact
}

func CompileMany(exprs []string) ([]Expr, error) {
	out := make([]Expr, 0, len(exprs))
	for _, e := range exprs {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		ex, err := compileOne(e)
		if err != nil {
			return nil, err
		}
		out = append(out, ex)
	}
	return out, nil
}

func compileOne(e string) (Expr, error) {
	segs := strings.Split(e, ".")
	parts := make([]segment, len(segs))
	for i, s := range segs {
		// translate globs to regex
		if strings.ContainsAny(s, "*?") {
			re := globToRegex(s)
			rgx, err := regexp.Compile("^" + re + "$")
			if err != nil {
				return Expr{}, err
			}
			parts[i] = segment{pattern: rgx}
		} else {
			parts[i] = segment{literal: s}
		}
	}
	return Expr{Raw: e, parts: parts}, nil
}

func globToRegex(s string) string {
	replacer := strings.NewReplacer(".", `\.`, "*", ".*", "?", ".")
	return replacer.Replace(s)
}

// ApplyToKeys filters/sorts keys per include/exclude expressions. If include is nil, keep all.
func ApplyToKeys(keys []string, include []Expr, exclude []Expr) []string {
	// preserve input order; when include provided, order by include-expr then input order
	added := map[string]struct{}{}
	out := make([]string, 0, len(keys))
	if len(include) == 0 {
		for _, k := range keys {
			// apply exclude filters on the fly
			if len(exclude) > 0 && matchesAny(k, exclude) {
				continue
			}
			if _, ok := added[k]; !ok {
				added[k] = struct{}{}
				out = append(out, k)
			}
		}
		return out
	}
	// include provided
	for _, inc := range include {
		for _, k := range keys {
			if _, ok := added[k]; ok {
				continue
			}
			if matchesAny(k, []Expr{inc}) {
				if len(exclude) > 0 && matchesAny(k, exclude) {
					continue
				}
				added[k] = struct{}{}
				out = append(out, k)
			}
		}
	}
	return out
}

func matchesAny(key string, exprs []Expr) bool {
	segs := strings.Split(key, ".")
	for _, ex := range exprs {
		if len(ex.parts) != len(segs) {
			continue
		}
		ok := true
		for i := range segs {
			p := ex.parts[i]
			if p.literal != "" {
				if p.literal != segs[i] {
					ok = false
					break
				}
			} else if p.pattern != nil {
				if !p.pattern.MatchString(segs[i]) {
					ok = false
					break
				}
			}
		}
		if ok {
			return true
		}
	}
	return false
}

// HeadersUnion returns the union of keys across rows in natural order of first occurrence.
func HeadersUnion(rows []flatten.FlatKV) []string {
	order := []string{}
	seen := map[string]struct{}{}
	for _, r := range rows {
		for k := range r {
			if _, ok := seen[k]; !ok {
				seen[k] = struct{}{}
				order = append(order, k)
			}
		}
	}
	return order
}

// MissingExpressions returns expressions with zero matches across provided keys.
func MissingExpressions(keys []string, include []Expr) []string {
	missing := []string{}
	for _, ex := range include {
		found := false
		for _, k := range keys {
			if matchesAny(k, []Expr{ex}) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, ex.Raw)
		}
	}
	return missing
}
