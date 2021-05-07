package codeforces

var (
	// LanguageID represents all available languages with ids.
	LanguageID = map[string]string{
		"Microsoft Visual C++ 2010":        "2",
		"Delphi 7":                         "3",
		"Free Pascal 3.0.2":                "4",
		"PHP 7.2.13":                       "6",
		"Python 2.7.18":                    "7",
		"C# Mono 6.8":                      "9",
		"Haskell GHC 8.10.1":               "12",
		"Perl 5.20.1":                      "13",
		"ActiveTcl 8.5":                    "14",
		"Io-2008-01-07 (Win32)":            "15",
		"Pike 7.8":                         "17",
		"Befunge":                          "18",
		"OCaml 4.02.1":                     "19",
		"Scala 2.12.8":                     "20",
		"OpenCobol 1.0":                    "22",
		"Factor":                           "25",
		"Secret_171":                       "26",
		"Roco":                             "27",
		"D DMD32 v2.091.0":                 "28",
		"Python 3.9.1":                     "31",
		"Go 1.15.6":                        "32",
		"Ada GNAT 4":                       "33",
		"JavaScript V8 4.8.0":              "34",
		"Java 1.8.0_241":                   "36",
		"Mysterious Language":              "38",
		"FALSE":                            "39",
		"PyPy 2.7 (7.3.0)":                 "40",
		"PyPy 3.7 (7.3.0)":                 "41",
		"GNU G++11 5.1.0":                  "42",
		"GNU GCC C11 5.1.0":                "43",
		"Picat 0.9":                        "44",
		"GNU C++11 5 ZIP":                  "45",
		"Java 8 ZIP":                       "46",
		"J":                                "47",
		"Kotlin 1.4.0":                     "48",
		"Rust 1.49.0":                      "49",
		"GNU G++14 6.4.0":                  "50",
		"PascalABC.NET 3.4.2":              "51",
		"Clang++17 Diagnostics":            "52",
		"GNU G++17 7.3.0":                  "54",
		"Node.js 12.6.3":                   "55",
		"Microsoft Q#":                     "56",
		"Text":                             "57",
		"Microsoft Visual C++ 2017":        "59",
		"Java 11.0.6":                      "60",
		"GNU G++17 9.2.0 (64 bit, msys 2)": "61",
		"UnknownX":                         "62",
		"C# 8, .NET Core 3.1":              "65",
		"Ruby 3.0.0":                       "67",
		"Secret 2021":                      "68",
	}

	// LanguageExtn corresponds to file extension of
	// given language source code.
	LanguageExtn = map[string]string{
		"GNU C11":               ".c",
		"Clang++17 Diagnostics": ".cpp",
		"GNU C++0x":             ".cpp",
		"GNU C++":               ".cpp",
		"GNU C++11":             ".cpp",
		"GNU C++14":             ".cpp",
		"GNU C++17":             ".cpp",
		"MS C++":                ".cpp",
		"MS C++ 2017":           ".cpp",
		"GNU C++17 (64)":        ".cpp",
		"Mono C#":               ".cs",
		"D":                     ".d",
		"Go":                    ".go",
		"Haskell":               ".hs",
		"Kotlin":                ".kt",
		"Ocaml":                 ".ml",
		"Delphi":                ".pas",
		"FPC":                   ".pas",
		"PascalABC.NET":         ".pas",
		"Perl":                  ".pl",
		"PHP":                   ".php",
		"Python 2":              ".py",
		"Python 3":              ".py",
		"PyPy 2":                ".py",
		"PyPy 3":                ".py",
		"Ruby":                  ".rb",
		"Rust":                  ".rs",
		"JavaScript":            ".js",
		"Node.js":               ".js",
		"Q#":                    ".qs",
		"Java":                  ".java",
		"Java 6":                ".java",
		"Java 7":                ".java",
		"Java 8":                ".java",
		"Java 9":                ".java",
		"Java 10":               ".java",
		"Java 11":               ".java",
		"Tcl":                   ".tcl",
		"F#":                    ".fs",
		"Befunge":               ".bf",
		"Pike":                  ".pike",
		"Io":                    ".io",
		"Factor":                ".factor",
		"Cobol":                 ".cbl",
		"Secret_171":            ".secret_171",
		"Ada":                   ".adb",
		"FALSE":                 ".f",
		"":                      ".txt",
	}
)
