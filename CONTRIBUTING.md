# Contributing to PufferPanel
We welcome any contributions to PufferPanel, whether in the form of bug reports, feature suggestions, and even code submissions. When submitting any of those, please follow this documentation to help speed up the process.

# Reporting Bugs
A bug is when the software does not behave in a way that is expected. It is **not** invalid configurations which render the panel broken.

If you believe you have located a bug, please report it to the [Bug Tracker](https://github.com/PufferPanel/PufferPanel/issues).

**Please make sure there is not an issue for your specific bug already!** If you find that someone else has reported a bug you have, please comment on that issue stating you have replicated that bug. Do not make a new issue.

When submitting those bugs, follow these standards:
* The title of the issue should **clearly** and **quickly** explain the issue. A good title would be "Cannot delete IPs from node if it has 2 or more ports".
* The description should contain the following information
  * A complete description of the problem. This should explain what you expect the panel to do and what the panel actually did.
  * Steps to reproduce the bug. It is hard to figure out what the bug truly is if we cannot do it ourselves.
  * OS, NodeJS version (node -v), npm version (npm -v)

# Submitting feature requests
If you have an idea for a new feature or enhancement, please suggest it on our [Community Forum](https://community.pufferpanel.com/forum/5-feature-requests/).

# How to Contribute
**When submitting new code to PufferPanel you must follow all standards outlined in this document including**:
* All Pull Requests must contain a reference to an **existing** issue. If there is no issue for your PR to reference, then create a new issue, following the guidelines above. You must reference the issue in your PR in order for it to be accepted.
  * *There is a Caveat to this*, if your PR is adding a new feature please link it to the proper thread on the community forums. If there is no thread, create a new thread explaining the feature request in detail.
* Pull Requests may only contain **one (1)** feature or enhancement. Kitchen sinks will be thrown out the window.
* Pull Requests should have **one** commit attached to them. Please squash your commits before creating a Pull Request. You may squash commits after creating the request if you forget, it will automatically update the PR.
* We utilize ESLint to keep track of basic code formatting in this project. Please ensure that your IDE can follow these ESLint standards. Your code will be run through a test after the Pull Request is created, we will not pull it into the repository if it fails tests.
* All submitted code **must** include 100% test coverage in order to be accepted.
* All submitted code **must** pass the entire test suite.
* You **must** follow the [style guidelines](http://hapijs.com/styleguide) that are provided by HapiJS.
