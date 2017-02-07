/*
Package wordexp exposes the libc function wordexp(3)
*/
package wordexp

/*
#include <stdlib.h>
#include <wordexp.h>

static size_t wordexp_wordcount(const wordexp_t *p)
{
	return p->we_wordc;
}

static char *wordexp_wordv(const wordexp_t *p, size_t n)
{
	return p->we_wordv[n];
}
*/
import "C"
import "fmt"

type wordexpError struct {
	val int
}

func (w *wordexpError) Error() string {
	switch w.val {
	case 0:
		return "not an error"
	case C.WRDE_BADCHAR:
		return "bad character"
	case C.WRDE_BADVAL:
		return "bad variable"
	case C.WRDE_CMDSUB:
		return "command execution not allowed"
	case C.WRDE_NOSPACE:
		return "not enough memory to store the result"
	case C.WRDE_SYNTAX:
		return "shell syntax error in words"
	default:
		return fmt.Sprintf("unrecognised error code")
	}
}

func newWordExpansionError(value int) error {
	return &wordexpError{val: value}
}

func wordexp(pattern string, allocator allocator) ([]string, error) {
	p, err := allocator.Alloc(C.sizeof_wordexp_t)
	if err != nil {
		return nil, newWordExpansionError(C.WRDE_NOSPACE)
	}
	defer allocator.Free(p)

	rc := C.wordexp(C.CString(pattern), p, C.WRDE_NOCMD|C.WRDE_UNDEF)

	if rc != 0 {
		return nil, newWordExpansionError(int(rc))
	}

	result := []string{}
	wordCount := C.wordexp_wordcount(p)

	for i := C.size_t(0); i < wordCount; i++ {
		result = append(result, C.GoString(C.wordexp_wordv(p, i)))
	}

	return result, nil
}

var _ wordExpander = (*foreignWordExpander)(nil)

// foreign to Go, because it is implemented in C.
type foreignWordExpander struct {
	allocator
}

func newWordExpander(a allocator) wordExpander {
	return &foreignWordExpander{allocator: a}
}

func (w *foreignWordExpander) Expand(pattern string) ([]string, error) {
	return wordexp(pattern, w.allocator)
}

type wordExpander interface {
	Expand(pattern string) ([]string, error)
}

// Expand performs posix-shell word expansion on pattern.
//
// This is a Cgo wrapper around wordexp(3) and explicitly passes
// WRDE_NOCMD|WRDE_UNDEF ("don't do command substitution", "treat
// undefined variables as errors") as flags to wordexp(3).
func Expand(pattern string) ([]string, error) {
	return newWordExpander(newNativeAllocator()).Expand(pattern)
}
