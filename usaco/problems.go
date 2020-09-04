package usaco

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	// SampleTest maps sample input to sample output.
	SampleTest struct {
		Input  string
		Output string
	}

	// Problem data is parsed to this struct.
	Problem struct {
		Name        string
		Contest     string
		InpStream   string
		OutStream   string
		SampleTests []SampleTest
		Arg         Args
	}
)

// ProblemPage returns link to problem page
func (arg Args) ProblemPage() (link string) {
	link = fmt.Sprintf("%v/index.php?page=viewproblem2&cpid=%v",
		hostURL, arg.Cpid)
	return
}

// GetProblem parses problem details along with sample tests.
func (arg Args) GetProblem() (Problem, error) {
	if arg.Cpid == "" {
		return Problem{}, ErrInvalidSpecifier
	}

	link := arg.ProblemPage()
	page, err := loadPage(link)
	if err != nil {
		return Problem{}, err
	}
	defer page.Close()

	doc, _ := goquery.NewDocumentFromReader(
		strings.NewReader(page.MustElement("html").MustHTML()))

	// to hold problem data
	var prob Problem
	// extract contest name
	prob.Contest = clean(doc.Find("h2").Eq(0).Text())
	// extract problem name (exclude Problem Index)
	if true {
		name := clean(doc.Find("h2").Eq(1).Text())
		re := regexp.MustCompile(`Problem \d+\. (.*)`)
		match := re.FindStringSubmatch(name)
		prob.Name = clean(match[1])
	}
	// extract input and output stream
	if true {
		body, _ := doc.Html()

		reInp := regexp.MustCompile(`\(file (.*\.in)\)`)
		matchInp := reInp.FindStringSubmatch(body)
		if len(matchInp) > 1 {
			prob.InpStream = matchInp[1]
		}

		reOut := regexp.MustCompile(`\(file (.*\.out)\)`)
		matchOut := reOut.FindStringSubmatch(body)
		if len(matchOut) > 1 {
			prob.OutStream = matchOut[1]
		}
	}
	// extract sample tests
	if true {
		inp, out := doc.Find("pre.in"), doc.Find("pre.out")
		sampleTests := make([]SampleTest, 0)
		for i := 0; i < inp.Length() && i < out.Length(); i++ {
			inpStr, _ := inp.Eq(i).Html()
			outStr, _ := out.Eq(i).Html()
			sampleTests = append(sampleTests, SampleTest{
				Input:  clean(inpStr) + "\n",
				Output: clean(outStr) + "\n",
			})
		}

		if len(sampleTests) == 0 {
			body, _ := doc.Html()

			// the old format. Use regex and extract
			reInp := regexp.MustCompile("SAMPLE INPUT .*?:\n\n([\\s\\S]+?\n)\n")
			reOut := regexp.MustCompile("SAMPLE OUTPUT .*?:\n\n([\\s\\S]+?\n)\n")

			inp := reInp.FindStringSubmatch(body)
			out := reOut.FindStringSubmatch(body)

			for i := 1; i < len(inp) && i < len(out); i++ {
				sampleTests = append(sampleTests, SampleTest{
					Input:  clean(inp[i]) + "\n",
					Output: clean(out[i]) + "\n",
				})
			}
		}

		prob.SampleTests = sampleTests
	}
	// this might come handy in the future
	prob.Arg = arg
	return prob, nil
}

// SubmitSolution submits source code to specified problem.
// View languages.go for valid langName values.
// file is the submission file to upload on the form.
//
// If submission completes successfully, returns nil.
func (arg Args) SubmitSolution(langName, file string) error {
	// problem not specified
	if arg.Cpid == "" {
		return ErrInvalidSpecifier
	}
	// invalid language specified
	if _, ok := LanguageID[langName]; !ok {
		return fmt.Errorf("Invalid language name")
	}

	link := arg.ProblemPage()
	page, err := loadPage(link)
	if err != nil {
		return err
	}
	defer page.Close()

	if !page.MustHas(`select[name="language"]`) {
		return fmt.Errorf("submission not possible")
	}

	// find previous submission status id
	prevSid := *page.MustElement(`#last-status`).MustAttribute(`data-sid`)

	page.MustElement(`select[name="language"]`).MustSelect(langName)
	page.MustElement(`input[name="sourcefile"]`).MustSetFiles(file)
	page.MustElement(`input#solution-submit`).MustClick()
	page.MustWaitLoad()

	currSid := *page.MustElement(`#last-status`).MustAttribute(`data-sid`)
	if prevSid == currSid {
		return fmt.Errorf("failed to submit solution")

	}
	return nil
}
