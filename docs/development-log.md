# 7zip Notes

* On Windows 10, it seems that when you install 7zip, the 7z.exe is NOT in the
  patch by default.
* 7zip (x64bit) version has been observed as installed in
  `C:\Program Files\7-Zip"`, 32bit will not be tested at this time.
    * When testing to see if 7zip is installed on a system, check for the
      output of `7z` and the exit code.
    * If the previous test fails on Windows, then check the default install
      path for the 64bit version on Windows. Quit and return false if exit code
      is a failure.
    * If the previous test fails on Linux/Mac, then quit and return false for
      failure to find.
* I was not able to see the documentation for 7zip on `https://www.7-zip.org/`
  however, after installing on Windows 10, I am able to review the CLI
  documentation in the `7-zip.chm` or Window help file.
* It looks like the command we need to start coding is:
  ```
  7z e archive.zip
  ```
  The documentation stats "extracts all files from archive archive.zip to
  the current directory."
* Since the template should be the only object at the root of the
  `archive.zip`; the program would have to work in the `cache/templates/`
  directory.
* The programs should specify where to extract the contents, this can be done
  with the `-e` flag.
* To have this work with passwords, we should allow the user to supply a
  password.
