* Add `git clone` support. Maybe use `go get` as it can handle caching.
* Add 7zip extract support.
  * Make use of 7zip extractor in main.
* Change `-a` to `-i` for answer short flag.
* .empty files must be removed from output
* Change excludes key to skipParsing in template.json
* Files in the "excludes" list must be included in the output without processing
* Remove answer file being required check, Make supplying an answer file optional.
* Format output when asking for variable input.
* Move all messaging to various arrays (big tedious job, but centralized text make easier to translate).
* Look for URL as the first argument and output path as the second.
* Remove use "IncludeFileExtensions" and "AllowedUrls" from config.
* Add cache dir to config.
* add setting sub-command to set the cache dir.
* Remove "running program tmpltoapp.exe" from verbosity output.
* Add version label to version verbosity output.
* Add output dir to verbosity output.
* Remove config from verbosity output.
* Make current directory the default parent output directory.
* Add setup sub-command to setup the config.
* Should append version to cached downloads.
* Show in verbosity if using cache or downloading a fresh copy.
