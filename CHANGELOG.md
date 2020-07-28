<a name="unreleased"></a>
## [Unreleased]


<a name="v0.2.0"></a>
## [v0.2.0] - 2020-07-28
### Build
- fix makefile test command
- fix build-zip command in makefile

### Doc
- updated docs
- add asciicast examples

### Feat
- improve consistency of the json reply format
- make massage paramters compatible with plzpy

### Misc
- minor dependency update

### BREAKING CHANGE

change the layout of json reply

The old format for historical data was [{year:value}, {year,value}, ..]
This change makes the output format more verbose but easier to deal with:
[{'year':year, 'count': value}...]
Also rename the 'count' label to 'total' for aggregated count

the cli massage paramters (--in,--out) are now called (--input,--output)


<a name="v0.1.0"></a>
## v0.1.0 - 2020-07-28
### Build
- add docker support
- add makefile and config for changelog gen

### Doc
- add readme

### Feat
- add command serve to expose  rest API endpoints
- add command to process src dataset

### Fix
- welcome banner broken using logging

### Test
- add testing for dataset ETL


[Unreleased]: https://github.com/noandrea/plz/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/noandrea/plz/compare/v0.1.0...v0.2.0
