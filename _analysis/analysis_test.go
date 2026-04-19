package analysis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gowordladder/generator"
	"gowordladder/solving"
	"gowordladder/words"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestAnalysisReport(t *testing.T) {
	f, err := os.Create("Analysis.md")
	require.NoError(t, err)
	defer f.Close()
	fmt.Fprint(f, "# Analysis Report\n\n")

	// collect statistics...
	stats := make([]map[int]int, 0)
	islands := [15]int{}
	doublets := [15]int{}
	counts := [15]int{}
	longests := [15]int{}
	bigMax := 0
	bigMaxWordLen := 0
	for wl := 2; wl <= 15; wl++ {
		d := words.LoadDictionary(wl)
		wds := d.Words()
		counts[wl-1] = len(wds)
		m := make(map[int]int)
		for _, w := range wds {
			if w.IsIsland() {
				islands[wl-1]++
			} else {
				mx := w.MaxSteps()
				for l := 2; l <= mx; l++ {
					m[l]++
				}
				if mx > bigMax {
					bigMax = mx
					bigMaxWordLen = wl
				}
				if w.IsDoublet() {
					doublets[wl-1]++
				}
				if mx > longests[wl-1] {
					longests[wl-1] = mx
				}
			}
		}
		stats = append(stats, m)
	}

	fmt.Fprint(f, "### Word statistics table\n\n")
	fmt.Fprint(f, "* Islands - are words that changing any letter will not form another word\n")
	fmt.Fprint(f, "* Doublets - are words that changing any letter will only form one other word\n")
	fmt.Fprint(f, "* LDS% (local decay smoothness) – the longest consecutive run of ladder lengths for which each word count is at least 95% of the previous ladder length’s count, expressed as a percentage of all ladder lengths in the dictionary\n")
	fmt.Fprint(f, "* Var – variance of consecutive count ratios, measuring how smoothly (low) or unevenly (high) word counts decay across ladder lengths\n")
	fmt.Fprint(f, "* Drop% – percentage decrease from ladder length 2 to 3, representing initial connectivity decay\n")
	fmt.Fprint(f, "* Longest - is the longest possible ladder length for the dictionary\n")
	fmt.Fprint(f, "* Numeric columns - are the number of words that can form that ladder length\n\n")
	fmt.Fprint(f, "| Letters | Words | Islands |   %  | Doublets |   %  | LDS% |  Var | Drop% | Longest |")
	for i := 2; i <= bigMax; i++ {
		fmt.Fprintf(f, "%6d |", i)
	}
	fmt.Fprint(f, "\n")
	fmt.Fprint(f, "|--------:|------:|--------:|-----:|---------:|-----:|-----:|-----:|------:|--------:|")
	for i := 2; i <= bigMax; i++ {
		fmt.Fprint(f, strings.Repeat("-", 6)+":|")
	}
	fmt.Fprint(f, "\n")
	for wl := 2; wl <= 15; wl++ {
		m := stats[wl-2]
		assert.Equal(t, counts[wl-1]-islands[wl-1], m[2])
		assert.Equal(t, counts[wl-1]-islands[wl-1]-doublets[wl-1], m[3])
		plateau := plateauSize(m, bigMax, 0.95)
		ll := longests[wl-1] - 1
		pp := (float64(plateau) / float64(ll)) * 100
		fc, fi, fd := float64(counts[wl-1]), float64(islands[wl-1]), float64(doublets[wl-1])
		v := ratioVariance(m, bigMax)
		drop := 100.0 - ((float64(m[3]) / float64(m[2])) * 100)
		fmt.Fprintf(f, "| %7d | %5d | %7d | %4.1f | %8d | %4.1f | %4.1f | %.2f | %5.1f | %7d |", wl, counts[wl-1], islands[wl-1], (fi/fc)*100, doublets[wl-1], (fd/fc)*100, pp, v, drop, longests[wl-1])
		for i := 2; i <= bigMax; i++ {
			if n := m[i]; n == 0 {
				fmt.Fprint(f, strings.Repeat(" ", 7)+"|")
			} else {
				fmt.Fprintf(f, "%6d |", n)
			}
		}
		fmt.Fprintf(f, "\n")
	}
	fmt.Fprint(f, "\nObservation notes:\n")
	fmt.Fprint(f, "1. word - islands = ladder length 2 words\n")
	fmt.Fprint(f, "2. word - islands - doublets = ladder length 3 words\n\n")
	// chart...
	fmt.Fprint(f, "\n![Chart]("+analysisChartFilename+")")
	analysisChart(stats, bigMax)

	// adjacents table...
	adjMax := 0
	adjacents := make([]map[int]int, 0)
	for wl := 2; wl <= 15; wl++ {
		d := words.LoadDictionary(wl)
		m := map[int]int{}
		for _, w := range d.Words() {
			n := len(w.LinkedWords())
			m[n]++
			if n > adjMax {
				adjMax = n
			}
		}
		adjacents = append(adjacents, m)
		assert.Equal(t, islands[wl-1], m[0])
	}
	fmt.Fprint(f, "\n\n### Adjacent Counts Table\n\n")
	fmt.Fprint(f, "This table shows the spread of adjacent word counts for each word in the dictionary\n\n")
	fmt.Fprint(f, "| Letters |")
	for i := 0; i <= adjMax; i++ {
		fmt.Fprintf(f, "%6d |", i)
	}
	fmt.Fprint(f, "\n")
	fmt.Fprint(f, "|--------:|")
	for i := 0; i <= adjMax; i++ {
		fmt.Fprint(f, strings.Repeat("-", 6)+":|")
	}
	fmt.Fprint(f, "\n")
	for wl := 2; wl <= 15; wl++ {
		m := adjacents[wl-2]
		fmt.Fprintf(f, "| %7d |", wl)
		for i := 0; i <= adjMax; i++ {
			if n := m[i]; n == 0 {
				fmt.Fprint(f, strings.Repeat(" ", 7)+"|")
			} else {
				fmt.Fprintf(f, "%6d |", n)
			}
		}
		fmt.Fprint(f, "\n")
	}
	// chart...
	fmt.Fprint(f, "\n![Chart]("+adjacentsChartFilename+")")
	adjacentsChart(adjacents, adjMax, counts)

	// longest ladders...
	fmt.Fprint(f, "\n\n### Longest Ladders\n\n")
	fmt.Fprintf(f, "%d-letter words yields the longest ladders (%d)\n\n", bigMaxWordLen, bigMax)
	d := words.LoadDictionary(bigMaxWordLen)
	wds := d.WordsWithSteps(bigMax)
	fmt.Fprint(f, "|")
	for _ = range len(wds) {
		fmt.Fprint(f, strings.Repeat(" ", bigMaxWordLen+2)+"|")
	}
	fmt.Fprint(f, "\n|")
	for _ = range len(wds) {
		fmt.Fprint(f, strings.Repeat("-", bigMaxWordLen+2)+"|")
	}
	solutions := make([]solving.Solution, 0)
	alts := make([]int, 0)
	for _, wd := range wds {
		sw := wd.ActualWord()
		puzzle, err := generator.GeneratePuzzle(bigMaxWordLen, bigMax, &sw, nil)
		require.NoError(t, err)
		solutions = append(solutions, puzzle.Solutions[0])
		alts = append(alts, len(puzzle.Solutions)-1)
	}
	for l := 0; l < bigMax; l++ {
		fmt.Fprint(f, "\n|")
		for _, s := range solutions {
			fmt.Fprintf(f, " `%s` |", s.Ladder()[l].ActualWord())
		}
	}
	fmt.Fprint(f, "\n|")
	for _, alt := range alts {
		fmt.Fprintf(f, " %d alternatives |", alt)
	}
}

func plateauSize(stats map[int]int, bigMax int, threshold float64) int {
	counts := make([]int, 0, bigMax-1)
	for i := 2; i <= bigMax; i++ {
		if n := stats[i]; n != 0 {
			counts = append(counts, n)
		}
	}
	if len(counts) < 2 {
		return 0
	}
	best := 0
	current := 0
	for i := 1; i < len(counts); i++ {
		prev := float64(counts[i-1])
		curr := float64(counts[i])
		if curr/prev >= threshold {
			if current == 0 {
				current = 2 // first matching pair starts a run of length 2
			} else {
				current++
			}
			if current > best {
				best = current
			}
		} else {
			current = 0
		}
	}
	return best
}

func ratioVariance(stats map[int]int, bigMax int) float64 {
	var ratios []float64
	var prev float64
	first := true
	for i := 2; i <= bigMax; i++ {
		n := stats[i]
		if n == 0 {
			continue
		}
		if first {
			prev = float64(n)
			first = false
			continue
		}
		curr := float64(n)
		ratios = append(ratios, curr/prev)
		prev = curr
	}
	if len(ratios) == 0 {
		return 0
	}
	// mean
	var sum float64
	for _, r := range ratios {
		sum += r
	}
	mean := sum / float64(len(ratios))
	// variance
	var varSum float64
	for _, r := range ratios {
		diff := r - mean
		varSum += diff * diff
	}
	return varSum / float64(len(ratios))
}

const (
	analysisChartFilename  = "analysis.png"
	adjacentsChartFilename = "adjacents.png"
)

func analysisChart(stats []map[int]int, bigMax int) {
	const (
		fontVariant = "Sans"
		title       = "Word Statistics"
		xAxisLabel  = "Ladder Lengths"
		yAxisLabel  = "Words"
	)
	p := plot.New()
	p.Legend.TextStyle.Font.Variant = fontVariant
	p.Title.Text = title
	p.Title.TextStyle.Font.Variant = fontVariant
	p.Title.TextStyle.Font.Weight = 3 //font2.WeightBold
	p.X.Label.Text = xAxisLabel
	p.X.Label.TextStyle.Font.Variant = fontVariant
	p.Y.Label.Text = yAxisLabel
	p.Y.Label.TextStyle.Font.Variant = fontVariant

	p.X.Tick.Label.Rotation = math.Pi / 2 //1.0472 // 60 degrees
	p.X.Tick.Label.XAlign = text.XRight
	p.X.Tick.Label.YAlign = text.YCenter
	p.X.Tick.Label.Font.Variant = fontVariant
	p.X.Padding = 5

	// make the legends appear above the plot area (by padding the title and shifting legends up)...
	legendsHt := vg.Length(14) * p.Legend.TextStyle.Height("X")
	p.Title.Padding = legendsHt
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.YOffs = legendsHt
	p.Legend.Padding = 3 // give some space between legend lines

	// x-axis...
	ticks := make([]plot.Tick, 0)
	for i := 2; i <= bigMax; i++ {
		lbl := strconv.Itoa(i) + "  "
		if i != 2 && i%5 != 0 {
			lbl = ""
		}
		ticks = append(ticks, plot.Tick{
			Value: float64(i),
			Label: lbl,
		})
	}
	p.X.Tick.Marker = plot.ConstantTicks(ticks)

	// plots...
	plots := make([]any, 0)
	for wl := 2; wl <= 15; wl++ {
		m := stats[wl-2]
		pts := make(plotter.XYs, 0)
		for i := 2; i <= bigMax; i++ {
			if n := m[i]; n != 0 {
				pts = append(pts, plotter.XY{X: float64(i), Y: float64(n)})
			}
		}
		plots = append(plots, fmt.Sprintf("%d-letter words", wl), pts)
	}
	err := plotutil.AddLinePoints(p, plots...)
	if err != nil {
		panic(err)
	}

	if err := p.Save(800, 500+legendsHt, analysisChartFilename); err != nil {
		panic(err)
	}
}

func adjacentsChart(stats []map[int]int, bigMax int, counts [15]int) {
	const (
		fontVariant = "Sans"
		title       = "Adjacent Words"
		xAxisLabel  = "Adjacent Counts"
		yAxisLabel  = "% Words"
	)
	p := plot.New()
	p.Legend.TextStyle.Font.Variant = fontVariant
	p.Title.Text = title
	p.Title.TextStyle.Font.Variant = fontVariant
	p.Title.TextStyle.Font.Weight = 3 //font2.WeightBold
	p.X.Label.Text = xAxisLabel
	p.X.Label.TextStyle.Font.Variant = fontVariant
	p.Y.Label.Text = yAxisLabel
	p.Y.Label.TextStyle.Font.Variant = fontVariant

	p.X.Tick.Label.Rotation = math.Pi / 2 //1.0472 // 60 degrees
	p.X.Tick.Label.XAlign = text.XRight
	p.X.Tick.Label.YAlign = text.YCenter
	p.X.Tick.Label.Font.Variant = fontVariant
	p.X.Padding = 5

	// make the legends appear above the plot area (by padding the title and shifting legends up)...
	legendsHt := vg.Length(14) * p.Legend.TextStyle.Height("X")
	p.Title.Padding = legendsHt
	p.Legend.Top = true
	p.Legend.Left = true
	p.Legend.YOffs = legendsHt
	p.Legend.Padding = 3 // give some space between legend lines

	// x-axis...
	ticks := make([]plot.Tick, 0)
	for i := 0; i <= bigMax; i++ {
		lbl := strconv.Itoa(i) + "  "
		if i%5 != 0 && i != bigMax {
			lbl = ""
		}
		ticks = append(ticks, plot.Tick{
			Value: float64(i),
			Label: lbl,
		})
	}
	p.X.Tick.Marker = plot.ConstantTicks(ticks)

	// plots...
	plots := make([]any, 0)
	for wl := 2; wl <= 15; wl++ {
		m := stats[wl-2]
		tot := float64(counts[wl-1])
		pts := make(plotter.XYs, 0)
		for i := 0; i <= bigMax; i++ {
			if n := m[i]; n != 0 {
				pts = append(pts, plotter.XY{X: float64(i), Y: (float64(n) / tot) * 100.0})
			}
		}
		plots = append(plots, fmt.Sprintf("%d-letter words", wl), pts)
	}
	err := plotutil.AddLinePoints(p, plots...)
	if err != nil {
		panic(err)
	}

	if err := p.Save(800, 500+legendsHt, adjacentsChartFilename); err != nil {
		panic(err)
	}
}
