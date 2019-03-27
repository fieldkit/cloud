package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/goadesign/goa"

	"github.com/conservify/sqlxcache"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/backend/ingestion"
	"github.com/fieldkit/cloud/server/data"
)

func ExportAllFiles(ctx context.Context, response *goa.ResponseData, download bool, iterator *FileIterator, exporter Exporter) error {
	log := Logger(ctx).Sugar()

	header := false

	for {
		cs, err := iterator.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if cs == nil {
			break
		}

		if !header {
			if download {
				fileName := exporter.FileName(&cs.IteratorFile)
				response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
				response.Header().Set("Content-Type", exporter.DownloadMimeType())
			} else {
				response.Header().Set("Content-Disposition", fmt.Sprintf("inline"))
				response.Header().Set("Content-Type", exporter.MimeType())
			}
			header = true
		}

		defer cs.Response.Body.Close()

		binaryReader := NewFkBinaryReader(exporter.ForFile(&cs.IteratorFile))

		err = binaryReader.Read(ctx, cs.Response.Body)
		if err != nil {
			log.Infow("Error reading stream", "error", err, "file_type_id", cs.FileID)
		}
	}

	exporter.Finish(ctx)

	return nil
}

type Exporter interface {
	ForFile(file *IteratorFile) FormattedMessageReceiver
	Finish(ctx context.Context) error
	FileName(file *IteratorFile) string
	DownloadMimeType() string
	MimeType() string
}

type SimpleJsonExporter struct {
	File    *IteratorFile
	Writer  io.Writer
	Records int
}

func NewSimpleJsonExporter(writer io.Writer) Exporter {
	return &SimpleJsonExporter{Writer: writer}
}

func (ce *SimpleJsonExporter) DownloadMimeType() string {
	return "application/json; charset=utf-8"
}

func (ce *SimpleJsonExporter) MimeType() string {
	return ce.DownloadMimeType()
}

func (ce *SimpleJsonExporter) FileName(file *IteratorFile) string {
	return fmt.Sprintf("%s.json", file.FileID)
}

func (ce *SimpleJsonExporter) ForFile(file *IteratorFile) FormattedMessageReceiver {
	ce.File = file
	return ce
}

func (ce *SimpleJsonExporter) HandleRecord(ctx context.Context, r *pb.DataRecord) error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if ce.Records > 0 {
		fmt.Fprintf(ce.Writer, ",\n")
	} else {
		fmt.Fprintf(ce.Writer, "[\n")
	}

	_, err = ce.Writer.Write(body)
	if err != nil {
		return err
	}

	ce.Records += 1

	return nil
}

func (ce *SimpleJsonExporter) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) (*ingestion.RecordChange, error) {
	return nil, nil
}

func (ce *SimpleJsonExporter) Finish(ctx context.Context) error {
	fmt.Fprintf(ce.Writer, "\n]\n")
	return nil
}

type SimpleCsvExporter struct {
	File          *IteratorFile
	Writer        io.Writer
	HeaderWritten bool
}

func NewSimpleCsvExporter(writer io.Writer) Exporter {
	return &SimpleCsvExporter{Writer: writer}
}

func (ce *SimpleCsvExporter) DownloadMimeType() string {
	return "text/csv; charset=utf-8"
}

func (ce *SimpleCsvExporter) MimeType() string {
	return "text/plain; charset=utf-8"
}

func (ce *SimpleCsvExporter) FileName(file *IteratorFile) string {
	return fmt.Sprintf("%s.csv", file.FileID)
}

func (ce *SimpleCsvExporter) ForFile(file *IteratorFile) FormattedMessageReceiver {
	ce.File = file
	return ce
}

func (ce *SimpleCsvExporter) HandleRecord(ctx context.Context, r *pb.DataRecord) error {
	if r.Log != nil {
		return ce.ExportLog(ctx, r)
	}

	return nil
}

func (ce *SimpleCsvExporter) ExportLog(ctx context.Context, r *pb.DataRecord) error {
	if !ce.HeaderWritten {
		fmt.Fprintf(ce.Writer, "%v,%v,%v,%v,%v,%v,%v,%v\n", "DeviceID", "FileID", "FileID", "Uptime", "Time", "Level", "Facility", "Message")
		ce.HeaderWritten = true
	}

	fmt.Fprintf(ce.Writer, "%v,%v,%v,%v,%v,%v,%v,%v\n", ce.File.DeviceID, ce.File.FileID, r.Log.Uptime, r.Log.Time, r.Log.Level, r.Log.Facility, strings.TrimSpace(r.Log.Message))

	return nil
}

func (ce *SimpleCsvExporter) HandleFormattedMessage(ctx context.Context, fm *ingestion.FormattedMessage) (*ingestion.RecordChange, error) {
	opaqueKeys := reflect.ValueOf(fm.MapValues).MapKeys()
	keys := make([]string, len(opaqueKeys))
	for i := 0; i < len(opaqueKeys); i++ {
		keys[i] = opaqueKeys[i].String()
	}
	sort.Strings(keys)

	if !ce.HeaderWritten {
		fmt.Fprintf(ce.Writer, "%v,%v,%v,%v,%v,%v,%v,%v", "Device", "File", "File", "Message", "Time", "Longitude", "Latitude", "Fixed")

		for _, key := range keys {
			fmt.Fprintf(ce.Writer, ",%v", key)
		}

		fmt.Fprintf(ce.Writer, "\n")

		ce.HeaderWritten = true
	}

	fmt.Fprintf(ce.Writer, "%v,%v,", ce.File.DeviceID, ce.File.FileID)

	fmt.Fprintf(ce.Writer, "%v,%v,", fm.MessageId, fm.Time)

	if fm.Location != nil && len(fm.Location) >= 2 {
		fmt.Fprintf(ce.Writer, "%v,%v", fm.Location[0], fm.Location[1])
	} else {
		fmt.Fprintf(ce.Writer, "%v,%v", 0.0, 0.0)
	}

	fmt.Fprintf(ce.Writer, ",%v", fm.Fixed)

	for _, key := range keys {
		fmt.Fprintf(ce.Writer, ",%v", fm.MapValues[key])
	}

	fmt.Fprintf(ce.Writer, "\n")

	return nil, nil
}

func (ce *SimpleCsvExporter) Finish(ctx context.Context) error {
	return nil
}

type FileIterator struct {
	Session   *session.Session
	S3Service *s3.S3
	Offset    int
	Limit     int
	SignUrls  bool
	Query     func(ctx context.Context, s *FileIterator) error
	Files     []*IteratorFile
	Index     int
}

type IteratorFile struct {
	FileID   string
	DeviceID string
	Size     int64
	URL      string
}

type OpenedFile struct {
	IteratorFile
	SignedURL string
	Response  *http.Response
}

func toIteratorFiles(dfs []*data.DeviceFile) []*IteratorFile {
	ifs := make([]*IteratorFile, len(dfs))
	for i, f := range dfs {
		ifs[i] = &IteratorFile{
			FileID:   f.FileID,
			DeviceID: f.DeviceID,
			Size:     f.Size,
			URL:      f.URL,
		}
	}
	return ifs
}

func LookupFileOnS3(ctx context.Context, session *session.Session, db *sqlxcache.DB, fileID string) (iterator *FileIterator, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("File", "file_id", fileID)

	iterator = &FileIterator{
		Offset:    0,
		Limit:     0,
		Index:     0,
		SignUrls:  session.Config.Credentials != nil,
		Session:   session,
		S3Service: s3.New(session),
		Query: func(ctx context.Context, s *FileIterator) error {
			fileRepository, err := NewFileRepository(session, "fk-streams")
			if err != nil {
				return err
			}

			fi, err := fileRepository.Info(ctx, fileID)
			if err != nil {
				return err
			}

			s.Files = make([]*IteratorFile, 1)
			s.Files[0] = &IteratorFile{
				FileID:   fi.Key,
				DeviceID: "",
				Size:     fi.Size,
				URL:      fi.URL,
			}

			return nil
		},
	}

	return
}

func LookupFile(ctx context.Context, session *session.Session, db *sqlxcache.DB, fileID string) (iterator *FileIterator, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("File", "file_id", fileID)

	iterator = &FileIterator{
		Offset:    0,
		Limit:     0,
		Index:     0,
		SignUrls:  session.Config.Credentials != nil,
		Session:   session,
		S3Service: s3.New(session),
		Query: func(ctx context.Context, s *FileIterator) error {
			var files []*data.DeviceFile
			if err := db.SelectContext(ctx, &files, `SELECT s.* FROM fieldkit.device_stream AS s WHERE s.id = $1`, fileID); err != nil {
				return err
			}
			s.Files = toIteratorFiles(files)
			return nil
		},
	}

	return
}

func LookupDeviceFiles(ctx context.Context, session *session.Session, db *sqlxcache.DB, deviceID string, fileTypeIDs []string) (iterator *FileIterator, err error) {
	log := Logger(ctx).Sugar()

	log.Infow("File", "device_id", deviceID, "file_type_ids", fileTypeIDs)

	iterator = &FileIterator{
		Offset:    0,
		Limit:     10,
		Index:     0,
		SignUrls:  session.Config.Credentials != nil,
		Session:   session,
		S3Service: s3.New(session),
		Query: func(ctx context.Context, s *FileIterator) error {
			var files []*data.DeviceFile
			if err := db.SelectContext(ctx, &files,
				`SELECT s.* FROM fieldkit.device_stream AS s WHERE (s.file_id = ANY($1)) AND (s.device_id = $2) ORDER BY time DESC LIMIT $3 OFFSET $4`,
				fileTypeIDs, deviceID, s.Limit, s.Offset); err != nil {
				return err
			}
			s.Files = toIteratorFiles(files)
			return nil
		},
	}

	return
}

func (iterator *FileIterator) Next(ctx context.Context) (cs *OpenedFile, err error) {
	log := Logger(ctx).Sugar()

	if iterator.Files != nil && iterator.Index >= len(iterator.Files) {
		iterator.Files = nil

		// Easy way to handle the single stream per batch.
		if iterator.Limit == 0 {
			log.Infow("No more batches")
			return nil, io.EOF
		}
	}

	if iterator.Files == nil {
		err = iterator.Query(ctx, iterator)
		if err != nil {
			return nil, err
		}

		log.Infow("Queried", "batch_size", len(iterator.Files))

		if len(iterator.Files) == 0 {
			log.Infow("No more batches")
			return nil, io.EOF
		}

		iterator.Index = 0
	}

	iteratorFile := iterator.Files[iterator.Index]

	log.Infow("File", "file_type_id", iteratorFile.FileID, "index", iterator.Index, "size", iteratorFile.Size, "url", iteratorFile.URL)

	iterator.Index += 1

	signed := iteratorFile.URL

	if iterator.SignUrls {
		signed, err = SignS3URL(iterator.S3Service, iteratorFile.URL)
		if err != nil {
			return nil, fmt.Errorf("Error signing stream URL: %v (%v)", iteratorFile.URL, err)
		}
		log.Infow("File", "file_type_id", iteratorFile.FileID, "signed_url", signed)
	}

	response, err := http.Get(signed)
	if err != nil {
		return nil, fmt.Errorf("Error opening stream URL: %v (%v)", iteratorFile.URL, err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Error retrieving stream: %v (status = %v)", signed, response.StatusCode)
	}

	cs = &OpenedFile{
		IteratorFile: *iteratorFile,
		SignedURL:    signed,
		Response:     response,
	}

	return
}
