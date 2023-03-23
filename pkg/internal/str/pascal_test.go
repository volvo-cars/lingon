/*
The MIT License (MIT)

Copyright (c) 2016 Segment
*/
package str

import (
	"testing"
)

func TestPascalCase(t *testing.T) {
	samples := []TT{
		{"sample text", "SampleText"},
		{"sample-text", "SampleText"},
		{"sample_text", "SampleText"},
		{"sample___text", "SampleText"},
		{"sampleText", "SampleText"},
		{"inviteYourCustomersAddInvites", "InviteYourCustomersAddInvites"},
		{"sample 2 Text", "Sample2Text"},
		{"   sample   2    Text   ", "Sample2Text"},
		{"   $#$sample   2    Text   ", "Sample2Text"},
		{"SAMPLE 2 TEXT", "Sample2Text"},
		{"___$$Base64Encode", "Base64Encode"},
		{"FOO:BAR$BAZ", "FooBarBaz"},
		{"FOO#BAR#BAZ", "FooBarBaz"},
		{"something.com", "SomethingCom"},
		{"$something%", "Something"},
		{"something.com", "SomethingCom"},
		{"•¶§ƒ˚foo˙∆˚¬", "Foo"},
	}

	for _, sample := range samples {
		if out := PascalCase(sample.str); out != sample.out {
			t.Errorf("got %q from %q, expected %q", out, sample.str, sample.out)
		}
	}
}

func BenchmarkPascalCase(t *testing.B) {
	for i := 0; i < t.N; i++ {
		PascalCase("some sample text here_noething:too$amazing")
	}
}
