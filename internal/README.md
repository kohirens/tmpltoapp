# Internal packages

Sometimes one wishes to have components that are not exported, for instance to
avoid acquiring clients of interfaces to code that is part of a public
repository but not intended for use outside the program to which it belongs.

To create such a package, place it in a directory named internal or in a
subdirectory of a directory named internal.

See [Internal packages] for details.

---

[Internal packages]: https://go.dev/doc/go1.4#internalpackages
