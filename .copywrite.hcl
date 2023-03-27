schema_version = 1

project {
  license        = "Apache-2.0"
  copyright_year = 2023
  copyright_holder = "Volvo Car Corporation"

  header_ignore = [
    # Ignore all testdata and generated files
    ".idea/**",
    ".vscode/**",
    "bin/**",
    "**/.lingon/**",
    "**/testdata/**",
    "**/out/**",
     "vendors/**",
  ]
}
