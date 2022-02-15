* Add 7zip extract support.
  * Make use of 7zip extractor in main.
* Add `git clone` support.
* Change `-a` to `-i` for answer short flag.
* .empty files must be removed from output
* Change excludes key to skipParsing in template.json
* Files in the "excludes" list must be included in the output without processing
* Remove answer file being required check
* Format output when asking for variable input.
* Move all messaging to various arrays (big tedious job, but centralized text make easier to translate).
