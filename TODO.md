* WIP: Add 7zip extract support.
  * Make use of 7zip extractor in main.
* Change excludes key to skipParsing in template.json
* Files in the "excludes" list must be included in the output without processing
* Format output when asking for variable input.
* WIP: Move all messaging to various arrays (big tedious job, but centralized text make easier to translate).
* Add version label to version verbosity output.
* Add output dir to verbosity output.
* Make current directory the default parent output directory.
* Should append version to cached downloads.
* Show in verbosity if using cache or downloading a fresh copy.
* Add basic (input regex, bool, int, unsigned) validation for placeholder.
* Add command to generate a template.json.
