# Contributing to cpt-lib
:confetti_ball::tada: Thank you for your interest in contributing to this project! :tada::confetti_ball:

There are some general rules and advices to keep in mind while making Pull Requests.

> The code is more what you call guidelines, than actual rules - *Cpt. Barbossa*

They are not intended to be strictly followed per se, so use your best judgement in making decisions.
Feel free to propose updates and changes to this document!



## Code of Conduct

Be kind, respectful and considerate towards the community. Take time to help others seeking advice. Any form of harassment will not be tolerated, and will be reported. Read the entire rules of conduct [here](CODE_OF_CONDUCT.md).



## Filling a bug report / feature request

1. Before creating a new issue, please check the existing issues to see if any similar one was already opened. Comment on existing ones, rather than creating duplicate issue reports.
2. If you think you've found a bug, please provide detailed steps of reproduction, the version of this library in use, the browser (and version) being automated, and any other useful metrics.
3. If you'd like to see a feature or enhancement, please open an issue with clear descriptions of what you'll like to have, and how its beneficial to the project.



## ELI5: How do I contribute?

First, you need to fork the repository, prior to submitting PRs. Then clone the fork to your computer:

```bash
git clone https://github.com/your_username/cpt-lib.git
cd cpt-lib
```

It is adviced to create a seperate feature branch (than making changes on the master branch):

```bash
git checkout -b my-feature-branch
```

Once you are done with the changes, add adequate tests (if applicable). Ensure each test covers a scenario of different cases, and try to reduce the overall number of tests, if it fetches data from the remote website (we don't want to be DDOS'ing the site, do we?).

If the new feature you are adding makes queries to the webiste, the tests **must** begin with a sleep timer of 10 seconds, to prevent unintended DDOS'ing. Use the following code to achieve the same:

```go
time.Sleep(time.Second * 10)
```


Once this is done, you will need to run your tests locally to check if your code works as expected. To do this, you will need to use **your personal account**, through which the tests will be conducted (all data parsed by your code must be public data, so that the testing account used doesn't run into content restriction).

In the corresponding website folder, creat a new `.env` file, and save your login credentials and other test details, in the following format (the .env files are untracked, preventing you from accidently committing your credentials):

```bash
# Keys should be all CAPS.

# The browser to use for testing.
BROWSER_BINARY=/usr/bin/google-chrome

# Login credentials. The prefix
# must be the folder name.
CODEFORCES_USERNAME=my-username
CODEFORCES_PASSWORD=my-password
```

Then you can run tests in one of the following ways:

- If you are using `VSCode`, use the CodeLens adornments (`Run Test | Debug Tests`) that appears at the top of the test function, to test your code.

- Otherwise, run the following command to test your code.

  ```bash
  # website refers to the website the changes are in.
  cd ./website
  # Test_MyFunction is the name of the test function.
  go test -v -run Test_MyFunction
  ```

The environment variables are automatically sourced by the tests, eliminating the need to source them prior to testing.

Once alls tests pass successfully, stage, commit and push changes using the commands:

```bash
git add .
git commit -m "Description of the changes"
git push origin my-feature-branch
```

Once the code is ready, create the Pull Request on GitHub and mark it for review. One of the maintainers will scrutinize changes in the commit, before running the testing workflow on them (disabled by default since login credentials passed can be misused by malicious code).

The reviewer(s) might suggest changes that should be done. Once satisfied, the PR will be merged, adding your name to the immortal Contributors Hall of Fame! :confetti_ball::confetti_ball: