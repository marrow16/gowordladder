package analysis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gowordladder/words"
	"slices"
	"strings"
	"testing"
)

func TestMaxDistances_Verify(t *testing.T) {
	for wl := 2; wl <= 15; wl++ {
		t.Run(fmt.Sprintf("%d-letter words", wl), func(t *testing.T) {
			d := words.NewDictionary(wl)
			wds := d.Words()
			slices.SortFunc(wds, func(a, b *words.Word) int {
				return strings.Compare(a.String(), b.String())
			})
			for _, w := range wds {
				aw := w.String()
				_ = aw
				wdm := words.NewWordDistanceMap(w, nil)
				mx := wdm.MaxDistance()
				assert.Equal(t, w.MaxSteps(), mx, "max distance %q = %d - should be %d", w, w.MaxSteps(), mx)
			}
		})
	}
}
