
# Algorithm Outline

1. Load settings from a file:
   1. Yes, then go to the next step.
   2. No, then load defaults, display message indicatin defaults were loaded, 
      then go to the next step.
2. Check required arguments are set:
   1. Yes, parse and go to the next step.
   2. No, stop and output an error.
3. Validate path/URL is allowed:
   1. Yes. go to next step.
   2. No, stop and output an error.
4. Detect if local path or URL:
   1. If local, then go to "Copy" step.
   2. If URL, then go to "Download" step.
5. Download template zip from URL:
   1. Download template zip to the download cache.
   2. Perform a checksum on the template zip.
   2. Extract the template zip to temp cache.
6. Copy template to the destination local path:
   1. Yes, go to the next step.
   2. No, stop and output error.
7. Recursively process all files in the app path:
   1. On success, output a friendly message indicating stats and let the user know they are all set.
   2. On failure, stop and output error.
