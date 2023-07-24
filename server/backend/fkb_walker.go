package backend

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/proto"

	pb "github.com/fieldkit/data-protocol"

	"github.com/fieldkit/cloud/server/common/logging"

	"github.com/fieldkit/cloud/server/files"
)

type WalkProgress struct {
	read int64
	file int64
}

type OnWalkProgress func(ctx context.Context, progress WalkProgress) error

type FkbWalker struct {
	verbose   bool
	handler   FkbHandler
	progress  OnWalkProgress
	files     files.FileArchive
	metrics   *logging.Metrics
	queue     []*parsedRecord
	keyRecord *keyRecord
	meta      *pb.DataRecord
}

type FkbHandler interface {
	OnSignedMeta(ctx context.Context, signedRecord *pb.SignedRecord, rawRecord *pb.DataRecord, bytes []byte) error
	OnMeta(ctx context.Context, recordNumber int64, rawMeta *pb.DataRecord, bytes []byte) error
	OnData(ctx context.Context, rawData *pb.DataRecord, rawMeta *pb.DataRecord, bytes []byte) error
	OnDone(ctx context.Context) error
}

func NewFkbWalker(files files.FileArchive, metrics *logging.Metrics, handler FkbHandler, progress OnWalkProgress, verbose bool) (ra *FkbWalker) {
	return &FkbWalker{
		verbose:   verbose,
		progress:  progress,
		files:     files,
		metrics:   metrics,
		handler:   handler,
		keyRecord: nil,
		meta:      nil,
	}
}

func (ra *FkbWalker) tryParseSignedRecord(sr *pb.SignedRecord, dataRecord *pb.DataRecord) (bytes []byte, err error) {
	if err := proto.Unmarshal(sr.Data, dataRecord); err == nil {
		return sr.Data, nil
	}

	buffer := proto.NewBuffer(sr.Data)
	size, err := buffer.DecodeVarint()
	if err != nil {
		return nil, err
	}

	sizeSize := uint64(proto.SizeVarint(size))

	if size+sizeSize != uint64(len(sr.Data)) {
		return nil, fmt.Errorf("bad length prefix in signed record: %d", size)
	}

	if err := buffer.Unmarshal(dataRecord); err != nil {
		return nil, err
	}

	bytes = sr.Data[sizeSize:]

	return bytes, nil
}

func (ra *FkbWalker) onSignedMeta(ctx context.Context, signedRecord *pb.SignedRecord, rawRecord *pb.DataRecord, bytes []byte) error {
	if err := ra.handler.OnSignedMeta(ctx, signedRecord, rawRecord, bytes); err != nil {
		return err
	}

	return nil
}

func (ra *FkbWalker) onMeta(ctx context.Context, recordNumber int64, rawRecord *pb.DataRecord, bytes []byte) error {
	if err := ra.handler.OnMeta(ctx, recordNumber, rawRecord, bytes); err != nil {
		return err
	}

	ra.meta = rawRecord

	return nil
}

func (ra *FkbWalker) onData(ctx context.Context, rawRecord *pb.DataRecord, bytes []byte) error {
	if err := ra.handler.OnData(ctx, rawRecord, ra.meta, bytes); err != nil {
		return err
	}

	return nil
}

func (ra *FkbWalker) flushQueue(ctx context.Context, required bool) (fatal error) {
	log := Logger(ctx).Sugar()

	// Do we have queued meta records in need of a record number?
	if len(ra.queue) == 0 {
		return nil
	}

	// If we're missing a key record, we can't import this file.
	if ra.keyRecord == nil {
		if !required {
			// Right now the only way this can happen is if we have a partial
			// file with no data and just meta records.
			log.Warnw("fkb:unrequired-flush no-key-record")
			return nil
		}

		return fmt.Errorf("fkb:flush error: no key record")
	}

	// So we can safetly process the queued meta records using their
	// file offset for the record number.
	for _, queued := range ra.queue {
		if err := ra.onMeta(ctx, ra.keyRecord.adjustRecordNumber(queued.FileOffset), queued.DataRecord, queued.Bytes); err != nil {
			return err
		}
	}

	log.Infow("fkb:flushed", "nrecords", len(ra.queue))

	ra.queue = make([]*parsedRecord, 0)

	return nil
}

func (ra *FkbWalker) handle(ctx context.Context, pr *parsedRecord) (warning error, fatal error) {
	log := Logger(ctx).Sugar()

	if pr.SignedRecord != nil {
		if err := ra.onSignedMeta(ctx, pr.SignedRecord, pr.DataRecord, pr.Bytes); err != nil {
			return nil, err
		}
	} else if pr.DataRecord != nil {
		// Older versions of firmware were using a nanopb version that wouldn't
		// allow us to elide empty child messages, so we need to check for a
		// field.
		if pr.DataRecord.Metadata != nil && pr.DataRecord.Metadata.DeviceId != nil {
			// Check to see if this record has a record number. Some older
			// firmware was never setting this which causes some headaches. What
			// we do is queue this up, and hope that we come across a data
			// record with a record number to anchor things. Notice that
			// non-single file meta records will come in via SignedMetaRecords.
			record := int64(pr.DataRecord.Metadata.Record)
			if record == 0 || len(ra.queue) > 0 {
				log.Infow("fkb:queue-meta")

				ra.queue = append(ra.queue, pr)
			} else {
				// I'm wasn't sure if this is entirely necessary and added this
				// in a better safe than sorry mentality and it seems to have
				// been a mistake. Right now, the only way to end up here is if
				// record is greater than 0, which will indicate a valid record
				// number. So, there's no need to adjust it.
				/*
					if ra.keyRecord != nil {
						if err := ra.onMeta(ctx, ra.keyRecord.adjustRecordNumber(record), pr.DataRecord, pr.Bytes); err != nil {
							return nil, err
						}
					} else {
						if err := ra.onMeta(ctx, record, pr.DataRecord, pr.Bytes); err != nil {
							return nil, err
						}
					}
				*/
				if err := ra.onMeta(ctx, record, pr.DataRecord, pr.Bytes); err != nil {
					return nil, err
				}
			}
		} else {
			// Are we need of a key record and does this record have one? Notice
			// that data readings can't be the first record, so looking for
			// non-zero readings here is safe.
			if ra.keyRecord == nil && pr.DataRecord.Readings.Reading > 0 {
				ra.keyRecord = &keyRecord{
					fileOffset:   pr.FileOffset,
					recordNumber: int64(pr.DataRecord.Readings.Reading),
				}
				log.Infow("key-record", "file_offset", ra.keyRecord.fileOffset, "record_number", ra.keyRecord.recordNumber)
			}

			if err := ra.flushQueue(ctx, true); err != nil {
				return nil, err
			}

			if err := ra.onData(ctx, pr.DataRecord, pr.Bytes); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

type ProgressReader struct {
	io.ReadCloser
	ctx    context.Context
	notify OnWalkProgress
	file   int64
}

func (r *ProgressReader) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	r.file += int64(n)

	if r.notify != nil {
		if err := r.notify(r.ctx, WalkProgress{
			read: int64(n),
			file: r.file,
		}); err != nil {
			return n, err
		}
	}

	return n, err
}

func (r *ProgressReader) Close() error {
	return r.ReadCloser.Close()
}

func (ra *FkbWalker) WalkUrl(ctx context.Context, url string) (info *WalkRecordsInfo, err error) {
	log := Logger(ctx).Sugar().With("url", url)

	opened, err := ra.files.OpenByURL(ctx, url)
	if err != nil {
		return nil, err
	}

	reader := &ProgressReader{opened.Body, ctx, ra.progress, 0}

	defer reader.Close()

	meta := false
	data := false

	totalRecords := 0
	dataProcessed := 0
	dataErrors := 0
	metaProcessed := 0
	metaErrors := 0
	records := 0
	recordNumber := int64(0)

	unmarshalFunc := UnmarshalFunc(func(b []byte) (proto.Message, error) {
		var unmarshalError error
		var dataRecord pb.DataRecord
		var signedRecord pb.SignedRecord

		totalRecords += 1

		if data || (!data && !meta) {
			err := proto.Unmarshal(b, &dataRecord)
			if err != nil {
				if data { // If we expected this record, return the error
					ra.metrics.DataErrorsParsing()
					return nil, fmt.Errorf("fkb:error parsing data record: %w", err)
				}
				unmarshalError = err
			} else {
				data = true
				warning, fatal := ra.handle(ctx, &parsedRecord{
					FileOffset: recordNumber,
					DataRecord: &dataRecord,
					Bytes:      b,
				})
				if fatal != nil {
					return nil, fatal
				}
				if warning == nil {
					dataProcessed += 1
					records += 1
				} else {
					dataErrors += 1
					if records > 0 {
						log.Infow("fkb:processed", "record_run", records)
						records = 0
					}
				}
			}
		}

		if meta || (!data && !meta) {
			err := proto.Unmarshal(b, &signedRecord)
			if err != nil {
				if meta { // If we expected this record, return the error
					ra.metrics.DataErrorsParsing()
					return nil, fmt.Errorf("fkb:error parsing signed record: %w", err)
				}
				unmarshalError = err
			} else {
				meta = true
				bytes, err := ra.tryParseSignedRecord(&signedRecord, &dataRecord)
				if err != nil {
					ra.metrics.DataErrorsParsing()
					return nil, fmt.Errorf("fkb:error parsing signed record: %w", err)
				}

				warning, fatal := ra.handle(ctx, &parsedRecord{
					FileOffset:   recordNumber,
					SignedRecord: &signedRecord,
					DataRecord:   &dataRecord,
					Bytes:        bytes,
				})
				if fatal != nil {
					return nil, fatal
				}
				if warning == nil {
					metaProcessed += 1
					records += 1
				} else {
					metaErrors += 1
					if records > 0 {
						log.Infow("fkb:processed", "record_run", records)
						records = 0
					}
				}
			}
		}

		// Only used with the new single file uploads.
		recordNumber += 1

		// Parsing either one failed, otherwise this would have been set. So return the error.
		if !meta && !data {
			return nil, unmarshalError
		}

		return nil, nil
	})

	if _, _, err = ReadLengthPrefixedCollection(ctx, MaximumDataRecordLength, reader, unmarshalFunc); err != nil {
		newErr := fmt.Errorf("fkb:error: %w", err)
		log.Errorw("fkb:error", "error", newErr)
		return nil, newErr
	}

	if err := ra.flushQueue(ctx, false); err != nil {
		return nil, err
	}

	withInfo := log.With("meta_processed", metaProcessed,
		"data_processed", dataProcessed,
		"meta_errors", metaErrors,
		"data_errors", dataErrors,
		"record_run", records,
	)

	withInfo.Infow("fkb:processed")

	if err := ra.handler.OnDone(ctx); err != nil {
		return nil, fmt.Errorf("fkb:done handler error: %w", err)
	}

	info = &WalkRecordsInfo{
		TotalRecords: int64(totalRecords),
		MetaRecords:  int64(metaProcessed),
		DataRecords:  int64(dataProcessed),
		MetaErrors:   int64(metaErrors),
		DataErrors:   int64(dataErrors),
	}

	withInfo.Infow("fkb:done")

	return
}

type WalkRecordsInfo struct {
	TotalRecords int64
	DataRecords  int64
	MetaRecords  int64
	MetaErrors   int64
	DataErrors   int64
}

type keyRecord struct {
	fileOffset   int64
	recordNumber int64
}

func (kr *keyRecord) adjustRecordNumber(fileOffset int64) int64 {
	return kr.recordNumber + (fileOffset - kr.fileOffset)
}

type parsedRecord struct {
	SignedRecord *pb.SignedRecord
	DataRecord   *pb.DataRecord
	Bytes        []byte
	FileOffset   int64
}

func prettyTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.String()
}
