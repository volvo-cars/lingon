# Security

## Vulnerability Scanning

### During Build

This project uses:
* [Go Vulnerability Database](https://go.dev/security/vuln/database)
* [OSV Vulnerability Database](https://osv.dev/)
* GitHub's [CodeQL](https://securitylab.github.com/tools/codeql) (hosted on GitHub)

to scan for vulnerabilities in the dependencies of this project.

### During Development

During a linting step is run to check for possible code vulnerabilities in this codebase
with the help of [gosec](https://github.com/securego/gosec)


## Secrets protection

This project uses GitHub security alerts containing a lists of supported secrets format 
to prevent fraudulent use of secrets that were committed accidentally.
For more information see [Secret scanning patterns](https://docs.github.com/en/code-security/secret-scanning/secret-scanning-patterns#supported-secrets-for-push-protection)

