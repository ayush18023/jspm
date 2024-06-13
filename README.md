# JSPM (Javascript package manager)

A javascript package installer that installs specified package along with the version. This installer is written in golang. Its a toy project and not for production. The sole purpose of this project was to learn golang :).


**Current Features:**
- **Install:** installs the package speified by package@version if version is null it takes latest.
```
jspm install express@4.18.2
```
- **Purge:** This command clears the node_module directory.
```
jspm purge
```

**What it does?:**
- Dependency resolution
- Custom lexer parser to parse version i.e. semver (eg: ">=1.0.2 <5.2.3 | ~1.0.5 | ^2.0.1")
- Package installation
- Terminal ui 
- Add bianry scripts<br/><br/>

It is faster compared to **npm**, but there are issues need to be resolved and testing required to make it compatible.

