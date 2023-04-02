schema_version = 1

project {
  license        = "Apache-2.0"
  copyright_year = 2023
  # Add the year into the copyright holder to force it being added
  copyright_holder = "2023 Volvo Car Corporation"

  header_ignore = [
    # Ignore all testdata and generated files
    ".idea/**",
    ".vscode/**",
    "bin/**",
    "**/.lingon/**",
    "**/.terra/**",
    "**/testdata/**",
    "**/out/**",
     "vendors/**",
  ]
}
