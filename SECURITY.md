# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

We take the security of Gip seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please DO NOT

- Open a public GitHub issue for security vulnerabilities
- Discuss the vulnerability in public forums or social media

### Please DO

**Report security vulnerabilities privately** by:

1. **GitHub Security Advisories** (Recommended): Use GitHub's private vulnerability reporting feature
   - Go to the [Security tab](https://github.com/iamHrithikRaj/gip/security)
   - Click "Report a vulnerability"
   - This ensures private disclosure and proper tracking

2. **Private Issue**: Open a security issue on GitHub marked as confidential

### What to Include

Please provide as much information as possible:

- **Type of vulnerability** (e.g., code injection, privilege escalation)
- **Full paths** of affected source files
- **Location** of the affected code (tag/branch/commit)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact** of the vulnerability
- **Potential remediation** if you have suggestions

### What to Expect

- **Acknowledgment**: We will acknowledge receipt within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Regular Updates**: We will keep you informed of our progress
- **Fix Timeline**: We aim to release fixes within 30 days for critical issues
- **Credit**: We will credit you in the security advisory (unless you prefer to remain anonymous)

## Security Best Practices

When using Gip:

### For Users

1. **Keep Gip Updated**: Always use the latest version
   ```bash
   cargo install gip
   # Or download from: https://github.com/iamHrithikRaj/gip/releases/latest
   ```

2. **Review Manifests**: Manifests are stored in `.gip/` - ensure they don't contain sensitive data
   ```bash
   # Add to .gitignore if manifests contain sensitive info
   echo ".gip/" >> .gitignore
   ```

3. **Validate Input**: Be cautious when running Gip on untrusted repositories

4. **Permissions**: Gip runs with your user permissions - be aware of what it can access

### For Contributors

1. **No Hardcoded Secrets**: Never commit API keys, passwords, or tokens
2. **Input Validation**: Always validate and sanitize user input
3. **Error Messages**: Don't leak sensitive information in error messages
4. **Dependencies**: Keep dependencies updated and review them for vulnerabilities
5. **Code Review**: Security-sensitive changes require thorough review

## Known Security Considerations

### Manifest Storage

- Manifests are stored in `.gip/manifest/` in the repository
- They are plain JSON files with commit metadata
- They may contain developer comments and descriptions
- Consider `.gip/` in your `.gitignore` for private repositories

### Git Integration

- Gip executes Git commands using `os/exec`
- It inherits Git's security model
- Malicious repository hooks could potentially affect Gip

### TOON Serialization

- TOON format is used for display in conflict markers
- It's a text-based format derived from JSON
- No execution of code or scripts from TOON data

## Security Updates

Security updates will be:
- Released as patch versions
- Announced in [GitHub Security Advisories](https://github.com/iamHrithikRaj/gip/security/advisories)
- Documented in release notes
- Communicated to users via GitHub

## Security Audit

This project has not yet undergone a formal security audit. If you're interested in sponsoring or conducting a security audit, please contact the maintainers.

## Attribution

We appreciate security researchers who help keep Gip secure. With your permission, we will:
- Credit you in the security advisory
- List you in our security acknowledgments
- Link to your website or GitHub profile

## Questions?

For security-related questions that are not vulnerabilities, please:
- Open a [GitHub Discussion](https://github.com/iamHrithikRaj/gip/discussions)
- Ask in an issue with the `security` label

---

**Thank you for helping keep Gip and its users safe!** 🔒
