# Test Coverage

- The primary test coverage should be acctest, unit test should only be used to cover complex logic with well defined input and output (for example: complex regular expression validation)
- Unit test should assist human reviewer in understanding the logic, if the logic is already obvious, avoid adding one
- Limit number of test cases to only scenarios likely to be encountered by users in real-world usage
- Avoid very long unit tests that is hard to read, for example: mocking http transport and invoking the resource lifecycle