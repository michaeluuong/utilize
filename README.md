# utilize
Provides common utilities that are popular with no one except myself.  
While not specifically relegated to macOS, no amount of testing has occurred on any other operating systems.

- package ***filing*** collects functions that deal with the file system
  - ***dir.go*** gathers operations related to directories  
  - ***file.go*** provides file related operations
- package ***jsoning*** collects convenience methods for JSON data
  - ***jsonify.go*** provides a Jsonify object with methods that operate on or with JSON data   
- package ***reflections*** collects functions that ponder the dynamic nature of content
  - ***reflections.go*** provides functions that provide help with dynamic determination
    - InitializeStruct(sPtr any) error initializes any maps in a struct   
- package ***routine*** helps with concurrency (i.e. go routines)
  - ***go_route.go*** provides a simple semaphor
- package ***stringy*** collects convenient functions for string data
  - ***distance.go*** provides Levenshtein string distance functions built on [github.com/agnivade/levenshtein](github.com/agnivade/levenshtein)
  - ***string_case.go*** provides case conversion (e.g. case folding, sentence case, title, lower & upper case)
  - ***strings.go*** provides common string operations
