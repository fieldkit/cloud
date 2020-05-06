package common

import (
	"strings"
)

const (
	FkDataBinaryContentType    = "application/vnd.fk.data+binary"
	FkDataBase64ContentType    = "application/vnd.fk.data+base64"
	MultiPartFormDataMediaType = "multipart/form-data"
	ContentTypeHeaderName      = "Content-Type"
	ContentLengthHeaderName    = "Content-Length"
	XForwardedForHeaderName    = "X-Forwarded-For"
	FkDeviceIdHeaderName       = "Fk-DeviceId"
	FkDeviceNameHeaderName     = "Fk-DeviceName"
	FkGenerationHeaderName     = "Fk-Generation"
	FkBlocksHeaderName         = "Fk-Blocks"
	FkFlagsIdHeaderName        = "Fk-Flags"
	FkTypeHeaderName           = "Fk-Type"

	// Deprecated
	FkFileIdHeaderName      = "Fk-FileId"
	FkVersionHeaderName     = "Fk-Version"
	FkFileNameHeaderName    = "Fk-FileName"
	FkFileVersionHeaderName = "Fk-FileVersion"
	FkBuildHeaderName       = "Fk-Build"
	FkFileOffsetHeaderName  = "Fk-FileOffset"
	FkCompiledHeaderName    = "Fk-Compiled"
	FkUploadNameHeaderName  = "Fk-UploadName"
)

func SanitizeMeta(m map[string]*string) map[string]*string {
	ci := make(map[string]*string)
	for key, value := range m {
		ci[strings.ToLower(key)] = value
	}
	newM := make(map[string]*string)
	for _, key := range []string{FkDeviceIdHeaderName, FkGenerationHeaderName, FkBlocksHeaderName, FkFlagsIdHeaderName, FkFileIdHeaderName, FkBuildHeaderName, FkFileNameHeaderName, FkVersionHeaderName} {
		if value, ok := ci[strings.ToLower(key)]; ok {
			newM[key] = value
		}
	}
	return newM
}
