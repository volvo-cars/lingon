// CODE GENERATED BY github.com/volvo-cars/lingon. DO NOT EDIT.

package local

import "github.com/volvo-cars/lingon/pkg/terra"

func NewDataFile(name string, args DataFileArgs) *DataFile {
	return &DataFile{
		Args: args,
		Name: name,
	}
}

var _ terra.DataResource = (*DataFile)(nil)

type DataFile struct {
	Name string
	Args DataFileArgs
}

func (f *DataFile) DataSource() string {
	return "local_file"
}

func (f *DataFile) LocalName() string {
	return f.Name
}

func (f *DataFile) Configuration() interface{} {
	return f.Args
}

func (f *DataFile) Attributes() dataFileAttributes {
	return dataFileAttributes{ref: terra.ReferenceDataResource(f)}
}

type DataFileArgs struct {
	// Filename: string, required
	Filename terra.StringValue `hcl:"filename,attr" validate:"required"`
}
type dataFileAttributes struct {
	ref terra.Reference
}

func (f dataFileAttributes) Content() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content"))
}

func (f dataFileAttributes) ContentBase64() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content_base64"))
}

func (f dataFileAttributes) ContentBase64Sha256() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content_base64sha256"))
}

func (f dataFileAttributes) ContentBase64Sha512() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content_base64sha512"))
}

func (f dataFileAttributes) ContentMd5() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content_md5"))
}

func (f dataFileAttributes) ContentSha1() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content_sha1"))
}

func (f dataFileAttributes) ContentSha256() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content_sha256"))
}

func (f dataFileAttributes) ContentSha512() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("content_sha512"))
}

func (f dataFileAttributes) Filename() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("filename"))
}

func (f dataFileAttributes) Id() terra.StringValue {
	return terra.ReferenceString(f.ref.Append("id"))
}
