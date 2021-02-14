package codeforces

import (
	"fmt"

	"github.com/go-rod/rod"
)

// CountdownPage returns link to countdown in contest.
func (arg Args) CountdownPage() (link string, err error) {
	if arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	switch arg.Class {
	case ClassGroup:
		if arg.Group == "" {
			return "", ErrInvalidSpecifier
		}

		link = fmt.Sprintf("%v/group/%v/contest/%v/countdown", hostURL, arg.Group, arg.Contest)

	case ClassContest, ClassGym:
		link = fmt.Sprintf("%v/%v/%v/countdown", hostURL, arg.Class, arg.Contest)

	default:
		return "", ErrInvalidSpecifier
	}

	return
}

// ContestsPage returns link to contests page of group/gym/contest.
func (arg Args) ContestsPage() (link string, err error) {

	switch arg.Class {
	case ClassGroup:
		if arg.Group == "" {
			return "", ErrInvalidSpecifier
		}

		// details of individual contest can't be parsed.
		// fallback to parsing all contests in group.
		link = fmt.Sprintf("%v/group/%v/contests?complete=true", hostURL, arg.Group)

	case ClassContest:
		if arg.Contest == "" {
			link = fmt.Sprintf("%v/contests?complete=true", hostURL)
			return
		}

		link = fmt.Sprintf("%v/contests/%v", hostURL, arg.Contest)

	case ClassGym:
		if arg.Contest == "" {
			link = fmt.Sprintf("%v/gyms?complete=true", hostURL)
			return
		}

		link = fmt.Sprintf("%v/contests/%v", hostURL, arg.Contest)

	default:
		return "", ErrInvalidSpecifier
	}

	return
}

// DashboardPage returns link to dashboard of contest.
func (arg Args) DashboardPage() (link string, err error) {
	if arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	switch arg.Class {
	case ClassGroup:
		if arg.Group == "" {
			return "", ErrInvalidSpecifier
		}

		link = fmt.Sprintf("%v/group/%v/contest/%v", hostURL, arg.Group, arg.Contest)

	case ClassContest, ClassGym:
		link = fmt.Sprintf("%v/%v/%v", hostURL, arg.Class, arg.Contest)

	default:
		return "", ErrInvalidSpecifier
	}

	return
}

// RegisterPage returns link to registration (not virtual contest registration)
// in contest.
func (arg Args) RegisterPage() (link string, err error) {
	if arg.Contest == "" || arg.Class != ClassContest {
		return "", ErrInvalidSpecifier
	}

	// gyms/groups don't support registration, do they!?
	link = fmt.Sprintf("%v/contestRegistration/%v",
		hostURL, arg.Contest)
	return
}

// ProblemsPage returns link to problem(s) page in contest.
func (arg Args) ProblemsPage() (link string, err error) {
	if arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	switch arg.Class {
	case ClassGroup:
		if arg.Group == "" {
			return "", ErrInvalidSpecifier
		}

		if arg.Problem == "" {
			link = fmt.Sprintf("%v/group/%v/contest/%v/problems", hostURL, arg.Group, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/group/%v/contest/%v/problem/%v", hostURL, arg.Group, arg.Contest, arg.Problem)
		}

	case ClassContest, ClassGym:
		if arg.Problem == "" {
			link = fmt.Sprintf("%v/%v/%v/problems", hostURL, arg.Class, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/%v/%v/problem/%v", hostURL, arg.Class, arg.Contest, arg.Problem)
		}

	default:
		return "", ErrInvalidSpecifier
	}

	return
}

// SubmissionsPage returns link to user submissions page.
func (arg Args) SubmissionsPage(handle string) (link string, err error) {
	// Contest not specified.
	if arg.Contest == "" {
		if handle == "" {
			// Extract handle from homepage.
			p, err := loadPage(hostURL)
			if err != nil {
				return "", ErrInvalidSpecifier
			}
			defer p.Close()
			p.WaitLoad()

			var elm *rod.Element
			if elm = p.MustElements(`#header a[href^="/profile/"]`).First(); elm == nil {
				return "", ErrInvalidSpecifier
			}

			handle = elm.MustText()
		}

		link = fmt.Sprintf("%v/submissions/%v", hostURL, handle)
		return
	}

	switch arg.Class {
	case ClassGroup:
		if handle != "" {
			// Fetching others submissions not possible.
			return "", ErrInvalidSpecifier
		}

		link = fmt.Sprintf("%v/group/%v/contest/%v/my", hostURL, arg.Group, arg.Contest)

	case ClassContest, ClassGym:
		if handle == "" {
			link = fmt.Sprintf("%v/%v/%v/my", hostURL, arg.Class, arg.Contest)
		} else {
			link = fmt.Sprintf("%v/submissions/%v/%v/%v", hostURL, handle, arg.Class, arg.Contest)
		}

	default:
		return "", ErrInvalidSpecifier
	}

	return
}

// SourceCodePage returns link to solution submission code.
func (sub Submission) SourceCodePage() (link string, err error) {
	if sub.ID == "" || sub.Arg.Contest == "" {
		return "", ErrInvalidSpecifier
	}

	switch sub.Arg.Class {
	case ClassGroup:
		link = fmt.Sprintf("%v/group/%v/contest/%v/submission/%v", hostURL, sub.Arg.Group, sub.Arg.Contest, sub.ID)

	case ClassContest, ClassGym:
		link = fmt.Sprintf("%v/%v/%v/submission/%v", hostURL, sub.Arg.Class, sub.Arg.Contest, sub.ID)

	default:
		return "", ErrInvalidSpecifier
	}

	return
}
