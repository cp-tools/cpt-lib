# cpt-lib
[![codecov](https://codecov.io/gh/cp-tools/cpt-lib/branch/master/graph/badge.svg?token=VMMMOHWT1L)](undefined) [![GoDoc](https://godoc.org/github.com/cp-tools/cpt-lib?status.svg)](https://godoc.org/github.com/cp-tools/cpt-lib) [![Go Report Card](https://goreportcard.com/badge/github.com/cp-tools/cpt-lib)](https://goreportcard.com/report/github.com/cp-tools/cpt-lib) ![GitHub](https://img.shields.io/github/license/cp-tools/cpt-lib)

Short for competitive programming tools library, `cpt-lib` is a collection of API wrappers to request and upload data to various competitive programming websites, enabling the extraction and processing of a myriad of data with relative ease.

Make sure to star the project if you found it useful. :smile:

<!--Or visit cpt-api for a command line interface-->

# Table of Contents

- [Overview](#overview)
- [Supported Websites](#supported-websites)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
- [FAQ](#faq)

# Overview

cpt-lib is a library that uses browser automation, to extract, process and automate various processes related to competitive programming websites. Built as an API wrapper, it can perform many different tasks, of which notable ones are:

- Fetching sample tests of problems.
- Submitting solution to remote judge.
- Returning dynamic status of submissions.
- Extracting public details of contests.
- Fetching submissions and its solution code.

Obviously, some websites have more features, while some have only a subset of these. To know all available functions, refer to the respective package documentation.



# Supported Websites

| Website                              | Support                        | Status            |                                                              |
| ------------------------------------ | ------------------------------ | ----------------- | ------------------------------------------------------------ |
| [CodeForces](https://codeforces.com) | :star::star::star::star::star: | Is supported      | [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/cp-tools/cpt-lib/Build%20and%20Test%20(codeforces)?label=Tests%20%28codeforces%29)](https://github.com/cp-tools/cpt-lib/actions) [![GoDoc](https://godoc.org/github.com/cp-tools/cpt-lib/codeforces?status.svg)](https://godoc.org/github.com/cp-tools/cpt-lib/codeforces) |
| [Atcoder](https://atcoder.jp)        | :star::star::star::star:       | Under development |                                                              |
| [USACO](https://usaco.org)           | :star::star::star:             | Under development |                                                              |
| [Codechef](https://codechef.com)     | :star::star::star:             | Future milestone  |                                                              |

#### Legend

- **Support**:

  - :star::star::star::star::star: high priority and long term support.
  - :star::star::star::star: medium priority with long term support.
  - :star::star::star: medium priority with long term support.
  - :star::star: low priority with only bug fixes.
  - :star: low priority, development stalled.

- **Status**:

  - Is supported - Package for website available in the latest release.
  - Under development - Website support in develop or feature branch.
  - Future milestone - Development underway in the near future.




# Getting Started

*For complete usage examples, view tests of the corresponding packages.*
*Refer`godoc` badges above for corresponding documentation.*

## Installation

Usage is simple. First use `go get` to install the latest version of the library.

```go
go get -u github.com/cp-tools/cpt-lib
```

Next, include cpt-lib in your application.

```go
import "github.com/cp-tools/cpt-lib"
```

## Usage

*The examples below use the codeforces module, for illustration purposes.*

The core functionalities of the library are achieved using browser automation, through the DevTools protocol. The package [rod](https://github.com/go-rod/rod) is used to control the automated browser.

To use the methods provided by the library, the automated browser must be initiated first. This can be done easily using the function provided in all sub packages - `Start().`

```go
func main(){
	// Initialization parameters.
	inHeadless := true
    browser := "google-chrome"
    browserProfile := "/home/<username>/.config/google-chrome/"
	codeforces.Start(inHeadless, browserProfile, browser)
	
	// Do parsing here...
}
```



At the root, each package implements a `Args` type. This holds metadata of a contest/problem group, on which the methods are provided. Instantiating a variable of this type is done using the provided `Parse()` function, which casts the provided specifiers to the variable.

Specifiers supported by `Parse()` varies between websites, but URLs to the contest/problem are supported by all packages.

```go
arg, err := codeforces.Parse("https://codeforces.com/contest/1234/problem/c1")
if err != nil{
	panic(err)
}

// Run methods on 'arg' next...
```

 

The returned variable can then be used to execute the different provided methods, using its metadata.
An example of fetching sample tests of the problem (specified in the previous snippet) is as follows:

```go
problems, err := arg.GetProblems()
if err != nil{
	panic(err)
}

// Display information of fetched problems.
for _, problem := range problems {
	fmt.Println("Problem Name:", problem.Name)
    fmt.Println("Time limit:", problem.TimeLimit)
	fmt.Println("Number of sample tests:", len(problem.SampleTests))
}
```



# FAQ

### Which browsers are supported?

As the project directly uses [rod](https://github.com/go-rod/rod) to control the automated browser, all browsers supported by it are supported by this package. Nevertheless, here the list of browser support:

- **Supported browsers**:
  - Google chrome (tested)
  - Chromium (tested)
  - Microsoft Edge (untested)
- **Unsupported browsers:**
  - **Firefox** (see issue [here](https://github.com/go-rod/rod/issues/193))
  - Safari
  - Opera
  - Internet Explorer



### Is sensitive data of my browser at risk? 

Short answer, No. The library doesn't access or modify any sensitive information, including browser cookies and login credentials.

The functioning of the the browser automation is as follows.

- Starts the specified browser, with the specified user data directory.
- Creates another new browser instance, with a **different** user data directory.
- Copies cookies data from the first browser instance to the second browser instance.
- Closes the former browser, and uses the latter browser profile to access various data.

This ensures nothing (history, cookies, bookmarks etc) of your specified user profile are modified.



### How do I make this library work in other languages?

Currently, there is no official support for the same.

We are working on a cross platform binary tool of the library, which will return the parsed data in `JSON` encoded format. Expect a beta version of the same to be available in the near future. 



### What are the benefits of using browser automation over fetching web pages with GET requests?

Here are the major plus sides of using browser automation over source code fetching:

1. Uses logged in session of specified browser, doing away with management of logged in sessions.
2. Improved security since no credentials/login information is given away.
3. Support for websites with dynamic loading of data (done through JavaScript).
4. No major difference in speed between the two methods (view test results for stats).
5. Support for monitoring and returning web socket controlled data (think submission status).
6. Anything that can be done manually, can be automated with automated browsers.

And thus, the pros of this method clearly outweigh the cons, making this the best method.

However, if you persist on using the older, HTTP request based fetching method, you can use version 1 of the library. **Note that, there will not be any future development or support for this version**.