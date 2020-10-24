# cpt-lib
[![codecov](https://codecov.io/gh/cp-tools/cpt-lib/branch/master/graph/badge.svg?token=VMMMOHWT1L)](undefined) [![GoDoc](https://godoc.org/github.com/cp-tools/cpt-lib?status.svg)](https://godoc.org/github.com/cp-tools/cpt-lib) [![Go Report Card](https://goreportcard.com/badge/github.com/cp-tools/cpt-lib)](https://goreportcard.com/report/github.com/cp-tools/cpt-lib) ![GitHub](https://img.shields.io/github/license/cp-tools/cpt-lib)

Short for competitive programming tools library, `cpt-lib` is a collection of API wrappers to request and upload data to various competitive programming websites, enabling the extraction of a myriad of data with ease.

<!--Or visit cpt-api for a command line interface-->

# Table of Contents

- [Overview](#overview)
- [Supported Websites](#supported-websites)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
- [FAQ](#faq)

# Overview

cpt-lib is a library that provides a dynamic API wrapper to perform various functions, using an instance of a headless browser. Some notable benefits (over other libraries and official API) include:

- Ability to fetch user exclusive data.
- No login credentials required (uses browser cookies)
- Optimised page fetching for faster performance.
- Ability to monitor dynamically updating data.



# Supported Websites

**Legend:** The number of symbols in **Support Status** signify the priority of development support.

| Website                              | Support Status                                               |                                                              |
| ------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| [CodeForces](https://codeforces.com) | :heavy_check_mark: :heavy_check_mark: :heavy_check_mark: Active support | [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/cp-tools/cpt-lib/Build%20and%20Test%20(codeforces)?label=Tests%20%28codeforces%29)](https://github.com/cp-tools/cpt-lib/actions) [![GoDoc](https://godoc.org/github.com/cp-tools/cpt-lib/codeforces?status.svg)](https://godoc.org/github.com/cp-tools/cpt-lib/codeforces) |
| [Atcoder](https://atcoder.jp)        | :white_check_mark: :white_check_mark: :white_check_mark: Active development |                                                              |
| [USACO](https://usaco.org)           | :white_check_mark:     Active development                    |                                                              |
| [Codechef](https://codechef.com)     | :black_medium_square: :black_medium_square:  ​Future milestone |                                                              |


# Getting Started

*For complete usage examples and documentation, view tests of the corresponding functions.*

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

The core module powering all website functions is headless browser automation, achieved using [rod](https://github.com/go-rod/rod).
Before fetching any data, the automated browser has to be initialised. This can be done easily using the `Start()` function provided:

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

*Note: Only Google Chrome/Chromium and Microsoft Edge is supported currently.*

Each function, at the base level, is a method to the type `Args`. This holds metadata of the contest and problem to extract data of, providing flexibility of usage. Converting appropriate links/descriptions can be done using the provided `Parse()` function.

```go
arg, err := codeforces.Parse("https://codeforces.com/contest/1234")
if err != nil{
	panic(err)
}

// Run methods on 'arg' next...
```

 

This specifier can then be used to perform a multitude of actions. For example, fetching sample test cases of the specified contest can be done easily, as given below:

```go
problems, err := arg.GetProblems()
if err != nil{
	panic(err)
}

// Display information of fetched problems.
for _, problem := range problems {
	fmt.Println("Problem Name:", problem.Name)
	fmt.Println("Number of sample tests:", len(problem.SampleTests))
}
```

 

# FAQ

### Which browsers are supported?

Since the project uses [go-rod/rod](https://github.com/go-rod/rod) as the headless browser controller, browsers that are supported by it are supported by this library. Nevertheless, is a gist of browser support:

- Chromium (and its derivative browsers) are supported and pass all tests.
- Due to issues in Firefox DevTools protocol implementation, **Firefox is currently not supported**. The related issue can be found [here](https://github.com/go-rod/rod/issues/193).
- Browsers without DevTools protocol are obviously not supported. Opera and Safari are two such examples.

- **Use unofficial forks of this project with extreme caution:** Maliciously modified code could steal sensitive information from your browser sessions. Therefore, it is recommended that you don't use any  unofficial version of this library.

### Are there any security issues?

No, the library doesn't access or modify any sensitive information, including login credentials.

However, that doesn't mean it isn't capable of doing the same. **Use unofficial versions of the project with extreme caution**, as it is very easy for malicious code to extract sensitive information from your browser cookies and logged in sessions.

### How do I use this library in other languages?

Currently, there is no support for the same.
However, we are working on a cross platform binary tool `cpt-api`, which will parse and return the specified information in `JSON` format. Expect it to be available very soon.

### What are the benefits of using browser automation over fetching web pages with GET requests?

After lot of consideration, here are the pros of the browser automation method:

- No hassle of login credential management. This makes it much more secure than the other method.
- Uses default browsers logged in session, doing away with reconfiguration (especially when you have multiple accounts :smirk:).
- Supports websites with single device login (looking at you - Codechef :smile:).
- Supports many more tasks (running custom JavaScript code on the page for example).
- Almost as fast as the previous method (see test execution details for proof).
- Dynamic data (delivered through web sockets) can be monitored and returned without the need for constantly fetching page to get results.

However, if you still wish to use the old method, you can use the last supported version of the same - `v1.5.1`. **Note that, there shall be no support and development for that version**.
