package analysis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gowordladder/generator"
	"gowordladder/solving"
	"gowordladder/words"
	"image/color"
	"io"
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
	fCsv1, err := os.Create("analysis.csv")
	require.NoError(t, err)
	defer fCsv1.Close()
	fCsv2, err := os.Create("adjacents.csv")
	require.NoError(t, err)
	defer fCsv2.Close()

	printMdHeader(f, 1, "Analysis Report")

	// collect statistics...
	stats := make([]map[int]int, 0)
	islands := [15]int{}
	doublets := [15]int{}
	counts := [15]int{}
	longests := [15]int{}
	bigMax := 0
	bigMaxWordLen := 0
	totalWords := 0
	for wl := 2; wl <= 15; wl++ {
		d := words.NewDictionary(wl)
		wds := d.Words()
		counts[wl-1] = len(wds)
		totalWords += len(wds)
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

	printMdHeader(f, 3, "Word Statistics")
	printMdNotes(f, "", false,
		"Islands - are words that changing any letter will not form another word",
		"Doublets - are words that changing any letter will only form one other word",
		"LDS% (local decay smoothness) – the longest consecutive run of ladder lengths for which each word count is at least 95% of the previous ladder length’s count, expressed as a percentage of all ladder lengths in the dictionary",
		"Var – variance of consecutive count ratios, measuring how smoothly (low) or unevenly (high) word counts decay across ladder lengths",
		"Drop% – percentage decrease from ladder length 2 to 3, representing initial connectivity decay",
		"Longest - is the longest possible ladder length for the dictionary",
		"Numeric columns - are the number of words that can form that ladder length",
	)
	printMdTableHeaders(f, []string{"Letters", "Words", "Islands", "  % ", "Doublets", "  % ", "LDS%", " Var", "Drop%", "Longest"}, 2, bigMax, 6)
	printCsvHeaders(fCsv1, []string{"Letters", "Words", "Islands", "Doublets", "Longest"}, 2, bigMax)
	totalPerms := 0
	for wl := 2; wl <= 15; wl++ {
		m := stats[wl-2]
		assert.Equal(t, counts[wl-1]-islands[wl-1], m[2])
		assert.Equal(t, counts[wl-1]-islands[wl-1]-doublets[wl-1], m[3])
		plateau := plateauSize(m, bigMax, 0.95)
		ll := longests[wl-1] - 1
		lds := (float64(plateau) / float64(ll)) * 100
		fc, fi, fd := float64(counts[wl-1]), float64(islands[wl-1]), float64(doublets[wl-1])
		v := ratioVariance(m, bigMax)
		drop := 100.0 - ((float64(m[3]) / float64(m[2])) * 100)
		_, _ = fmt.Fprintf(f, "| %7d | %5d | %7d | %4.1f | %8d | %4.1f | %4.1f | %.2f | %5.1f | %7d |", wl, counts[wl-1], islands[wl-1], (fi/fc)*100, doublets[wl-1], (fd/fc)*100, lds, v, drop, longests[wl-1])
		_, _ = fmt.Fprintf(fCsv1, "%d,%d,%d,%d,%d", wl, counts[wl-1], islands[wl-1], doublets[wl-1], longests[wl-1])
		for i := 2; i <= bigMax; i++ {
			if n := m[i]; n == 0 {
				_, _ = fmt.Fprint(f, strings.Repeat(" ", 7)+"|")
				_, _ = fmt.Fprint(fCsv1, ",")
			} else {
				_, _ = fmt.Fprintf(f, "%6d |", n)
				_, _ = fmt.Fprintf(fCsv1, ",%d", n)
				if i > 2 {
					totalPerms += n
				}
			}
		}
		_, _ = fmt.Fprintf(f, "\n")
		_, _ = fmt.Fprint(fCsv1, "\n")
	}
	fmt.Printf("Words: %d, Perms: %d\n", totalWords, totalPerms)
	printMdNotes(f, "\nObservation notes:", true,
		"word - islands = ladder length 2 words",
		"word - islands - doublets = ladder length 3 words",
	)
	// chart...
	printMdChartLink(f, analysisChartFilename)
	analysisChart(stats, bigMax)

	// adjacents table...
	adjMax := 0
	adjacents := make([]map[int]int, 0)
	for wl := 2; wl <= 15; wl++ {
		d := words.NewDictionary(wl)
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
	printMdHeader(f, 3, "Adjacent Words Counts")
	printMdNotes(f, "This table shows the spread of adjacent word counts for each word in the dictionary.", false, "Words are considered adjacent if changing just one letter in one word forms the other word.")

	printMdTableHeaders(f, []string{"Letters"}, 0, adjMax, 6)
	printCsvHeaders(fCsv2, []string{"Letters"}, 0, adjMax)
	for wl := 2; wl <= 15; wl++ {
		m := adjacents[wl-2]
		_, _ = fmt.Fprintf(f, "| %7d |", wl)
		for i := 0; i <= adjMax; i++ {
			if n := m[i]; n == 0 {
				_, _ = fmt.Fprint(f, strings.Repeat(" ", 7)+"|")
			} else {
				_, _ = fmt.Fprintf(f, "%6d |", n)
			}
		}
		_, _ = fmt.Fprint(f, "\n")
		_, _ = fmt.Fprintf(fCsv2, "%d", wl)
		for i := 0; i <= adjMax; i++ {
			if n := m[i]; n == 0 {
				_, _ = fmt.Fprint(fCsv2, ",")
			} else {
				_, _ = fmt.Fprintf(fCsv2, ",%d", n)
			}
		}
		_, _ = fmt.Fprint(fCsv2, "\n")
	}
	// chart...
	printMdChartLink(f, adjacentsChartFilename)
	adjacentsChart(adjacents, adjMax, counts)

	// longest ladders...
	printMdHeader(f, 3, "Longest Ladders")
	printMdNotes(f, fmt.Sprintf("%d-letter words yields the longest ladders (%d)\n", bigMaxWordLen, bigMax), false)
	d := words.NewDictionary(bigMaxWordLen)
	wds := d.WordsWithSteps(bigMax)
	_, _ = fmt.Fprint(f, "|")
	for range len(wds) {
		_, _ = fmt.Fprint(f, strings.Repeat(" ", bigMaxWordLen+4)+"|")
	}
	_, _ = fmt.Fprint(f, "\n|")
	for range len(wds) {
		_, _ = fmt.Fprint(f, strings.Repeat("-", bigMaxWordLen+4)+"|")
	}
	solutions := make([]*solving.Solution, 0)
	alts := make([]int, 0)
	for _, wd := range wds {
		sw := wd.String()
		puzzle, err := generator.GeneratePuzzle(bigMaxWordLen, bigMax, &sw, nil)
		require.NoError(t, err)
		solutions = append(solutions, puzzle.Solutions[0])
		alts = append(alts, len(puzzle.Solutions)-1)
	}
	for l := 0; l < bigMax; l++ {
		_, _ = fmt.Fprint(f, "\n|")
		for _, s := range solutions {
			_, _ = fmt.Fprintf(f, " `%s` |", s.Ladder()[l])
		}
	}
	_, _ = fmt.Fprint(f, "\n|")
	for _, alt := range alts {
		_, _ = fmt.Fprintf(f, " %d alternatives |", alt)
	}
}

func printMdHeader(f *os.File, level int, header string) {
	if level > 1 {
		_, _ = fmt.Fprint(f, "\n")
	}
	_, _ = fmt.Fprint(f, strings.Repeat("#", level))
	_, _ = fmt.Fprint(f, " "+header+"\n\n")
}

func printMdNotes(f io.Writer, title string, numbered bool, notes ...string) {
	if title != "" {
		_, _ = fmt.Fprint(f, title+"\n")
	}
	for i, n := range notes {
		if numbered {
			_, _ = fmt.Fprintf(f, "%d. ", i+1)
		} else {
			_, _ = fmt.Fprint(f, "* ")
		}
		_, _ = fmt.Fprint(f, n+"\n")
	}
	if len(notes) > 0 {
		_, _ = fmt.Fprint(f, "\n")
	}
}

func printMdTableHeaders(f io.Writer, hdrs []string, numMin, numMax, numLen int) {
	_, _ = fmt.Fprint(f, "|")
	for _, h := range hdrs {
		_, _ = fmt.Fprint(f, " "+h+" |")
	}
	numFmt := "%" + strconv.Itoa(numLen) + "d |"
	for i := numMin; i <= numMax; i++ {
		_, _ = fmt.Fprintf(f, numFmt, i)
	}
	_, _ = fmt.Fprint(f, "\n|")
	for _, h := range hdrs {
		_, _ = fmt.Fprint(f, strings.Repeat("-", len(h)+1)+":|")
	}
	for i := numMin; i <= numMax; i++ {
		_, _ = fmt.Fprint(f, strings.Repeat("-", numLen)+":|")
	}
	_, _ = fmt.Fprint(f, "\n")
}

func printMdChartLink(f io.Writer, link string) {
	_, _ = fmt.Fprint(f, "\n![Chart]("+link+")\n")
}

func printCsvHeaders(f io.Writer, hdrs []string, numMin, numMax int) {
	for i, h := range hdrs {
		if i > 0 {
			_, _ = fmt.Fprint(f, ",")
		}
		_, _ = fmt.Fprint(f, h)
	}
	for i := numMin; i <= numMax; i++ {
		_, _ = fmt.Fprintf(f, ",%d", i)
	}
	_, _ = fmt.Fprint(f, "\n")
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
	p := newPlot("Word Statistics", "Ladder Lengths", "Words")

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
	for wl := 2; wl <= 15; wl++ {
		m := stats[wl-2]
		pts := make(plotter.XYs, 0)
		for i := 2; i <= bigMax; i++ {
			if n := m[i]; n != 0 {
				pts = append(pts, plotter.XY{X: float64(i), Y: float64(n)})
			}
		}
		if l, s, err := plotter.NewLinePoints(pts); err == nil {
			c, d := colorAndDashes(wl)
			l.Color = c
			l.Dashes = d
			s.Color = c
			s.Shape = nil
			p.Add(l, s)
			p.Legend.Add(fmt.Sprintf("%d-letter words", wl), l, s)
		} else {
			panic(err)
		}
	}
	if err := p.Save(800, 500+legendsHt, analysisChartFilename); err != nil {
		panic(err)
	}
}

func adjacentsChart(stats []map[int]int, bigMax int, counts [15]int) {
	p := newPlot("Adjacent Words Counts", "Number of Adjacent Words", "% of Words")

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
	for wl := 2; wl <= 15; wl++ {
		m := stats[wl-2]
		tot := float64(counts[wl-1])
		pts := make(plotter.XYs, 0)
		for i := 0; i <= bigMax; i++ {
			if n := m[i]; n != 0 {
				pts = append(pts, plotter.XY{X: float64(i), Y: (float64(n) / tot) * 100.0})
			}
		}
		if l, s, err := plotter.NewLinePoints(pts); err == nil {
			c, d := colorAndDashes(wl)
			l.Color = c
			l.Dashes = d
			s.Color = c
			s.Shape = nil
			p.Add(l, s)
			p.Legend.Add(fmt.Sprintf("%d-letter words", wl), l, s)
		} else {
			panic(err)
		}
	}
	if err := p.Save(800, 500+legendsHt, adjacentsChartFilename); err != nil {
		panic(err)
	}
}

func newPlot(title, xAxisLabel, yAxisLabel string) *plot.Plot {
	const fontVariant = "Sans"
	p := plot.New()
	p.Legend.TextStyle.Font.Variant = fontVariant
	p.Title.Text = title
	p.Title.TextStyle.Font.Variant = fontVariant
	p.Title.TextStyle.Font.Weight = 3 //font2.WeightBold

	p.Y.Label.Text = yAxisLabel
	p.Y.Label.TextStyle.Font.Variant = fontVariant
	p.Y.Tick.Label.Font.Variant = fontVariant

	p.X.Label.Text = xAxisLabel
	p.X.Label.TextStyle.Font.Variant = fontVariant
	p.X.Tick.Label.Font.Variant = fontVariant
	p.X.Tick.Label.Rotation = math.Pi / 2 // 90 degrees
	p.X.Tick.Label.XAlign = text.XRight
	p.X.Tick.Label.YAlign = text.YCenter
	p.X.Padding = 5
	return p
}

func colorAndDashes(wordLen int) (color color.RGBA, dashes []vg.Length) {
	i := wordLen - 2
	color = defaultColors[i%len(defaultColors)]
	if n := i / 5; n > 0 {
		dashes = defaultDashes[n-1]
	}
	return
}

var defaultColors = []color.RGBA{
	{31, 119, 180, 255}, // blue
	{255, 127, 14, 255}, // orange
	{44, 160, 44, 255},  // green
	{214, 39, 40, 255},  // red
	{23, 190, 207, 255}, // cyan
}

var defaultDashes = [][]vg.Length{
	{vg.Points(6), vg.Points(2)},
	{vg.Points(2), vg.Points(2)},
}
